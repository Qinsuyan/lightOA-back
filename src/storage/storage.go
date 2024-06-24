package storage

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"lightOA-end/src/entity"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var passphrase []byte
var dir string

func Init(path string, pass string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	combined := hostname + pass
	hash := sha256.Sum256([]byte(combined))
	passphrase = hash[:]
	//检查文件夹是否存在
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			//创建文件夹
			err = os.MkdirAll(absPath, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if !info.IsDir() {
		return errors.New("file path is not a path")
	}
	//检查权限
	testFile := filepath.Join(absPath, ".test")
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()
	os.Remove(testFile)
	dir = absPath
	return nil
}

func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	zipped, err := compress(ciphertext)
	if err != nil {
		return nil, err
	}
	return zipped, nil
}

func decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	unzipped, err := decompress(data)
	if err != nil {
		return nil, err
	}
	iv := unzipped[:aes.BlockSize]
	ciphertext := unzipped[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

func SaveFile(fileHeader *multipart.FileHeader, allowTypes map[string]bool) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(fileData)
	if !allowTypes[fileType] {
		return "invalid format", errors.New("file type not allowed")
	}
	encryptedData, err := encrypt(fileData)
	if err != nil {
		return "", err
	}

	uuid := uuid.New()
	fullPath := fmt.Sprintf("%s/%s", dir, uuid)
	err = os.WriteFile(fullPath, encryptedData, 0644)
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress decompresses data using gzip.
func decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}

func GetFile(info *entity.File) ([]byte, error) {
	fullPath := fmt.Sprintf("%s/%s", dir, info.UUID)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return decrypt(data)
}
