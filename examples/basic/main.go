package main

import (
	"fmt"
	"os"

	"github.com/ganlvtech/go-kahla-api/cryptojs"
	"github.com/ganlvtech/go-kahla-api/kahla"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("go run main.go user@example.com password")
		return
	}
	email := os.Args[1]
	password := os.Args[2]

	c := kahla.NewClient("https://staging.server.kahla.app", "https://oss.aiursoft.com")

	response1, _, _ := c.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{
		Email:    email,
		Password: password,
	})
	fmt.Println(response1)

	response2, _, _ := c.Auth.Me()
	fmt.Println(response2)

	response3, _, _ := c.Conversation.All(&kahla.Conversation_AllRequest{
		Skip: 0,
		Take: 15,
	})
	fmt.Println(response3)

	// find first group
	var conversationId uint32
	for _, v := range response3.Items {
		if v.Discriminator == kahla.Conversation_Discriminator_name[kahla.Conversation_Discriminator_GroupConversation] {
			conversationId = v.ConversationId
		}
	}

	response4, _, _ := c.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})
	fmt.Println(response4)

	file, _ := os.Open("main.go")
	defer file.Close()
	response5, _, _ := c.Files.UploadFile(&kahla.Files_UploadFileRequest{
		File:           file,
		ConversationId: conversationId,
	})
	fmt.Println(response5)

	// content := "Send By protobuf rpc compiled go-kahla-api client"
	content := fmt.Sprintf("[file]%d-%s-%d B", response5.FileKey, response5.SavedFileName, response5.FileSize)
	content, _ = cryptojs.AesEncrypt(content, response4.Value.AesKey)
	response6, _, _ := c.Conversation.SendMessage(&kahla.Conversation_SendMessageRequest{
		At:      []string{response4.Value.Users[0].User.Id},
		Content: content,
		Id:      conversationId,
	})
	fmt.Println(response6)

	response7, _, _ := c.Friendship.Mine()
	fmt.Println(response7)

	nickName := "Ganlv"
	bio := "Hello everyone. I'm a ... i don't know."
	response8, _, _ := c.Auth.UpdateInfo(&kahla.Auth_UpdateInfoRequest{
		NickName: &nickName,
		Bio:      &bio,
	})
	fmt.Println(response8)

	response9, _, _ := c.Devices.MyDevices()
	fmt.Println(response9)
}
