.PHONY: proto run clean setup

proto:
	protoc --go_out=. --go-grpc_out=. ./api/proto/aggregator.proto;

run:
	go run cmd/server/main.go

clean:
	rm -rf api/proto/aggregator/v1