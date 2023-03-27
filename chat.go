package main

import "github.com/google/uuid"

type ChatRoom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var ChatRooms []ChatRoom
var ChatRoomToUser map[ChatRoom]User
var UserToChatRoom map[User]ChatRoom

func CreateChatRoom(name string) ChatRoom {
	id := uuid.New().String()
	chatRoom := ChatRoom{
		ID:   id,
		Name: name,
	}
	ChatRooms = append(ChatRooms, chatRoom)
	return chatRoom
}

func JoinChatRoom(user User, chatRoom ChatRoom) {
	ChatRoomToUser[chatRoom] = user
	UserToChatRoom[user] = chatRoom
}

func LeaveChatRoom(user User) {
	chatRoom := UserToChatRoom[user]
	delete(ChatRoomToUser, chatRoom)
	delete(UserToChatRoom, user)
}
