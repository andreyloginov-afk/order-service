package httph

import (
	"encoding/json"
	"errors"
	"net/http"
)

type httpCoder interface {
	error
	HTTPStatus() int
}

func SendJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

func SendEmpty(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorApply(r, err)

	var hc httpCoder
	if errors.As(err, &hc) {
		sendError(w, hc.HTTPStatus(), hc)
		return
	}

	sendError(w, http.StatusInternalServerError, err)
}

func sendError(w http.ResponseWriter, statusCode int, err error) {
	SendJSON(w, statusCode, map[string]string{"error": err.Error()})
}
