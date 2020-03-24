package api

import "github.com/labstack/echo"

//Module ...
type Module struct {
	Method   string
	Path     string
	Function echo.HandlerFunc
}

// LoadModules ...
func LoadModules() []*Module {
	return []*Module{
		&Module{
			Method: "POST",
			Path:   "/login",
		},
		&Module{
			Method: "POST",
			Path:   "/logout",
		},
		&Module{
			Method: "GET",
			Path:   "/status",
		},
		&Module{
			Method: "POST",
			Path:   "/upload",
		},
	}
}
