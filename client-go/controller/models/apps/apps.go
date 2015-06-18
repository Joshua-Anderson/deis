package apps

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List lists apps on controller
func List(c client.Client) ([]api.App, error) {
	body, status, err := c.BasicRequest("GET", "/v1/apps/", []byte(""))

	if err != nil {
		return []api.App{}, err
	}

	if status != 200 {
		return []api.App{}, errors.New(body)
	}

	apps := api.Apps{}
	err = json.Unmarshal([]byte(body), &apps)

	if err != nil {
		return []api.App{}, err
	}

	return apps.Apps, nil
}

// Create app on controller
func Create(c client.Client, id string) (api.App, error) {
	body := []byte{}

	var err error
	if id != "" {
		req := api.AppCreateRequest{ID: id}
		body, err = json.Marshal(req)

		if err != nil {
			return api.App{}, err
		}
	}

	resBody, status, err := c.BasicRequest("POST", "/v1/apps/", body)

	if err != nil {
		return api.App{}, err
	}

	if status != 201 {
		return api.App{}, errors.New(resBody)
	}

	app := api.App{}
	err = json.Unmarshal([]byte(resBody), &app)

	if err != nil {
		return api.App{}, err
	}

	return app, nil
}

// Get app details from controller
func Get(c client.Client, appID string) (api.App, error) {
	URL := fmt.Sprintf("/v1/apps/%s/", appID)

	body, status, err := c.BasicRequest("GET", URL, []byte{})

	if err != nil {
		return api.App{}, err
	}

	if status != 200 {
		return api.App{}, errors.New(body)
	}

	app := api.App{}

	err = json.Unmarshal([]byte(body), &app)

	if err != nil {
		return api.App{}, err
	}

	return app, nil
}

// Logs retrieves logs from an app
func Logs(c client.Client, appID string, lines int) (string, error) {
	URL := fmt.Sprintf("/v1/apps/%s/logs", appID)

	if lines > 0 {
		URL += "?log_lines=" + strconv.Itoa(lines)
	}

	body, status, err := c.BasicRequest("GET", URL, []byte{})

	if err != nil {
		return "", err
	}

	if status != 200 {
		return body, errors.New(body)
	}

	return strings.Trim(body, `"`), nil
}

// Run one time command in app
func Run(c client.Client, appID string, command string) (api.AppRunResponse, error) {
	req := api.AppRunRequest{Command: command}
	body, err := json.Marshal(req)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	URL := fmt.Sprintf("/v1/apps/%s/run", appID)

	resBody, status, err := c.BasicRequest("POST", URL, body)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	if status != 200 {
		return api.AppRunResponse{}, errors.New(resBody)
	}

	out := make([]interface{}, 2)

	err = json.Unmarshal([]byte(resBody), &out)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	return api.AppRunResponse{Output: out[1].(string), ReturnCode: int(out[0].(float64))}, nil
}

// Delete app
func Delete(c client.Client, appID string) error {
	URL := fmt.Sprintf("/v1/apps/%s/", appID)

	body, status, err := c.BasicRequest("DELETE", URL, []byte{})

	if err != nil {
		return err
	}

	if status != 204 {
		return errors.New(body)
	}

	return nil
}
