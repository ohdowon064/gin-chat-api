package main

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"net/http"
)

func BroadcastToClient() {
	for {
		// broadcast 채널로부터 메시지를 받아서 모든 클라이언트에게 전송
		for roomId, client := range clients {
			msg := <-broadcast[roomId]
			fmt.Println("broadcast message:", msg)
			fmt.Println(client)

			for _, conn := range client {
				fmt.Println(conn)
				err := wsutil.WriteServerMessage(*conn, ws.OpText, msg)
				if err != nil {
					fmt.Println(err)
					(*conn).Close()
					for i, c := range client {
						if c == conn {
							clients[roomId] = append(clients[roomId][:i], clients[roomId][i+1:]...)
							break
						}
					}
				}
			}
		}
	}
}

func ReceiveFromClient(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query()["rooms"][0]

	conn, _, _, err := ws.UpgradeHTTP(r, w)
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
