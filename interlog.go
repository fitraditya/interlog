package interlog

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const (
	TextFormat = 0
	JSONFormat = 1
)

// Interlog is a server interceptor for gRPC logging
type Interlog struct {
	// Env: "staging", "production", or whatever environment you defined
	Env string
	// Logger: logrus instance
	Logger *logrus.Logger
}

// New() returns a new logging interceptor
func New() *Interlog {
	log := logrus.New()
	// Default config
	// TO DO: Add method for override default config
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	return &Interlog{
		Env:    "development",
		Logger: log,
	}
}

// Unary() returns a server interceptor method to logging unary gRPC call
func (interlog *Interlog) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ts := time.Now()
		ip := "0.0.0.0"

		if pr, ok := peer.FromContext(ctx); ok {
			if tcp, ok := pr.Addr.(*net.TCPAddr); ok {
				ip = tcp.IP.String()
			} else {
				ip = pr.Addr.String()
			}
		}

		resp, err = handler(ctx, req)
		end := time.Since(ts)

		interlog.Logger.WithFields(logrus.Fields{
			"env":          interlog.Env,
			"remote_addr":  ip,
			"full_method":  info.FullMethod,
			"status":       grpc.Code(err),
			"request_time": end.Milliseconds(),
		}).Infof(grpc.Code(err).String())

		return
	}
}
