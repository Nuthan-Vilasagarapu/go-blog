// session := make(map[string]string)
package stores

import (
	"errors"
	"go-blog/interfaces"
)

var (
	Session = make(map[string]string)
	Blogs   = []interfaces.IBlog{}
	Users   = []interfaces.IUser{}
)

func SetSession(token, user_id string) {
	Session[token] = user_id
}

func GetSession(token string) (*string, error) {
	user_id := Session[token]
	if user_id == "" {
		return nil, errors.New("No Session found")
	}
	return &user_id, nil
}
