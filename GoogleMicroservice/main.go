package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"microE/GMSCode"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var port = "8080"


func main() {
	var (
		httpAddr = flag.String("http", ":"+port, "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := GMSCode.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := GMSCode.Endpoints{
		FilesEndpoint:    GMSCode.MakeFilesEndpoint(srv),
		UploadEndpoint:   GMSCode.MakeUploadEndpoint(srv),
		DownloadEndpoint: GMSCode.MakeDownloadEndpoint(srv),
		GetAuthCodeEndpoint: GMSCode.MakeGetAuthCodeEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("GoogleMicroservice is listening on port:", *httpAddr)
		handler := GMSCode.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
