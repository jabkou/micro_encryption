package UXMSCode

import (
	"context"
)

type Service interface {
	Template(ctx context.Context) (string, error)
	EncryptAndUpload(ctx context.Context, password string, route string, fileName string) (string, error)
}

type uxService struct{}

func NewService() Service {
	return uxService{}
}

func (uxService) Template(ctx context.Context) (string, error) {

	return "template", nil
}

func (uxService) EncryptAndUpload(ctx context.Context, password string, route string, fileName string) (string, error) {
	return "template", nil
}
