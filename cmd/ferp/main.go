package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/satimoto/go-ferp/internal/converter"
	"github.com/satimoto/go-ferp/internal/exchange"
	"github.com/satimoto/go-ferp/internal/rest"
	"github.com/satimoto/go-ferp/internal/rpc"
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "ferp",
	Short: "Run the Fiat Exchange Rate Provider",
	Long:  "Run the Fiat Exchange Rate Provider",
	Run:   startFerp,
}

func main() {
	configFile, err := os.UserHomeDir()

	if err == nil {
		configFile = configFile + "/.ferp/"
	}

	configFile = configFile + "ferp.conf"

	runCommand.Flags().StringP("configfile", "C", configFile, "Config")
	runCommand.Execute()
}

func startFerp(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("configfile")

	godotenv.Load(configFile)

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

	sigtermChan := make(chan os.Signal, 1)
	signal.Notify(sigtermChan, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigtermChan

	log.Printf("Shutting down FERP server")

	cancelFunc()
	waitGroup.Wait()

	log.Printf("FERP server shut down")
}
