package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vniche/users-microservice/datastore"
	"github.com/vniche/users-microservice/graphql"
	"github.com/vniche/users-microservice/protocol"
)

func serve(ctx context.Context) (err error) {
	// starts datastore client
	datastore.Start()
	defer datastore.Close()

	go func() {
		if err := protocol.Serve(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		if err = graphql.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+v\n", err)
		}
	}()

	log.Printf("server started")

	<-ctx.Done()

	log.Printf("server stopped")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

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
		log.Printf("failed to serve: %+v\n", err)
	}
}
