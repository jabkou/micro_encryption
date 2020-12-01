package GMSCode

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	FilesEndpoint		endpoint.Endpoint
	UploadEndpoint		endpoint.Endpoint
	DownloadEndpoint	endpoint.Endpoint
	GetAuthCodeEndpoint	endpoint.Endpoint
}


func MakeFilesEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(filesRequest)
		f, err := srv.Files(ctx)
		if err != nil {
			return filesResponse{f}, err
		}
		return filesResponse{f}, nil
	}
}

func MakeUploadEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uploadRequest)
		f, err := srv.Upload(ctx, req.Upload, req.Route)
		if err != nil {
			return uploadResponse{f}, err
		}
		return uploadResponse{f}, nil
	}
}

func MakeDownloadEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(downloadRequest)
		f, err := srv.Download(ctx, req.Download, req.Route)
		if err != nil {
			return downloadResponse{f}, err
		}
		return downloadResponse{f}, nil
	}
}

func MakeGetAuthCodeEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAuthCodeRequest)
		f, err := srv.GetAuthCode(ctx, req.AuthCode)
		if err != nil {
			return getAuthCodeResponse{f}, err
		}
		return getAuthCodeResponse{f}, nil
	}
}

func (e Endpoints) Files(ctx context.Context) ([2][]string, error) {
	req := filesRequest{}
	resp, err := e.FilesEndpoint(ctx, req)
	if err != nil {
		return [2][]string{}, err
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

	return uploadResp.Response, nil
}

func (e Endpoints) Download(ctx context.Context) (string, error) {
	req := downloadRequest{}
	resp, err := e.DownloadEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	downloadResp := resp.(downloadResponse)

	return downloadResp.Response, nil
}

func (e Endpoints) GetAuthCode(ctx context.Context) (string, error) {
	req := getAuthCodeRequest{}
	resp, err := e.GetAuthCodeEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	downloadResp := resp.(getAuthCodeResponse)

	return downloadResp.Response, nil
}
