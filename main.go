package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"net/http"
)

var clients = make(map[*net.Conn]bool) // 연결된 클라이언트들
var broadcast = make(chan []byte)      // 모든 클라이언트에게 전송할 메시지 채널

func main() {
	// gin 웹서버
	r := gin.Default()
	r.LoadHTMLGlob("templates/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// handshake error: bad "Upgrade" header 에러 때문에 추가
	// Upgrade 헤더설정

	// websocket 핸들러 함수 설정
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		// 연결된 클라이언트를 clients 맵에 추가
		clients[&conn] = true

		for {
			// 클라이언트로부터 메시지를 받고, broadcast 채널에 전송
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				fmt.Println(err)
				break
			}
			broadcast <- msg
			fmt.Printf("Received message %s (opcode: %d)\n", msg, op)
		}

		// 연결된 클라이언트를 clients 맵에서 제거
		delete(clients, &conn)
	})

	// 브로드캐스트 핸들러 함수 설정
	go func() {
		for {
			// broadcast 채널로부터 메시지를 받아서 모든 클라이언트에게 전송
			msg := <-broadcast
			for conn := range clients {
				err := wsutil.WriteServerMessage(*conn, ws.OpText, msg)
				if err != nil {
					fmt.Println(err)
					(*conn).Close()
					delete(clients, conn)
				}
			}
		}
	}()

	r.Run()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
