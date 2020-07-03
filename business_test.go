package trustpilot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	testBBaseURLPath = "/v1"
	serviceReviews   = map[string]interface{}{
		"reviews": []map[string]interface{}{
			{
				"language": "da",
				"links": []map[string]interface{}{
					{
						"href":   "<Url for the resource>",
						"method": "<Http method for the resource>",
						"rel":    "<Description of the relation>",
					},
				},
				"title": "My review",
				"businessUnit": map[string]interface{}{
					"displayName": "Trustpilot",
					"id":          "507f191e810c19729de860ea",
					"links": []map[string]interface{}{
						{
							"href":   "<Url for the resource>",
							"method": "<Http method for the resource>",
							"rel":    "<Description of the relation>",
						},
					},
				},
				"name": map[string]interface{}{
					"referring": []interface{}{
						"trustpilot.com",
						"www.trustpilot.com",
					},
					"identifying": "trustpilot.com",
				},
				"text":      "This shop is great.",
				"updatedAt": "2013-09-07T13:37:00",
				"createdAt": "2013-09-07T13:37:00",
				"invitation": map[string]interface{}{
					"businessUnitId": "507f191e810c19729de860ea",
				},
				"location": map[string]interface{}{
					"urlFormattedName": "Pilestraede58",
					"id":               "43f51215-a1fc-4c60-b6dd-e4afb6d7b831",
					"name":             "Pilestraede 58",
				},
				"stars": 5,
				"companyReply": map[string]interface{}{
					"text":      "This is our reply.",
					"createdAt": "2013-09-07T13:37:00",
				},
				"consumer": map[string]interface{}{
					"profileUrl": "http://www.trustpilot.com/users/55cc4f3b0000fe0002c4f125",
					"profileImage": map[string]interface{}{
						"image35x35": map[string]interface{}{
							"url":    "<Url for the image>",
							"width":  "<Image width>",
							"height": "<Image height>",
						},
						"image24x24": map[string]interface{}{
							"url":    "<Url for the image>",
							"width":  "<Image width>",
							"height": "<Image height>",
						},
						"image73x73": map[string]interface{}{
							"url":    "<Url for the image>",
							"width":  "<Image width>",
							"height": "<Image height>",
						},
						"image64x64": map[string]interface{}{
							"url":    "<Url for the image>",
							"width":  "<Image width>",
							"height": "<Image height>",
						},
						"displayName": "John Doe",
						"id":          "507f191e810c19729de860ea",
						"links": []map[string]interface{}{
							{
								"href":   "<Url for the resource>",
								"method": "<Http method for the resource>",
								"rel":    "<Description of the relation>",
							},
						},
					},
				},
				"id":         "507f191e810c19729de860ea",
				"isVerified": true,
			},
		},
	}
	respStr = `{
		"reviews": [
		  {
			"language": "da",
			"links": [
			  {
				"href": "<Url for the resource>",
				"method": "<Http method for the resource>",
				"rel": "<Description of the relation>"
			  }
			],
			"title": "My review",
			"businessUnit": {
			  "displayName": "Trustpilot",
			  "id": "507f191e810c19729de860ea",
			  "links": [
				{
				  "href": "<Url for the resource>",
				  "method": "<Http method for the resource>",
				  "rel": "<Description of the relation>"
				}
			  ],
			  "name": {
				"referring": [
				  "trustpilot.com",
				  "www.trustpilot.com"
				],
				"identifying": "trustpilot.com"
			  }
			},
			"text": "This shop is great.",
			"updatedAt": "2013-09-07T13:37:00",
			"createdAt": "2013-09-07T13:37:00",
			"invitation": {
			  "businessUnitId": "507f191e810c19729de860ea"
			},
			"location": {
			  "urlFormattedName": "Pilestraede58",
			  "id": "43f51215-a1fc-4c60-b6dd-e4afb6d7b831",
			  "name": "Pilestraede 58"
			},
			"stars": 5,
			"companyReply": {
			  "text": "This is our reply.",
			  "createdAt": "2013-09-07T13:37:00"
			},
			"consumer": {
			  "profileUrl": "http://www.trustpilot.com/users/55cc4f3b0000fe0002c4f125",
			  "profileImage": {
				"image35x35": {
				  "url": "<Url for the image>",
				  "width": "<Image width>",
				  "height": "<Image height>"
				},
				"image24x24": {
				  "url": "<Url for the image>",
				  "width": "<Image width>",
				  "height": "<Image height>"
				},
				"image73x73": {
				  "url": "<Url for the image>",
				  "width": "<Image width>",
				  "height": "<Image height>"
				},
				"image64x64": {
				  "url": "<Url for the image>",
				  "width": "<Image width>",
				  "height": "<Image height>"
				}
			  },
			  "displayName": "John Doe",
			  "id": "507f191e810c19729de860ea",
			  "links": [
				{
				  "href": "<Url for the resource>",
				  "method": "<Http method for the resource>",
				  "rel": "<Description of the relation>"
				}
			  ]
			},
			"id": "507f191e810c19729de860ea",
			"isVerified": true
		  }
		]
	  }
	`
)

func TestBusiness_getCrededential(t *testing.T) {
	client, mux, _, teardown := bsetup()
	defer teardown()
	mux.HandleFunc("/business-units/find", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryParams(t, r, "name", "Trustpilot")
		testHeader(t, r, "Authorization", "eHh4eHh4eDp4eHh4eHh4")
		fmt.Fprint(w, `{"displayName": "Trustpilot","id": "507f191e810c19729de860ea"}`)
	})
	client.CTX = ctx
	client.ClientID = "xxxxxxx"
	client.ClientSecret = "xxxxxxx"
	got, err := client.Business.GetBusinessCredentials("eHh4eHh4eDp4eHh4eHh4", "Trustpilot")
	if err != nil {
		t.Errorf("TestBusiness_getCrededential returned error: %v", err)
	}
	want := &Business{DisplayName: String("Trustpilot"), ID: String("507f191e810c19729de860ea")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TestBusiness_getCrededential returned auth %+v, want %+v", got, want)
	}
}

func TestBusiness_getServiceReviews(t *testing.T) {
	client, mux, _, teardown := bsetup()
	defer teardown()
	mux.HandleFunc("/reviews/latest", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Authorization", "xxxxxxx")
		fmt.Fprint(w, respStr)
	})
	client.CTX = ctx
	client.ClientID = "xxxxxxx"
	client.ClientSecret = "xxxxxxx"
	got, err := client.Business.GetServiceReviews(3)
	if err != nil {
		t.Errorf("TestBusiness_getServiceReviews returned error: %v", err)
	}
	wantSR := new(ServiceReviews)
	bytes, _ := json.Marshal(&serviceReviews)
	_ = json.Unmarshal(bytes, &wantSR)
	log.Printf("%#v", got)
	if !reflect.DeepEqual(got, wantSR) {
		t.Errorf("TestBusiness_getServiceReviews returned auth %+v, want %+v", got, wantSR)
	}
}

func bsetup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	log.Println("testBaseURLPath == ", testBBaseURLPath)
	apiHandler.Handle(testBBaseURLPath+"/", http.StripPrefix(testBBaseURLPath, mux))
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
	url, err := url.Parse(server.URL + testBBaseURLPath + "/")
	if err != nil {
		panic(err)
	}
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testQueryParams(t *testing.T, r *http.Request, q, want string) {
	t.Helper()
	got := r.URL.Query().Get(q)
	if got != want {
		t.Errorf("Query params is %s, want %s", got, want)
	}
}
