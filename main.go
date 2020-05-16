// package main implements a server for GrpcMicroservice service
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/vniche/users-microservice/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

// server is used to implement proto.GrpcMicroservice
type server struct {
	pb.UnimplementedGrpcMicroserviceServer
}

func (s *server) Method(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Method not implemented")
}

func (s *server) NoRequestMethod(ctx context.Context, req *pb.Empty) (*pb.Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoRequestMethod not implemented")
}

func main() {
	// creates a new TCP listener
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// creates a new gRPC server instance and register microservice to it
	srv := grpc.NewServer()
	pb.RegisterGrpcMicroserviceServer(srv, &server{})

	// starts gRPC server with the TCP listener
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
