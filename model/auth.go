package model

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"net/url"
	"os"
)

//AuthType Authentication Type of each Provider
var AuthType = map[string]string{
	"twitter": "OAuth1.0a",
}
var config1 = map[string]*oauth1.Config{
	"twitter": {
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		CallbackURL:    os.Getenv("TWITTER_CALLBACK_URL"),
		Endpoint:       twitter.AuthorizeEndpoint,
	},
}

//GetAuthURL Get Authorization URL
func GetAuthURL(provider string) (*url.URL, error) {
	switch AuthType[provider] {
	case "OAuth1.0a":
		requestToken, _, err := config1[provider].RequestToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get request token: %v", err)
		}
		return config1[provider].AuthorizationURL(requestToken)
	}
	return nil, fmt.Errorf("unknown provider")
}

//Login Get Data and Generate Token
func Login(provider string, query url.Values) (string, error) {
	switch AuthType[provider] {
	case "OAuth1.0a":
		accessToken, _, err := config1[provider].AccessToken(query.Get("oauth_token"), "", query.Get("oauth_verifier"))
		if err != nil {
			return "", fmt.Errorf("failed to get access token: %v", err)
		}
		return accessToken, nil
	}
	return "", fmt.Errorf("unknown provider")
}
