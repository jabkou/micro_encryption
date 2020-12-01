package UXMSCode

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

//type templateRequest struct {}

type encryptAndUploadRequest struct {
	Route string `json:"route"`
	Filename string `json:"filename"`
	Password string `json:"password"`
}

type decryptAndDownloadRequest struct {
	Route string `json:"route"`
	FileId string `json:"fileId"`
	Password string `json:"password"`
}

//type templateResponse struct {
//	Template string `json:"template"`
//}

type encryptAndUploadResponse struct {
	Response string `json:"response"`
}

type decryptAndDownloadResponse struct {
	Response string `json:"response"`
}


//func decodeTemplateRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
//	var req templateRequest
//	return req, nil
//}

func decodeEncryptAndUploadRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req encryptAndUploadRequest

	route, ok := r.URL.Query()["route"]
	if !ok || len(route[0]) < 1 {
		log.Println("Url Param 'route' is missing")
	}

	filename, ok := r.URL.Query()["filename"]
	if !ok || len(filename[0]) < 1 {
		log.Println("Url Param 'filename' is missing")
	}

	password, ok := r.URL.Query()["password"]
	if !ok || len(password[0]) < 1 {
		log.Println("Url Param 'password' is missing")
	}

	req.Route = route[0]
	req.Filename = filename[0]
	req.Password = password[0]

	return req, nil
}

func decodeDecryptAndDownload(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req decryptAndDownloadRequest

	route, ok := r.URL.Query()["route"]
	if !ok || len(route[0]) < 1 {
		log.Println("Url Param 'route' is missing")
	}

	fileId, ok := r.URL.Query()["fileId"]
	if !ok || len(fileId[0]) < 1 {
		log.Println("Url Param 'fileId' is missing")
	}

	password, ok := r.URL.Query()["password"]
	if !ok || len(password[0]) < 1 {
		log.Println("Url Param 'password' is missing")
	}

	req.Route = route[0]
	req.FileId = fileId[0]
	req.Password = password[0]

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
