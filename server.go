package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn *websocket.Conn) {
	fmt.Printf("Подсоединился новый клиент: %s \n", conn.RemoteAddr().String())
	c := &client{
		nick:     "Аноним",
		conn:     conn,
		commands: s.commands,
	}
	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}
	c.nick = args[1]
	c.msg(fmt.Sprintf("Теперь твое имя: %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := strings.Join(args[1:], " ")
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	_, ok = r.members[c.conn.RemoteAddr()]
	if !ok {
		r.members[c.conn.RemoteAddr()] = c
		s.quitCurrentRoom(c)
		c.room = r
		c.room.broadCast(c, fmt.Sprintf("%s зашел в комнату", c.nick))
		c.msg(fmt.Sprintf("Вы зашли в комнату: %s", r.name))
	}

}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("Доступные комнаты: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("вы не находитесь в комнате"))
	} else {
		c.room.broadCast(c, c.nick+": "+strings.Join(args[1:], " "))
	}
}

func (s *server) quit(c *client) {
	fmt.Printf("client disconnected %s", c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)
	c.msg("Пока :(")
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadCast(c, fmt.Sprintf("%s вышел из комнаты: %s", c.nick, c.room.name))
	}
}
