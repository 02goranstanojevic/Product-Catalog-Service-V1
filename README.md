# Product Catalog Service

A Go microservice implementing a Product Catalog with DDD, Clean Architecture, CQRS, and transactional outbox pattern. Built on Google Cloud Spanner with gRPC transport.

## Architecture

```
cmd/server/          -> Entry point
internal/
  app/product/
    domain/          -> Pure business logic (aggregates, value objects, events)
    usecases/        -> Command handlers (write operations)
    queries/         -> Query handlers (read operations)
    contracts/       -> Repository & read model interfaces
    repo/            -> Spanner implementations
  models/            -> Database model structs & field constants
  transport/grpc/    -> gRPC handlers, mappers, error mapping
  services/          -> DI container
  pkg/               -> Shared utilities (clock, committer)
proto/               -> Protocol Buffer definitions
migrations/          -> Spanner DDL
```

### Design Decisions

**Domain Purity**: The domain layer has zero external dependencies - no `context.Context`, no database imports, no proto types. All business logic is expressed as pure Go.

**Golden Mutation Pattern**: Every write operation follows: Load aggregate -> Execute domain logic -> Build CommitPlan -> Apply atomically. Repositories return `*spanner.Mutation` but never apply them. The use case interactor owns the transaction boundary.

**CQRS**: Commands go through domain aggregates for validation and event capture. Queries bypass the domain, reading directly from the database via a read model interface with DTOs.

**Change Tracking**: The `ChangeTracker` on each aggregate tracks dirty fields, enabling targeted `UPDATE` mutations that only touch modified columns.

**Transactional Outbox**: Domain events are captured as intent structs during business operations. The use case enriches and serializes them into the `outbox_events` table within the same atomic transaction as the aggregate mutation.

**Money as `*big.Rat`**: All monetary values use `math/big.Rat` (stored as numerator/denominator pairs) for lossless decimal arithmetic.

### Trade-offs

- **Proto types are hand-written stubs** - avoids requiring `protoc` tooling for local dev. gRPC method descriptors and handler routing are wired manually (matching protoc-generated patterns). Run `make proto` to generate from `.proto` files when protoc is available; the proto codec requires protoc-generated types for on-the-wire serialization.
- **Custom CommitPlan instead of `github.com/Vektor-AI/commitplan`** - The `commitplan` module is internal to the organization and not publicly available.  The internal `pkg/committer` package follows the same pattern: `NewPlan()` -> `plan.Add(mutation)` -> `committer.Apply(ctx, plan)`, providing identical semantics with Spanner's `Apply` for atomic transactions.
- The outbox processor is not implemented - events are stored but not dispatched. This is by design per the spec.
- No authentication, metrics, or REST gateway.

## Prerequisites

- Go 1.21+
- Docker (for Spanner emulator)
- `gcloud` CLI (for migrations)
- `make` (optional, mostly for Linux/macOS convenience)

### Windows Note

On Windows, use `scripts/run.ps1` as the primary workflow. `make` is optional and not required.

## Quick Start

For Windows, prefer the `scripts/run.ps1` flow in **Local Run Order (Client Review)** below.

```bash
# Start Spanner emulator
make emulator-up

# Create instance, database, and apply schema
make migrate

# Build and run the gRPC server
make run

# Run unit tests (no emulator needed)
make test-unit

# Run all tests including E2E (emulator required)
make test
```

## Local Run Order (Client Review)

Use the commands below in this exact order for a clean local backend verification.

### Option A: PowerShell script (Windows, fastest)

```powershell
# 1) Start emulator
.\scripts\run.ps1 -Action up

# 2) Apply schema
.\scripts\run.ps1 -Action migrate

# 3) Run full test suite (unit + e2e)
.\scripts\run.ps1 -Action test

# 4) Start backend server
.\scripts\run.ps1 -Action run
```

### Option B: Makefile commands

Use this if `make` is installed on your system.

```bash
# 1) Start emulator
make emulator-up

# 2) Apply schema
make migrate

# 3) Run full test suite (unit + e2e)
make test

# 4) Start backend server
make run
```

### Stop local infrastructure

```powershell
.\scripts\run.ps1 -Action down
```

or

```bash
make emulator-down
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `50051` | gRPC server port |
| `SPANNER_DATABASE` | `projects/test-project/instances/test-instance/databases/test-db` | Spanner database path |
| `SPANNER_EMULATOR_HOST` | - | Set to `localhost:9010` for local development |

## API

The service exposes a gRPC `ProductService` with:

**Commands**: `CreateProduct`, `UpdateProduct`, `ActivateProduct`, `DeactivateProduct`, `ApplyDiscount`, `RemoveDiscount`

**Queries**: `GetProduct`, `ListProducts` (with category filter and pagination)

See [proto/product/v1/product_service.proto](proto/product/v1/product_service.proto) for the full API definition.

## Testing

**Unit tests** cover domain logic in isolation - product state machine, money arithmetic, discount validation, pricing calculations.

**E2E tests** run against the Spanner emulator, exercising the full flow from use case interactors through repositories to the database and back through queries.

### Test modes (local)

- Unit only (no Docker required):

PowerShell (Windows):

```powershell
go test ./internal/app/product/domain/... -v -count=1
```

Makefile:

```bash
make test-unit
```

- E2E only (requires emulator):

```bash
SPANNER_EMULATOR_HOST=localhost:9010 go test ./tests/e2e -v -count=1
```

- All backend tests:

PowerShell (Windows):

```powershell
.\scripts\run.ps1 -Action test
```

Makefile:

```bash
make test
```

On Windows PowerShell for E2E:

```powershell
$env:SPANNER_EMULATOR_HOST='localhost:9010'; go test ./tests/e2e -v -count=1
```

### Important note

- This repository is backend-only (no frontend is included by design).
- Main readiness goal is: build succeeds, migrations apply, tests pass, and gRPC server starts successfully.

```bash
# Unit tests only
make test-unit

# All tests (requires emulator)
make test
```
