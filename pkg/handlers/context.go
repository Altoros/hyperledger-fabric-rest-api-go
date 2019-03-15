package handlers

import (
	"fabric-rest-api-go/pkg/api"
	"github.com/labstack/echo/v4"
)

type ApiContext struct{
	echo.Context

	fsc *api.FabricSdkClient
}

func (c *ApiContext) Fsc() *api.FabricSdkClient {
	return c.fsc
}

func (c *ApiContext) SetFsc(fsc *api.FabricSdkClient) {
	c.fsc = fsc
}
