package EMSCode

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Service interface {
	Template(ctx context.Context) (string, error)
	Encrypt(ctx context.Context, route string, filename string, password string) (string, error)
	Decrypt(ctx context.Context, route string, filename string, password string) (string, error)
}

type googService struct{}

func NewService() Service {
	return googService{}
}

func (googService) Template(ctx context.Context) (string, error) {

	return "template", nil
}

func DeriveKey(password, salt []byte, route string) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
		saltPath := route+"/salt"
		_ = ioutil.WriteFile(saltPath, salt, 0777)
	}
	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}

func (googService) Encrypt(ctx context.Context, route string, filename string, password string) (string, error) {

	fullPath := route+"/"+filename
	infile, err := os.Open(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	//key, err := ioutil.ReadFile("key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	pass := []byte(password)
	key, _, err := DeriveKey(pass, nil, route)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	outfilePath := route+"/"+filename+".bin"
	outfile, err := os.OpenFile(outfilePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	// Append the IV
	outfile.Write(iv)

	return "template", nil
}

func (googService) Decrypt(ctx context.Context, route string, filename string, password string) (string, error) {

	fullPath := route+"/"+filename

	infile, err := os.Open(fullPath+".bin")
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)

	pass := []byte(password)
	saltPath := route+"/salt"
	dat, _ := ioutil.ReadFile(saltPath)
	key, _, err := DeriveKey(pass, dat, route)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	fi, err := infile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		log.Fatal(err)
	}

	outfile, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}

	return "template", nil
}