package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/spanner"
	"github.com/product-catalog-service/internal/services"
	pb "github.com/product-catalog-service/proto/product/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	spannerDB := os.Getenv("SPANNER_DATABASE")
	if spannerDB == "" {
		spannerDB = "projects/test-project/instances/test-instance/databases/test-db"
	}

	spannerEmulator := os.Getenv("SPANNER_EMULATOR_HOST")

	var opts []option.ClientOption
	if spannerEmulator != "" {
		log.Printf("Using Spanner emulator at %s", spannerEmulator)
	}

	client, err := spanner.NewClient(ctx, spannerDB, opts...)
	if err != nil {
		log.Fatalf("Failed to create Spanner client: %v", err)
	}
	defer client.Close()

	container := services.NewContainer(client)

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, container.ProductHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
}
