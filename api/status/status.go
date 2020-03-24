package status

import (
	"net/http"

	"github.com/labstack/echo"
)

func StatusResponse(c echo.Context) error {
	return c.String(http.StatusOK, "Status: bruh")
}
