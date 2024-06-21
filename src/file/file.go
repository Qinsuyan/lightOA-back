package file

import (
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
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
