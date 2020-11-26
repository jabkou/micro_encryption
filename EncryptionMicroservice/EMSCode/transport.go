package EMSCode

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

//type templateRequest struct {}

type encryptionRequest struct {
	Route string `json:"route"`
	Filename string `json:"filename"`
	Password string `json:"password"`
}

type decryptionRequest struct {
	Route string `json:"route"`
	Filename string `json:"filename"`
	Password string `json:"password"`
}

//type templateResponse struct {
//	Template string `json:"template"`
//}

type encryptionResponse struct {
	Encryption string `json:"encryption"`
}

type decryptionResponse struct {
	Decryption string `json:"decryption"`
}

//func decodeTemplateRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
//	var req templateRequest
//	return req, nil
//}

func decodeEncryptionRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req encryptionRequest

	route, ok := r.URL.Query()["route"]
	if !ok || len(route[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	filename, ok := r.URL.Query()["filename"]
	if !ok || len(filename[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	password, ok := r.URL.Query()["password"]
	if !ok || len(password[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	req.Route = route[0]
	req.Filename = filename[0]
	req.Password = password[0]

	return req, nil
}

func decodeDecryptionRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req decryptionRequest

	route, ok := r.URL.Query()["route"]
	if !ok || len(route[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	filename, ok := r.URL.Query()["filename"]
	if !ok || len(filename[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	password, ok := r.URL.Query()["password"]
	if !ok || len(password[0]) < 1 {
		log.Println("Url Param 'key' is missing")
	}

	req.Route = route[0]
	req.Filename = filename[0]
	req.Password = password[0]

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
