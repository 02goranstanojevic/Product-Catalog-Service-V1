package e2e

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	databasepb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/product-catalog-service/internal/app/product/domain"
	"github.com/product-catalog-service/internal/app/product/queries/get_product"
	"github.com/product-catalog-service/internal/app/product/queries/list_products"
	"github.com/product-catalog-service/internal/app/product/repo"
	"github.com/product-catalog-service/internal/app/product/usecases/activate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/apply_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/archive_product"
	"github.com/product-catalog-service/internal/app/product/usecases/create_product"
	"github.com/product-catalog-service/internal/app/product/usecases/deactivate_product"
	"github.com/product-catalog-service/internal/app/product/usecases/remove_discount"
	"github.com/product-catalog-service/internal/app/product/usecases/update_product"
	"github.com/product-catalog-service/internal/pkg/clock"
	"github.com/product-catalog-service/internal/pkg/committer"
)

const (
	projectID  = "test-project"
	instanceID = "test-instance"
)

var (
	spannerClient *spanner.Client
	testDBPath    string
)

func TestMain(m *testing.M) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		fmt.Println("SPANNER_EMULATOR_HOST not set, skipping e2e tests")
		os.Exit(0)
	}

	fmt.Printf("[e2e] emulator at %s\n", os.Getenv("SPANNER_EMULATOR_HOST"))

	setupCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	dbID := fmt.Sprintf("testdb-%s", uuid.New().String()[:8])
	testDBPath = fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbID)
	fmt.Printf("[e2e] test database: %s\n", testDBPath)

	fmt.Println("[e2e] setting up instance...")
	if err := setupInstance(setupCtx); err != nil {
		fmt.Printf("[e2e] FATAL instance setup: %v\n", err)
		cancel()
		os.Exit(1)
	}
	fmt.Println("[e2e] instance OK")

	fmt.Println("[e2e] setting up database...")
	if err := setupDatabase(setupCtx, dbID); err != nil {
		fmt.Printf("[e2e] FATAL database setup: %v\n", err)
		cancel()
		os.Exit(1)
	}
	cancel()
	fmt.Println("[e2e] database OK")

	fmt.Println("[e2e] creating spanner client...")
	var err error
	spannerClient, err = spanner.NewClient(context.Background(), testDBPath)
	if err != nil {
		fmt.Printf("[e2e] FATAL spanner client: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[e2e] ready - running tests")

	code := m.Run()
	done := make(chan struct{})
	go func() { spannerClient.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	os.Exit(code)
}

// closeAdmin closes a client with a short deadline so it never blocks the process.
func closeAdmin(c interface{ Close() error }) {
	done := make(chan struct{})
	go func() {
		_ = c.Close()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
}

func setupInstance(ctx context.Context) error {
	fmt.Println("[e2e]   creating instance admin client...")
	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("instance admin client: %w", err)
	}

	fmt.Println("[e2e]   sending CreateInstance...")
	op, err := adminClient.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/emulator-config", projectID),
			DisplayName: "Test Instance",
			NodeCount:   1,
		},
	})
	if err != nil {
		// Instance already exists - that's fine.
		fmt.Println("[e2e]   instance exists, reusing")
		closeAdmin(adminClient)
		return nil
	}
	fmt.Println("[e2e]   waiting for instance operation...")
	_, _ = op.Wait(ctx) // ignore error - instance is usable
	closeAdmin(adminClient)
	return nil
}

func setupDatabase(ctx context.Context, dbID string) error {
	fmt.Println("[e2e]   creating database admin client...")
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("database admin client: %w", err)
	}

	fmt.Printf("[e2e]   sending CreateDatabase %s...\n", dbID)
	op, err := adminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", dbID),
		ExtraStatements: []string{
			`CREATE TABLE products (
				product_id STRING(36) NOT NULL,
				name STRING(255) NOT NULL,
				description STRING(MAX),
				category STRING(100) NOT NULL,
				base_price_numerator INT64 NOT NULL,
				base_price_denominator INT64 NOT NULL,
				discount_percent NUMERIC,
				discount_start_date TIMESTAMP,
				discount_end_date TIMESTAMP,
				status STRING(20) NOT NULL,
				created_at TIMESTAMP NOT NULL,
				updated_at TIMESTAMP NOT NULL,
				archived_at TIMESTAMP,
			) PRIMARY KEY (product_id)`,
			`CREATE TABLE outbox_events (
				event_id STRING(36) NOT NULL,
				event_type STRING(100) NOT NULL,
				aggregate_id STRING(36) NOT NULL,
				payload JSON NOT NULL,
				status STRING(20) NOT NULL,
				created_at TIMESTAMP NOT NULL,
				processed_at TIMESTAMP,
			) PRIMARY KEY (event_id)`,
			`CREATE INDEX idx_outbox_status ON outbox_events(status, created_at)`,
			`CREATE INDEX idx_products_category ON products(category, status)`,
		},
	})
	if err != nil {
		closeAdmin(adminClient)
		return fmt.Errorf("create database: %w", err)
	}
	fmt.Println("[e2e]   waiting for database operation...")
	if _, err := op.Wait(ctx); err != nil {
		closeAdmin(adminClient)
		return fmt.Errorf("database wait: %w", err)
	}
	closeAdmin(adminClient)
	return nil
}

