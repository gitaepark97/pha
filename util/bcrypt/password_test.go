package bcrypt

import (
	"testing"

	"github.com/gitaepark/pha/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := util.CreateRandomString(10)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := util.CreateRandomString(10)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hashedPassword2, _ := HashPassword(password)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
