package UXMSCode

import (
	"context"
	"encoding/json"
	"net/http"
)

type templateRequest struct {}

type templateResponse struct {
	Template string `json:"template"`
}


func decodeTemplateRequest(ctx context.Context, r *http.Request) (interface{}, error ) {
	var req templateRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
