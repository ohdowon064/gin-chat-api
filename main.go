package main

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

var clients = make(map[string][]*net.Conn)   // 연결된 클라이언트들
var broadcast = make(map[string]chan []byte) // 모든 클라이언트에게 전송할 메시지 채널

func main() {
	// ServeMux 생성
	broadcast["1"] = make(chan []byte)
	broadcast["2"] = make(chan []byte)

	// gin 웹서버
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/rooms/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		c.HTML(http.StatusOK, "chat.html", gin.H{"RoomId": roomId})
	})

	r.POST("/rooms", func(c *gin.Context) {
		// 채팅방 생성
		name := c.PostForm("name")
		chatRoom := CreateChatRoom(name)
		broadcast[chatRoom.ID] = make(chan []byte)
		c.JSON(http.StatusOK, chatRoom)
	})

	r.GET("/rooms", func(c *gin.Context) {
		c.JSON(http.StatusOK, ChatRooms)
	})

	// websocket 핸들러 함수 설정
	r.GET("/ws", ReceiveFromClient)

	// 브로드캐스트 핸들러 함수 설정
	//go BroadcastToClient()

	r.Run(":8080")

}
