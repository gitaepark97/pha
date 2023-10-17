package repository

import (
	"context"
	"encoding/binary"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/corpix/uarand"
	"github.com/gitaepark/pha/util/jwt"
	"github.com/stretchr/testify/require"
)

var (
	userAgent string
	clientIp  string
)

func init() {
	userAgent = uarand.GetRandom()
	buf := make([]byte, 4)
	ip := rand.Uint32()
	binary.LittleEndian.PutUint32(buf, ip)
	clientIp = net.IP(buf).To4().String()
}

func TestCreateSession(t *testing.T) {
	user := getRandomUser(t)

	createRandomSession(t, user)
}

func TestGetSesssion(t *testing.T) {
	user := getRandomUser(t)
	refreshToken, refreshPayload := createRandomSession(t, user)

	session, err := testQueries.GetSession(context.Background(), refreshPayload.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, session.ID, refreshPayload.ID)
	require.Equal(t, session.UserID, user.ID)
	require.Equal(t, session.RefreshToken, refreshToken)
	require.Equal(t, session.UserAgent, userAgent)
	require.Equal(t, session.ClientIp, clientIp)
	require.False(t, session.IsBlocked)
	require.WithinDuration(t, session.ExpiredAt, refreshPayload.ExpiredAt, time.Second)
}

func createRandomSession(t *testing.T, user User) (string, *jwt.Payload) {
	refreshToken, refreshPayload, _ := jwt.CreateToken(user.ID, testConfig.JWTSecret, testConfig.RefreshTokenDuration)

	arg := CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)

	return refreshToken, refreshPayload
}
