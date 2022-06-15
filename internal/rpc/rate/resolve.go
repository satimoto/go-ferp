package rate

import (
	"context"

	"github.com/satimoto/go-ferp/internal/converter"
)

type RpcRateResolver struct {
	ConverterService converter.Converter
	Shutdown         context.CancelFunc
	shutdownCtx      context.Context
}

func NewResolver(converterService converter.Converter) *RpcRateResolver {
	shutdownCtx, shutdown := context.WithCancel(context.Background())

	return &RpcRateResolver{
		ConverterService: converterService,
		Shutdown: shutdown,
		shutdownCtx: shutdownCtx,
	}
}
