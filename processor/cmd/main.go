package main

import (
	"fmt"
	"log"
	"net"

	collector "github.com/homka122/Gomka122/processor/internal/adapter"
	"github.com/homka122/Gomka122/processor/internal/config"
	controller "github.com/homka122/Gomka122/processor/internal/controller/grpc"
	"github.com/homka122/Gomka122/processor/internal/usecase"
	pbProcessor "github.com/homka122/Gomka122/proto/processor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient(cfg.CollectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	collector := collector.NewCollector(conn)
	usecase := usecase.NewRepositoryUseCase(collector)
	server := controller.NewServer(usecase)

	grpcServer := grpc.NewServer()
	pbProcessor.RegisterProcessorServiceServer(grpcServer, server)

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
