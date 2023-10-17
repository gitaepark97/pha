package jwt

import (
	"testing"
	"time"

	"github.com/gitaepark/pha/util"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	userID := util.CreateRandomInt64(1, 10)
	secret := util.CreateRandomString(32)

	token, payload1, err := CreateToken(userID, secret, time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload1)

	require.Equal(t, payload1.UserID, userID)
	require.WithinDuration(t, payload1.ExpiredAt, time.Now(), time.Second)

	payload2, err := VerifyToken(token, secret)
	require.NoError(t, err)
	require.Equal(t, payload2.ID, payload1.ID)
	require.Equal(t, payload2.UserID, payload1.UserID)
	require.WithinDuration(t, payload2.ExpiredAt, payload1.ExpiredAt, time.Second)

	_, err = VerifyToken(util.CreateRandomString(50), secret)
	require.ErrorIs(t, err, ErrInvalidToken)

	time.Sleep(time.Second)

	_, err = VerifyToken(token, secret)
	require.ErrorIs(t, err, ErrExpiredToken)
}
