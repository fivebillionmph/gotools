package gotools

import (
	"net/http"
	"time"
	"io/ioutil"
)

func Get_request(url string) (result []byte, err error) {
	http_client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil { return }
	res, err := http_client.Do(req)
	if err != nil { return }
	if res.Body != nil {
		defer res.Body.Close()
	}
	result, err = ioutil.ReadAll(res.Body)

	return
}
