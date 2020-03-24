package main

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	modules := api.LoadModules()
	fmt.Println("== URLs Loaded == ")
	for _, mod := range modules {
		switch mod.Method {
		case "GET":
			fmt.Println("\tGET:\t" + mod.Path)
			e.GET(mod.Path, mod.Function)
			break
		case "POST":
			fmt.Println("\tPOST:\t" + mod.Path)
			e.POST(mod.Path, mod.Function)
			break
		}
	}

	e.Logger.Fatal(e.Start(":9999"))
}
