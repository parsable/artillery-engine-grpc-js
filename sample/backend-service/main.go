package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	backend_services_v1 "github.com/kenju/artillery-engine-grpc/sample/backend-service/backend/services/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

const (
	defaultAddr = ":8080"
)

func init() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}

func main() {
	port := getEnv("ADDR", defaultAddr)
	listenPort, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	severOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	}

	server := grpc.NewServer(severOpts...)

	backend_services_v1.RegisterHelloServiceServer(server, newBackendServer())

	reflection.Register(server)
	server.Serve(listenPort)
}

//--------------------------------
// utility
//--------------------------------

func getEnv(key, defaultVal string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return defaultVal
	}
	return v
}

//--------------------------------
// backend application server
//--------------------------------
type backendServer struct {
}

func newBackendServer() *backendServer {
	return &backendServer{}
}

func (bs *backendServer) Hello(
	ctx context.Context,
	req *backend_services_v1.HelloRequest,
) (*backend_services_v1.HelloResponse, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Info("Hello()")

	// NOTE: sleep for rondom milliseconds for benchmarking
	r := rand.Intn(100) // up to 100 msec
	time.Sleep(time.Duration(r) * time.Millisecond)

	return &backend_services_v1.HelloResponse{
		Message: fmt.Sprintf("world (code=%d)", codes.OK),
	}, nil
}

func (bs *backendServer) Bye(
	ctx context.Context,
	req *backend_services_v1.ByeRequest,
) (*backend_services_v1.ByeResponse, error) {
	log.WithFields(log.Fields{
		"request": req,
	}).Info("Bye()")

	// NOTE: sleep for rondom milliseconds for benchmarking
	r := rand.Intn(500) // up to 500 msec
	time.Sleep(time.Duration(r) * time.Millisecond)

	return &backend_services_v1.ByeResponse{
		Message: fmt.Sprintf("bye (code=%d)", codes.OK),
	}, nil
}
