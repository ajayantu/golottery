package helpers

import (
	"encoding/json"
	"net/http"
)

type FailStruct struct {
	Message    string `json:"message"`
	ErrorField string `json:"error_field,omitempty"`
}

type FailResponse struct {
	Success bool         `json:"success"`
	Errors  []FailStruct `json:"errors,omitempty"`
}

type SuccessResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

const (
	InternalError = "Internal server error"
	ContentType   = "Content-Type"
	JsonType      = "application/json"
)

func Fail(w http.ResponseWriter, status int, errors []FailStruct) {
	res := &FailResponse{
		Success: false,
		Errors:  errors,
	}
	out, err := json.Marshal(res)
	if err != nil {
		http.Error(
			w,
			InternalError,
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(ContentType, JsonType)
	w.WriteHeader(status)
	w.Write(out)
}

func Success(w http.ResponseWriter, status int, result interface{}) {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(
			w,
			InternalError,
			http.StatusInternalServerError,
		)
		return
	}
	r := &SuccessResponse{
		Success: true,
		Data:    res,
	}

	out, err := json.Marshal(r)
	if err != nil {
		http.Error(
			w,
			InternalError,
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set(ContentType, JsonType)
	w.WriteHeader(status)
	w.Write(out)
}
