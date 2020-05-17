package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/vniche/users-microservice/datastore"
	"github.com/vniche/users-microservice/graphql"
	pb "github.com/vniche/users-microservice/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// server is used to implement proto.GrpcMicroservice
type server struct {
	pb.UnimplementedUsersServer
}

func (s *server) SignUp(ctx context.Context, req *pb.User) (*pb.Created, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}

func (s *server) List(ctx context.Context, req *pb.Empty) (*pb.UserList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func main() {
	// GRPC Server
	grpcPort := ":5000"
	if os.Getenv("GRPC_PORT") != "" {
		grpcPort = os.Getenv("GRPC_PORT")
	}

	// creates a new TCP listener
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	datastore.Start()
	defer datastore.Client.Close()

	// creates a new gRPC server instance and register microservice to it
	srv := grpc.NewServer()
	pb.RegisterUsersServer(srv, &server{})
	reflection.Register(srv)

	// GraphQL Server
	graphqlPort := "3000"
	if os.Getenv("GRAPHQL_PORT") != "" {
		graphqlPort = os.Getenv("GRAPHQL_PORT")
	}

	http.Handle("/graphql/playground", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(
		graphql.Config{Resolvers: &graphql.Resolver{}}),
	))

	log.Printf("connect to http://localhost:%s/graphql/playground for GraphQL playground", graphqlPort)
	log.Fatal(http.ListenAndServe(":"+graphqlPort, nil))

	// starts gRPC server with the TCP listener
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
