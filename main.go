package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/vniche/users-microservice/datastore"
	"github.com/vniche/users-microservice/entities"
	"github.com/vniche/users-microservice/graphql"
	pb "github.com/vniche/users-microservice/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement proto.GrpcMicroservice
type server struct {
	pb.UnimplementedUsersServer
}

func (s *server) SignUp(ctx context.Context, req *pb.NewUser) (*pb.Created, error) {
	uid, err := entities.SignUp(&entities.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to sign up user %s: %s", req.FirstName, err.Error())
	}

	return &pb.Created{Uid: uid}, nil
}

func (s *server) List(ctx context.Context, req *pb.Empty) (*pb.UserList, error) {
	users, err := entities.List()
	if err != nil {
		return nil, fmt.Errorf("Unable to list users: %s", err.Error())
	}

	parsed := make([]*pb.User, len(users))
	for index, user := range users {
		parsed[index] = &pb.User{
			Uid:       user.UID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt.String(),
		}
	}

	return &pb.UserList{Users: parsed}, nil
}

func serve(ctx context.Context) (err error) {
	// GRPC Server
	grpcPort := ":5000"
	if os.Getenv("GRPC_PORT") != "" {
		grpcPort = os.Getenv("GRPC_PORT")
	}

	datastore.Start()
	defer datastore.Client.Close()

	// creates a new gRPC server instance and register microservice to it
	grpcServer := grpc.NewServer()
	pb.RegisterUsersServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// creates a new TCP listener
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		// starts gRPC server with the TCP listener
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// GraphQL Server
	graphqlPort := ":3000"
	if os.Getenv("GRAPHQL_PORT") != "" {
		graphqlPort = os.Getenv("GRAPHQL_PORT")
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql/playground", handler.Playground("GraphQL playground", "/query"))
	mux.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(
		graphql.Config{Resolvers: &graphql.Resolver{}}),
	))

	graphqlServer := &http.Server{
		Addr:    graphqlPort,
		Handler: mux,
	}

	go func() {
		// starts graphql http server with the TCP listener
		log.Printf("connect to http://localhost:%s/graphql/playground for GraphQL playground", graphqlPort)
		if err = graphqlServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	grpcServer.Stop()

	if err = graphqlServer.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("unable to shutdown graphql server:%+s\n", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
