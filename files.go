package gotools

import (
	"io/ioutil"
	"os"
)

func Write_to_file(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func File_exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
