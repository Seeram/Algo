package handles

import (
	"algo/models"
	"algo/utils"
	"algo/web"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

type IdGeneratorConfiguration struct {
	Charset string
	Length  int
}

type Environment struct {
	Conf   IdGeneratorConfiguration
	Editor interface {
		Get(string) (models.Editor, error)
		Create(string)
		All() []string
		Save(string, models.Editor)
	}
	Executioner interface {
		Execute(string, chan string)
	}
	ClientHub interface {
		CreateClient(string, echo.Context) *models.Client
		GetClientOutput(string) chan string
		Run(*models.Client, interface {
			Save(string, models.Editor)
		})
	}
}

func (e Environment) Playground(c echo.Context) error {
	// http://localhost:8080/slkdj
	//                         ^
	editId := strings.Split(c.Request().URL.Path, "/")[1]

	editor, err := e.Editor.Get(editId)

	if err != nil {
		log.Printf("Error getting page %v reason %v", editId, err)

		return c.String(http.StatusNotFound, "Page not found")
	}

	allPlaygrounds := e.Editor.All()

	// redirect after newEdit enforces new page
	c.Response().Header().Set("Cache-Control", "no-cache")

	return c.Render(http.StatusOK, "playground.html", web.Playground{
		Editor:         editor,
		AllPlaygrounds: allPlaygrounds,
	})
}

func (e Environment) WebSocket(c echo.Context) error {
	// http://localhost:8080/ws/slkdj
	//                            ^
	editId := strings.Split(c.Request().URL.Path, "/")[2]

	_, err := e.Editor.Get(editId)

	if err != nil {
		log.Printf("Error getting page %v reason %v", editId, err)

		return c.String(http.StatusNotFound, "Page not found")
	}

	client := e.ClientHub.CreateClient(editId, c)

	e.ClientHub.Run(client, e.Editor)

	return nil
}

func (e Environment) Execute(c echo.Context) error {
	// http://localhost:8080/slkdj
	//                         ^
	editId := strings.Split(c.Request().URL.Path, "/")[2]

	editor, err := e.Editor.Get(editId)

	if err != nil {
		log.Printf("Error getting page %v reason %v", editId, err)

		return c.String(http.StatusNotFound, "Page not found")
	}

	go e.Executioner.Execute(editor.Content, e.ClientHub.GetClientOutput(editId))

	return nil
}

func (e Environment) NewPlayground(c echo.Context) error {
	editorId := utils.StringWithCharset(e.Conf.Length, e.Conf.Charset)

	log.Printf("Created new playground: %v\n", editorId)

	e.Editor.Create(editorId)

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("http://localhost:8080/%v", editorId))
}
