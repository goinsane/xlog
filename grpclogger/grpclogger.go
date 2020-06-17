package grpclogger

import (
	"github.com/goinsane/xlog"
)

type GrpcLogger struct {
	*xlog.Logger
}

func (g *GrpcLogger) V(l int) bool {
	return g.Logger.V(xlog.Verbose(l)) != nil
}
