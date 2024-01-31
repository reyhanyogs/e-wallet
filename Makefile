postgresinit:
	docker run --name e-wallet-postgres -p 5434:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16-alpine

createdb:
	docker exec -it e-wallet-postgres createdb --username=postgres --owner=postgres e-wallet

dropdb: 
	docker exec -it e-wallet-postgres dropdb e-wallet

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5434/e-wallet?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5434/e-wallet?sslmode=disable" -verbose down

.PHONY: postgresinit createdb dropdb migrateup migratedown