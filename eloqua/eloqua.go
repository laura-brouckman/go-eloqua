package eloqua

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var _ = fmt.Printf
var _ = ioutil.ReadAll

// The client manages communications with the Eloqua API
type Client struct {
	client *http.Client

	// The base URL for the eloqua instance
	BaseURL string

	// Eloqua login company name
	companyName string
	// Eloqua login user name
	userName string
	// Eloqua login password
	password string
	// Basic auth header value
	authHeader string

	// Various components of the API
	Emails *EmailService
}

// NewClient creates a new instance of an Eloqua HTTP client
// used to interface with the Eloqua API.
func NewClient(baseURL string, companyName string, userName string, password string) *Client {

	authString := companyName + "\\" + userName + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	c := &Client{
		client:      http.DefaultClient,
		BaseURL:     strings.Trim(baseURL, " /"),
		companyName: companyName,
		userName:    userName,
		password:    password,
		authHeader:  "Basic " + encodedAuth,
	}

	// Create services
	c.Emails = &EmailService{client: c}

	return c
}

// Custom eloqua response type
type Response struct {
	*http.Response

	// Variables used in listing operations.
	// Will remain zero-valued for other operations

	// The main body containing the request entities
	Elements json.RawMessage `json:"elements,omitempty"`
	// The current page of the response
	Page int `json:"page,omitempty"`
	// The page size of the response
	PageSize int `json:"pageSize,omitempty"`
	// The total entities found in the query
	Total int `json:"total,omitempty"`
}

// newReponse creates a new custom Response for the given http.Response
func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

// Options for listing requests
type ListOptions struct {
	// Level of detail returned from request
	// Values: "minimal", "partial", "complete"
	Depth string `url:"depth,omitempty"`
	// Number of entities to return
	Count int `url:"count,omitempty"`
	// The page count of entities to return, Starting at 1
	Page int `url:"page,omitempty"`
	// A term for searching through entities
	Search string `url:"search,omitempty"`
	// The property on which to sort the returned data
	Sort string `url:"sort,omitempty"`
	// The direction of the applied sort
	SortDir string `url:"dir,omitempty"`
	// The field on which to order results
	OrderBy string `url:"orderBy,omitempty"`
	// A minimum last updated timestamp
	LastUpdatedAt int `url:"lastUpdatedAt,omitempty"`
}

// Perform a request to the Eloqua API
// Flexible to allow any use of any endpoint but it only returns a
// simple respose.
func (c *Client) RestRequest(endpoint string, method string, jsonData string) (*Response, error) {
	url := c.BaseURL + "/api/rest/2.0/" + strings.Trim(endpoint, " /")
	// fmt.Println(jsonData)
	jsonStr := []byte(jsonData)
	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return newResponse(resp), nil
}

// Performs a GET request and decodes the response into the provided interface
func (c *Client) getRequestDecode(endpoint string, v interface{}) (*Response, error) {
	resp, err := c.RestRequest(endpoint, "GET", "")
	defer resp.Body.Close()
	// TODO - check the response further for common eloqua responses
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

// Performs a GET request for a listing endpoint and decodes the response into the provided interface
func (c *Client) getRequestListDecode(endpoint string, v interface{}, opts *ListOptions) (*Response, error) {

	// Create our options if not set
	if opts == nil {
		opts = &ListOptions{}
	}
	// Set a default minimal depth
	if opts.Depth == "" {
		opts.Depth = "minimal"
	}

	encoder, _ := query.Values(opts)
	endpoint += "?" + encoder.Encode()

	resp, err := c.RestRequest(endpoint, "GET", "")

	defer resp.Body.Close()
	// TODO - check the response further for common eloqua responses
	if err != nil {
		return resp, err
	}

	// Decode response
	if resp != nil {
		err = json.NewDecoder(resp.Body).Decode(resp)
	}

	// Decode elements onto model
	if v != nil {
		err = json.Unmarshal(resp.Elements, v)
	}

	return resp, err
}

// Performs a HTTP request using the given method
// and decodes the response into the provided interface
func (c *Client) requestDecode(endpoint string, method string, v interface{}) (*Response, error) {

	postBody := ""

	if v != nil {
		jsonString, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		postBody = string(jsonString)
	}

	resp, err := c.RestRequest(endpoint, strings.ToUpper(method), postBody)
	defer resp.Body.Close()
	// TODO - check the response further for common eloqua responses
	if err != nil {
		return resp, err
	}

	// content, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(content))

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
	}

	return resp, err
}

// Performs a POST request and decodes the response into the provided interface
func (c *Client) postRequestDecode(endpoint string, v interface{}) (*Response, error) {
	return c.requestDecode(endpoint, "POST", v)
}

// Performs a PUT request and decodes the response into the provided interface
func (c *Client) putRequestDecode(endpoint string, v interface{}) (*Response, error) {
	return c.requestDecode(endpoint, "PUT", v)
}

// Performs a PUT request and decodes the response into the provided interface
func (c *Client) deleteRequest(endpoint string, v interface{}) (*Response, error) {
	postBody := ""

	if v != nil {
		jsonString, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		postBody = string(jsonString)
	}

	return c.RestRequest(endpoint, "DELETE", postBody)
}