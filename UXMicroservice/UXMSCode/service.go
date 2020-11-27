package UXMSCode

import (
	"context"
	"log"
	"net/http"
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


	req1, err := http.NewRequest("GET", "localhost:8070/encrypt?route="+route+"&filename="+fileName+"&password="+password, nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
	}

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	println(resp1)

	req2, err := http.NewRequest("GET", "localhost:8080/upload?route="+route+"&filename="+fileName+"&password="+password, nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
	}

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	println(resp2)




	return "template", nil
}
