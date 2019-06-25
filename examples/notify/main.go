package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ganlvtech/go-kahla-api/cryptojs"
	"github.com/ganlvtech/go-kahla-api/kahla"
	"github.com/ganlvtech/go-kahla-api/pusher"
)

func pusherEventHandler(i interface{}) {
	// You should do anything lasting for a long time in another goroutine.
	switch v := i.(type) {
	case *pusher.Pusher_NewMessageEvent:
		content, err := cryptojs.AesDecrypt(v.Content, v.AesKey)
		if err != nil {
			log.Println(err)
		} else {
			title := v.Sender.NickName
			message := content
			log.Println(title, ":", message)
		}
	case *pusher.Pusher_NewFriendRequestEvent:
		title := "Friend request"
		message := "You have got a new friend request!"
		log.Println(title, ":", message, "nick name:", v.Requester.NickName, "id:", v.Requester.Id)
	case *pusher.Pusher_WereDeletedEvent:
		title := "Were deleted"
		message := "You were deleted by one of your friends from his friend list."
		log.Println(title, ":", message, "nick name:", v.Trigger.NickName, "id:", v.Trigger.Id)
	case *pusher.Pusher_FriendAcceptedEvent:
		title := "Friend request"
		message := "Your friend request was accepted!"
		log.Println(title, ":", message, "nick name:", v.Target.NickName, "id:", v.Target.Id)
	case *pusher.Pusher_TimerUpdatedEvent:
		title := "Self-destruct timer updated!"
		message := fmt.Sprintf("Your current message life time is: %d", v.NewTimer)
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	case *pusher.Pusher_NewMemberEvent:
		title := "New member"
		message := fmt.Sprintf("%s has join the group.", v.NewMember.NickName)
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	case *pusher.Pusher_SomeoneLeftEvent:
		title := "Someone left"
		message := fmt.Sprintf("%s has successfully left the group.", v.LeftUser.NickName)
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	case *pusher.Pusher_DissolveEvent:
		title := "Group Dissolved"
		message := "A group dissolved"
		log.Println(title, ":", message, "conversation id:", v.ConversationId)
	}
}

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("go run main.go user@example.com password")
		return
	}
	email := os.Args[1]
	password := os.Args[2]
	c := kahla.NewClient("https://staging.server.kahla.app", "https://oss.aiursoft.com")

	authByPasswordResponse, httpResponse, err := c.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	if httpResponse.StatusCode != http.StatusOK {
		log.Fatal("Response status code not 200")
	}
	if authByPasswordResponse.Code != 0 {
		log.Fatal(authByPasswordResponse.Code, ":", authByPasswordResponse.Message)
	}

	initPusherResponse, httpResponse, err := c.Auth.InitPusher()
	if err != nil {
		log.Fatal(err)
	}
	if httpResponse.StatusCode != http.StatusOK {
		log.Fatal("Response status code not 200")
	}
	if initPusherResponse.Code != 0 {
		log.Fatal(initPusherResponse.Code, ":", initPusherResponse.Message)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	interrupt2 := make(chan struct{})
	go func() {
		// cast channel type
		<-interrupt
		close(interrupt2)
	}()

	p := pusher.NewPusher(initPusherResponse.ServerPath, pusherEventHandler)
	err = p.Connect(interrupt2)
	if err != nil {
		log.Fatal(err, "Connect to pusher failed.")
	}
}
