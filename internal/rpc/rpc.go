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

func NewRpc() Rpc {
	return &RpcService{
		server:          grpc.NewServer(),
		RpcRateResolver: rate.NewResolver(),
	}
}

func (rs *RpcService) StartRpc(ctx context.Context, waitGroup *sync.WaitGroup) {
	log.Printf("Starting Rpc service")
	waitGroup.Add(1)

	go rs.listenAndServe()

	go func() {
		<-ctx.Done()
		log.Printf("Shutting down Rpc service")

		rs.shutdown()

		log.Printf("Rpc service shut down")
		waitGroup.Done()
	}()
}

func (rs *RpcService) GetRateService() *rate.RpcRateResolver {
	return rs.RpcRateResolver
}

func (rs *RpcService) listenAndServe() {
	rpcPort := os.Getenv("RPC_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcPort))

	if err != nil {
		util.LogOnError("FERP003", "Error creating network address", err)
		log.Printf("FERP003: RpcPort=%v", rpcPort)
	}

	ferprpc.RegisterRateServiceServer(rs.server, rs.RpcRateResolver)
	err = rs.server.Serve(listener)

	if err != nil {
		util.LogOnError("FERP004", "Error in Rpc service", err)
	}
}

func (rs *RpcService) shutdown() {
	rs.server.GracefulStop()
}
