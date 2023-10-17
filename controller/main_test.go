package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/pha/service"
	"github.com/gitaepark/pha/util"
	"github.com/stretchr/testify/require"
)

var testConfig = util.Config{
	JWTSecret:            util.CreateRandomString(32),
	AccessTokenDuration:  time.Minute,
	RefreshTokenDuration: time.Minute,
}

func newTestController(t *testing.T, service service.Service) *Controller {
	return NewController(testConfig, service)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

type responseBody struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	Data interface{}
}

func getResponseBody(t *testing.T, body *bytes.Buffer) responseBody {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponseBody responseBody
	err = json.Unmarshal(data, &gotResponseBody)
	require.NoError(t, err)

	return gotResponseBody
}
