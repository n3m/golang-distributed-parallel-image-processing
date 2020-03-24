package status

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func StatusResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/status")
	return c.String(http.StatusOK, "Status: bruh")
}
