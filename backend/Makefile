# Adfluence backend — quick commands (run from backend/)
# Usage: make help

.PHONY: help up down restart reset logs psql migrate migrate-security migrate-phone-verification migrate-portfolio seed-categories status tables seed-admin

help:
	@echo "Adfluence database commands"
	@echo ""
	@echo "  make up        Start PostgreSQL + MinIO (Docker)"
	@echo "  make down      Stop PostgreSQL"
	@echo "  make restart   Restart container"
	@echo "  make reset     Wipe DB + recreate + migrate (DESTROYS DATA)"
	@echo "  make migrate         Run ALL migrations (fresh DB only)"
	@echo "  make migrate-security  Add rate-limit tables to existing DB"
	@echo "  make migrate-phone-verification  Add phone OTP table to existing DB"
	@echo "  make migrate-portfolio     Add influencer_portfolio_items table"
	@echo "  make seed-categories     Seed niche categories (Beauty, Tech, …)"
	@echo "  make logs      Show postgres logs"
	@echo "  make psql      Open psql shell"
	@echo "  make status    Container + table count"
	@echo "  make seed-admin  Insert dev admin (admin@adfluence.com / admin1111)"
	@echo "  make tables    List all tables"

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart postgres

reset:
	docker compose down -v
	docker compose up -d
	@echo "Waiting for migrations..."
	@sleep 10
	@$(MAKE) status

migrate:
	./scripts/migrate.sh

migrate-security:
	chmod +x ./scripts/migrate_security.sh
	./scripts/migrate_security.sh

migrate-phone-verification:
	chmod +x ./scripts/migrate_phone_verification.sh
	./scripts/migrate_phone_verification.sh

migrate-portfolio:
	chmod +x ./scripts/migrate_portfolio.sh
	./scripts/migrate_portfolio.sh

seed-categories:
	PGPASSWORD=adfluence123 psql -h localhost -p 5433 -U adfluence -d adfluence_db -v ON_ERROR_STOP=1 -f migration/seed_categories.sql

logs:
	docker compose logs -f postgres

psql:
	PGPASSWORD=adfluence123 psql -h localhost -p 5433 -U adfluence -d adfluence_db

status:
	@docker compose ps
	@PGPASSWORD=adfluence123 psql -h localhost -p 5433 -U adfluence -d adfluence_db -t -c \
		"SELECT 'Tables: ' || count(*) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';" 2>/dev/null || echo "DB not ready — run: make up"

tables:
	PGPASSWORD=adfluence123 psql -h localhost -p 5433 -U adfluence -d adfluence_db -c \
		"SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY tablename;"

seed-admin:
	chmod +x ./scripts/seed_admin.sh
	./scripts/seed_admin.sh
