package ws

// Hub 维护所有的活动客户端，并向客户端广播消息。
type Hub struct {
	// 注册了的客户端名单 (类似于一个 Set)
	Clients map[*Client]bool

	// 广播通道：将来你的 MQTT 收到数据，只要丢进这里，所有人都能收到！
	Broadcast chan []byte

	// 客户端上线通道
	Register chan *Client

	// 客户端下线通道
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run 启动总控室，死循环监听各种事件
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true // 用户上线，加入名单
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client) // 用户下线，踢出名单
				close(client.Send)
			}
		case message := <-h.Broadcast:
			// 核心机制：收到广播指令，循环发给名单里的所有人
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}