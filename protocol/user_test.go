package protocol

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vniche/users-microservice/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	RegisterUsersServer(srv, &server{})
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSignUp(t *testing.T) {
	// starts datastore client
	datastore.Start()

	// creates context to be used for connections
	ctx := context.Background()

	// tries to create a connection to gRPC server
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	// creates a client with the created connection
	client := NewUsersClient(conn)

	// tries to sign up a user
	resp, err := client.SignUp(ctx, &NewUser{FirstName: "John", LastName: "Doe"})
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}

	// ensure response is an UUID generated for the created user
	uuid.MustParse(resp.Uid)

	t.Cleanup(func() {
		if err := datastore.Close(); err != nil {
			t.Error("error resetting:", err)
		}
	})
}

func TestList(t *testing.T) {
	// starts datastore client
	datastore.Start()

	// creates context to be used for connections
	ctx := context.Background()

	// tries to create a connection to gRPC server
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	// creates a client with the created connection
	client := NewUsersClient(conn)

	// tries to sign up a user
	resp, err := client.List(ctx, &Empty{})
	if err != nil {
		t.Fatalf("SignUp failed: %v", err)
	}

	// ensure response users list size is not empty
	assert.NotEmpty(t, resp.Users, 0)

	t.Cleanup(func() {
		if err := datastore.Close(); err != nil {
			t.Error("error resetting:", err)
		}
	})
}
