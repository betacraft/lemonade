package framework

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

func MD5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GenerateOtp() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, 5)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
