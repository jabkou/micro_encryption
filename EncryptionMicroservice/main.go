package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"microE/EncryptionMicroservice/EMSCode"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8088", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := EMSCode.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := EMSCode.Endpoints{
		TemplateEndpoint: EMSCode.MakeTemplateEndpoint(srv),
		EncryptEndpoint:  EMSCode.MakeEncryptionEndpoint(srv),
		DecryptEndpoint:  EMSCode.MakeDecryptionEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("EncryptionMicroservice is listening on port:", *httpAddr)
		handler := EMSCode.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
