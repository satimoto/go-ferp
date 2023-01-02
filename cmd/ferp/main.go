package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/satimoto/go-ferp/internal/converter"
	"github.com/satimoto/go-ferp/internal/exchange"
	metrics "github.com/satimoto/go-ferp/internal/metric"
	"github.com/satimoto/go-ferp/internal/rest"
	"github.com/satimoto/go-ferp/internal/rpc"
)

func main() {
	log.Printf("Starting up FERP server")
	shutdownCtx, cancelFunc := context.WithCancel(context.Background())
	waitGroup := &sync.WaitGroup{}

	exchangeService := exchange.NewService()
	exchangeService.Start(shutdownCtx, waitGroup)

	converterService := converter.NewService(exchangeService)
	converterService.Start(shutdownCtx, waitGroup)
	
	metricsService := metrics.NewMetrics()
	metricsService.StartMetrics(shutdownCtx, waitGroup)

	restService := rest.NewRest()
	restService.StartRest(shutdownCtx, waitGroup)

	rpcService := rpc.NewRpc(converterService)
	rpcService.StartRpc(shutdownCtx, waitGroup)

	sigtermChan := make(chan os.Signal, 1)
	signal.Notify(sigtermChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigtermChan

	log.Printf("Shutting down FERP server")

	cancelFunc()
	waitGroup.Wait()

	log.Printf("FERP server shut down")
}
