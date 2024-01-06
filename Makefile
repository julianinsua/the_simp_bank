postgres:
	docker run --name postgres_db --network simp_net -p 5432:5432 -e POSTGRES_USER=postgres POSTGRES_PASSWORD=password -d postgres:14.1-alpine 

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

serverrun:
	go run main.go

mock:
	mockgen --build_flags=--mod=mod -destination db/mock/store.go -package mock_db github.com/julianinsua/the_simp_bank/internal/database Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc serverrun mock
	
