package bootstrap

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func StartGRPCServer(port int, register func(*grpc.Server)) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	register(grpcServer)
	log.Printf("gRPC server started on %s", addr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
