package ws

import (
	"bytes"
	"fmt"
	"goRunShengXiao/internal/service"
	"log"
	"net/http"
	"time" // 引入时间包

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client 代表一个连接到 WebSocket 的单个前端网页
type Client struct {
	Hub        *Hub
	Conn       *websocket.Conn
	Send       chan []byte // 专属于这个网页的发送通道
	CurrentCmd string      // 新增：记录当前客户端订阅的指令
}

// ReadPump：负责听取客户端的指令并切换“订阅状态”
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("用户断开: %v", err)
			break
		}
		cmd := string(bytes.TrimSpace(message))
		log.Printf("收到订阅请求: %s", cmd)

		// 核心逻辑：收到指令后不直接返回，而是修改订阅状态
		// 这样定时器就会根据新的状态开始推送数据
		c.CurrentCmd = cmd
	}
}

// SubscriptionPump：新增！专门负责每 3 秒检查一次订阅状态并推送
func (c *Client) SubscriptionPump() {
	ticker := time.NewTicker(3 * time.Second) // 创建 3 秒定时器
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			// 如果当前有订阅指令，就调用 service 生成数据并发送
			if c.CurrentCmd != "" {
				reply := service.ProcessCommand(c.CurrentCmd)
				// 发送到 Send 通道，由 WritePump 最终发出
				select {
				case c.Send <- []byte(reply):
				default:
					fmt.Printf("发送失败，用户可能已断开: %s\n", c.CurrentCmd) // 发送失败，用户可能已断开 后台直接断开 防止程序崩溃
					//log.Printf("发送失败，用户可能已断开: %s", c.CurrentCmd) // 发送失败，用户可能已断开 后台直接断开 防止程序崩溃
					return
				}
			}
		case <-c.Hub.Unregister:
			fmt.Printf("用户已下线，停止推送:%s\n", c.CurrentCmd) // 监听用户下线事件，停止推送
			//log.Printf("用户已下线，停止推送: %s", c.CurrentCmd)
			return
		}
	}
}

// WritePump 专门负责数据推送
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

// ServeWs 是 HTTP 升级为 WebSocket 的入口函数
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级 WS 失败:", err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client
	// 多个协程发送 一个接受 一个发送消息互不干扰
	go client.WritePump()
	go client.ReadPump()
	go client.SubscriptionPump() // 新增：专门处理订阅推送
}
