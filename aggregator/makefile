.PHONY: proto run clean setup

proto:
	protoc --go_out=. --go-grpc_out=. ./api/proto/aggregator.proto;

run-server:
	go run main.go

run-mock-node:
	go run cmd/mock/node.go

clean:
	rm -rf api/proto/aggregator/v1