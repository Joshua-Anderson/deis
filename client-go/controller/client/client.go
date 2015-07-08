package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

// Client oversees the interaction between the client and controller
type Client struct {
	// HTTP client used to communicate with the API.
	HTTPClient *http.Client
	SSLVerify  bool

	// URL used to communicate with the controller.
	URL url.URL

	// Token is used to authenticate the request against the API.
	Token string

	// Username is the name of the user performing requests against the API.
	Username string
}

type settingsFile struct {
	Username   string `json:"username"`
	SslVerify  bool   `json:"ssl_verify"`
	Controller string `json:"controller"`
	Token      string `json:"token"`
}

// Load settings from file
func Load() (*Client, error) {
	filename := locateSettingsFile()

	_, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("Not logged in. Use 'deis login' or 'deis register' to get started.")
		}

		return nil, err
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	settings := settingsFile{}
	err = json.Unmarshal(contents, &settings)

	if err != nil {
		return nil, err
	}

	URL, err := url.Parse(settings.Controller)
	if err != nil {
		return nil, err
	}

	return &Client{HTTPClient: CreateHTTPClient(settings.SslVerify), SSLVerify: settings.SslVerify,
		URL: *URL, Token: settings.Token, Username: settings.Username}, nil
}

// Save settings to a file
func (c Client) Save() error {
	settings := settingsFile{Username: c.Username,
		SslVerify:  c.SSLVerify,
		Controller: c.URL.String(), Token: c.Token}

	settingsContents, err := json.Marshal(settings)

	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(os.Getenv("HOME"), "/.deis/"), 0775)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(locateSettingsFile(), settingsContents, 0775)

	if err != nil {
		return err
	}

	return nil
}
