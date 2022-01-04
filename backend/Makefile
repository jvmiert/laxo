include .env
MIGRATE_VERSION_TAG=v4.15.1
SQLC_VERSION_TAG=v1.11.0

.PHONY: install migrate upgrade generate

install:
	# Installing migrate tool to manage migrations
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION_TAG)
	#
	# sqlc is used to easily create reusable queries
	#
	go install github.com/kyleconroy/sqlc/cmd/sqlc@${SQLC_VERSION_TAG}
	go mod tidy

generate:
	sqlc generate -f laxo/sqlc/config/sqlc.yaml

migrate:
	@read -p "Enter the name of the new migration: " name; \
	migrate create -ext sql -dir laxo/migrations -seq $$name;

upgrade:
	migrate -database ${POSTGRESQL_URL} -path laxo/migrations up

