package main

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code)
	assert.Equal(t, response.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	assert.Equal(t, true, strings.Contains(response.Body.String(), "Fabric REST Api"))
}
