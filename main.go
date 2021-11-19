package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type consulta struct {
	peruana    bool
	embarazada bool
	hijos      bool
	trabaja    bool
	edad       bool
	casada     bool
	estudia    bool
	seguro     bool
	distrito   string
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Metodo GET para
	e.GET("/", func(c echo.Context) error {

		//FUNCION DONDE SE MANDA CONSULTA

		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
