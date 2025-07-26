package shared

import (
	"encoding/json"
	"net/http"
)

func SuccessResp(w http.ResponseWriter, body any) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(body)
}

type FailResp struct {
	Message string `json:"message"`
}

func FailedResp(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	message := FailResp{
		Message: msg,
	}
	json.NewEncoder(w).Encode(message)
}
