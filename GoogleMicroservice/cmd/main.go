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

	"microE"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := microE.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := microE.Endpoints{
		FilesEndpoint:		microE.MakeFilesEndpoint(srv),
		UploadEndpoint:		microE.MakeUploadEndpoint(srv),
		DownloadEndpoint:	microE.MakeDownloadEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("microE is listening on port:", *httpAddr)
		handler := microE.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
