package handlers

import (
	"io"
	"net/http"
)

// Generic wrapper for GET handlers
func GetHandlerWrapper(w http.ResponseWriter, _ *http.Request, handlerFunc func() (string, error)) {
	jsonString, err := handlerFunc()
	GetJsonOutputWrapper(w, jsonString, err)
}

func GetJsonOutputWrapper(w http.ResponseWriter, jsonString string, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, jsonString)
	if err != nil {
		panic(err)
	}
}

