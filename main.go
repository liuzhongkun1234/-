package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	// 引入你刚才写的 ws 包
	"goRunShengXiao/internal/ws"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// 1. 启动 Hub（总控室）
	hub := ws.NewHub()
	go hub.Run() // 后台独立运行

	// 2. 注册路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	// 3. 启动服务器
	port := ":8082"
	fmt.Printf("已启动，监听端口 %s/ws\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
