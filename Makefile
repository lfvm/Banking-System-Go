postgres:
	docker run --name finance-postgres --network bank-network  -p 5432:5432  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -d postgres:12-alpine 

createdb:
	docker exec -it finance-postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it finance-postgres dropdb simple_bank


migrateup:
	migrate -path db/migrations -database "postgresql://postgres:Gy5RBAD3GB6twmZi33Ge@simplebank.cv3c7y8izzg1.us-east-1.rds.amazonaws.com:5432/simple_bank" -verbose up 

migrateuplocal:
	migrate -path db/migrations -database  "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable"  -verbose up 

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go
devServer:
	air -c .air.toml

mockdb: 
	mockgen -package mockdb  -destination db/mock/store.go github.com/lfvm/simplebank/db/sqlc Store

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mockdb migrateup1  migratedown1 devServer migrateuplocal new_migration
