package GMSCode

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	FilesEndpoint		endpoint.Endpoint
	UploadEndpoint		endpoint.Endpoint
	DownloadEndpoint	endpoint.Endpoint
}


func MakeFilesEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(filesRequest)
		f, err := srv.Files(ctx)
		if err != nil {
			return filesResponse{f}, nil
		}
		return filesResponse{f}, nil
	}
}

func MakeUploadEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadRequest)
		f, err := srv.Upload(ctx, req.Upload)
		if err != nil {
			return uploadResponse{f}, nil
		}
		return uploadResponse{f}, nil
	}
}

func MakeDownloadEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(downloadRequest)
		f, err := srv.Download(ctx, req.Download)
		if err != nil {
			return downloadResponse{f}, nil
		}
		return downloadResponse{f}, nil
	}
}

func (e Endpoints) Files(ctx context.Context) ([2][]string, error) {
	req := filesRequest{}
	resp, err := e.FilesEndpoint(ctx, req)
	if err != nil {
	}
	filesResp := resp.(filesResponse)

	return filesResp.Files, nil
}

func (e Endpoints) Upload(ctx context.Context) (string, error) {
	req := uploadRequest{}
	resp, err := e.UploadEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	uploadResp := resp.(uploadResponse)

	return uploadResp.Upload, nil
}

func (e Endpoints) Download(ctx context.Context) (string, error) {
	req := downloadRequest{}
	resp, err := e.DownloadEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	downloadResp := resp.(downloadResponse)

	return downloadResp.Download, nil
}
