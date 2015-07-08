package client

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CreateRemote adds a git remote in the current directory
func (c Client) CreateRemote(remote, appID string) error {
	_, err := exec.Command("git", "remote", "add", remote, c.RemoteURL(appID)).Output()

	if err != nil {
		return err
	}

	fmt.Printf("Git remote %s added\n", remote)

	return nil
}

// DeleteRemote removes a git remote in the current directory
func (c Client) DeleteRemote(appID string) error {
	name, err := remoteNameFromAppID(appID)

	fmt.Printf("'%s'\n", name)

	if err != nil {
		return err
	}

	_, err = exec.Command("git", "remote", "remove", name).Output()

	if err != nil {
		return err
	}

	fmt.Printf("Git remote %s removed\n", name)

	return nil
}

func remoteNameFromAppID(appID string) (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	for _, line := range strings.Split(cmd, "\n") {
		if strings.Contains(line, appID) {
			return strings.Split(line, "\t")[0], nil
		}
	}

	return "", errors.New("Could not find remote matching app in 'git remote -v'")
}

// DetectApp detects if there is deis remote in git
func (c Client) DetectApp() (string, error) {
	remote, err := c.findRemote()

	if err != nil {
		return "", err
	}

	ss := strings.Split(remote, "/")
	return strings.Split(ss[len(ss)-1], ".")[0], nil
}

func (c Client) findRemote() (string, error) {
	out, err := exec.Command("git", "remote", "-v").Output()

	if err != nil {
		return "", err
	}

	cmd := string(out)

	for _, line := range strings.Split(cmd, "\n") {
		for _, remote := range strings.Split(line, " ") {
			if strings.Contains(remote, c.URL.Host) {
				return strings.Split(remote, "\t")[1], nil
			}
		}
	}

	return "", errors.New("Could not find deis remote in 'git remote -v'")
}

// RemoteURL returns he git URL of app
func (c Client) RemoteURL(appID string) string {
	return fmt.Sprintf("ssh://git@%s:2222/%s.git", c.URL.Host, appID)
}
