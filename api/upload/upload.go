package upload

import (
	"fmt"

	"github.com/labstack/echo"
)

func UploadResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/upload")
	return nil
}
