.PHONY: build run dev test test-unit test-integration test-e2e test-all test-coverage clean
.PHONY: migrate-up migrate-down migrate-create migrate-force migrate-version migrate-install
.PHONY: docker-up docker-down fmt vet lint deps

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=devprep
MAIN_PATH=./cmd/devprep/main.go

DB_URL=postgres://devprep_postgres:password@localhost:5432/devprep_db?sslmode=disable
TEST_DB_URL=postgres://postgres:password@localhost:5432/devprep_test?sslmode=disable

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

dev:
	$(GOCMD) run $(MAIN_PATH)

test:
	$(GOTEST) -v -race ./...

test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -race -short ./test/unit/...

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -race ./test/integration/...

test-e2e:
	@echo "Running E2E tests..."
	$(GOTEST) -v -race ./test/e2e/...

test-all: test-unit test-integration test-e2e

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt:
	$(GOCMD) fmt ./...

vet:
	$(GOCMD) vet ./...

lint:
	golangci-lint run

deps:
	$(GOMOD) tidy
	$(GOMOD) verify

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

migrate-install:
	$(GOCMD) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir database/migrations -seq $$name

migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path database/migrations -database "$(DB_URL)" force $$version

migrate-version:
	migrate -path database/migrations -database "$(DB_URL)" version
