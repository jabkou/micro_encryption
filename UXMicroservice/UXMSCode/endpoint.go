package UXMSCode

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	//TemplateEndpoint			endpoint.Endpoint
	EncryptAndUploadEndpoint	endpoint.Endpoint
	DecryptAndDownloadEndpoint	endpoint.Endpoint

}

//func MakeTemplateEndpoint(srv Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//		_ = request.(templateRequest)
//		f, err := srv.Template(ctx)
//		if err != nil {
//			return templateResponse{f}, nil
//		}
//		return templateResponse{f}, nil
//	}
//}

func MakeEncryptAndUploadEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(encryptAndUploadRequest)
		f, err := srv.EncryptAndUpload(ctx, req.Password, req.Route, req.Filename)
		if err != nil {
			return encryptAndUploadResponse{f}, err
		}
		return encryptAndUploadResponse{f}, nil
	}
}

func MakeDecryptAndDownloadEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(decryptAndDownloadRequest)
		f, err := srv.DecryptAndDownload(ctx, req.Password, req.Route, req.FileId)
		if err != nil {
			return decryptAndDownloadResponse{f}, err
		}
		return decryptAndDownloadResponse{f}, nil
	}
}

//func (e Endpoints) Template(ctx context.Context) (string, error) {
//	req := templateRequest{}
//	resp, err := e.TemplateEndpoint(ctx, req)
//	if err != nil {
//	}
//	templateResp := resp.(templateResponse)
//
//	return templateResp.Template, nil
//}

func (e Endpoints) EncryptAndUpload(ctx context.Context) (string, error) {
	req := encryptAndUploadRequest{}
	resp, err := e.EncryptAndUploadEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	response := resp.(encryptAndUploadResponse)

	return response.Response, nil
}

func (e Endpoints) DecryptAndDownload(ctx context.Context) (string, error) {
	req := decryptAndDownloadRequest{}
	resp, err := e.DecryptAndDownloadEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	response := resp.(decryptAndDownloadResponse)

	return response.Response, nil
}