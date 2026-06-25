// DeleteUser
// UpdateUser

package repository

import (
	"errors"
	"go-blog/interfaces"
	"go-blog/stores"
	"time"

	"github.com/google/uuid"
)

func GetUsers() []interfaces.IUser {
	return stores.Users
}

func GetUserById(id string) *interfaces.IUser {
	for _, author := range GetUsers() {
		if author.Id == id {
			return &author
		}
	}
	return nil
}

func GetUserByUsername(user_name string) *interfaces.IUser {
	for _, author := range GetUsers() {
		if author.UserName == user_name {
			return &author
		}
	}
	return nil
}

func CreateUser(username, passwordHash string) (*interfaces.IUser, error) {
	id := uuid.New()

	user := interfaces.IUser{
		Id:           id.String(),
		PasswordHash: passwordHash,
		UserName:     username,
		CreatedAt:    time.Now(),
	}
	for _, userExist := range GetUsers() {
		if userExist.UserName == username {
			return nil, errors.New("User Already Exists")
		}
	}
	stores.Users = append(stores.Users, user)
	WriteUsersToDb(&stores.DB, stores.Users)
	return &user, nil
}

func UpdateUser(id string, new_user interfaces.IUser) bool {
	for i, user := range stores.Users {
		if user.Id == id {
			if new_user.UserName != "" {
				user.UserName = new_user.UserName
			}
			if !new_user.LastLoginAt.IsZero() {
				user.LastLoginAt = new_user.LastLoginAt
			}
			user.UpdateAt = time.Now()
			stores.Users[i] = user
			WriteUsersToDb(&stores.DB, stores.Users)
			return true
		}
	}
	return false
}
