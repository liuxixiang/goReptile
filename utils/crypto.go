package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(input string) (hash string) {
	h := md5.New()
	h.Write([]byte(input))
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

func Sha1(input string) (hash string) {
	h := sha1.New()
	h.Write([]byte(input))
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}
