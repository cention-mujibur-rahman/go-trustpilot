package trustpilot

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	ctx             = context.Background()
	testBaseURLPath = "/v1/oauth/oauth-business-users-for-applications"
)

func TestAuthorizationsService_Code(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if r.URL.Query().Get("client_id") != "xxxxxxx" {
			fmt.Fprint(w, `blah/blah?code=wrong`)
			return
		}
		if r.URL.Query().Get("redirect_uri") != "blah/ohai" {
			fmt.Fprint(w, `blah/blah?code=wrong`)
			return
		}
		fmt.Fprint(w, `blah/blah?code=612576sadxas`)
	})
	client.CTX = ctx
	client.ClientID = "xxxxxxx"
	got, err := client.Authorizations.AuthorizationCode("blah/ohai")
	if err != nil {
		t.Errorf("Authorizations.AuthorizationCode returned error: %v", err)
	}
	want := "blah/blah?code=612576sadxas"
	if got != want {
		t.Errorf("Authorizations.AuthorizationCode returned auth %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_RetriveToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/accesstoken", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"grant_type":"authorization_code","code":"612576sadxas","redirect_uri":"blah/ohai"}`+"\n")
		//testHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		testHeader(t, r, "Authorization", "Basic eHh4eHh4eDp4eHh4eHh4")
		fmt.Fprint(w, `{"access_token":"12345abc","refresh_token":"deadmeat"}`)
	})
	client.CTX = ctx
	client.ClientID = "xxxxxxx"
	client.ClientSecret = "xxxxxxx"
	got, err := client.Authorizations.RetrieveAccessToken("612576sadxas", "blah/ohai")
	if err != nil {
		t.Errorf("TestAuthorizationsService_RetriveTokenTestAuthorizationsService_RetriveToken returned error: %v", err)
	}
	want := &Authorization{AccessToken: String("12345abc"), RefreshToken: String("deadmeat")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TestAuthorizationsService_RetriveToken returned auth %+v, want %+v", got, want)
	}
}

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(testBaseURLPath+"/", http.StripPrefix(testBaseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the trustpilot client being tested and is
	// configured to use test server.
	client = NewClient(nil)
	url, err := url.Parse(server.URL + testBaseURLPath + "/")
	if err != nil {
		panic(err)
	}
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	t.Helper()
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}
