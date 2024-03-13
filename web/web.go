package web

import (
	"algo/execution"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

var upgrader = websocket.Upgrader{}

type TemplateRenderer struct {
	templates *template.Template
}

type Code struct {
	Text string
}

var goFile = Code{}
var ws *websocket.Conn
var MsgType int

func WebSocketHandler(c echo.Context) error {
	ws, _ = upgrader.Upgrade(c.Response(), c.Request(), nil)

	defer ws.CloseHandler()

	for {
		msgType, msg, err := ws.ReadMessage()

		if err != nil {
			fmt.Println("Closing connection")
			break
		}

		MsgType = msgType
		goFile = Code{Text: string(msg)}
	}

	return nil
}

func ExecuteHandler(c echo.Context) error {
	executeStdout := execution.ExecuteGo(goFile.Text)

	if err := ws.WriteMessage(MsgType, []byte(executeStdout)); err != nil {
		fmt.Println("Connection closed")
	}

	return nil
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func GetRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob("web/*.html")),
	}
}
