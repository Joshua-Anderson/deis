package keys

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List keys
func List(c client.Client) ([]api.Key, error) {
	body, status, err := c.BasicRequest("GET", "/v1/keys/", []byte(""))

	if err != nil {
		return []api.Key{}, err
	}

	if status != 200 {
		return []api.Key{}, errors.New(body)
	}

	keys := api.Keys{}
	err = json.Unmarshal([]byte(body), &keys)

	if err != nil {
		return []api.Key{}, err
	}

	return keys.Keys, nil
}

// Create key
func Create(c client.Client, id string, pubKey string) (api.Key, error) {
	req := api.KeyCreateRequest{ID: id, Public: pubKey}
	body, err := json.Marshal(req)

	resBody, status, err := c.BasicRequest("POST", "/v1/keys/", body)

	if err != nil {
		return api.Key{}, err
	}

	if status != 201 {
		return api.Key{}, errors.New(resBody)
	}

	key := api.Key{}
	err = json.Unmarshal([]byte(resBody), &key)

	if err != nil {
		return api.Key{}, err
	}

	return key, nil
}

// Delete key
func Delete(c client.Client, keyID string) error {
	URL := fmt.Sprintf("/v1/keys/%s", keyID)

	body, status, err := c.BasicRequest("DELETE", URL, []byte{})

	if err != nil {
		return err
	}

	if status != 204 {
		return errors.New(body)
	}

	return nil
}
