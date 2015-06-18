package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/deis/deis/client-go/controller/client"
	"golang.org/x/crypto/ssh/terminal"
)

// Register a user with the controller
func Register(controller string, username string, password string, email string,
	sslVerify bool) error {

	url, err := url.Parse(controller)

	if err != nil {
		return err
	}

	url, err = chooseScheme(*url)

	if err != nil {
		return err
	}

	err = client.CheckConection(client.CreateHTTPClient(sslVerify), *url)

	if err != nil {
		return err
	}

	if username == "" {
		fmt.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Printf("\npassword (confirm): ")
		passwordConfirm, err := readPassword()
		fmt.Println()

		if err != nil {
			return err
		}

		if password != passwordConfirm {
			return errors.New("Password mismatch, aborting registration.")
		}
	}

	if email == "" {
		fmt.Print("email: ")
		fmt.Scanln(&email)
	}

	return client.Register(url, username, password, email, sslVerify, true)
}

// Login to a controller
func Login(controller string, username string, password string, sslVerify bool) error {
	url, err := url.Parse(controller)

	if err != nil {
		return err
	}

	url, err = chooseScheme(*url)

	if err != nil {
		return err
	}

	err = client.CheckConection(client.CreateHTTPClient(sslVerify), *url)

	if err != nil {
		return err
	}

	if username == "" {
		fmt.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	return client.Login(url, username, password, sslVerify)
}

// Logout from a controller
func Logout() error {
	return client.Logout()
}

// Passwd changes a users password
func Passwd(username string, password string, newPassword string) error {
	var err error

	if password == "" {
		fmt.Print("current password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	if newPassword == "" {
		fmt.Print("new password: ")
		newPassword, err = readPassword()
		fmt.Printf("\nnew password (confirm): ")
		passwordConfirm, err := readPassword()

		fmt.Println()

		if err != nil {
			return err
		}

		if newPassword != passwordConfirm {
			return errors.New("Password mismatch, not changing.")
		}
	}

	return client.Passwd(username, password, newPassword)
}

// Whoami prints the logged in user
func Whoami() error {
	return client.Whoami(os.Stdout)
}

// Cancel deletes a users account
func Cancel(username string, password string, yes bool) error {
	c, err := client.Load()

	if err != nil {
		return err
	}

	fmt.Println("Please log in again in order to cancel this account")

	err = Login(c.URL.String(), username, password, c.SSLVerify)

	if err != nil {
		return err
	}

	if yes == false {
		confirm := ""

		c, err = client.Load()

		if err != nil {
			return err
		}

		fmt.Printf("cancel account %s at %s? (y/N): ", c.Username, c.URL.String())
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) == "y" {
			yes = true
		}
	}

	if yes == false {
		fmt.Println("Account not changed")
		return nil
	}

	return client.Cancel()
}

// Regenerate regenenerates a user's token
func Regenerate(username string, all bool) error {
	return client.Regenerate(username, all)
}

func readPassword() (string, error) {
	password, err := terminal.ReadPassword(0)

	return string(password), err
}

func chooseScheme(URL url.URL) (*url.URL, error) {
	if URL.Scheme == "" {
		URL.Scheme = "http"
		URL, err := URL.Parse(URL.String())

		if err != nil {
			return &url.URL{}, err
		}

		return URL, nil
	}

	return &URL, nil
}
