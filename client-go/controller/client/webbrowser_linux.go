package client

import (
	"os/exec"
)

// Webbrowser opens a url with the default browser
func Webbrowser(URL string) (err error) {
	_, err = exec.Command("xdg-open", URL).Output()
	return
}
