package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

func PostInvokeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, err.Error())
		return
	}

	fcn := r.FormValue("fcn")
	args := strings.Split(r.FormValue("args"), ",")

	if fcn == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Fcn is required")
		return
	}

	resultString, err := api.Invoke(&api.FscInstance, vars["channelId"], vars["chaincodeId"], fcn, args)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, resultString))
	if err != nil {
		panic(err)
	}
}