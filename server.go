package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP      string
	Port    int
	UserMap map[string]*User
	maplock sync.Mutex
	msg     chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:      ip,
		Port:    port,
		UserMap: make(map[string]*User),
		msg:     make(chan string),
	}
	return server
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.msg
		s.maplock.Lock()
		for _, user := range s.UserMap {
			user.c <- msg
		}
		s.maplock.Unlock()
	}
}
func (s *Server) Broadcast(msg string) {
	s.msg <- msg
}
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("处理连接")

	user := NewUser(conn, s)

	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("读取错误")
				return
			}

			umsg := string(buf[:n-1])

			user.DoMessage(umsg)

			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:
		case <-time.After(time.Second * 200):
			user.SendMessage("你被T了")
			close(user.c)
			conn.Close()
			return
		}
	}

}
func (s *Server) Start() {
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("Listener错误")
		return
	}
	defer Listener.Close()

	go s.ListenMessage()

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("listener accept错误")
			continue
		}
		go s.Handler(conn)
	}
}
