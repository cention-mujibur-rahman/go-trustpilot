package trustpilot

import (
	"encoding/json"
	"fmt"
	"log"
)

var (
	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

// AuthorizationsService handles communication with the authorization related
// methods of the Trustpilot API.
//
// This service requires HTTP Basic Authentication; it cannot be accessed using
// an OAuth token.
//
// GitHub API docs: https://developers.trustpilot.com/authentication
type AuthorizationsService service

// Authorization represents an individual GitHub authorization.
type Authorization struct {
	AccessToken  *string `json:"access_token,omitmpty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	ExpiresIN    *string `json:"expires_in,omitempty"`
}

func (a Authorization) String() string {
	return Stringify(a)
}

// AuthorizationApp represents an individual trustpilot app (in the context of authorization).
type AuthorizationApp struct {
	AccessToken  *string `json:"access_token,omitmpty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	ExpiresIN    *string `json:"expires_in,omitempty"`
}

func (a AuthorizationApp) String() string {
	return Stringify(a)
}

//AuthorizationCode The user is redirected to a website owned by Trustpilot in order to be authorized. After the authorization succeeds,
//Trustpilot redirects the user back to the client site with a code parameter containing the authorization code.
//
//https://developers.trustpilot.com/authentication
func (a *AuthorizationsService) AuthorizationCode(redirectURL string) (string, error) {
	u := fmt.Sprintf("%s/authenticate?response_type=code&client_id=%v&redirect_uri=%v", fakeURL, a.client.ClientID, redirectURL)
	if isTEST {
		u = fmt.Sprintf("%s?response_type=code&client_id=%v&redirect_uri=%v", "/v1/oauth/oauth-business-users-for-applications/authenticate", a.client.ClientID, redirectURL)
	}
	//https://authenticate.trustpilot.com?client_id=APIKey&redirect_uri=https://www.clientsSite.com&response_type=code
	req, err := a.client.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Err %v", err)
		return "", err
	}
	resp, err := a.client.Do(a.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return "", err
	}
	//Redirects back to: https://www.clientsSite.com/?code=Code as the response
	return string(resp), nil
}

//RetrieveAccessToken Convert an authorization code to an access token
func (a *AuthorizationsService) RetrieveAccessToken(code string, redirectURL string) (*Authorization, error) {
	//https://api.trustpilot.com/v1/oauth/oauth-business-users-for-applications/accesstoken
	// send the request
	//grant_type=authorization_code&code=Code&redirect_uri=https://www.clientsSite.com
	reqBody := &struct {
		GrantType   string `json:"grant_type"`
		Code        string `json:"code"`
		RedirectURI string `json:"redirect_uri"`
	}{GrantType: "authorization_code", Code: code, RedirectURI: redirectURL}
	if isTEST {
		accessTokenURL = fmt.Sprintf("%s", "/v1/oauth/oauth-business-users-for-applications/accesstoken")
	} else {
		accessTokenURL = fakeURL + "accesstoken"
	}
	req, err := a.client.NewRequest("POST", accessTokenURL, reqBody)
	if err != nil {
		return nil, err
	}
	auth := new(Authorization)
	bc := base64Encode(fmt.Sprintf("%s:%s", a.client.ClientID, a.client.ClientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+bc)
	resp, err := a.client.Do(a.client.CTX, req)
	if err != nil {
		return auth, err
	}
	err = json.Unmarshal(resp, &auth) // convert the response data to json

	if err != nil {
		return nil, err
	}
	return auth, nil
}