type testEnv struct {
	ctx               context.Context
	clock             *clock.MockClock
	createProduct     *create_product.Interactor
	updateProduct     *update_product.Interactor
	activateProduct   *activate_product.Interactor
	deactivateProduct *deactivate_product.Interactor
	archiveProduct    *archive_product.Interactor
	applyDiscount     *apply_discount.Interactor
	removeDiscount    *remove_discount.Interactor
	getProduct        *get_product.Query
	listProducts      *list_products.Query
}

func newTestEnv() *testEnv {
	ctx := context.Background()
	clk := clock.NewMock(time.Date(2026, 2, 16, 12, 0, 0, 0, time.UTC))
	comm := committer.New(spannerClient)
	productRepo := repo.NewProductRepo(clk)
	outboxRepo := repo.NewOutboxRepo(clk)
	readModel := repo.NewProductReadModel(clk)

	return &testEnv{
		ctx:               ctx,
		clock:             clk,
		createProduct:     create_product.New(productRepo, outboxRepo, comm, clk),
		updateProduct:     update_product.New(productRepo, outboxRepo, comm, spannerClient),
		activateProduct:   activate_product.New(productRepo, outboxRepo, comm, spannerClient),
		deactivateProduct: deactivate_product.New(productRepo, outboxRepo, comm, spannerClient),
		archiveProduct:    archive_product.New(productRepo, outboxRepo, comm, spannerClient, clk),
		applyDiscount:     apply_discount.New(productRepo, outboxRepo, comm, spannerClient, clk),
		removeDiscount:    remove_discount.New(productRepo, outboxRepo, comm, spannerClient),
		getProduct:        get_product.New(readModel, spannerClient),
		listProducts:      list_products.New(readModel, spannerClient),
	}
}

func (e *testEnv) createActiveProduct(t *testing.T) string {
	t.Helper()
	id, err := e.createProduct.Execute(e.ctx, create_product.Request{
		Name:        "Test Product",
		Description: "A test product",
		Category:    "electronics",
		Numerator:   1999,
		Denominator: 100,
	})
	require.NoError(t, err)

	err = e.activateProduct.Execute(e.ctx, activate_product.Request{ProductID: id})
	require.NoError(t, err)
	return id
}

