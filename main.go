package main

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api"
	"golang-distributed-parallel-image-processing/api/helpers"
	"golang-distributed-parallel-image-processing/controller"
	"golang-distributed-parallel-image-processing/models"
	"golang-distributed-parallel-image-processing/scheduler"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sonyarouje/simdb/db"
)

func main() {
	/* Setup Folders*/
	_ = os.MkdirAll("public/download", 0755)
	_ = os.MkdirAll("public/results", 0755)

	/* Setup Variables */
	ControllerConnectionURL := "tcp://localhost:40899"
	APIPort := ":8080"
	DBName := "workers"
	currentWorkers := map[string]models.Worker{}
	jobs := make(chan scheduler.Job)
	workloadsFileNumbers := make(map[string]int64)

	/* DB Setup */
	db, err := db.New(DBName)
	if err != nil {
		panic(err)
	}

	/* Scheduler Setup */
	go scheduler.Start(jobs, currentWorkers)

	/* Controller Setup */
	go controller.Start(ControllerConnectionURL, currentWorkers, db)
	log.Printf("[SETUP] Controller Connection URL: %+v", ControllerConnectionURL)

	/* API EndPoint Setup */
	e := echo.New()

	/* Custom Context */
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &helpers.CustomContext{c, currentWorkers, workloadsFileNumbers, jobs}
			return next(cc)
		}
	})

	modules := api.LoadModules()
	fmt.Println("== URLs Loaded == ")
	/* Static Server Routes */

	fmt.Println("\tSTATIC:\t" + "/public")
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: true,
	}))
	e.Use(middleware.Static("public"))

	for _, mod := range modules {
		switch mod.Method {
		case "GET":
			fmt.Println("\tGET:\t" + mod.Path)
			if mod.Middleware != nil {
				e.GET(mod.Path, mod.Function, *mod.Middleware)
			} else {
				e.GET(mod.Path, mod.Function)
			}
			break
		case "POST":
			fmt.Println("\tPOST:\t" + mod.Path)
			if mod.Middleware != nil {
				e.POST(mod.Path, mod.Function, *mod.Middleware)
			} else {
				e.POST(mod.Path, mod.Function)
			}
			break
		}
	}

	go e.Logger.Fatal(e.Start(APIPort))

}
