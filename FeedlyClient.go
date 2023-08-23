package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// FeedlyClient represents the client to interact with Feedly's APIs.
type FeedlyClient struct {
	ClientID          string
	ClientSecret      string
	Sandbox           bool
	ServiceHost       string
	AdditionalHeaders map[string]string
	Token             string
	Secret            string
}

// NewFeedlyClient initializes a new Feedly client with the given options.
func NewFeedlyClient(options map[string]interface{}) *FeedlyClient {
	client := &FeedlyClient{
		ClientID:          options["client_id"].(string),
		ClientSecret:      options["client_secret"].(string),
		Sandbox:           options["sandbox"].(bool),
		AdditionalHeaders: options["additional_headers"].(map[string]string),
		Token:             options["token"].(string),
		Secret:            options["secret"].(string),
	}

	if client.Sandbox {
		client.ServiceHost = "sandbox.feedly.com"
	} else {
		client.ServiceHost = "cloud.feedly.com"
	}

	return client
}

// GetCodeURL constructs the URL for OAuth authentication with Feedly.
func (client *FeedlyClient) GetCodeURL(callbackURL string) string {
	scope := "https://cloud.feedly.com/subscriptions"
	responseType := "code"
	endpoint := client.getEndpoint("v3/auth/auth")
	requestURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=%s",
		endpoint, client.ClientID, callbackURL, scope, responseType)
	return requestURL
}

// GetAccessToken retrieves the access token using the authentication code.
func (client *FeedlyClient) GetAccessToken(redirectURI, code string) (map[string]interface{}, error) {
	params := url.Values{
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
		"code":          {code},
	}
	requestURL := client.getEndpoint("v3/auth/token")
	res, err := http.PostForm(requestURL, params)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

// RefreshAccessToken refreshes the access token using the refresh token.
func (client *FeedlyClient) RefreshAccessToken(refreshToken string) (map[string]interface{}, error) {
	params := url.Values{
		"refresh_token": {refreshToken},
		"client_id":     {client.ClientID},
		"client_secret": {client.ClientSecret},
		"grant_type":    {"refresh_token"},
	}
	requestURL := client.getEndpoint("v3/auth/token")
	res, err := http.PostForm(requestURL, params)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

// GetUserProfile fetches the user's profile information.
func (client *FeedlyClient) GetUserProfile(accessToken string) (map[string]interface{}, error) {
	headers := map[string]string{
		"Authorization": "OAuth " + accessToken,
	}
	requestURL := client.getEndpoint("v3/user")
	req, _ := http.NewRequest("GET", requestURL, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

// GetUserSubscriptions fetches the list of user's subscriptions.
func (client *FeedlyClient) GetUserSubscriptions(accessToken string) ([]map[string]interface{}, error) {
	headers := map[string]string{
		"Authorization": "OAuth " + accessToken,
	}
	requestURL := client.getEndpoint("v3/subscriptions")
	req, _ := http.NewRequest("GET", requestURL, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

// GetFeedContent fetches the content of a specific feed.
func (client *FeedlyClient) GetFeedContent(accessToken, streamID string, unreadOnly bool, newerThan int64) (map[string]interface{}, error) {
	headers := map[string]string{
		"Authorization": "OAuth " + accessToken,
	}
	params := url.Values{
		"streamId":    {streamID},
		"unreadOnly":  {strconv.FormatBool(unreadOnly)},
		"newerThan":   {strconv.FormatInt(newerThan, 10)},
	}
	requestURL := client.getEndpoint("v3/streams/contents") + "?" + params.Encode()
	req, _ := http.NewRequest("GET", requestURL, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

// MarkArticleRead marks one or multiple articles as read.
func (client *FeedlyClient) MarkArticleRead(accessToken string, entryIds []string) (*http.Response, error) {
	headers := map[string]string{
		"content-type":  "application/json",
		"Authorization": "OAuth " + accessToken,
	}
	params := map[string]interface{}{
		"action":  "markAsRead",
		"type":    "entries",
		"entryIds": entryIds,
	}
	jsonData, _ := json.Marshal(params)
	requestURL := client.getEndpoint("v3/markers")
	req, _ := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := http.DefaultClient.Do(req)
	return res, err
}

// SaveForLater saves one or multiple articles for later reading.
func (client *FeedlyClient) SaveForLater(accessToken, userID string, entryIds []string) (*http.Response, error) {
	headers := map[string]string{
		"content-type":  "application/json",
		"Authorization": "OAuth " + accessToken,
	}
	requestURL := client.getEndpoint("v3/tags") + "/user%2F" + userID + "%2Ftag%2Fglobal.saved"
	params := map[string]interface{}{
		"entryIds": entryIds,
	}
	jsonData, _ := json.Marshal(params)
	req, _ := http.NewRequest("PUT", requestURL, bytes.NewBuffer(jsonData))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := http.DefaultClient.Do(req)
	return res, err
}

// getEndpoint constructs the URL for the given endpoint path.
func (client *FeedlyClient) getEndpoint(path string) string {
	return "https://" + client.ServiceHost + "/" + path
}
