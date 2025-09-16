.PHONY: build run dev test clean migrate-up migrate-down migrate-create migrate-force migrate-version migrate-install docker-up docker-down

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=devprep
MAIN_PATH=./cmd/devprep/main.go

DB_URL=postgres://postgres:password@localhost:5432/auth_db?sslmode=disable

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

test:
	$(GOTEST) ./...

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
