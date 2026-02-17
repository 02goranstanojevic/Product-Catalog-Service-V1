.PHONY: all build run test proto migrate clean emulator-up emulator-down

PROJECT_ID := test-project
INSTANCE_ID := test-instance
DATABASE_ID := test-db

SPANNER_DB := projects/$(PROJECT_ID)/instances/$(INSTANCE_ID)/databases/$(DATABASE_ID)

all: build

build:
	go build -o bin/server ./cmd/server

run: build
	SPANNER_EMULATOR_HOST=localhost:9010 \
	SPANNER_DATABASE=$(SPANNER_DB) \
	./bin/server

test:
	SPANNER_EMULATOR_HOST=localhost:9010 \
	go test ./... -v -count=1

test-unit:
	go test ./internal/app/product/domain/... -v -count=1

emulator-up:
	docker-compose up -d

emulator-down:
	docker-compose down

migrate:
	@echo "Creating Spanner instance and database on emulator..."
	gcloud config configurations create emulator --no-activate 2>/dev/null || true
	SPANNER_EMULATOR_HOST=localhost:9010 gcloud spanner instances create $(INSTANCE_ID) \
		--config=emulator-config --description="Test Instance" --nodes=1 \
		--project=$(PROJECT_ID) 2>/dev/null || true
	SPANNER_EMULATOR_HOST=localhost:9010 gcloud spanner databases create $(DATABASE_ID) \
		--instance=$(INSTANCE_ID) \
		--project=$(PROJECT_ID) \
		--ddl-file=migrations/001_initial_schema.sql 2>/dev/null || true
	@echo "Migration complete."

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/product/v1/product_service.proto

clean:
	rm -rf bin/
