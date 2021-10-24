package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fernandoocampo/tracing/internal/servicea"
)

func main() {
	serviceBClient := servicea.NewServiceBClient("http://localhost:8087")
	service := servicea.NewService(serviceBClient)
	endpoints := servicea.NewEndpoints(service)

	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// HTTP user
	go func() {
		log.Println("msg", "starting http server service a", "http: 8890")
		handler := servicea.NewHTTPServer(endpoints)
		errChan <- http.ListenAndServe(":8890", handler)
	}()

	errmsg := <-errChan
	fmt.Println("Ending service a:", errmsg)
}
