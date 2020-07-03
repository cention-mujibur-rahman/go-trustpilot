package trustpilot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	baseURL        = "https://api.trustpilot.com/v1/" // api domain
	authURL        = "https://authenticate.trustpilot.com"
	accessTokenURL = "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken"
	refreshURL     = "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/refresh"
	revokeURL      = "https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/revoke"

	//Due to unavailability of the api access, I used a fake server which behaving same as like trustpilot API server
	fakeURL = "http://localhost:8005/trustpilot/"

	isTEST = false
)

// Client manages communication with the trustpilot API.
type Client struct {
	CTX      context.Context
	clientMu sync.Mutex // clientMu protects the client during calls that modify the CheckRedirect func.
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// UserAgent agent used when communicating with Trustpilot API.
	UserAgent string

	// Application client_id
	ClientID string

	// Application client_secret
	ClientSecret string

	// ResponseType is type of response from trustpilot e.g., code, password, implicit
	ResponseType string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the API.
	Authorizations *AuthorizationsService
	Business       *BusinessService
	Product        *ProductService

	// Temporary Response
	Response *Response
}

type service struct {
	client *Client
}

// NewClient returns a new GitHub API client. If a nil httpClient is
// provided, a new http.Client will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	bURL, _ := url.Parse(baseURL)
	c := &Client{client: httpClient, BaseURL: bURL}
	c.common.client = c
	c.Authorizations = (*AuthorizationsService)(&c.common)
	c.Business = (*BusinessService)(&c.common)
	c.Product = (*ProductService)(&c.common)
	return c
}

//Response represents the raw http response and rate limit
type Response struct {
	*http.Response

	// Explicitly specify the Rate type so Rate's String() receiver doesn't
	// propagate to Response.
	Rate Rate
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request) ([]byte, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
