package GMSCode

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Pair struct {
	name string
	id string
}
var authCode string = ""

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "../data/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	tokenSource := config.TokenSource(oauth2.NoContext, tok)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	if newToken.AccessToken != tok.AccessToken {
		saveToken(tokFile, newToken)
		log.Println("Saved new token:", newToken.AccessToken)
	}


	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)


	//if _, err := log.Scan(&authCode); err != nil {
	//	log.Fatalf("Unable to read authorization code %v", err)
	//}

	var i = 0
	for ; i < 100; {
		i += 1
		if authCode == "" {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if  authCode == "" {
		log.Fatalf("Unable to read authorization code ")
	}


	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//validate user
func validateUser() (*oauth2.Config, error){
	b, err := ioutil.ReadFile("../data/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config, err
}

//getUser
func getUser() (*drive.Service, error){
	config, err := validateUser()

	client := getClient(config)
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return srv, err
}

func createDir(service *drive.Service, name string, parentId string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := service.Files.Create(d).Do()

	if err != nil {
		log.Println("Could not create dir: " + err.Error())
		return nil, err
	}

	return file, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func getService() (*drive.Service, error) {
	b, err := ioutil.ReadFile("../data/credentials.json")
	if err != nil {
		log.Printf("Unable to read credentials.json file. Err: %v\n", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)

	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, err := drive.New(client)

	if err != nil {
		log.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return service, err
}

type Service interface {
	Files(ctx context.Context) ([2][]string, error)
	Upload(ctx context.Context, fileName string, route string) (string, error)
	Download(ctx context.Context, fileId string, route string) (string, error)
	GetAuthCode(ctx context.Context, authCode string) (string, error)
}

type googService struct{}

func NewService() Service {
	return googService{}
}

func (googService) Files(ctx context.Context) ([2][]string, error) {

	var temp [2][]string
	srv, err := getUser()
	if err != nil {
		log.Println("Error while getting user.\n[ERRO] -", err)
		return temp, err
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Println("Unable to retrieve files: %v", err)
		return temp, err
	}
	log.Println("Files:")
	if len(r.Files) == 0 {
		log.Println("No files found.")
	} else {
		for _, i := range r.Files {
			log.Printf("%s (%s)\n", i.Name, i.Id)
			temp[0] = append(temp[0], i.Name)
			temp[1] = append(temp[1], i.Id)
		}
	}
	return temp, nil
}

func (googService) Upload(ctx context.Context, fileName string, route string) (string, error) {

	fullPath := route+"/encrypted/"+fileName
	// Step 1. Open the file
	f, err := os.Open(fullPath)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	defer f.Close()
	buffer := make([]byte, 512)

	// Step 2. Get the Google Drive service
	service, err := getService()
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	// Step 3. Create the directory
	dir, err := createDir(service, "My Folder", "root")

	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	// Step 4. Create the file and upload its content

	file, err := createFile(service, fileName, http.DetectContentType(buffer), f, dir.Id)

	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	log.Printf("File '%s' successfully uploaded in '%s' directory", file.Name,  dir.Name)
	return "OK", nil
}

func (googService) Download(ctx context.Context, fileId string, route string) (string, error) {


	srv, err := getUser()
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	file, err := srv.Files.Get(fileId).Do()
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}


	log.Println("Downloading file...")

	fullPath := route+"/downloaded/"
	f, err := os.Create(fullPath+file.Name)
	if err != nil {
		log.Printf("create file: %v", err)
		return "Error: "+err.Error(), err
	}
	defer f.Close()

	tok, err := tokenFromFile("../data/token.json")
	if err != nil {
		log.Printf("create file: %v", err)
		return "Error: "+err.Error(), err
	}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/drive/v3/files/" + fileId + "?alt=media", nil)
	if err != nil {
		log.Println("Error on request.\n[ERRO] -", err)
		return "Error: "+err.Error(), err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", tok.TokenType + " " + tok.AccessToken)
	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return "Error: "+err.Error(), err
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}
	log.Println(resp.Status)

	return file.Name, nil
}

func (googService) GetAuthCode(ctx context.Context, authCodes string) (string, error) {
	authCode = authCodes
	return "OK", nil
}