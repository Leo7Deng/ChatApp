include .env

create_migration:
	migrate create -ext=sql -dir=postgres/migrations -seq init

migrate_up:
	migrate -path=postgres/migrations -database "postgresql://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_DBNAME}?sslmode=${PSQL_SSLMODE}" -verbose up

migrate_down:
	migrate -path=postgres/migrations -database "postgresql://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_DBNAME}?sslmode=${PSQL_SSLMODE}" -verbose down

migrate_force:
	migrate -path=postgres/migrations -database "postgresql://${PSQL_USER}:${PSQL_PASSWORD}@${PSQL_HOST}:${PSQL_PORT}/${PSQL_DBNAME}?sslmode=${PSQL_SSLMODE}" force 2

.PHONY: create_migration migrate_up migrate_down migrate_force