package api

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/login"
	"golang-distributed-parallel-image-processing/api/logout"
	"golang-distributed-parallel-image-processing/api/status"
	"golang-distributed-parallel-image-processing/api/upload"
	"golang-distributed-parallel-image-processing/api/workloads"
	"golang-distributed-parallel-image-processing/models"
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
		{
			Method:   "GET",
			Path:     "/",
			Function: rootResponse,
		},
		{
			Method:   "GET",
			Path:     "/login",
			Function: login.LoginResponse,
		},
		{
			Method:     "GET",
			Path:       "/logout",
			Function:   logout.LogoutResponse,
			Middleware: &IsLoggedIn,
		},
		{
			Method:     "GET",
			Path:       "/status",
			Function:   status.StatusResponse,
			Middleware: &IsLoggedIn,
		},
		{
			Method:     "GET",
			Path:       "/status/:worker",
			Function:   status.StatusWorkerResponse,
			Middleware: &IsLoggedIn,
		},
		{
			Method:     "POST",
			Path:       "/upload",
			Function:   upload.UploadResponse,
			Middleware: &IsLoggedIn,
		},
		{
			Method:     "GET",
			Path:       "/workloads/test",
			Function:   workloads.WorkloadsResponse,
			Middleware: &IsLoggedIn,
		},
		{
			Method:     "POST",
			Path:       "/workloads/filter",
			Function:   workloads.WorkloadsFilterResponse,
			Middleware: &IsLoggedIn,
		},
	}
}

func rootResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/")
	return c.JSON(http.StatusForbidden, &models.Message{Message: "You're not allowed to do this. [AM - Nothing here to see]"})
}
