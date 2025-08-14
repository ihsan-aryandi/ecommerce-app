# Load .env file if exists
-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

# ========== VARIABLES ==========
# Migrations folder path
MIGRATIONS_PATH=migrations

# Database connection
DATABASE=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Wire file path
WIRE_PATH=./internal/app # Must contain "./"; otherwise, Google Wire considers it as a package.

# ========== GOLANG MIGRATE ==========
# Create migration file
migrate-create:
	@echo Creating migration file...
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-force:
	@echo Running migrations...
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE) force $(version)

migrate-drop:
	@echo Running migrations...
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE) drop

migrate-up:
	@echo Running migrations...
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE) up

migrate-down:
	@echo Running migrations...
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE) down

# ========== GOOGLE WIRE ==========
# Generate wire_gen.go
wire:
	wire $(WIRE_PATH)