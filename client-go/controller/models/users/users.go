package users

import (
	"encoding/json"
	"errors"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List users registered with the controller
func List(c client.Client) ([]api.User, error) {
	body, status, err := c.BasicRequest("GET", "/v1/users/", []byte(""))

	if err != nil {
		return []api.User{}, err
	}

	if status != 200 {
		return []api.User{}, errors.New(body)
	}

	users := api.Users{}
	err = json.Unmarshal([]byte(body), &users)

	if err != nil {
		return []api.User{}, err
	}

	return users.Users, nil
}
