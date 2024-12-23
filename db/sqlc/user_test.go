package db

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/siddheshRajendraNimbalkar/GO-BANK/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	var name string = util.GenerateRandomName()

	userName := strings.Replace(name, " ", "_", 1)
	userName = userName + strconv.Itoa(int(util.GenerateRandomBalance(10, 99)))

	namePart := strings.Split(name, " ")
	email := userName + "@" + namePart[0] + "." + "com"

	password := strconv.Itoa(int(util.GenerateRandomBalance(10000000, 99999999)))

	arg := CreateUserParams{
		Username:       userName,
		FullName:       name,
		Email:          email,
		HashedPassword: password,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func TestCreateRandomUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetRandomUser(t *testing.T) {
	user := createRandomUser(t)
	getUser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, getUser)
	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.FullName, getUser.FullName)
	require.Equal(t, user.HashedPassword, getUser.HashedPassword)
	require.Equal(t, user.Email, getUser.Email)
	require.Equal(t, user.PasswordChangedAt, getUser.PasswordChangedAt)

}
