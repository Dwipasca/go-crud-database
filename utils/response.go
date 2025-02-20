package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJson(w http.ResponseWriter, code int, status string, data interface{}, message string) {
	var req = Response{
		Code:    code,
		Status:  status,
		Data:    data,
		Message: message,
	}

	res, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}