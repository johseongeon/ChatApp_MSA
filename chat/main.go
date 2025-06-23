package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/johseongeon/chat_package"
)

// upgrade to websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var RoomMgr = &chat_package.RoomManager{}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	var initMsg struct {
		Username string `json:"username"`
		ChatID   string `json:"chat_id"`
	}
	log.Println("WebSocket connection attempt")
	err = conn.ReadJSON(&initMsg)
	if err != nil {
		log.Println("Failed to read init message:", err)
		return
	}

	client := &chat_package.Client{
		Username: initMsg.Username,
		Conn:     conn,
		Rooms:    make(map[string]*chat_package.ChatRoom),
	}

	chatroom := RoomMgr.GetRoom(initMsg.ChatID)
	client.Rooms[initMsg.ChatID] = chatroom
	RoomMgr.ConnectToRoom(client, chatroom)
	log.Printf("User %s joined chat %s", client.Username, initMsg.ChatID)

	for {
		var msg struct {
			Message string `json:"message"`
			RoomID  string `json:"room_id"`
		}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		roomID := initMsg.ChatID
		if msg.RoomID != "" {
			roomID = msg.RoomID
		}

		chatMsg := chat_package.ChatMessage{
			Username:  client.Username,
			Message:   msg.Message,
			RoomID:    roomID,
			Timestamp: time.Now(),
		}
		if err := chat_package.MessageLog.LogMessage(chatMsg); err != nil {
			log.Printf("Failed to log message: %v", err)
		}

		client.BroadcastToRoom(roomID, map[string]string{
			"from":    client.Username,
			"message": msg.Message,
		})
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server started on :8080")
	client, err := chat_package.ConnectMongoDB()
	if err != nil {
		log.Fatal("Failed to connect MongoDB:", err)
	}
	chat_package.MessageLog.Client = client
	RoomMgr.Client = client
	chat_package.LoadRoomsFromDB(RoomMgr)

	// RoomManager 동기화
	go func() {
		for {
			chat_package.LoadWhileRunning(RoomMgr)
			time.Sleep(3 * time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
