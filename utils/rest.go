package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type JsonMessage struct {
	Message any `json:"message"`
}

func WriteResponse(w http.ResponseWriter, status int, message any) {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(status)
	w.Write(marshalledMessage)
}

func ParseJSON(body io.ReadCloser, v any) error {
	defer body.Close()
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}

func ValidateStruct(s any) error {
	validate := validator.New()
	return validate.Struct(s)
}

func WriteErrorResponse(w http.ResponseWriter, status int, errMsg string) {
	WriteResponse(w, status, JsonMessage{Message: errMsg})
}
