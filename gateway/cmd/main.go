package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/homka122/Gomka122/gateway/internal/adapter/collector"
	"github.com/homka122/Gomka122/gateway/internal/config"
	controller "github.com/homka122/Gomka122/gateway/internal/controller/http"
	"github.com/homka122/Gomka122/gateway/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	conn, error := grpc.NewClient(cfg.CollectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}
	defer conn.Close()

	collector := collector.NewCollector(conn)

	repositoryUseCase := usecase.NewRepositoryUseCase(collector)

	handler := controller.NewHandler(repositoryUseCase)

	http.HandleFunc("/", handler.GetRepository)
	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
