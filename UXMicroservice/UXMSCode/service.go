package UXMSCode

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Service interface {
	//Template(ctx context.Context) (string, error)
	EncryptAndUpload(ctx context.Context, password string, route string, fileName string) (string, error)
	DecryptAndDownload(ctx context.Context, password string, route string, fileId string) (string, error)
}

type uxService struct{}

func NewService() Service {
	return uxService{}
}

//func (uxService) Template(ctx context.Context) (string, error) {
//
//	return "template", nil
//}

func (uxService) EncryptAndUpload(ctx context.Context, password string, route string, fileName string) (string, error) {


	req1, err := http.NewRequest("GET", "http://localhost:8070/encrypt?route="+route+"&filename="+fileName+"&password="+password, nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
		return "Error on request: "+err.Error(), err
	}

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil || resp1 == nil || resp1.StatusCode != 200 {
		log.Println("Error on response.\n[ERRO] -", err)
		err = errors.New("Error")
		return "Error on response", err
	}

	req2, err := http.NewRequest("GET", "http://localhost:8080/upload?upload="+fileName+".bin&route="+route, nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
		return "Error on request: "+err.Error(), err
	}

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil || resp2 == nil || resp2.StatusCode != 200 {
		log.Println("Error on response.\n[ERRO] -", err)
		return "Error on response", err
	}

	return "OK", nil
}

func (uxService) DecryptAndDownload(ctx context.Context, password string, route string, fileId string) (string, error) {


	req1, err := http.NewRequest("GET", "http://localhost:8080/download?download="+fileId+"&route="+route, nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
		return "Error on request: "+err.Error(), err
	}

	client1 := &http.Client{}
	resp1, err := client1.Do(req1)
	if err != nil || resp1 == nil || resp1.StatusCode != 200 {
		log.Println("Error on response.\n[ERRO] -", err)
		return "Error on response: ", err
	}

	defer resp1.Body.Close()
	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		log.Println("Error on reading body.\n[ERRO] -", err)
		return "Error: "+err.Error(), err
	}
	body2 := string(body)

	r, _ := regexp.Compile("(?::)\".+\"")

	body3 := r.FindString(body2)

	//fileName2 := strings.TrimLeft(body3, ":\"")
	//fileName := strings.TrimRight(fileName2, ".bin\"")

	fileName := strings.Trim(body3, ":\" .bin\"")


	req2, err := http.NewRequest("GET", "http://localhost:8070/decrypt?route="+route+"&filename="+fileName+"&password="+password, nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
		return "Error: "+err.Error(), err
	}

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil || resp2 == nil || resp2.StatusCode != 200 {
		log.Println("Error on response.\n[ERRO] -", err)
		return "Error on response", err
	}
	println(resp2.Status)

	return "OK", nil
}