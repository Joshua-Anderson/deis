package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/deis/deis/version"
)

// Register Expected
const rE string = `{"username":"test","password":"opensesame","email":"test@example.com"}`
const lE string = `{"username":"test","password":"opensesame"}`
const pE string = `{"username":"test","password":"old","new_password":"new"}`
const rAE string = `{"all":true}`
const rUE string = `{"username":"test"}`

type fakeAuthHTTPServer struct {
	regenBodyEmpty    bool
	regenBodyAll      bool
	regenBodyUsername bool
}

func (f fakeAuthHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/auth/register/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
		}

		if string(body) != rE {
			fmt.Printf("Expected '%s', Got '%s'\n", rE, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte{})
		return
	}

	if req.URL.Path == "/v1/auth/login/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
		}

		if string(body) != lE {
			fmt.Printf("Expected '%s', Got '%s'\n", lE, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
			return
		}

		res.Write([]byte(`{"token":"abc"}`))
		return
	}

	if req.URL.Path == "/v1/auth/passwd/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
		}

		if string(body) != pE {
			fmt.Printf("Expected '%s', Got '%s'\n", lE, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
			return
		}

		res.Write([]byte{})
		return
	}

	if req.URL.Path == "/v1/auth/tokens/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
		}

		if string(body) == rAE && !f.regenBodyAll {
			f.regenBodyAll = true
			res.Write([]byte{})
			return
		} else if string(body) == rUE && !f.regenBodyUsername {
			f.regenBodyUsername = true
			res.Write([]byte{})
			return
		} else if !f.regenBodyEmpty {
			f.regenBodyEmpty = true
			res.Write([]byte(`{"token":"abc"}`))
			return
		}

		fmt.Printf("Expected '%s', Got '%s'\n", lE, body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte{})
		return
	}

	if req.URL.Path == "/v1/auth/cancel/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte{})
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte{})
}

func TestRegister(t *testing.T) {
	t.Parallel()

	handler := fakeAuthHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	URL, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	err = Register(URL, "test", "opensesame", "test@example.com", false, false)

	if err != nil {
		t.Error(err)
	}
}

func TestLogin(t *testing.T) {
	err := createTempProfile("")

	handler := fakeAuthHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	controllerURL, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	URL := *controllerURL

	err = Login(&URL, "test", "opensesame", false)

	if err != nil {
		t.Error(err)
	}

	client, err := Load()

	if err != nil {
		t.Fatal(err)
	}

	if client.URL.String() != controllerURL.String() {
		t.Error(fmt.Errorf("Expected %s, Got %s", controllerURL.String(), client.URL.String()))
	}

	expected := "test"
	if client.Username != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Username))
	}

	expected = "abc"
	if client.Token != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Token))
	}

	expectedB := false
	if client.SSLVerify != expectedB {
		t.Error(fmt.Errorf("Expected %t, Got %t", expectedB, client.SSLVerify))
	}
}

func TestWhoAmI(t *testing.T) {
	err := createTempProfile(sFile)

	if err != nil {
		t.Fatal(err)
	}

	testWriter := bytes.Buffer{}

	err = Whoami(&testWriter)

	if err != nil {
		t.Fatal(err)
	}

	expected := "You are t at http://d.t\n"
	actual := testWriter.String()

	if expected != actual {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, actual))
	}
}

func TestLogout(t *testing.T) {
	err := createTempProfile(sFile)

	if err != nil {
		t.Fatal(err)
	}

	err = Logout()

	if err != nil {
		t.Fatal(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Error(fmt.Errorf("File %s exists, supposed to have been deleted.", file))
	}
}

func TestPasswd(t *testing.T) {
	handler := fakeAuthHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	sF := fmt.Sprintf(`{"username":"t","ssl_verify":false,"controller":"%s","token":"a"}`, server.URL)
	err := createTempProfile(sF)

	if err != nil {
		t.Fatal(err)
	}

	err = Passwd("test", "old", "new")

	if err != nil {
		t.Error(err)
	}
}

func TestCancel(t *testing.T) {
	handler := fakeAuthHTTPServer{regenBodyEmpty: false, regenBodyAll: false,
		regenBodyUsername: false}
	server := httptest.NewServer(handler)
	defer server.Close()

	sF := fmt.Sprintf(`{"username":"t","ssl_verify":false,"controller":"%s","token":"a"}`, server.URL)
	err := createTempProfile(sF)

	if err != nil {
		t.Fatal(err)
	}

	err = Regenerate("", true)

	if err != nil {
		t.Error(err)
	}

	err = Cancel()

	if err != nil {
		t.Error(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Error(fmt.Errorf("File %s exists, supposed to have been deleted.", file))
	}
}

func TestRegenerate(t *testing.T) {
	handler := fakeAuthHTTPServer{regenBodyEmpty: false, regenBodyAll: false,
		regenBodyUsername: false}
	server := httptest.NewServer(handler)
	defer server.Close()

	sF := fmt.Sprintf(`{"username":"t","ssl_verify":false,"controller":"%s","token":"a"}`, server.URL)
	err := createTempProfile(sF)

	if err != nil {
		t.Fatal(err)
	}

	err = Regenerate("", true)

	if err != nil {
		t.Error(err)
	}

	err = Regenerate("test", false)

	if err != nil {
		t.Error(err)
	}

	err = Regenerate("", false)

	if err != nil {
		t.Error(err)
	}

	client, err := Load()

	if err != nil {
		t.Error(err)
	}

	expected := "abc"
	if client.Token != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, client.Token))
	}
}
