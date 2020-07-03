package trustpilot

import (
	"encoding/json"
	"fmt"
	"log"
)

// ProductService handles communication with the product review related
// methods of the trustpilot API.
//
// Trustpilot Product Review API docs: https://developers.trustpilot.com/product-reviews-api
type ProductService service

// Product represents an individual product under a business
type Product struct{}

func (p ProductService) String() string {
	return Stringify(p)
}

//ProductReviews ...
type ProductReviews struct {
	Reviews []*SingleProductReview `json:"productReviews"`
}

//SingleProductReview ...
type SingleProductReview struct {
	Content     *string       `json:"content"`
	CreatedAt   *string       `json:"createdAt"`
	Stars       *int          `json:"stars"`
	ID          *string       `json:"id"`
	Consumer    *BusinessUnit `json:"consumer"`
	Links       []*Links
	Attachments []*Attachments
}

//Attachments product review attachment
type Attachments struct {
	State         *string `json:"state"`
	ID            *string `json:"id"`
	ProcessedFile []*TFiles
}

//TFiles product review files
type TFiles struct {
	MimeType  *string `json:"mimeType"`
	URL       *string `json:"url"`
	Dimension *string `json:"dimension"`
}

//GetProductReviews gets product reviews
//This method allows you to get business unit product reviews for SKUs and / or productUrls.
//Note: Even though productUrl and sku are listed as optional parameters at least one of them must be specified.
//It includes review content, date of creation of review, individual star rating, id and display name of the consumer who wrote the review.
//Pagination and filtering reviews by language is also possible.
//
//https://developers.trustpilot.com/product-reviews-api#get-product-reviews
func (p *ProductService) GetProductReviews(businessUnitID string) (*ProductReviews, error) {
	u := fmt.Sprintf("%s/product-reviews/business-units/%s/reviews", fakeURL, businessUnitID)
	if isTEST {
		u = "/product-reviews/business-units/507f191e810c19729de860ea/reviews"
	}
	pr := new(ProductReviews)
	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Err %v", err)
		return pr, err
	}

	req.Header.Add("Authorization", ""+p.client.ClientID)
	resp, err := p.client.Do(p.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return pr, err
	}
	err = json.Unmarshal(resp, &pr)

	if err != nil {
		return pr, err
	}
	return pr, nil
}

//GetProductPrivateReviews Get private product review
//Given a list of SKUs or product urls return a list of product reviews. This method includes private information such as consumer e-mail and reference id.
//By default only published reviews are returned. To get reviews with other states, provide a list in the state field. Pagination is used to retrieve all results.
//
//https://developers.trustpilot.com/product-reviews-api#get-product-reviews
func (p *ProductService) GetProductPrivateReviews(token, businessUnitID string) (*ProductReviews, error) {
	u := fmt.Sprintf("%s/private/product-reviews/business-units/%s/reviews", fakeURL, businessUnitID)
	if isTEST {
		u = "/private/product-reviews/business-units/507f191e810c19729de860ea/reviews"
	}
	pr := new(ProductReviews)
	req, err := p.client.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("Err %v", err)
		return pr, err
	}

	req.Header.Add("Authorization", ""+token)
	resp, err := p.client.Do(p.client.CTX, req)
	if err != nil {
		log.Printf("Err1 %v", err)
		return pr, err
	}
	err = json.Unmarshal(resp, &pr)

	if err != nil {
		return pr, err
	}
	return pr, nil
}
