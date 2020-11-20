package encryptionMicroservice

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	TemplateEndpoint		endpoint.Endpoint
}


func MakeTemplateEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(templateRequest)
		f, err := srv.Template(ctx)
		if err != nil {
			return templateResponse{f}, nil
		}
		return templateResponse{f}, nil
	}
}

func (e Endpoints) Template(ctx context.Context) (string, error) {
	req := templateRequest{}
	resp, err := e.TemplateEndpoint(ctx, req)
	if err != nil {
	}
	templateResp := resp.(templateResponse)

	return templateResp.Template, nil
}
