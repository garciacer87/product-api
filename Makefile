PRINT= printf "%b"
BLUE=\033[0;94m
NC=\033[0m

create-postgres:
	@$(PRINT) "$(BLUE)Starting PostgreSQL$(NC)\n"
	docker run --name api-db -e POSTGRES_USER=productapi -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:13-alpine

create-db-test:
	@$(PRINT) "$(BLUE)Creating DB for tests$(NC)\n"
	docker exec api-db createdb --username=productapi --owner=productapi productapitest

drop-db-test:
	@$(PRINT) "$(BLUE)Dropping tests DB$(NC)\n"
	docker exec api-db dropdb --username=productapi productapitest

drop-db:
	@$(PRINT) "$(BLUE)Dropping DB$(NC)\n"
	docker exec api-db dropdb --username=productapi productapi

drop-postgres:
	@$(PRINT) "$(BLUE)Stopping and removing postgresql container$(NC)\n"
	docker stop api-db
	docker container rm api-db

migrate-up:
	@$(PRINT) "$(BLUE)Running Database migration...$(NC)\n"
	migrate -path sql/postgresql/ -database "postgresql://productapi:password@localhost:5432/productapi?sslmode=disable" -verbose up

migrate-down:
	@$(PRINT) "$(BLUE)Undoing Database migration...$(NC)\n"
	migrate -path sql/postgresql/ -database "postgresql://productapi:password@localhost:5432/productapi?sslmode=disable" -verbose down -all

migrate-drop:
	@$(PRINT) "$(BLUE)Dropping migration...$(NC)\n"
	migrate -path sql/postgresql/ -database "postgresql://productapi:password@localhost:5432/productapi?sslmode=disable" -verbose drop -f

go-build:
	go build -o ./target/product-api ./cmd/api

docker-build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./target/product-api ./cmd/api
	docker build -f build/docker/dockerfile -t product-api:1.0.0 .

docker-run:
	docker run -e PORT=8080 -e DATABASE_URI="postgres://productapi:password@localhost:5432/productapi" -d -p 8080:8080 product-api:1.0.0

test:
	go test -cover ./...

vet:
	@$(PRINT) "$(BLUE)Vetting the source code...$(NC)\n"
	go vet ./...

revive:
	@$(PRINT) "$(BLUE)Running revive...$(NC)\n"
	revive -formatter friendly ./...

staticcheck:
	@$(PRINT) "$(BLUE)Running Static Code Analysys...$(NC)\n"
	staticcheck ./...

check: vet revive staticcheck

run:
	./target/product-api

.PHONY: build run check