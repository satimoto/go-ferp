package rpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/satimoto/go-datastore/pkg/util"
	"github.com/satimoto/go-ferp/ferprpc"
	"github.com/satimoto/go-ferp/internal/converter"
	metrics "github.com/satimoto/go-ferp/internal/metric"
	"github.com/satimoto/go-ferp/internal/rpc/rate"
	"google.golang.org/grpc"
)

type Rpc interface {
	StartRpc(context.Context, *sync.WaitGroup)
	GetRateService() *rate.RpcRateResolver
}

type RpcService struct {
	server          *grpc.Server
	RpcRateResolver *rate.RpcRateResolver
}

func NewRpc(converterService converter.Converter) Rpc {
	return &RpcService{
		server:          grpc.NewServer(),
		RpcRateResolver: rate.NewResolver(converterService),
	}
}

func (rs *RpcService) StartRpc(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	log.Printf("Starting Rpc service")
	waitGroup.Add(1)

	go rs.listenAndServe()
	go rs.waitForShutdown(shutdownCtx, waitGroup)
}

func (rs *RpcService) GetRateService() *rate.RpcRateResolver {
	return rs.RpcRateResolver
}

func (rs *RpcService) listenAndServe() {
	rpcPort := os.Getenv("RPC_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcPort))
	util.PanicOnError("FERP003", "Error creating network address", err)

	ferprpc.RegisterRateServiceServer(rs.server, rs.RpcRateResolver)
	err = rs.server.Serve(listener)

	if err != nil {
		metrics.RecordError("FERP004", "Error in Rpc service", err)
	}
}

func (rs *RpcService) waitForShutdown(shutdownCtx context.Context, waitGroup *sync.WaitGroup) {
	<-shutdownCtx.Done()
	log.Printf("Shutting down Rpc service")

	rs.RpcRateResolver.Shutdown()
	rs.server.GracefulStop()

	log.Printf("Rpc service shut down")
	waitGroup.Done()
}
