package encryptionMicroservice

import (
	"context"
	"encoding/json"
	"net/http"
)

//type templateRequest struct {}

type encryptionRequest struct {}

type decryptionRequest struct {}

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
	return req, nil
}

func decodeDecryptionRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req decryptionRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
