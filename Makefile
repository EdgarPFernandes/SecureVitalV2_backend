
APP=api
MIGRATIONS=./internal/adapters/postgresql/migrations
DB_URL=postgres://postgres:postgres@localhost:5432/svdb?sslmode=disable

# =====================
# GO APP
# =====================
run:
	go run .\cmd\main.go .\cmd\api.go

build:
	go build -o bin/$(APP) ./cmd/$(APP).go

# =====================
# DOCKER
# =====================
up:
	docker compose up

down:
	docker compose down

down-v:
	docker compose down -v

# =====================
# MIGRATIONS (GOOSE)
# =====================
migrate-up:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" status

# create migration (FIXED so you DON'T need to cd)
create-migration:
	goose -dir $(MIGRATIONS) create $(name) sql

# =====================
# SQLC
# =====================
sqlc:
	docker run --rm -v ${PWD}:/src -w /src sqlc/sqlc generate

# =====================
# FULL DEV FLOW
# =====================
dev:
	docker compose up -d
	make migrate-up
	make run