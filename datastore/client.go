package datastore

import (
	"context"
	"log"
	"os"

	ceteProto "github.com/vniche/cete/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataStoreClient struct {
	KVSClient ceteProto.KVSClient
	Conn      *grpc.ClientConn
}

var Client *DataStoreClient

func Start() {
	// cete grpc port
	ceteEndpoint := ":9000"
	if os.Getenv("CETE_ENDPOINT") != "" {
		ceteEndpoint = os.Getenv("CETE_ENDPOINT")
	}

	dialOpts := append([]grpc.DialOption{},
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.EmptyCallOption{}))
	conn, err := grpc.Dial(ceteEndpoint, dialOpts...)
	if err != nil {
		log.Fatal(err)
	}

	Client = &DataStoreClient{
		KVSClient: ceteProto.NewKVSClient(conn),
		Conn:      conn,
	}
}

func (client *DataStoreClient) Node() (*ceteProto.NodeResponse, error) {
	return client.KVSClient.Node(context.Background(), &emptypb.Empty{}, grpc.EmptyCallOption{})
}

func Close() error {
	return Client.Conn.Close()
}
