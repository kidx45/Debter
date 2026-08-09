migrateup:
	migrate -path ./internal/repository/migration -database "postgresql://root:secret@localhost:5434/debter?sslmode=disable" -verbose up
migratedown:
	migrate -path ./internal/repository/migration -database "postgresql://root:secret@localhost:5434/debter?sslmode=disable" -verbose down

.PHONY: migrateup migratedown

