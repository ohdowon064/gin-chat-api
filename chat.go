package main

import "github.com/google/uuid"

type ChatRoom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var ChatRooms []ChatRoom

func CreateChatRoom(name string) ChatRoom {
	id := uuid.New().String()
	chatRoom := ChatRoom{
		ID:   id,
		Name: name,
	}
	ChatRooms = append(ChatRooms, chatRoom)
	return chatRoom
}
