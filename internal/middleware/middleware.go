package middleware

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Response struct {
	Data interface{} `json:"data"`
	Err  *Error      `json:"error"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Details interface{} `json:"data"`
}

func WriteSuccessData(w http.ResponseWriter, data interface{}) {
	_ = jsoniter.NewEncoder(w).Encode(Response{
		Data: data,
	})

	w.WriteHeader(200)
}

func WriteErrorResponse(w http.ResponseWriter, errCode int, msg string) {
	_ = jsoniter.NewEncoder(w).Encode(Response{
		Err: &Error{
			Code:    errCode,
			Message: msg,
		},
	})

	w.WriteHeader(200)
}
