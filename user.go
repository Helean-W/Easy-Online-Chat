package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	c      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		c:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

func (u *User) Online() {
	u.server.maplock.Lock()
	u.server.UserMap[u.Name] = u
	u.server.maplock.Unlock()
	message := "[" + u.Addr + "]" + u.Name + "已上线"
	u.server.Broadcast(message)
}

func (u *User) Offline() {
	u.server.maplock.Lock()
	delete(u.server.UserMap, u.Name)
	u.server.maplock.Unlock()
	message := "[" + u.Addr + "]" + u.Name + "已下线"
	u.server.Broadcast(message)
}

func (u *User) SendMessage(msg string) {
	u.conn.Write([]byte(msg))
}

func (u *User) DoMessage(msg string) {
	if msg == "who" {
		for _, user := range u.server.UserMap {
			msg := "[" + user.Name + "]" + "ONLINE" + "\n"
			u.SendMessage(msg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := u.server.UserMap[newName]; ok {
			u.SendMessage("name repeated!")
		} else {
			u.server.maplock.Lock()
			delete(u.server.UserMap, u.Name)
			u.server.UserMap[newName] = u
			u.server.maplock.Unlock()
			u.Name = newName
			u.SendMessage("更新成功")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMessage("格式为to|接收者|发送的信息。。。")
			return
		}
		remoteUser, ok := u.server.UserMap[remoteName]
		if !ok {
			u.SendMessage("用户名不存在")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMessage("无消息内容")
			return
		}
		remoteUser.SendMessage(u.Name + "对你说:" + content + "\n")
	} else {
		msg = "[" + u.Addr + "]" + u.Name + ":" + msg
		u.server.Broadcast(msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.c
		u.conn.Write([]byte(msg + "\n"))
	}
}
