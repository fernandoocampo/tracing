package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/fernandoocampo/tracing/internal/items"
	"github.com/fernandoocampo/tracing/internal/servicec"
	"github.com/fernandoocampo/tracing/internal/tracers"
	"google.golang.org/grpc"
)

func main() {
	service := servicec.NewService()
	endpoints := servicec.NewEndpoints(service)
	grpcServer := servicec.NewGRPCServer(endpoints)

	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println("during", "Listen", "err", err)
		os.Exit(1)
	}

	// HTTP user
	go func() {
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(tracers.GRPCServerTrace()))
		pb.RegisterItemCodeServiceServer(baseServer, grpcServer)
		log.Println("msg", "starting grpc server service b", "grpc: 50051")
		errChan <- baseServer.Serve(grpcListener)
	}()

	errmsg := <-errChan
	fmt.Println("Ending service a:", errmsg)
}
