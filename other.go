package gotools

import (
	"crypto/md5"
	"fmt"
	"time"
)

func Md5(s string) string {
	b := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(b))
}

func Unix_timestamp() int {
	return int(time.Now().Unix())
}
