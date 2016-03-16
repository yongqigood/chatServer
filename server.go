package main

import (
	"bytes"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
)
/***********holds data for a chat room.*************/

type Room struct {
	roomName    string
	msg string
	users   map[string] *User
}
/***********holds data for a user.*************/
type User struct {
	UserName    string
	RoomName    string
}

//Available chat rooms.
var chats map[string] *Room = make(map[string] *Room)
//connected users
var allUsers map[string] *User = make(map[string] *User)
//user chat channels
var userInfo map[string] chan string = make(map[string] chan string)

func PushHandler(w http.ResponseWriter, req *http.Request) {
	r := req.FormValue("room")
	room := chats[r]
	b, _ := json.Marshal(room.users)
	
	fmt.Println(string(b))
	for _, value := range room.users{
		ch := userInfo[value.UserName]
		if ch == nil {
		 	ch = make (chan string)
		 	userInfo[value.UserName] = ch  
		}
		ch <- string(b)
	}	
}

// poll information from server
func PollHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(400)
		return		
	}
	rcpt := req.FormValue("rcpt")
	ch := userInfo[rcpt]
	fmt.Println("good")
	if ch == nil {
		ch = make (chan string)
		userInfo[rcpt] = ch
	}
	io.WriteString(w, <-ch)
}
// user login handler
func LoginHandler(w http.ResponseWriter, req *http.Request) {
	
	rcpt := req.FormValue("rcpt")
	if _, ok := allUsers[rcpt]; ok {
    	fmt.Println(rcpt)
  		io.WriteString(w, "user reg")
  		return 
	}
	var buffer bytes.Buffer
	for _,rooms := range chats {
  		buffer.WriteString(rooms.roomName)
  		buffer.WriteString(",")
	}
	user := NewUser(rcpt, "not select rooms")
	allUsers[rcpt] = user
	fmt.Println(buffer.String())
	io.WriteString(w, buffer.String())
}
// poll information from server to user who login
func JoinRoomHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(400)
		return		
	}
	from := req.FormValue("from")
	ch := userInfo[from]
	if ch == nil {
		ch = make (chan string)
		userInfo[from] = ch
	}
	io.WriteString(w, <-ch)
	
}
// push notification that user leave this room
func LeaveRoomHandler(w http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	room := req.FormValue("room")
	users := chats[room].users
	delete(users, from)
	var buffer bytes.Buffer
 	if _, ok := chats[room]; ok {
		for _, user := range users{
			buffer.WriteString(user.UserName)
  			buffer.WriteString(",")
		}
		buffer.WriteString(from)
  		buffer.WriteString(",leave")
		for _, user := range users{
			rcpt := user.UserName
			fmt.Println(rcpt)
			ch := userInfo[rcpt]
			if ch == nil {
				ch = make (chan string)
		 		userInfo[rcpt] = ch  
			}
			ch <- buffer.String()
		}
	}else{
		fmt.Println(room)
  		io.WriteString(w, "No this room")
  		return 
	}
}

// user logout handler
func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	
	from := req.FormValue("from")
	if _, ok := allUsers[from]; ok{
		delete(allUsers, from)
		io.WriteString(w, "success")		
	}else{
  		io.WriteString(w, "No this user")
  		return 
	}
}
// push information to users who join this room
func NotifyHandler(w http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	room := req.FormValue("room")
	rcpt := req.FormValue("rcpt")
	user := allUsers[from]
	users := chats[room].users
	users[from] = user
	var buffer bytes.Buffer
	if _, ok := chats[room]; ok {
		for _, user := range users{
			buffer.WriteString(user.UserName)
  			buffer.WriteString(",")
		}
		buffer.WriteString(from)
		if req.Method == "GET" {
			buffer.WriteString(",come")		
		}else{
			body, _ := ioutil.ReadAll(req.Body)
			buffer.WriteString(",")
			buffer.WriteString(string(body))
		}
  		if(rcpt == "" || rcpt == "All"){
			for _, user := range users{
				rcpt1 := user.UserName
				fmt.Println(rcpt1)
				ch := userInfo[rcpt1]
				if ch == nil {
					ch = make (chan string)
		 			userInfo[rcpt1] = ch  
				}
				ch <- buffer.String()
			}
		}else{
			fmt.Println(rcpt)
			ch := userInfo[rcpt]
			if ch == nil {
				ch = make (chan string)
		 		userInfo[rcpt] = ch  
			}
			ch <- buffer.String()
		}
	}else{
		fmt.Println(room)
  		io.WriteString(w, "No this room")
  		return 
	}
}
// constructor of User
func NewUser(userName1 string, roomName1 string) *User {
	user := new(User)
	user.UserName = userName1
	user.RoomName = roomName1
	return user
}
//constructor of Room
func NewChatRoom(roomName string) *Room {
	chatRoom := &Room{
		users:  make(map[string] *User),
	}
	chatRoom.roomName = roomName
	return chatRoom
}
// main function
func main() {
	chatRoom1 := NewChatRoom("room1")
	chatRoom2 := NewChatRoom("room2")
	chats["room1"] = chatRoom1;
	chats["room2"] = chatRoom2;
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/poll", PollHandler)
	http.HandleFunc("/push", PushHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/joinroom", JoinRoomHandler)
	http.HandleFunc("/notify", NotifyHandler)
	http.HandleFunc("/leaveroom", LeaveRoomHandler)
	http.HandleFunc("/logout", LogoutHandler)
    err := http.ListenAndServe("localhost:8005", nil)
    if err != nil {
        //log.Fatal("ListenAndServe: ", err.String())
	}
}