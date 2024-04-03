package util

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 自定义token生成方法
func FormToken(username string) string {
	u := uuid.NewV4().String()
	hash := sha256.New()
	hash.Write([]byte(u + time.Now().Format(time.RFC1123) + username))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}
