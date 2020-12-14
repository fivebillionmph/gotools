package server

import (
	"net/http"
	"encoding/json"
)


func Send_json_response(w http.ResponseWriter, status_code int, data interface{}) {
	j_data, err := json.Marshal(data)
	if err != nil {
		Send_text_message(w, 500, "Unexpected error")
		return
	}
	w.WriteHeader(status_code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j_data)
}

func Send_text_message(w http.ResponseWriter, status_code int, message string) {
	w.WriteHeader(status_code)
	w.Write([]byte(message))
}
