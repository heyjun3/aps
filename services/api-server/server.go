package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"api-server/handler"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
	}))

	e.GET("/api/list", handler.GetFilenames)
	e.GET("/api/counts", handler.GetStatusCounts)
	e.GET("/api/chart_list/:filename", handler.GetCharts)
	e.DELETE("/api/deleteFile/:filename", handler.DeleteProducts)
	e.Logger.Fatal(e.Start(":5000"))
}
