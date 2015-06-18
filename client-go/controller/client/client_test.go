package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"testing"
)

const sFile string = `{"username":"t","ssl_verify":false,"controller":"http://d.t","token":"a"}`

func createTempProfile(contents string) error {
	name, err := ioutil.TempDir("", "client")

	if err != nil {
		return err
	}

	os.Unsetenv("DEIS_PROFILE")
	os.Setenv("HOME", name)
	folder := path.Join(name, "/.deis/")
	err = os.Mkdir(folder, 0755)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(folder, "client.json"), []byte(contents), 0775)

	if err != nil {
		return err
	}

	return nil
}

func TestLoadSave(t *testing.T) {
	err := createTempProfile(sFile)

	if err != nil {
		t.Fatal(err)
	}

	client, err := Load()

	if err != nil {
		t.Fatal(err)
	}

	expectedB := false
	if client.SSLVerify != expectedB {
		t.Error(fmt.Errorf("Expected %t, Got %t", expectedB, client.SSLVerify))
	}

	expected := "a"
	if client.Token != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Token))
	}

	expected = "t"
	if client.Username != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Username))
	}

	expected = "http://d.t"
	if client.URL.String() != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.URL.String()))
	}

	client.SSLVerify = true
	client.Token = "b"
	client.Username = "c"

	URL, err := url.Parse("http://deis.test")

	if err != nil {
		t.Fatal(err)
	}

	client.URL = *URL

	err = client.Save()

	if err != nil {
		t.Fatal(err)
	}

	client, err = Load()

	expectedB = true
	if client.SSLVerify != expectedB {
		t.Error(fmt.Errorf("Expected %t, Got %t", expectedB, client.SSLVerify))
	}

	expected = "b"
	if client.Token != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Token))
	}

	expected = "c"
	if client.Username != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Username))
	}

	expected = "http://deis.test"
	if client.URL.String() != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.URL.String()))
	}
}
