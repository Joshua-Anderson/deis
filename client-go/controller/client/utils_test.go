package client

import (
	"fmt"
	"os"
	"testing"
)

func TestChooseSettingsFileLocation(t *testing.T) {
	os.Unsetenv("DEIS_PROFILE")
	os.Setenv("HOME", "/home/test")
	expected := "/home/test/.deis/client.json"

	actual := locateSettingsFile()

	if actual != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, actual))
	}
}

func TestChooseSettingsFileUsingProfile(t *testing.T) {
	os.Setenv("DEIS_PROFILE", "testing")
	os.Setenv("HOME", "/home/test")
	expected := "/home/test/.deis/testing.json"

	actual := locateSettingsFile()

	if actual != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, actual))
	}
}

func TestDeleteSettings(t *testing.T) {
	err := createTempProfile("")

	if err != nil {
		t.Fatal(err)
	}

	err = deleteSettings()

	if err != nil {
		t.Fatal(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Error(fmt.Errorf("File %s exists, supposed to have been deleted.", file))
	}
}
