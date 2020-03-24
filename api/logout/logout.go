package logout

import (
	"fmt"

	"github.com/labstack/echo"
)

func LogoutResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/logout")
	return nil
}
