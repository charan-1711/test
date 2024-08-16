include .env

migrate_up: 
	migrate -path database/migrations/ -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" -verbose up
migrate_down: 
	migrate -path database/migrations/ -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" -verbose down
migrate_fix:
	migrate -path database/migrations/ -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" force 1`