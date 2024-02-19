postgresinit:
	docker run --name e-wallet-postgres -p 5434:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16-alpine

createdb:
	docker exec -it e-wallet-postgres createdb --username=postgres --owner=postgres e-wallet

dropdb: 
	docker exec -it e-wallet-postgres dropdb e-wallet

redisinit:
	docker run --name e-wallet-redis -p 6378:6379 -d redis:7.2-alpine

redis:
	docker exec -it e-wallet-redis redis-cli

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5434/e-wallet?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5434/e-wallet?sslmode=disable" -verbose down

composeup:
	docker compose up

composedown:
	docker compose down

.PHONY: postgresinit createdb dropdb redisinit redis new_migration migrateup migratedown composeup composedown