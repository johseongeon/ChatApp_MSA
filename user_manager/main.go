package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/johseongeon/chat_package"
)

var UserManagerInstance = &chat_package.UserManager{}

var Collection = &chat_package.MessageCollection{}

func main() {

	// connect to MongoDB
	client, err := chat_package.ConnectMongoDB()
	if err != nil {
		log.Fatal("MongoDB 연결 실패:", err)
	}

	// Initialize UserManager and RoomManager
	userManager := &chat_package.UserManager{Client: client}
	RoomMgr := &chat_package.RoomManager{Client: client}
	// Load users from DB
	chat_package.LoadRoomsFromDB(RoomMgr)

	// RoomManager 동기화
	go func() {
		for {
			chat_package.LoadWhileRunning(RoomMgr)
			time.Sleep(3 * time.Second)
		}
	}()

	Collection.MessageCol = client.Database("ChatDB").Collection("users")

	//register
	http.HandleFunc("/register", chat_package.RegisterServer(client))

	//addFriend
	http.HandleFunc("/addFriend", chat_package.Add_friend(client, userManager))

	//getFriends
	http.HandleFunc("/getFriends", chat_package.GetFriends(client, userManager))

	//getRooms
	http.HandleFunc("/getRooms", chat_package.GetRooms(client, userManager))

	//createRoom
	http.HandleFunc("/createRoom", chat_package.CreateRoom(client, RoomMgr))

	//joinUser
	http.HandleFunc("/joinUser", chat_package.JoinUser(client, RoomMgr))

	fmt.Println("Server started on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
