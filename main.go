package main

import (
	"algo/handles"
	"algo/middleware"
	"algo/web"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	middleware.InitMiddleware(e)
	e.Renderer = web.GetRenderer()

	e.GET("/", handles.Index)
	e.GET("/edit", handles.Edit)
	e.GET("/wsCodeEditor", handles.WebSocket)
	e.GET("/execute", handles.Execute)

	err := e.Start(":8080")

	if err != nil {
		fmt.Errorf("failed to initiate server '%v'", err)
	}
}
