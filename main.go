package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

const PORT = ":7081"

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func chatStart(s *server, w http.ResponseWriter, r *http.Request) {
	upgrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go s.newClient(ws)
}

func homePage(c *gin.Context) {
	tmp, _ := template.ParseFiles("index.html")
	err := tmp.Execute(c.Writer, "")
	if err != nil {
		log.Println("Ошибка")
	}
}

func main() {
	var s = newServer()
	go s.run()
	fmt.Printf("HTTP Сервер запущен на порту %s\n", PORT)
	g := gin.Default()
	g.GET("/", homePage)
	g.GET("/ws", func(c *gin.Context) {
		chatStart(s, c.Writer, c.Request)
	})
	err := g.Run(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}
