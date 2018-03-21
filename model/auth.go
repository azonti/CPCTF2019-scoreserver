package model

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"gopkg.in/resty.v1"
	"net/url"
	"os"
)

var authType = map[string]string{
	"twitter": "OAuth1.0a",
}
var oauth1Config = map[string]*oauth1.Config{
	"twitter": {
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("TWITTER_CALLBACK_URL"),
		Endpoint:       twitter.AuthorizeEndpoint,
	},
}

//GetAuthoURL Get an Authorization URL
func GetAuthoURL(provider string) (*url.URL, error) {
	switch authType[provider] {
	case "OAuth1.0a":
		requestToken, _, err := oauth1Config[provider].RequestToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get a request token: %v", err)
		}
		return oauth1Config[provider].AuthorizationURL(requestToken)
	}
	return nil, fmt.Errorf("an unknown provider")
}

//GetAuthedUserID Get the Authenticated User ID
func GetAuthedUserID(provider string, query *url.Values) (string, error) {
	switch authType[provider] {
	case "OAuth1.0a":
		requestToken, verifier := query.Get("oauth_token"), query.Get("oauth_verifier")
		accessToken, accessTokenSecret, err := oauth1Config[provider].AccessToken(requestToken, "", verifier)
		if err != nil {
			return "", fmt.Errorf("failed to get an access token: %v", err)
		}
		httpClient := oauth1Config[provider].Client(oauth1.NoContext, oauth1.NewToken(accessToken, accessTokenSecret))
		client := resty.New().SetTransport(httpClient.Transport)
		switch provider {
		case "twitter":
			data := &struct {
				IDStr string `json:"id_str"`
			}{}
			if _, err := client.R().SetResult(data).Get("https://api.twitter.com/1.1/account/verify_credentials.json"); err != nil {
				return "", err
			}
			if data.IDStr == "" {
				return "", fmt.Errorf("failed for unknown reason")
			}
			return data.IDStr, nil
		}
	}
	return "", fmt.Errorf("an unknown provider")
}
