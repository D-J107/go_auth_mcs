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
	// settings, err := godotenv.Read("../.env")
	// os.Setenv("DATABASE_URL", settings["DATABASE_URL"])
	// os.Setenv("PORT", settings["PORT"])
	// if err != nil {
	// 	panic("cant read .env file!")
	// }

	ctx := context.Background()

	// 1) init db
	db, err := NewDB(ctx)
	if err != nil {
		panic("cant init db: " + err.Error())
	}

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
