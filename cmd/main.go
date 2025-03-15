package main

import (
	"log"
	"net"

	"github.com/tanmaygupta069/auth-service-go/config"
	pb "github.com/tanmaygupta069/auth-service-go/generated"
	"github.com/tanmaygupta069/auth-service-go/internal/auth"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("error: %v", err.Error())
	}
	// router:=router.GetRouter()
	// router.Run(":"+cfg.ServerConfig.Port)
	sayHelloController := auth.NewAuthController()
	// _, err = credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}
	listener, err := net.Listen("tcp4",":"+cfg.GrpcServerConfig.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, sayHelloController)
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s", cfg.GrpcServerConfig.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
