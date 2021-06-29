package grpclogger

import (
	"github.com/goinsane/xlog"
)

type GrpcLogger struct {
	*xlog.Logger
}

func (g *GrpcLogger) V(v int) bool {
	return g.Logger.V(xlog.Verbose(v)) != nil
}
