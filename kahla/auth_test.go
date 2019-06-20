package kahla_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ganlvtech/go-kahla-api/kahla"
)

var testDataAuthBio string
var testDataAuthHeadImgKey uint32
var testDataAuthHideMyEmail bool
var testDataAuthNickName string
var testDataAuthThemeId uint32
var testDataAuthEnableEmailNotification bool

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

func TestAuthService_AuthByPassword(t *testing.T) {
	out, resp, err := c.Auth.AuthByPassword(&kahla.Auth_AuthByPasswordRequest{
		Email:    email,
		Password: password,
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
	assert.Equal(t, testDataAuthBio, out.Value.Bio)
	assert.Equal(t, testDataAuthHeadImgKey, out.Value.HeadImgFileKey)
	assert.Equal(t, testDataAuthHideMyEmail, !out.Value.MakeEmailPublic)
	assert.Equal(t, testDataAuthNickName, out.Value.NickName)
	assert.Equal(t, testDataAuthThemeId, out.Value.ThemeId)
	assert.Equal(t, testDataAuthEnableEmailNotification, out.Value.EnableEmailNotification)
}

func TestAuthService_ChangePassword(t *testing.T) {
	out, resp, err := c.Auth.ChangePassword(&kahla.Auth_ChangePasswordRequest{
		OldPassword:    password,
		NewPassword:    password,
		RepeatPassword: password,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully changed your password!", out.Message)
}

func TestAuthService_InitPusher(t *testing.T) {
	out, resp, err := c.Auth.InitPusher()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, KahlaResponseCodeOK, out.Code)
	assert.Equal(t, "Successfully get your channel.", out.Message)
	assert.NotEmpty(t, out.ChannelId)
	assert.Len(t, out.ConnectKey, 32)
	assert.Equal(t, out.ServerPath, fmt.Sprintf("wss://stargate.aiursoft.com/Listen/Channel?Id=%d&Key=%s", out.ChannelId, out.ConnectKey))
}

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

func TestAuthService_OAuth(t *testing.T) {
	resp, err := c.Auth.OAuth()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "api.aiursoft.com", resp.Request.URL.Host)
	assert.Equal(t, "/oauth/authorize", resp.Request.URL.Path)
}

// TODO RegisterKahla
// TODO AuthResult
// TODO SendEmail
