package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"microE/EncryptionMicroservice"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8088", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := encryptionMicroservice.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := encryptionMicroservice.Endpoints{
		TemplateEndpoint:		encryptionMicroservice.MakeTemplateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("EncryptionMicroservice is listening on port:", *httpAddr)
		handler := encryptionMicroservice.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
