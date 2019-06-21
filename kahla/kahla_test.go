package kahla_test

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/malisit/kolpa"
	"github.com/stretchr/testify/assert"

	"github.com/ganlvtech/go-kahla-api/cryptojs"
	"github.com/ganlvtech/go-kahla-api/kahla"
)

const (
	KahlaResponseCodeOK int32 = 0
)

var testDataInitkahlaServer string
var testDataInitOssServer string
var testDataInitEmail string
var testDataInitPassword string
var c *kahla.Client
var k kolpa.Generator
var r *rand.Rand
var testDataAuthBio string
var testDataAuthHeadImgKey uint32
var testDataAuthHideMyEmail bool
var testDataAuthUserId string
var testDataAuthNickName string
var testDataAuthThemeId uint32
var testDataAuthEnableEmailNotification bool
var testDataConversationConversationId uint32
var testDataConversationAesKey string

func assertDateTime(t *testing.T, dateTime string) {
	assert.Regexp(t, `^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.\d+Z$`, dateTime)
}

func assertEmail(t *testing.T, email string) {
	assert.Regexp(t, `^([0-9a-zA-Z]([-\.\w]*[0-9a-zA-Z])*@([0-9a-zA-Z][-\w]*[0-9a-zA-Z]\.)+[a-zA-Z]{2,9})$`, email)
}

func assertAesKey(t *testing.T, aesKey string) {
	assert.Regexp(t, `^[0-9a-f]{32}$`, aesKey)
}

func assertUserId(t *testing.T, userId string) {
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, userId)
}

func assertDiscriminator(t *testing.T, discriminator string) {
	_, ok := kahla.Conversation_Discriminator_value[discriminator]
	assert.True(t, ok)
}

func TestMain(m *testing.M) {
	testDataInitkahlaServer = os.Getenv("KAHLA_SERVER")
	testDataInitOssServer = os.Getenv("OSS_SERVER")
	testDataInitEmail = os.Getenv("EMAIL")
	testDataInitPassword = os.Getenv("PASSWORD")
	c = kahla.NewClient(testDataInitkahlaServer, testDataInitOssServer)
	k = kolpa.C()
	// k.SetLanguage("zh_CN")
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	os.Exit(m.Run())
}

// TODO TestAuthService_RegisterKahla
// TODO TestAuthService_AuthResult
// TODO TestAuthService_SendEmail

func TestAuthService_Index(t *testing.T) {
	out, resp, err := c.Auth.Index()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Welcome to Aiursoft Kahla API! Running in Production mode.", out.Message)
}

func TestAuthService_Version(t *testing.T) {
	out, resp, err := c.Auth.Version()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get the latest version number for Kahla App and Kahla.CLI.", out.Message)
	assert.Contains(t, out.LatestVersion, ".")
	assert.Contains(t, out.LatestCLIVersion, ".")
	assert.Equal(t, "https://www.kahla.app", out.DownloadAddress)
}

// Must test OAuth before AuthByPassword.
func TestAuthService_OAuth(t *testing.T) {
	resp, err := c.Auth.OAuth()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "api.aiursoft.com", resp.Request.URL.Host)
	assert.Equal(t, "/oauth/authorize", resp.Request.URL.Path)
}

