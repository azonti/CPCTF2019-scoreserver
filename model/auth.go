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
		CallbackURL:    os.Getenv("DEPROY_URL") + os.Getenv("API_URL_PREFIX") + "/auth/twitter/callback",
		Endpoint:       twitter.AuthorizeEndpoint,
	},
}

//ErrUnknownProvider an Error due to an Unknown Provider
var ErrUnknownProvider = fmt.Errorf("an unknown provider")

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
	return nil, ErrUnknownProvider
}

//GetAuthedUserID Get the Authenticated User's ID
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
				return "", fmt.Errorf("failed to get the user's information: %v", err)
			}
			if data.IDStr == "" {
				return "", fmt.Errorf("failed for unknown reason")
			}
			return provider + "_" + data.IDStr, nil
		}
	}
	return "", ErrUnknownProvider
}
