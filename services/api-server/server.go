package main

import (
	"github.com/labstack/echo/v4"

	"api-server/handler"
)



func main() {
	e := echo.New()
	e.GET("/api/list", handler.GetFilenames)
	e.DELETE("/api/deletefile/:filename", handler.DeleteProducts)
	e.Logger.Fatal(e.Start(":1323"))
}