func TestAuthService_AuthByPassword(t *testing.T) {
	out, resp, err := c.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{
		Email:    testDataInitEmail,
		Password: testDataInitPassword,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Auth success.", out.Message)
}

func TestAuthService_SignInStatus(t *testing.T) {
	out, resp, err := c.Auth.SignInStatus()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get your signin status.", out.Message)
	assert.True(t, out.Value)
}

func TestAuthService_UpdateInfo(t *testing.T) {
	testDataAuthBio = k.LoremSentence()
	if len(testDataAuthBio) > 80 {
		testDataAuthBio = testDataAuthBio[:80]
	}
	testDataAuthHeadImgKey = uint32(r.Intn(10000))
	testDataAuthHideMyEmail = r.Intn(2) > 0
	testDataAuthNickName = k.FirstName()
	out, resp, err := c.Auth.UpdateInfo(&kahla.Auth_UpdateInfoRequest{
		Bio:         &testDataAuthBio,
		HeadImgKey:  &testDataAuthHeadImgKey,
		HideMyEmail: &testDataAuthHideMyEmail,
		NickName:    &testDataAuthNickName,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully set your personal info.", out.Message)
}

func TestAuthService_UpdateClientSetting_ThemeId(t *testing.T) {
	testDataAuthThemeId = uint32(r.Intn(6))
	out, resp, err := c.Auth.UpdateClientSetting(&kahla.Auth_UpdateClientSettingRequest{
		ThemeId: &testDataAuthThemeId,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully update your client setting.", out.Message)
}

func TestAuthService_UpdateClientSetting_EnableEmailNotification(t *testing.T) {
	testDataAuthEnableEmailNotification = r.Intn(2) > 0
	out, resp, err := c.Auth.UpdateClientSetting(&kahla.Auth_UpdateClientSettingRequest{
		EnableEmailNotification: &testDataAuthEnableEmailNotification,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully update your client setting.", out.Message)
}

func TestAuthService_Me(t *testing.T) {
	out, resp, err := c.Auth.Me()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get your information.", out.Message)
	assert.NotNil(t, out.Value)
	assertDateTime(t, out.Value.AccountCreateTime)
	assert.Equal(t, testDataAuthBio, out.Value.Bio)
	assert.Equal(t, testDataInitEmail, out.Value.Email)
	assert.Equal(t, testDataAuthHeadImgKey, out.Value.HeadImgFileKey)
	assertUserId(t, out.Value.Id)
	testDataAuthUserId = out.Value.Id
	assert.Equal(t, testDataAuthHideMyEmail, !out.Value.MakeEmailPublic)
	assert.Equal(t, testDataAuthNickName, out.Value.NickName)
	assert.Equal(t, testDataAuthThemeId, out.Value.ThemeId)
	assert.Equal(t, testDataAuthEnableEmailNotification, out.Value.EnableEmailNotification)
}

func TestAuthService_ChangePassword(t *testing.T) {
	out, resp, err := c.Auth.ChangePassword(&kahla.Auth_ChangePasswordRequest{
		OldPassword: testDataInitPassword,
		// TODO change password to a new one
		NewPassword:    testDataInitPassword,
		RepeatPassword: testDataInitPassword,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully changed your password!", out.Message)
}

func TestAuthService_SendEmail(t *testing.T) {
	out, resp, err := c.Auth.SendEmail(&kahla.Auth_SendEmailRequest{
		Email: testDataInitEmail,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	switch out.Code {
	case KahlaResponseCodeOK:
		assert.Equal(t, "Successfully sent the validation email.", out.Message)
	case -6:
		assert.Equal(t, fmt.Sprintf("The email :%s was already validated!", testDataInitEmail), out.Message)
	default:
		assert.Fail(t, "out.Code not equals 0 or -6")
	}
}

func TestAuthService_InitPusher(t *testing.T) {
	out, resp, err := c.Auth.InitPusher()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get your channel.", out.Message)
	assert.NotEmpty(t, out.ChannelId)
	assert.Len(t, out.ConnectKey, 32)
	assert.Equal(t, out.ServerPath, fmt.Sprintf("wss://stargate.aiursoft.com/Listen/Channel?id=%d&key=%s", out.ChannelId, out.ConnectKey))
}

func TestConversationService_All(t *testing.T) {
	out, resp, err := c.Conversation.All(&kahla.Conversation_AllRequest{
		Take: 15,
		Skip: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get all your friends.", out.Message)
	if len(out.Items) <= 0 {
		fmt.Println("Response item list empty.")
	}
	for _, v := range out.Items {
		assert.NotEmpty(t, v.DisplayName)
		assert.NotEmpty(t, v.DisplayImageKey)
		assert.NotEmpty(t, v.LatestMessage)
		assertDateTime(t, v.LatestMessageTime)
		assert.NotEmpty(t, v.ConversationId)
		assertDiscriminator(t, v.Discriminator)
		switch v.Discriminator {
		case kahla.Conversation_Discriminator_name[kahla.Conversation_Discriminator_PrivateConversation]:
			assertUserId(t, v.UserId)
		case kahla.Conversation_Discriminator_name[kahla.Conversation_Discriminator_GroupConversation]:
			assert.Empty(t, v.UserId)
		}
		assertAesKey(t, v.AesKey)
		content, err := cryptojs.AesDecrypt(v.LatestMessage, v.AesKey)
		assert.NotEmpty(t, content)
		assert.NoError(t, err)
		if testDataConversationConversationId == 0 {
			testDataConversationConversationId = v.ConversationId
			testDataConversationAesKey = v.AesKey
		}
	}
}

func TestConversationService_GetMessage(t *testing.T) {
	if testDataConversationConversationId == 0 {
		fmt.Println("No conversation. Skip.")
		return
	}
	out, resp, err := c.Conversation.GetMessage(&kahla.Conversation_GetMessageRequest{
		Id:       testDataConversationConversationId,
		Take:     15,
		SkipTill: -1,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get all your messages.", out.Message)
	if len(out.Items) <= 0 {
		fmt.Println("Response item list empty.")
	}
	for _, v := range out.Items {
		assert.NotEmpty(t, v.Id)
		assert.Equal(t, testDataConversationConversationId, v.ConversationId)
		for _, v1 := range v.Ats {
			assertUserId(t, v1.TargetUserId)
		}
		assertUserId(t, v.SenderId)
		assert.NotEmpty(t, v.Sender)
		assert.Equal(t, v.SenderId, v.Sender.Id)
		if v.Sender.MakeEmailPublic {
			assertEmail(t, v.Sender.Email)
		} else {
			assert.Empty(t, v.Sender.Email)
		}
		assert.NotEmpty(t, v.Sender.NickName)
		assert.NotEmpty(t, v.Sender.HeadImgFileKey)
		assertDateTime(t, v.Sender.AccountCreateTime)
		assertDateTime(t, v.SendTime)
		assert.NotEmpty(t, v.Content)
		assert.NotEmpty(t, v.ConversationId)
		s, err := cryptojs.AesDecrypt(v.Content, testDataConversationAesKey)
		assert.NotEmpty(t, s)
		assert.NoError(t, err)
	}
}

func TestConversationService_SendMessage(t *testing.T) {
	if testDataConversationConversationId == 0 {
		fmt.Println("No conversation. Skip.")
		return
	}
	content, err := cryptojs.AesEncrypt(k.LoremSentence(), testDataConversationAesKey)
	assert.NoError(t, err)
	out, resp, err := c.Conversation.SendMessage(&kahla.Conversation_SendMessageRequest{
		Id:      testDataConversationConversationId,
		Content: content,
		At:      []string{},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Your message has been sent.", out.Message)
}

func TestConversationService_ConversationDetail(t *testing.T) {
	if testDataConversationConversationId == 0 {
		fmt.Println("No conversation. Skip.")
		return
	}
	out, resp, err := c.Conversation.ConversationDetail(&kahla.Conversation_ConversationDetailRequest{
		Id: testDataConversationConversationId,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get target conversation.", out.Message)
	assert.NotEmpty(t, out.Value.Id)
	assertDiscriminator(t, out.Value.Discriminator)
	switch out.Value.Discriminator {
	case kahla.Conversation_Discriminator_name[kahla.Conversation_Discriminator_PrivateConversation]:
		assertUserId(t, out.Value.RequesterId)
		assert.NotEmpty(t, out.Value.RequestUser)
		u1 := out.Value.RequestUser
		assert.Equal(t, out.Value.RequesterId, u1.Id)
		if u1.MakeEmailPublic {
			assertEmail(t, u1.Email)
		} else {
			assert.Empty(t, u1.Email)
		}
		assert.NotEmpty(t, u1.NickName)
		assert.NotEmpty(t, u1.HeadImgFileKey)
		assertDateTime(t, u1.AccountCreateTime)
		assertUserId(t, out.Value.TargetId)
		assert.NotEmpty(t, out.Value.TargetUser)
		u2 := out.Value.TargetUser
		assert.Equal(t, out.Value.TargetId, u2.Id)
		if u2.MakeEmailPublic {
			assertEmail(t, u2.Email)
		} else {
			assert.Empty(t, u2.Email)
		}
		assert.NotEmpty(t, u2.NickName)
		assert.NotEmpty(t, u2.HeadImgFileKey)
		assertDateTime(t, u2.AccountCreateTime)
		assertUserId(t, out.Value.AnotherUserId)
		switch out.Value.AnotherUserId {
		case out.Value.RequesterId:
			assert.Equal(t, testDataAuthUserId, out.Value.RequesterId)
		case out.Value.TargetId:
			assert.Equal(t, testDataAuthUserId, out.Value.TargetId)
		default:
			assert.Fail(t, "AnotherUserId not equals	RequesterId or TargetId.")
		}
	case kahla.Conversation_Discriminator_name[kahla.Conversation_Discriminator_GroupConversation]:
		assert.NotEmpty(t, out.Value.Users)
		assert.NotEmpty(t, out.Value.GroupImageKey)
		assert.NotEmpty(t, out.Value.GroupName)
		assertUserId(t, out.Value.OwnerId)
	}
	assertAesKey(t, out.Value.AesKey)
	assertDateTime(t, out.Value.ConversationCreateTime)
	assert.NotEmpty(t, out.Value.DisplayName)
	assert.NotEmpty(t, out.Value.DisplayImage)
}

func TestConversationService_UpdateMessageLifeTime(t *testing.T) {
	if testDataConversationConversationId == 0 {
		fmt.Println("No conversation. Skip.")
		return
	}
	out, resp, err := c.Conversation.UpdateMessageLifeTime(&kahla.Conversation_UpdateMessageLifeTimeRequest{
		Id:          testDataConversationConversationId,
		NewLifeTime: math.MaxInt32,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	switch out.Code {
	case KahlaResponseCodeOK:
		assert.Contains(t, out.Message, "Successfully updated your life time. Your current message life time is:")
	case -8:
		assert.Equal(t, "You are not the owner of that group.", out.Message)
	default:
		assert.Fail(t, "out.Code not equals 0 or -8")
	}
}

// Must log off at the end
func TestAuthService_LogOff(t *testing.T) {
	// TODO use a real DeviceId instead of -1
	out, resp, err := c.Auth.LogOff(&kahla.Auth_LogOffRequest{
		DeviceId: -1,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// TODO -3 means "Successfully logged you off, but we did not find device with id: -1". Try to remove it and assert response code only 0.
	assert.True(t, out.Code == KahlaResponseCodeOK || out.Code == -3)
	assert.Contains(t, out.Message, "Successfully logged you off")
}
