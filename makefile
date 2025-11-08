.PHONY: proto run clean setup db-up db-down db-migrate db-info db-clean db-reset db-connect db-tables

proto:
	cd aggregator && protoc --go_out=. --go-grpc_out=. ./api/proto/aggregator.proto;

run-server:
	cd aggregator && go run main.go

run-mock-node:
	cd aggregator && go run cmd/mock/node.go

clean:
	rm -rf aggregator/api/proto/aggregator/v1

# Database commands
db-up:
	docker-compose up -d db

db-down:
	docker-compose down

db-migrate:
	docker-compose up flyway

db-info:
	docker-compose run --rm flyway -configFiles=/flyway/conf/flyway.conf info

db-clean:
	docker-compose run --rm flyway -configFiles=/flyway/conf/flyway.conf clean

db-reset: db-clean db-migrate

# Start everything (database + migrations)
db-setup: db-up db-migrate