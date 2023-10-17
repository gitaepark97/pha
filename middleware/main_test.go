package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/util"
)

var testConfig = util.Config{
	JWTSecret:           util.CreateRandomString(32),
	AccessTokenDuration: time.Minute,
}

type Server struct {
	router *gin.Engine
}

func newServer() Server {

	return Server{
		router: gin.Default(),
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
