package main

import (
	"algo/handles"
	"algo/middleware"
	"algo/models"
	"algo/web"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strings"
	"time"
)

type ConfigurationFile struct {
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	Upgrader struct {
		ReadBufferSize  int
		WriteBufferSize int
	}
	ClientConfiguration struct {
		WriteDeadlineSeconds int
		ReadDeadlineSeconds  int
		PingPeriodSeconds    int
	}
	IdGenerator struct {
		Charset string
		Length  int
	}
}

var (
	env *handles.Environment
)

func init() {
	var conf ConfigurationFile

	file, err := os.Open("configuration.json")
	if err != nil {
		log.Fatal(err)
	}
	var decoder = json.NewDecoder(file)
	err = decoder.Decode(&conf)

	log.Printf("Configuration: \n%+v", conf)

	if err != nil {
		log.Fatalf("Error decoding configuration file: %v", err)
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     strings.Join([]string{conf.Redis.Host, conf.Redis.Port}, ":"),
		Password: conf.Redis.Password,
		DB:       0,
	})

	env = &handles.Environment{
		Conf: handles.IdGeneratorConfiguration{
			Charset: conf.IdGenerator.Charset,
			Length:  conf.IdGenerator.Length,
		},
		Executioner: models.ExecutionerModel{},
		Editor: models.EditorModel{
			Redis: redisClient,
		},
		ClientHub: models.ClientModel{
			Conf: models.ClientConfiguration{
				WriteDeadlineSeconds: time.Duration(conf.ClientConfiguration.WriteDeadlineSeconds),
				ReadDeadlineSeconds:  time.Duration(conf.ClientConfiguration.ReadDeadlineSeconds),
				PingPeriodSeconds:    time.Duration(conf.ClientConfiguration.PingPeriodSeconds),
			},
			Clients: make(map[string]*models.Client),
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  conf.Upgrader.ReadBufferSize,
				WriteBufferSize: conf.Upgrader.WriteBufferSize,
			},
		},
	}
}

func main() {
	server := echo.New()

	middleware.InitMiddleware(server)
	server.Renderer = web.GetRenderer()

	server.File("/favicon.ico", "web/files/favicon.ico")

	server.GET("/", env.NewPlayground)
	server.GET("/*", env.Playground)
	server.GET("/ws/*", env.WebSocket)
	server.GET("/execute/*", env.Execute)

	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Failed to initiate server '%v'", err)
	}
}
