// DeleteUser
// UpdateUser

package repository

import (
	"errors"
	"go-blog/interfaces"
	"time"

	"github.com/google/uuid"
)

func GetUsers() []interfaces.DBUserFmt {
	return ReadUsersFromDb()
}

func GetUserById(id string) *interfaces.DBUserFmt {
	for _, author := range GetUsers() {
		if author.Id == id {
			return &author
		}
	}
	return nil
}

func GetUserByUsername(user_name string) *interfaces.DBUserFmt {
	for _, author := range GetUsers() {
		if author.Data.Username == user_name {
			return &author
		}
	}
	return nil
}

func CreateUser(username, passwordHash string) (*interfaces.DBUserFmt, error) {
	users := GetUsers()
	id := uuid.New()
	user := interfaces.DBUserFmt{
		Id: id.String(),
		Data: interfaces.DBUserDataFmt{
			PasswordHash: passwordHash,
			Username:     username,
			CreatedAt:    time.Now().Unix(),
			LastLoginAt:  time.Now().Unix(),
			UpdatedAt:    int(time.Now().Unix()),
		},
	}
	for _, userExist := range GetUsers() {
		if userExist.Data.Username == username {
			return nil, errors.New("User Already Exists")
		}
	}
	users = append(users, user)
	WriteUsersToDb(users)
	return &user, nil
}

func UpdateUser(id string, new_user interfaces.DBUserDataFmt) bool {
	users := GetUsers()
	for i, user := range users {
		if user.Id == id {
			if new_user.Username != "" {
				user.Data.Username = new_user.Username
			}
			if new_user.LastLoginAt != 0 {
				user.Data.LastLoginAt = new_user.LastLoginAt
			}
			user.Data.UpdatedAt = int(time.Now().Unix())
			users[i] = user

			WriteUsersToDb(users)
			return true
		}
	}
	return false
}
