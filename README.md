SUMUNARPOS CORE

<!-- migration -->

source .env
migrate create -ext sql -dir db/migrations -seq create_users_table
migrate -path db/migrations -database "$DB_URL" up

migrate -database "$DB_URL" -path db/migrations force 0
