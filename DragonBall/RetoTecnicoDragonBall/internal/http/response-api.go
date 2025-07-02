package internal

import (
	"encoding/json"
	"net/http"
)

type ResponseAPI[T any] struct {
	ErrorCode string `json:"errorCode"`
	Data      T      `json:"data"`
	IsSuccess bool   `json:"isSuccess"`
	IsWarning bool   `json:"isWarning"`
	Message   string `json:"message"`
}

func (ResponseAPI[T]) GetIsSuccess(result T) *ResponseAPI[T] {
	return &ResponseAPI[T]{
		Data:      result,
		ErrorCode: "",
		IsSuccess: true,
		IsWarning: false,
		Message:   "",
	}
}

func (r *ResponseAPI[T]) GetIsWarning(message string) *ResponseAPI[T] {
	var data T
	return &ResponseAPI[T]{
		Data:      data,
		ErrorCode: "",
		IsSuccess: true,
		IsWarning: true,
		Message:   message,
	}
}

type Response struct {
	Data   interface{} `json:"result,omitempty"`
	Status int         `json:"status,omitempty"`
}

func HttpResult(result interface{}, status int) *Response {
	return &Response{
		Status: status,
		Data:   result,
	}
}
func (r *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	return json.NewEncoder(w).Encode(r.Data)
}