func TestProductCreationFlow(t *testing.T) {
	env := newTestEnv()

	productID, err := env.createProduct.Execute(env.ctx, create_product.Request{
		Name:        "Test Product",
		Description: "A test product",
		Category:    "electronics",
		Numerator:   1999,
		Denominator: 100,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, productID)

	product, err := env.getProduct.Execute(env.ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, "electronics", product.Category)
	assert.Equal(t, "19.99", product.BasePrice)
	assert.Equal(t, "DRAFT", product.Status)

	events := getOutboxEvents(t, env.ctx, productID)
	require.GreaterOrEqual(t, len(events), 1)
	assert.Equal(t, "product.created", events[0].EventType)
}

func TestProductUpdateFlow(t *testing.T) {
	env := newTestEnv()
	productID := env.createActiveProduct(t)

	err := env.updateProduct.Execute(env.ctx, update_product.Request{
		ProductID:   productID,
		Name:        "Updated Product",
		Description: "Updated description",
		Category:    "gadgets",
	})
	require.NoError(t, err)

	product, err := env.getProduct.Execute(env.ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Product", product.Name)
	assert.Equal(t, "Updated description", product.Description)
	assert.Equal(t, "gadgets", product.Category)
}

func TestDiscountApplicationFlow(t *testing.T) {
	env := newTestEnv()
	productID := env.createActiveProduct(t)

	now := env.clock.Now()
	err := env.applyDiscount.Execute(env.ctx, apply_discount.Request{
		ProductID:  productID,
		Percentage: 20,
		StartDate:  now.Add(-time.Hour),
		EndDate:    now.Add(24 * time.Hour),
	})
	require.NoError(t, err)

	product, err := env.getProduct.Execute(env.ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "15.99", product.EffectivePrice)
	assert.NotNil(t, product.DiscountPercent)

	events := getOutboxEvents(t, env.ctx, productID)
	hasDiscountEvent := false
	for _, e := range events {
		if e.EventType == "discount.applied" {
			hasDiscountEvent = true
			break
		}
	}
	assert.True(t, hasDiscountEvent)
}

func TestRemoveDiscountFlow(t *testing.T) {
	env := newTestEnv()
	productID := env.createActiveProduct(t)

	now := env.clock.Now()
	_ = env.applyDiscount.Execute(env.ctx, apply_discount.Request{
		ProductID:  productID,
		Percentage: 20,
		StartDate:  now.Add(-time.Hour),
		EndDate:    now.Add(24 * time.Hour),
	})

	err := env.removeDiscount.Execute(env.ctx, remove_discount.Request{ProductID: productID})
	require.NoError(t, err)

	product, err := env.getProduct.Execute(env.ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, product.BasePrice, product.EffectivePrice)
}

func TestProductActivationDeactivation(t *testing.T) {
	env := newTestEnv()

	productID, err := env.createProduct.Execute(env.ctx, create_product.Request{
		Name:        "Status Test",
		Description: "desc",
		Category:    "test",
		Numerator:   500,
		Denominator: 100,
	})
	require.NoError(t, err)

	err = env.activateProduct.Execute(env.ctx, activate_product.Request{ProductID: productID})
	require.NoError(t, err)

	product, _ := env.getProduct.Execute(env.ctx, productID)
	assert.Equal(t, "ACTIVE", product.Status)

	err = env.deactivateProduct.Execute(env.ctx, deactivate_product.Request{ProductID: productID})
	require.NoError(t, err)

	product, _ = env.getProduct.Execute(env.ctx, productID)
	assert.Equal(t, "INACTIVE", product.Status)
}

func TestProductArchiveFlow(t *testing.T) {
	env := newTestEnv()
	productID := env.createActiveProduct(t)

	err := env.archiveProduct.Execute(env.ctx, archive_product.Request{ProductID: productID})
	require.NoError(t, err)

	product, err := env.getProduct.Execute(env.ctx, productID)
	require.NoError(t, err)
	assert.Equal(t, "ARCHIVED", product.Status)
	assert.NotEmpty(t, product.UpdatedAt)

	err = env.activateProduct.Execute(env.ctx, activate_product.Request{ProductID: productID})
	assert.ErrorIs(t, err, domain.ErrProductArchived)
}

func TestBusinessRuleValidation_DiscountOnInactiveProduct(t *testing.T) {
	env := newTestEnv()

	productID, _ := env.createProduct.Execute(env.ctx, create_product.Request{
		Name:        "Draft Product",
		Description: "desc",
		Category:    "test",
		Numerator:   1000,
		Denominator: 100,
	})

	now := env.clock.Now()
	err := env.applyDiscount.Execute(env.ctx, apply_discount.Request{
		ProductID:  productID,
		Percentage: 10,
		StartDate:  now.Add(-time.Hour),
		EndDate:    now.Add(24 * time.Hour),
	})
	assert.ErrorIs(t, err, domain.ErrProductNotActive)
}

func TestListProducts(t *testing.T) {
	env := newTestEnv()

	for i := 0; i < 3; i++ {
		id, _ := env.createProduct.Execute(env.ctx, create_product.Request{
			Name:        fmt.Sprintf("List Product %d", i),
			Description: "desc",
			Category:    "list-test",
			Numerator:   int64(1000 + i*100),
			Denominator: 100,
		})
		_ = env.activateProduct.Execute(env.ctx, activate_product.Request{ProductID: id})
	}

	result, err := env.listProducts.Execute(env.ctx, list_products.Request{
		Category: "list-test",
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.Products), 3)
}

func TestConcurrentUpdates(t *testing.T) {
	env := newTestEnv()
	productID := env.createActiveProduct(t)

	requests := []update_product.Request{
		{
			ProductID:   productID,
			Name:        "Concurrent Name A",
			Description: "Concurrent Desc A",
			Category:    "concurrency",
		},
		{
			ProductID:   productID,
			Name:        "Concurrent Name B",
			Description: "Concurrent Desc B",
			Category:    "concurrency",
		},
	}

	errCh := make(chan error, len(requests))
	var wg sync.WaitGroup

	for _, req := range requests {
		wg.Add(1)
		go func(r update_product.Request) {
			defer wg.Done()
			errCh <- env.updateProduct.Execute(env.ctx, r)
		}(req)
	}

	wg.Wait()
	close(errCh)

	var successCount int
	for err := range errCh {
		if err == nil {
			successCount++
		}
	}

	require.GreaterOrEqual(t, successCount, 1)

	product, err := env.getProduct.Execute(env.ctx, productID)
	require.NoError(t, err)
	assert.Contains(t, []string{"Concurrent Name A", "Concurrent Name B"}, product.Name)
	assert.Contains(t, []string{"Concurrent Desc A", "Concurrent Desc B"}, product.Description)
}

type outboxEvent struct {
	EventType   string
	AggregateID string
}

func getOutboxEvents(t *testing.T, ctx context.Context, aggregateID string) []outboxEvent {
	t.Helper()

	stmt := spanner.Statement{
		SQL:    "SELECT event_type, aggregate_id FROM outbox_events WHERE aggregate_id = @id ORDER BY created_at",
		Params: map[string]interface{}{"id": aggregateID},
	}

	tx := spannerClient.Single()
	defer tx.Close()

	iter := tx.Query(ctx, stmt)
	defer iter.Stop()

	var events []outboxEvent
	for {
		row, err := iter.Next()
		if err != nil {
			break
		}
		var e outboxEvent
		if err := row.Columns(&e.EventType, &e.AggregateID); err != nil {
			continue
		}
		events = append(events, e)
	}
	return events
}
