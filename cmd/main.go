package main

import (
	"log"
	"net"

	"github.com/tanmaygupta069/auth-service/config"
	"github.com/tanmaygupta069/auth-service/internal/auth"
	"google.golang.org/grpc"
	pb "github.com/tanmaygupta069/auth-service/generated"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/credentials"
)

func main(){
	cfg,err := config.GetConfig()
	if err!=nil{
		log.Printf("error: %v",err.Error());
	}
	// router:=router.GetRouter()
	// router.Run(":"+cfg.ServerConfig.Port)
	sayHelloController:=auth.NewAuthController()
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}
	listener, err := net.Listen("tcp", ":"+cfg.GrpcServerConfig.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterAuthServiceServer(grpcServer, sayHelloController)
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on port %s",cfg.GrpcServerConfig.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}