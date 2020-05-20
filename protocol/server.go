package protocol

import (
	context "context"
	fmt "fmt"
	"log"
	"net"
	"os"

	"github.com/vniche/users-microservice/entities"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Serve() error {
	// creates a new gRPC server instance and register microservice to it
	grpcServer := grpc.NewServer()
	RegisterUsersServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	// gRPC Server port
	grpcPort := ":5000"
	if os.Getenv("GRPC_PORT") != "" {
		grpcPort = os.Getenv("GRPC_PORT")
	}

	// creates a new TCP listener
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listen.Close()

	// starts gRPC server with the TCP listener
	log.Printf("gRPC is available at tcp://localhost%s", grpcPort)
	return grpcServer.Serve(listen)
}

// server is used to implement proto.GrpcMicroservice
type server struct {
	UnimplementedUsersServer
}

func (s *server) SignUp(ctx context.Context, req *NewUser) (*Created, error) {
	uid, err := entities.SignUp(&entities.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to sign up user %s: %s", req.FirstName, err.Error())
	}

	return &Created{Uid: uid}, nil
}

func (s *server) List(ctx context.Context, req *Empty) (*UserList, error) {
	users, err := entities.List()
	if err != nil {
		return nil, fmt.Errorf("Unable to list users: %s", err.Error())
	}

	parsed := make([]*User, len(users))
	for index, user := range users {
		parsed[index] = &User{
			Uid:       user.UID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt.String(),
		}
	}

	return &UserList{Users: parsed}, nil
}
