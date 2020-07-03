package trustpilot

import (
	"encoding/json"
	"fmt"
	"log"
)

// BusinessService handles communication with the business related
// methods of the trustpilot API.
//
// GitHub API docs: https://developers.trustpilot.com/business-units-api#get-a-list-of-all-business-units
type BusinessService service

// Business represents an individual business
type Business struct {
	DisplayName *string `json:"displayName,omitmpty"`
	ID          *string `json:"id,omitempty"`
	Links       []Links `json:"links,omitempty"`
}

//Links represents the business links
type Links struct {
	HREF   *string
	Method *string
	REL    *string
}

func (b BusinessService) String() string {
	return Stringify(b)
}

//GetBusinessCredentials Returns the business unit given by the provided name
func (b *BusinessService) GetBusinessCredentials(token, search string) (*Business, error) {
	u := fmt.Sprintf("%s/v1/business-units/find?name=%s", fakeURL, search)
	if isTEST {
		u = fmt.Sprintf("/v1/business-units/find?name=%s", search)
	}
	bs := new(Business)
	req, err := b.client.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Err %v", err)
		return bs, err
	}

	req.Header.Add("Authorization", ""+token)
	resp, err := b.client.Do(b.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return bs, err
	}
	err = json.Unmarshal(resp, &bs) // convert the response data to json

	if err != nil {
		return bs, err
	}
	return bs, nil
}

//ServiceReviews ...
type ServiceReviews struct {
	Reviews []*SingleServiceReview
}

//SingleServiceReview ...
type SingleServiceReview struct {
	Title        *string       `json:"title"`
	Text         *string       `json:"text"`
	UpdatedAt    *string       `json:"updatedAt"`
	CreatedAt    *string       `json:"createdAt"`
	Stars        *int          `json:"stars"`
	BusinessUnit *BusinessUnit `json:"businessUnit"`
	ID           *string       `json:"id"`
	Consumer     *BusinessUnit `json:"consumer"`
}

//BusinessUnit ...
type BusinessUnit struct {
	DisplaName *string `json:"displayName"`
	ID         *string `json:"id"`
}

//GetServiceReviews gets latest reviews by language
//This method gets the latest reviews written in a specfic language.
//
//https://developers.trustpilot.com/service-reviews-api#get-latest-reviews-by-language
func (b *BusinessService) GetServiceReviews(count int) (*ServiceReviews, error) {
	u := fmt.Sprintf("%s/reviews/latest", fakeURL)
	if isTEST {
		u = fmt.Sprintf("/v1/reviews/latest")
	}
	sr := new(ServiceReviews)
	req, err := b.client.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Err %v", err)
		return sr, err
	}

	req.Header.Add("Authorization", ""+b.client.ClientID)
	resp, err := b.client.Do(b.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return sr, err
	}
	err = json.Unmarshal(resp, &sr) // convert the response data to json

	if err != nil {
		return sr, err
	}
	return sr, nil
}

//GetServicePrivateReview Get private review
//This method gets the reviews's basic public information but also some private information (referenceEmail and referenceId)
//and status as either "active" or "reported".
//
//https://developers.trustpilot.com/service-reviews-api#get-private-review
func (b *BusinessService) GetServicePrivateReview(token, reviewID string) (*SingleServiceReview, error) {
	u := fmt.Sprintf("%s/private/reviews/%s", fakeURL, reviewID)
	if isTEST {
		u = fmt.Sprintf("/private/reviews/507f191e810c19729de860ea")
	}
	sr := new(SingleServiceReview)
	req, err := b.client.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Err %v", err)
		return sr, err
	}

	req.Header.Add("Authorization", ""+token)
	resp, err := b.client.Do(b.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return sr, err
	}
	err = json.Unmarshal(resp, &sr) // convert the response data to json

	if err != nil {
		return sr, err
	}
	return sr, nil
}

//ServiceReviewResp decode the response from api
type ServiceReviewResp struct{}

//SendServiceReviews Reply to a review.
//This method will post a reply to a review.
//
//https://developers.trustpilot.com/service-reviews-api#reply-to-a-review-
func (b *BusinessService) SendServiceReviews(token, reviewID, message string) (*ServiceReviewResp, error) {
	u := fmt.Sprintf("%s/private/reviews/%s/reply", fakeURL, reviewID)
	if isTEST {
		u = fmt.Sprintf("/private/reviews/507f191e810c19729de860ea/reply")
	}
	sr := new(ServiceReviewResp)
	reqBody := &struct {
		Message string `json:"message"`
	}{Message: message}

	req, err := b.client.NewRequest("POST", u, reqBody)
	if err != nil {
		log.Printf("Err %v", err)
		return sr, err
	}

	req.Header.Add("Authorization", ""+token)
	resp, err := b.client.Do(b.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return sr, err
	}
	err = json.Unmarshal(resp, &sr)

	if err != nil {
		return sr, err
	}
	return sr, nil
}
