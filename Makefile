migrateup:
	migrate -path ./internal/migration -database ${DB_URL} -verbose up
migratedown:
	migrate -path ./internal/migration -database ${DB_URL} -verbose down

.PHONY: migrateup migratedown

