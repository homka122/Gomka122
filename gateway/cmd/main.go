package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/homka122/Gomka122/gateway/internal/adapter/collector"
	"github.com/homka122/Gomka122/gateway/internal/config"
	controller "github.com/homka122/Gomka122/gateway/internal/controller/http"
	"github.com/homka122/Gomka122/gateway/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//	@title		Gomka122 API
//	@version	0.1

//	@contact.name	Homka122
//	@contact.url	t.me/homka122
//	@contact.email	kombaev02@gmail.com

//	@license.name	MIT
//	@license.url	https://mit-license.org/

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	cfg := config.Load()

	conn, error := grpc.NewClient(cfg.ProcessorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}
	defer conn.Close()

	collector := collector.NewCollector(conn)

	repositoryUseCase := usecase.NewRepositoryUseCase(collector)

	handler := controller.NewHandler(repositoryUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/docs/swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger/doc.json", cfg.Port))))
	mux.HandleFunc("/repo/{owner}/{repo}", handler.GetRepository)

	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux))
}
