package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fernandoocampo/tracing/internal/serviceb"
)

func main() {
	serviceCClient, err := serviceb.NewServicCGRPCClient()
	if err != nil {
		panic(err)
	}
	service := serviceb.NewService(serviceCClient)
	endpoints := serviceb.NewEndpoints(service)

	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// HTTP user
	go func() {
		log.Println("msg", "starting http server service b", "http: 8087")
		handler := serviceb.NewHTTPServer(endpoints)
		errChan <- http.ListenAndServe(":8087", handler)
	}()

	errmsg := <-errChan
	fmt.Println("Ending service a:", errmsg)
}
