include .env
MIGRATE_VERSION_TAG=v4.15.1
SQLC_VERSION_TAG=996a73a

.PHONY: install migrate upgrade generate test testrace translate serve migrate-version

serve:
	go run ./cmd/laxo/main.go

translate:
	go generate ./laxo/translations/translations.go

test:
	migrate -database ${POSTGRESQL_TEST_URL} -path laxo/migrations drop -f
	migrate -database ${POSTGRESQL_TEST_URL} -path laxo/migrations up
	go test -v ./laxo_tests

testrace:
	migrate -database ${POSTGRESQL_TEST_URL} -path laxo/migrations drop -f
	migrate -database ${POSTGRESQL_TEST_URL} -path laxo/migrations up
	go test -race -v ./cmd/laxo

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

migrate-version:
	migrate -database ${POSTGRESQL_URL} -path laxo/migrations version

migrate-force:
	@read -p "Enter the migration version number to force to: " number; \
	migrate -database ${POSTGRESQL_URL} -path laxo/migrations force $$number;
