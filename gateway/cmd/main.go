package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/homka122/Gomka122/gateway/internal/adapter/collector"
	"github.com/homka122/Gomka122/gateway/internal/adapter/processor"
	"github.com/homka122/Gomka122/gateway/internal/adapter/subscriber"
	"github.com/homka122/Gomka122/gateway/internal/config"
	controller "github.com/homka122/Gomka122/gateway/internal/controller/http"
	"github.com/homka122/Gomka122/gateway/internal/usecase"
	"github.com/homka122/Gomka122/internal/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title		Gomka122 API
//	@version	0.1

//	@contact.name	Homka122
//	@contact.url	t.me/homka122
//	@contact.email	kombaev02@gmail.com

//	@license.name	MIT
//	@license.url	https://mit-license.org/

//	@host		localhost:8080
//	@BasePath	/

func main() {
	cfg := config.Load()
	logger := logger.Load("DEBUG")

	processor := processor.NewProcessor(cfg, logger)
	collector := collector.NewCollector(cfg, logger)
	subscriber := subscriber.NewSubscriber(cfg, logger)

	repositoryUseCase := usecase.NewRepositoryUseCase(processor, logger)
	pingUseCase := usecase.NewPingUsecase(map[string]usecase.Pinger{
		"collector":  collector,
		"processor":  processor,
		"subscriber": subscriber,
	}, logger)
	subscribeUseCase := usecase.NewSubscribeUseCase(subscriber, logger)

	handler := controller.NewHandler(repositoryUseCase, pingUseCase, subscribeUseCase, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/docs/swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger/doc.json", cfg.Port))))
	mux.HandleFunc("GET /api/repositories/info", handler.GetRepository)
	mux.HandleFunc("GET /api/ping", handler.PingServices)
	mux.HandleFunc("POST /api/subscriptions", handler.Subscribe)
	mux.HandleFunc("DELETE /api/subscriptions/{owner}/{repo}", handler.Unsubscribe)
	mux.HandleFunc("GET /api/subscriptions", handler.GetSubscriptions)
	mux.HandleFunc("GET /api/subscriptions/info", handler.GetSubscribedRepositories)

	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux))
}
