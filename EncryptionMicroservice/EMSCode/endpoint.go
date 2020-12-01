package EMSCode

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	//TemplateEndpoint		endpoint.Endpoint
	EncryptEndpoint		endpoint.Endpoint
	DecryptEndpoint		endpoint.Endpoint
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

func MakeEncryptionEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(encryptionRequest)
		f, err := srv.Encrypt(ctx, req.Route, req.Filename, req.Password)
		if err != nil {
			return encryptionResponse{f}, err
		}
		return encryptionResponse{f}, nil
	}
}

func MakeDecryptionEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(decryptionRequest)
		f, err := srv.Decrypt(ctx, req.Route, req.Filename, req.Password)
		if err != nil {
			return decryptionResponse{f}, err
		}
		return decryptionResponse{f}, nil
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

func (e Endpoints) Encrypt(ctx context.Context) (string, error) {
	req := encryptionRequest{}
	resp, err := e.EncryptEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	response := resp.(encryptionResponse)

	return response.Response, nil
}

func (e Endpoints) Decrypt(ctx context.Context) (string, error) {
	req := decryptionRequest{}
	resp, err := e.DecryptEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	response := resp.(decryptionResponse)

	return response.Response, nil
}
