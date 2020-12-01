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
	//Template(ctx context.Context) (string, error)
	Encrypt(ctx context.Context, route string, filename string, password string) (string, error)
	Decrypt(ctx context.Context, route string, filename string, password string) (string, error)
}

type encryptionService struct{}

func NewService() Service {
	return encryptionService{}
}

//func (encryptionService) Template(ctx context.Context) (string, error) {
//
//	return "template", nil
//}

func DeriveKey(password, salt []byte, route string, fileName string) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
		saltPath := route+"/salts/salt-"+fileName
		_ = ioutil.WriteFile(saltPath, salt, 0777)
	}
	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}

func (encryptionService) Encrypt(ctx context.Context, route string, filename string, password string) (string, error) {

	fullPath := route+"/"+filename
	infile, err := os.Open(fullPath)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}
	defer infile.Close()

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	//key, err := ioutil.ReadFile("key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	pass := []byte(password)
	key, _, err := DeriveKey(pass, nil, route, filename)
	if err != nil {
		return "Error: "+err.Error(), err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	outfilePath := route+"/encrypted/"+filename+".bin"
	outfile, err := os.OpenFile(outfilePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
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

	return "OK", nil
}

func (encryptionService) Decrypt(ctx context.Context, route string, filename string, password string) (string, error) {

	fullPath := route+"/downloaded/"+filename

	infile, err := os.Open(fullPath+".bin")
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}
	defer infile.Close()

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)

	pass := []byte(password)
	saltPath := route+"/salts/salt-"+filename
	dat, err := ioutil.ReadFile(saltPath)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}
	key, _, err := DeriveKey(pass, dat, route, filename)
	if err != nil {
		return "Error: "+err.Error(), err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	fi, err := infile.Stat()
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
	}

	outPath := route+"/decrypted/"+filename
	outfile, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return "Error: "+err.Error(), err
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

	return "OK", nil
}