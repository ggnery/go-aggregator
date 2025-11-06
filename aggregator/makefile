.PHONY: proto run clean setup

proto:
	protoc --go_out=. --go-grpc_out=. ./proto/aggregator.proto;

run:
	go run main.go

clean:
	rm -rf proto/aggregator