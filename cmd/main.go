package main

import (
	"context"
	"fmt"
	auth "go_jwt_mcs/gen/go"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// 1) init db
	db := NewDB(ctx)

	// 2 init  grpc server and auth handler
	authHandler := NewAuthHandler(db)
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, authHandler)

	// 3) listen and serve
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT variable not configured in .env")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic("cant establish tcp connection:" + err.Error())
	}
	fmt.Println("gRPC server listening on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
