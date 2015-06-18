package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/deis/deis/client-go/controller/api"
)

// Register a user with the controller
func Register(url *url.URL, username string, password string, email string, sslVerify bool,
	loginAfter bool) error {
	client := CreateHTTPClient(sslVerify)

	// Create a JSON body of request from user api
	user := api.AuthRegisterRequest{Username: username, Password: password, Email: email}
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	url.Path = "/v1/auth/register/"

	headers := http.Header{}

	controllerClient, err := Load()

	if err == nil {
		headers.Add("Authorization", "token "+controllerClient.Token)
	}

	headers.Add("Content-Type", "application/json")
	addUserAgent(&headers)

	res, err := rawRequest(client, "POST", url.String(), bytes.NewBuffer(body), headers, 201)

	if err != nil {
		return err
	}

	res.Body.Close()

	fmt.Printf("Registered %s\n", username)

	if loginAfter {
		return Login(url, username, password, sslVerify)
	}

	return nil
}

// Login to the controller
func Login(url *url.URL, username string, password string, sslVerify bool) error {
	client := CreateHTTPClient(sslVerify)

	controllerURL := *url

	// Create a JSON body of request from api
	user := api.AuthLoginRequest{Username: username, Password: password}
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	url.Path = "/v1/auth/login/"

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	addUserAgent(&headers)

	res, err := rawRequest(client, "POST", url.String(), bytes.NewBuffer(body), headers, 200)

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)

	token := api.AuthLoginResponse{}

	err = json.Unmarshal([]byte(resBody), &token)
	res.Body.Close()

	if err != nil {
		return err
	}

	controllerClient := Client{Username: username, SSLVerify: sslVerify,
		URL: controllerURL, Token: token.Token}

	err = controllerClient.Save()

	if err != nil {
		return nil
	}

	fmt.Printf("Logged in as %s\n", username)
	return nil
}

// Logout from a controller
func Logout() error {
	err := deleteSettings()

	if err != nil {
		return err
	}

	fmt.Println("Logged out")
	return nil
}

// Passwd changes a users password
func Passwd(username string, password string, newPassword string) error {
	client, err := Load()

	if err != nil {
		return err
	}

	req := api.AuthPasswdRequest{Password: password, NewPassword: newPassword}

	if username != "" {
		req.Username = username
	}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	resBody, status, err := client.BasicRequest("POST", "/v1/auth/passwd/", body)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("Password change failed: %s", resBody)
	}

	fmt.Println("Password change succeeded.")
	return nil
}

// Whoami prints logged in user
func Whoami(writer io.Writer) error {
	client, err := Load()

	if err != nil {
		return err
	}

	fmt.Fprintf(writer, "You are %s at %s\n", client.Username, client.URL.String())
	return nil
}

// Cancel deletes a users account
func Cancel() error {
	client, err := Load()

	if err != nil {
		return err
	}

	body, status, err := client.BasicRequest("DELETE", "/v1/auth/cancel/", []byte{})

	if status != 204 {
		return fmt.Errorf("Cancelation failed: %s", body)
	}

	err = deleteSettings()

	if err != nil {
		return err
	}

	fmt.Println("Account cancelled")
	return nil
}

// Regenerate regenenerates a user's token
func Regenerate(username string, all bool) error {
	client, err := Load()

	if err != nil {
		return err
	}

	var body []byte

	if all == true {
		body, err = json.Marshal(api.AuthRegenerateRequest{All: all})
	} else if username != "" {
		body, err = json.Marshal(api.AuthRegenerateRequest{Name: username})
	} else {
		body = []byte{}
	}

	if err != nil {
		return err
	}

	resBody, status, err := client.BasicRequest("POST", "/v1/auth/tokens/", body)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("Token regeneration failed: %s", resBody)
	}

	// Update the user's token
	if username == "" && all == false {
		token := api.AuthRegenerateResponse{}
		err = json.Unmarshal([]byte(resBody), &token)

		if err != nil {
			return nil
		}

		client.Token = token.Token

		err = client.Save()

		if err != nil {
			return nil
		}
	}

	fmt.Println("Token Regenerated")
	return nil
}
