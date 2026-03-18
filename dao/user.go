package dao

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
