package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIP:   ip,
		ServerPort: port,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	return client
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1-公聊模式")
	fmt.Println("2-私聊模式")
	fmt.Println("3-更新用户名")
	fmt.Println("0-退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的范围")
		return false
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println("请输入新的名字：")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.write error:", err)
		return false
	}
	return true
}

func (c *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			msg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("conn.write error:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}

}

func (c *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.write error:", err)
		return
	}

}

func (c *Client) PrivateChat() {
	var remoteUser string
	var chatContent string
	c.SelectUser()
	fmt.Println("请输入用户名，exit退出")
	fmt.Scanln(&remoteUser)
	for remoteUser != "exit" {

		fmt.Println("请输入聊天内容, exit退出")
		fmt.Scanln(&chatContent)

		for chatContent != "exit" {
			if len(chatContent) != 0 {
				sendMsg := "to|" + remoteUser + "|" + chatContent + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.write error:", err)
					break
				}
			}
			chatContent = ""
			fmt.Println("请输入聊天内容，exit退出")
			fmt.Scanln(&chatContent)
		}

		remoteUser = ""
		c.SelectUser()
		fmt.Println("请输入用户名，exit退出")
		fmt.Scanln(&remoteUser)
	}

}

func (c *Client) DealResponse() {
	//一旦c.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, c.conn)

	// for {
	// 	buf := make([]byte, 4096)
	// 	c.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}
		//不需要break，不会向下执行
		switch c.flag {
		case 1:
			fmt.Println("选择公聊模式")
			c.PublicChat()
		case 2:
			fmt.Println("选择私聊模式")
			c.PrivateChat()
		case 3:
			fmt.Println("选择更改用户名")
			c.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器IP地址")
	flag.IntVar(&serverPort, "port", 8888, "服务器端口号")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	go client.DealResponse()
	if client == nil {
		fmt.Println("连接失败")
		return
	}
	fmt.Println("连接成功")
	client.Run()
}
