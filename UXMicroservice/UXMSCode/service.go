package UXMSCode

import (
	"context"
)

type Service interface {
	Template(ctx context.Context) (string, error)
}

type googService struct{}

func NewService() Service {
	return googService{}
}

func (googService) Template(ctx context.Context) (string, error) {

	return "template", nil
}
