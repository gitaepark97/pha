package repository

import (
	"context"
	"testing"

	"github.com/gitaepark/pha/util"
	"github.com/gitaepark/pha/util/bcrypt"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	getRandomUser(t)
}

func createRandomUser(t *testing.T) (string, string) {
	phoneNumber := util.CreateRandomPhoneNumber()
	password := util.CreateRandomString(10)
	hashedPassword, _ := bcrypt.HashPassword(password)

	arg := CreateUserParams{
		PhoneNumber:    phoneNumber,
		HashedPassword: hashedPassword,
	}

	err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	return phoneNumber, password
}

func getRandomUser(t *testing.T) User {
	phoneNumber, password := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), phoneNumber)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.Equal(t, user.PhoneNumber, phoneNumber)
	require.NoError(t, bcrypt.CheckPassword(password, user.HashedPassword))
	require.NotZero(t, user.CreatedAt)

	return user
}
