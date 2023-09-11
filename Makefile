postgres:
	docker run --name postgres_db -p 5432:5432 -e POSTGRES_USER=postgres POSTGRES_PASSWORD=password -d postgres:14.1-alpine 

createdb:
	docker exec -it postgres_db createdb --username=postgres --owner=postgres simp_bank

dropdb:
	docker exec -it postgres_db dropdb --username=postgres simp_bank

migrateup:
	goose postgres postgresql://postgres:password@localhost:5432/simp_bank up

migratedown:
	goose postgres postgresql://postgres:password@localhost:5432/simp_bank down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
	

