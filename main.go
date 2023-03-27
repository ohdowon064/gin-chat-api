package main

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

var clients = make(map[*net.Conn]bool) // 연결된 클라이언트들
var broadcast = make(chan []byte)      // 모든 클라이언트에게 전송할 메시지 채널

func main() {
	// ServeMux 생성
	mux := http.NewServeMux()

	// gin 웹서버
	r := gin.Default()
	r.LoadHTMLGlob("templates/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// websocket 핸들러 함수 설정
	mux.HandleFunc("/ws", ReceiveFromClient)

	// 브로드캐스트 핸들러 함수 설정
	go BroadcastToClient()

	mux.Handle("/", r)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
