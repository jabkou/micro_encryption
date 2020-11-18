package microE

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type filesRequest struct {}

type filesResponse struct {
	Files string `json:"fil"`
}

type uploadRequest struct {
	Upload string `json:"upload"`
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

//func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
//	var req validateRequest
//	err := json.NewDecoder(r.Body).Decode(&req)
//	if err != nil {
//		return nil, err
//	}
//	return req, nil
//}

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

	key := upload[0]
	req.Upload = key
	return req, nil
}

func decodeDownloadRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req downloadRequest

	download, ok := r.URL.Query()["download"]

	if !ok || len(download[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	key := download[0]
	req.Download = key
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
