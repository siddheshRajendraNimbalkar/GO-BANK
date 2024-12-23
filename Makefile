.PHONY: postgres createdb dropdb migrateup migratedown sqlc Test Server migratedown1 migrateup1

postgres:
	docker run --name goBank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it goBank createdb --username=root --owner=root  simple_bank

dropdb:
	docker exec -it goBank dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
	
migratedown1:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

Test:
	go test ./...

Server: 
	go run main.go