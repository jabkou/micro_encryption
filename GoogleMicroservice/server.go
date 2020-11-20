package microE

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux" //implements a request router and dispatcher for matching incoming requests to their respective handler.
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware) // @see https://stackoverflow.com/a/51456342

	r.Methods("GET").Path("/files").Handler(httptransport.NewServer(
		endpoints.FilesEndpoint,
		decodeFilesRequest,
		encodeResponse))

	r.Methods("GET").Path("/upload").Handler(httptransport.NewServer(
		endpoints.UploadEndpoint,
		decodeUploadRequest,
		encodeResponse))

	r.Methods("GET").Path("/download").Handler(httptransport.NewServer(
		endpoints.DownloadEndpoint,
		decodeDownloadRequest,
		encodeResponse))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
