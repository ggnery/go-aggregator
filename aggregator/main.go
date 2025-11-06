package main

import (
	"log"
	"net" // network operations
	"google.golang.org/grpc"
	"aggregator/controller"
)

func main() {
	// Define the port to listen on
	port := ":50051"
	
	// Create a TCP listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	controller.RegisterGRPCServers(grpcServer)	

	log.Printf("gRPC server listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}