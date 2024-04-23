package models

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

type ClientConfiguration struct {
	WriteDeadlineSeconds time.Duration
	ReadDeadlineSeconds  time.Duration
	PingPeriodSeconds    time.Duration
}

type Client struct {
	conf     ClientConfiguration
	id       string
	conn     *websocket.Conn
	sendChan chan string
}

type ClientModel struct {
	Conf     ClientConfiguration
	Upgrader *websocket.Upgrader
	Clients  map[string]*Client
}

func (cm ClientModel) CreateClient(editId string, c echo.Context) *Client {
	ws, err := cm.Upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		log.Fatalf("Error upgrading to websocket %v", err.Error())
	}

	return &Client{
		id:       editId,
		conn:     ws,
		sendChan: make(chan string),
		conf:     cm.Conf,
	}
}

func (cm ClientModel) Run(c *Client, editor interface{ Save(string, Editor) }) {
	cm.Clients[c.id] = c

	go c.reader(editor)
	go c.writer()
}

func (cm ClientModel) GetClientOutput(editId string) chan string {
	return cm.Clients[editId].sendChan
}

func (c Client) reader(editor interface{ Save(string, Editor) }) {
	c.conn.SetReadDeadline(time.Now().Add(c.conf.ReadDeadlineSeconds * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.conf.ReadDeadlineSeconds * time.Second))
		return nil
	})

	defer func() {
		log.Printf("Closing websocket connection %v", c.id)
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket unexpectly closed: %v %v", err, c.id)
			}
			break
		}

		editor.Save(c.id, Editor{Content: string(msg)})
	}
}

func (c Client) writer() {

	for {
		select {
		case msg := <-c.sendChan:
			c.conn.SetWriteDeadline(time.Now().Add(c.conf.WriteDeadlineSeconds * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Printf("Error sending message to ws: ", err.Error(), msg)
			}

		case <-time.After(c.conf.PingPeriodSeconds * time.Second):
			c.conn.SetWriteDeadline(time.Now().Add(c.conf.WriteDeadlineSeconds * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
