package trello

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	BASE_URL        = "https://trello.com"
	BASE_API_URL    = "https://api.trello.com"
	CALLBACK_METHOD = "fragment" // Always fragment

	ALL StatusEnum = iota
	CLOSED
	NONE
	OPEN
)

type (
	StatusEnum int

	Client struct {
		returnUrl      string
		expiration     string
		apiKey         string
		token          string
		appName        string
		organizationID string
		scopes         []string
	}
)

func (t StatusEnum) GetValue() (string, error) {
	switch t {
	case 0:
		return "ALL", nil
	case 1:
		return "CLOSED", nil
	case 2:
		return "NONE", nil
	case 3:
		return "OPEN", nil
	default:
		return "", fmt.Errorf("INVALID VALUE %d", t)
	}
}

func NewClient(options ...func(*Client)) *Client {
	c := &Client{}
	for _, f := range options {
		f(c)
	}
	return c
}

func WithOrganizationID(id string) func(*Client) {
	return func(c *Client) {
		c.organizationID = id
	}
}

func WithAppName(appName string) func(*Client) {
	return func(c *Client) {
		c.appName = appName
	}
}

func WithReturnURL(returnUrl string) func(*Client) {
	return func(c *Client) {
		c.returnUrl = returnUrl
	}
}

func WithScopes(scopes []string) func(*Client) {
	return func(c *Client) {
		c.scopes = scopes
	}
}

func WithExpiration(expiration string) func(*Client) {
	return func(c *Client) {
		c.expiration = expiration
	}
}

func WithAPIKey(apiKey string) func(*Client) {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func WithToken(token string) func(*Client) {
	return func(c *Client) {
		c.token = token
	}
}

func (c *Client) OrganizationID() string {
	return c.organizationID
}

func (c *Client) authURI() string {
	return fmt.Sprintf(
		"%s/1/authorize?expiration=%s&name=%s&scope=%s&response_type=%s&key=%s&callback_method=%s&return_url=%s",
		BASE_URL,
		c.expiration,
		c.appName,
		strings.Join(c.scopes, ","),
		"token",
		c.apiKey,
		CALLBACK_METHOD,
		c.returnUrl,
	)
}

func (c *Client) NewBoard(req CreateBoardReq) (BoardRes, error) {
	params := fmt.Sprintf("%s&key=%s&token=%s", req.String(), c.apiKey, c.token)
	var body BoardRes = BoardRes{}

	uri := fmt.Sprintf("%s/%s?%s", BASE_API_URL, "1/boards/", params)
	res, err := http.Post(uri, "Content-Type/json", nil)

	if err != nil {
		return body, err
	}

	if res.StatusCode != http.StatusOK {
		// I should return the complete body here... whichi I'm kinda doing
		return body, fmt.Errorf("failed to create board - statusCode: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}

func (c *Client) Boards() ([]BoardRes, error) {
	baseUrl := fmt.Sprintf("%s/%s?key=%s&token=%s", BASE_API_URL, fmt.Sprintf("1/organizations/%s/boards", c.organizationID), c.apiKey, c.token)
	urlParts := []string{baseUrl}
	url := strings.Join(urlParts, "&")
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get all boards")
	}

	var body []BoardRes
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) BoardByName(name string) (BoardRes, error) {
	boards, err := c.Boards()
	if err != nil {
		return BoardRes{}, err
	}

	for _, b := range boards {
		if b.Name == name {
			return b, nil
		}
	}

	return BoardRes{}, fmt.Errorf(" board with name: %s not found", name)
}

func (c *Client) Invite(email, boardID string) error {
	url := fmt.Sprintf("%s/%s?key=%s&token=%s&email=%s", BASE_API_URL, fmt.Sprintf("1/boards/%s/members", boardID), c.apiKey, c.token, email)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Add("Content-Type", "Application/json")
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to get all the lists for the board with id: %++v", req)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}
	fmt.Println("%++v", body)
	return nil
}

func (c *Client) Lists(req GetListReq) ([]ListRes, error) {
	params, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	qParams := params.String()
	baseUrl := fmt.Sprintf("%s/%s?key=%s&token=%s", BASE_API_URL, fmt.Sprintf("1/boards/%s/lists", req.ID), c.apiKey, c.token)
	urlParts := []string{baseUrl}

	if qParams != "" {
		urlParts = append(urlParts, qParams)
	}

	url := strings.Join(urlParts, "&")

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get all the lists for the board with id: %s", req.ID)
	}

	var body []ListRes
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) NewCard(req NewCardReq) (NewCardRes, error) {
	name := strings.Replace(strings.Replace(req.Name, " ", "%20", -1), "\n", "%0A", -1)
	desc := strings.Replace(strings.Replace(req.Description, " ", "%20", -1), "\n", "%0A", -1)
	qParams := fmt.Sprintf("name=%s&desc=%s&idList=%s", name, desc, req.IDList)

	url := fmt.Sprintf("%s/%s?key=%s&token=%s&%s", BASE_API_URL, "1/cards", c.apiKey, c.token, qParams)
	var body NewCardRes

	res, err := http.Post(url, "Content-Type: Application/json", nil)
	if err != nil {
		return body, err
	}

	if res.StatusCode != http.StatusOK {
		return body, fmt.Errorf("Error While creating a card - Status code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return body, err
	}
	return body, nil
}
