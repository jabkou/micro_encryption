package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"microE/UXMicroservice/UXMSCode"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var port string = "8090"

func main() {
	var (
		httpAddr = flag.String("http", ":"+port, "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := UXMSCode.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := UXMSCode.Endpoints{
		TemplateEndpoint: UXMSCode.MakeTemplateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("UXMicroservice is listening on port:", *httpAddr)
		handler := UXMSCode.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
