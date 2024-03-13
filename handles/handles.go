package handles

import (
	"algo/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func Edit(c echo.Context) error {
	return c.Render(http.StatusOK, "edit.html", nil)
}

func WebSocket(c echo.Context) error {
	return web.WebSocketHandler(c)
}

func Execute(c echo.Context) error {
	return web.ExecuteHandler(c)
}
