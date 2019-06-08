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

	var conversationId uint32
	for _, v := range response3.Items {
		if v.Discriminator == kahla.Conversation_Discriminator_name[int32(kahla.Conversation_Discriminator_GroupConversation)] {
			conversationId = v.ConversationId
		}
	}

	response4, _, _ := c.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: conversationId,
	})
	fmt.Println(response4)

	file, _ := os.Open("main.go")
	defer file.Close()
	file2 := kahla.RequestFile(file)
	response9, _, _ := c.Files.UploadFile(&kahla.Files_UploadFileRequest{
		File:           &file2,
		ConversationId: conversationId,
	})
	fmt.Println(response9)

	// content := "Send By protobuf rpc compiled go-kahla-api client"
	content := fmt.Sprintf("[file]%d-%s-%d B", response9.FileKey, response9.SavedFileName, response9.FileSize)
	content, _ = cryptojs.AesEncrypt(content, response4.Value.AesKey)
	response5, _, _ := c.Conversation.SendMessage(&kahla.Conversation_SendMessageRequest{
		At:      []string{response4.Value.Users[0].User.Id},
		Content: content,
		Id:      conversationId,
	})
	fmt.Println(response5)

	// response6, _, _ := c.Friendship.Mine()
	// fmt.Println(response6)
	//
	// response7, _, _ := c.Auth.UpdateInfo(&kahla.Auth_UpdateInfoRequest{
	// 	NickName: "Ganlv",
	// 	Bio:      "Hello everyone. I'm a ... i don't know.",
	// })
	// fmt.Println(response7)
	//
	// response8, _, _ := c.Devices.MyDevices()
	// fmt.Println(response8)
}
