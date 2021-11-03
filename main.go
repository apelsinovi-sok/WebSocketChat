package main

import (
	"fmt"
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

func chatStart(w http.ResponseWriter, r *http.Request) {
	upgrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		//log.Println(err)
		return
	}
	go s.newClient(ws)
}



func homePage(w http.ResponseWriter, r *http.Request)  {
	tmp, _ := template.ParseFiles("index.html")
	err := tmp.Execute(w,"")
	if err != nil {
		log.Println("Ошибка")
	}
}

var s = newServer()

func main() {
	http.HandleFunc("/", homePage)
	go s.run()
	fmt.Printf("HTTP Сервер запущен на порту %s\n", PORT)
	http.HandleFunc("/ws", chatStart)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Println(err)
		return
	}

}
