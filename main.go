package main

import (
	"customer/config"
	"customer/injector"
	"customer/pb"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// create a connection to the database
	db := config.NewConn().Database("db_customer")

	// get the port number from the environment
	port := os.Getenv("PORT")

	// listen for incoming connections on the specified port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	// create a new gRPC server
	GRPServer := grpc.NewServer()

	// enable reflection on the server
	reflection.Register(GRPServer)

	// create a new CustomerService handler with the database connection
	handler := injector.NewCustomerInjector(db)

	// register the handler with the gRPC server
	pb.RegisterCustomerServiceServer(GRPServer, handler)

	// log that the server is ready
	fmt.Printf("⚡️[server]: gRPC Server is running on port %s\n", port)

	// start serving requests
	if err := GRPServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
