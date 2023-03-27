package main

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net/http"
)

func BroadcastToClient() {
	for {
		// broadcast 채널로부터 메시지를 받아서 모든 클라이언트에게 전송
		for _, chatRoom := range ChatRooms {
			msg := <-broadcast[chatRoom.ID]
			for conn := range clients {
				if clients[conn] != chatRoom.ID {
					continue
				}
				err := wsutil.WriteServerMessage(*conn, ws.OpText, msg)
				if err != nil {
					fmt.Println(err)
					(*conn).Close()
					delete(clients, conn)
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
	clients[&conn] = roomId

	for {
		// 클라이언트로부터 메시지를 받고, broadcast 채널에 전송
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			fmt.Println(err)
			break
		}
		broadcast[roomId] <- msg
		fmt.Printf("Received message %s (opcode: %d)\n", msg, op)
	}

	// 연결된 클라이언트를 clients 맵에서 제거
	delete(clients, &conn)
}
