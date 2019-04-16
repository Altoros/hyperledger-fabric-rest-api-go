package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type userCredentials struct {
	OrgName  string `json:"orgName"`
	Username string `json:"username"`
}

func PostUsersHandler(c echo.Context) error {
	userCredentials := new(userCredentials)
	if err := c.Bind(userCredentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = userCredentials.Username
	claims["org"] = userCredentials.OrgName
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
