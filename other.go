package gotools

import (
	"crypto/md5"
	"fmt"
)

func Md5(s string) string {
	b := []byte{s}
	return fmt.Sprintf("%x", md5.Sum(b))
}
