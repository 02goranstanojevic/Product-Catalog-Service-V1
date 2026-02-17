.PHONY: all build run test test-unit proto migrate init clean emulator-up emulator-down

ifeq ($(OS),Windows_NT)
POWERSHELL := powershell.exe -NoProfile -ExecutionPolicy Bypass
endif

PROJECT_ID := test-project
INSTANCE_ID := test-instance
DATABASE_ID := test-db

SPANNER_DB := projects/$(PROJECT_ID)/instances/$(INSTANCE_ID)/databases/$(DATABASE_ID)

all: build

ifeq ($(OS),Windows_NT)
build:
	go build -o bin/server.exe ./cmd/server

run:
	@$(POWERSHELL) -File scripts/run.ps1 -Action run

test:
	@$(POWERSHELL) -File scripts/run.ps1 -Action test

test-unit:
	go test ./internal/app/product/domain/... -v -count=1
else
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
endif

emulator-up:
	docker-compose up -d

emulator-down:
	docker-compose down

ifeq ($(OS),Windows_NT)
migrate:
	@$(POWERSHELL) -File scripts/run.ps1 -Action migrate
else
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
endif

init:
	$(MAKE) emulator-up
	$(MAKE) migrate

proto:
	go run github.com/bufbuild/buf/cmd/buf@v1.34.0 generate

clean:
ifeq ($(OS),Windows_NT)
	@if exist bin rmdir /s /q bin
else
	rm -rf bin/
endif
