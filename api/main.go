package api

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/login"
	"golang-distributed-parallel-image-processing/api/logout"
	"golang-distributed-parallel-image-processing/api/status"
	"golang-distributed-parallel-image-processing/api/upload"
	"net/http"

	"github.com/labstack/echo"
)

//Module ...
type Module struct {
	Method     string
	Path       string
	Function   echo.HandlerFunc
	Middleware *echo.MiddlewareFunc
}

type Message struct {
	Message string `json:"message"`
}

// LoadModules ...
func LoadModules() []*Module {
	return []*Module{
		&Module{
			Method:   "GET",
			Path:     "/",
			Function: rootResponse,
		},
		&Module{
			Method:   "POST",
			Path:     "/login",
			Function: login.LoginResponse,
		},
		&Module{
			Method:   "POST",
			Path:     "/logout",
			Function: logout.LogoutResponse, //TODO Add a function response for logout
		},
		&Module{
			Method:   "GET",
			Path:     "/status",
			Function: status.StatusResponse, //TODO Add a function response for status
		},
		&Module{
			Method:   "POST",
			Path:     "/upload",
			Function: upload.UploadResponse, //TODO Add a function response for upload
		},
	}
}

func rootResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/")
	return c.JSON(http.StatusForbidden, &Message{Message: "You're not allowed to do this. [AM - Nothing here to see]"})
}
