package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
)

func broadcastToClient(roomId string, clientList []*net.Conn) {
	msg := <-broadcast[roomId]
	newStr := make([]byte, len(msg)) // 새로운 byte 배열을 만듦
	copy(newStr, msg)                // 기존 문자열의 값을 새로운 배열에 복사

	fmt.Println("broadcast message:", msg)

	for _, conn := range clientList {
		fmt.Println(conn)
		err := wsutil.WriteServerMessage(*conn, ws.OpText, newStr)
		if err != nil {
			fmt.Println(err)
			(*conn).Close()
			for i, c := range clientList {
				if c == conn {
					clients[roomId] = append(clients[roomId][:i], clients[roomId][i+1:]...)
					break
				}
			}
		}
	}
}

func BroadcastToClient() {
	for {
		// broadcast 채널로부터 메시지를 받아서 모든 클라이언트에게 전송
		for roomId, client := range clients {
			go broadcastToClient(roomId, client)
		}
	}
}

func ReceiveFromClient(c *gin.Context) {
	roomId := c.Query("roomId")

	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	fmt.Println("websocket connected:", conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 연결된 클라이언트를 clients 맵에 추가
	if clients[roomId] == nil {
		clients[roomId] = make([]*net.Conn, 0)
	}

	clients[roomId] = append(clients[roomId], &conn)

	for {
		// 클라이언트로부터 메시지를 받고, broadcast 채널에 전송
		msg, op, err := wsutil.ReadClientData(conn)
		fmt.Println("received message:", msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		broadcast[roomId] <- msg
		fmt.Printf("Received message %s (opcode: %d)\n", msg, op)
	}

	// 연결된 클라이언트를 clients 맵에서 제거
	for i, c := range clients[roomId] {
		if c == &conn {
			clients[roomId] = append(clients[roomId][:i], clients[roomId][i+1:]...)
			break
		}
	}
}
