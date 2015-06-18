package config

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/deis/deis/client-go/controller/client"
	"github.com/deis/deis/version"
)

const configSetExpected string = `{"values":{"FOO":"bar","TEST":"yes"}}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/apps/example-go/config" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
		}

		if string(body) != configSetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configSetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte{})
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`[0,"hi\n"]`))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write([]byte{})
}

func TestConfigSet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	URL, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, URL: *URL, Token: "abc"}

	testVars := map[string]string{
		"TEST": "yes",
		"FOO":  "bar",
	}

	err = Set(client, "example-go", testVars)

	if err != nil {
		t.Error(err)
	}
}
