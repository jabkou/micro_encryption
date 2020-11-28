package GMSCode

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type filesRequest struct {}

type filesResponse struct {
	Files [2][]string `json:"files"`
}

type uploadRequest struct {
	Upload string `json:"upload"`
	Route  string `json:"route"`
}

type uploadResponse struct {
	Upload string `json:"uploadResponse"`
}

type downloadRequest struct {
	Download string `json:"download"`
}

type downloadResponse struct {
	Download string `json:"downloadResponse"`
}

func decodeFilesRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req filesRequest
	return req, nil
}

func decodeUploadRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req uploadRequest

	upload, ok := r.URL.Query()["upload"]

	if !ok || len(upload[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	req.Upload = upload[0]

	route, ok := r.URL.Query()["route"]

	if !ok || len(route[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	req.Route = route[0]

	return req, nil
}

func decodeDownloadRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req downloadRequest

	download, ok := r.URL.Query()["download"]

	if !ok || len(download[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	req.Download = download[0]
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
