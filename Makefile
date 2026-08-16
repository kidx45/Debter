migrateup:
	migrate -path ./internal/migration -database ${DB_URL} -verbose up
migratedown:
	migrate -path ./internal/migration -database ${DB_URL} -verbose down
migratecreate:
	migrate create -ext sql -dir ./internal/migration -seq ${SCHEMA}

.PHONY: migrateup migratedown

