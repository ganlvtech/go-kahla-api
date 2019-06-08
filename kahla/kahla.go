package kahla

import (
	"net/http"

	cookiejar "github.com/ganlvtech/go-exportable-cookiejar"
)

type Client struct {
	client       *http.Client
	baseUrl      string
	Auth         *AuthService
	Conversation *ConversationService
	Devices      *DevicesService
	Files        *FilesService
	Friendship   *FriendshipService
	Groups       *GroupsService
	Oss          *OssService
}

func NewClient(baseUrl string, ossUrl string) *Client {
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)
	return &Client{
		client:client,
		baseUrl: baseUrl,
		Auth: &AuthService{client: client, baseUrl: baseUrl},
		Conversation: &ConversationService{client: client, baseUrl: baseUrl},
		Devices: &DevicesService{client: client, baseUrl: baseUrl},
		Files: &FilesService{client: client, baseUrl: baseUrl},
		Friendship: &FriendshipService{client: client, baseUrl: baseUrl},
		Groups: &GroupsService{client: client, baseUrl: baseUrl},
		Oss: &OssService{client: client, baseUrl: ossUrl},
	}
}
