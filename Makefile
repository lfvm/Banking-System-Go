postgres:
	docker run --name finance-postgres --network bank-network  -p 5432:5432  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine 

createdb:
	docker exec -it finance-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it finance-postgres dropdb simple_bank


migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 

migrateup1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

mockdb: 
	mockgen -package mockdb  -destination db/mock/store.go github.com/lfvm/simplebank/db/sqlc Store


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mockdb migrateup1 migratedown1
