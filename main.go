package main

import (
	"golang-distributed-parallel-image-processing/api"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	modules := api.LoadModules()
	for _, mod := range modules {
		switch mod.Method {
		case "GET":
			e.GET(mod.Path, mod.Function)
			break
		case "POST":
			e.POST(mod.Path, mod.Function)
			break
		}
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
