package service

import (
	"testing"
	"time"

	"github.com/gitaepark/pha/repository"
	"github.com/gitaepark/pha/util"
)

var testConfig = util.Config{
	JWTSecret:            util.CreateRandomString(32),
	AccessTokenDuration:  time.Minute,
	RefreshTokenDuration: time.Minute,
}

func newTestService(t *testing.T, repository repository.Repository) Service {
	return NewService(testConfig, repository)
}
