package main

import (
	"fabric-rest-api-go/pkg/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWelcome(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, handlers.WelcomeHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Fabric REST Api")
	}
}
