package main

import (
	"ballot/controller"
	"ballot/data"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	e := echo.New()
	e.Group("/api").POST("/reverse", controller.MessageReverseHandler)
	data.LoadData()
	log.Fatal(e.Start(":2201"))
}
