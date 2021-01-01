package server

import (
	"encoding/json"
	"net/http"
)

func Retrieve_json_request(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(s)
	return err
}
