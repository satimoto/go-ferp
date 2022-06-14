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
	"github.com/satimoto/go-ferp/internal/rest"
	"github.com/satimoto/go-ferp/internal/rpc"
)

func main() {
	log.Printf("Starting up FERP server")
	shutdownCtx, cancelFunc := context.WithCancel(context.Background())
	waitGroup := &sync.WaitGroup{}

	restService := rest.NewRest()
	restService.StartRest(shutdownCtx, waitGroup)

	rpcService := rpc.NewRpc()
	rpcService.StartRpc(shutdownCtx, waitGroup)

	exchangeService := exchange.NewService(rpcService)
	exchangeService.Start(shutdownCtx, waitGroup)

	converterService := converter.NewService(exchangeService, rpcService)
	converterService.Start(shutdownCtx, waitGroup)

	sigtermChan := make(chan os.Signal)
	signal.Notify(sigtermChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigtermChan

	log.Printf("Shutting down FERP server")

	cancelFunc()
	waitGroup.Wait()

	log.Printf("FERP server shut down")
}
