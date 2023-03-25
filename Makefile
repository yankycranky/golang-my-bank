postgres:
	docker run --name bank-postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it bank-postgres createdb --username=postgres --owner=postgres my_bank
dropdb:
	docker exec -it bank-postgres dropdb --username=postgres my_bank

migrate-up:
	migrate -path db/migration -database postgres://postgres:secret@localhost:5432/my_bank?sslmode=disable -verbose up

migrate-down:
	migrate -path db/migration -database postgres://postgres:secret@localhost:5432/my_bank?sslmode=disable -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrate-up migrate-down sqlc