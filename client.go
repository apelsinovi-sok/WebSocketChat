package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
)

type client struct {
	nick     string
	conn     *websocket.Conn
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		_, msgByte, err := c.conn.ReadMessage()
		msg := string(msgByte)
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
				args:   args,
			}
		default:
			c.err(errors.New(fmt.Sprintf("Неверная команда: %s", cmd)))
		}
	}
}

func (c *client) err(err error) {
	err = c.conn.WriteMessage(1, []byte("Ошибка: " + err.Error()))
	if err != nil {
		return
	}
}

func (c *client) msg(msg string) {
	err := c.conn.WriteMessage(1, []byte(msg))
	if err != nil {
		return
	}
}
