include .env
MIGRATE_VERSION_TAG=v4.15.1

.PHONY: install migrate upgrade

install:
	# Installing migrate tool to manage migrations
	go get -u -d github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION_TAG)
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION_TAG)

migrate:
	@read -p "Enter the name of the new migration: " name; \
	migrate create -ext sql -dir laxo/migrations -seq $$name;

upgrade:
	migrate -database ${POSTGRESQL_URL} -path laxo/migrations up

