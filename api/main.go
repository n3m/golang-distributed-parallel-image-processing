package api

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/login"
	"golang-distributed-parallel-image-processing/api/logout"
	"golang-distributed-parallel-image-processing/api/status"
	"golang-distributed-parallel-image-processing/api/upload"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

var IsLoggedIn = checkIfLoggedIn()

func checkIfLoggedIn() echo.MiddlewareFunc {
	data := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})
	return data
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
			Method:     "POST",
			Path:       "/logout",
			Function:   logout.LogoutResponse,
			Middleware: &IsLoggedIn,
		},
		&Module{
			Method:     "GET",
			Path:       "/status",
			Function:   status.StatusResponse,
			Middleware: &IsLoggedIn,
		},
		&Module{
			Method:     "POST",
			Path:       "/upload",
			Function:   upload.UploadResponse,
			Middleware: &IsLoggedIn,
		},
	}
}

func rootResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/")
	return c.JSON(http.StatusForbidden, &Message{Message: "You're not allowed to do this. [AM - Nothing here to see]"})
}
