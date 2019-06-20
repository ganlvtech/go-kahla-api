package kahla_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/malisit/kolpa"

	"github.com/ganlvtech/go-kahla-api/kahla"
)

const (
	KahlaResponseCodeOK int32 = 0
)

var kahlaServer string
var ossServer string
var email string
var password string
var c *kahla.Client
var k kolpa.Generator
var r *rand.Rand

func TestMain(m *testing.M) {
	kahlaServer = os.Getenv("KAHLA_SERVER")
	ossServer = os.Getenv("OSS_SERVER")
	email = os.Getenv("EMAIL")
	password = os.Getenv("PASSWORD")
	c = kahla.NewClient(kahlaServer, ossServer)
	k = kolpa.C()
	// k.SetLanguage("zh_CN")
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	os.Exit(m.Run())
}
