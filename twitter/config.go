package twitter

import (
	"errors"
	"os"
)

type TwitterConfig struct {
	Username       string
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func NewTwitterConfig() *TwitterConfig {
	return &TwitterConfig{}
}

//Load loads config from env.
func (t *TwitterConfig) Load() {
	t.Username = os.Getenv("TWITTER_USERNAME")
	t.ConsumerKey = os.Getenv("CONSUMER_KEY")
	t.ConsumerSecret = os.Getenv("CONSUMER_SECRET")
	t.AccessToken = os.Getenv("ACCESS_TOKEN")
	t.AccessSecret = os.Getenv("ACCESS_SECRET")

}

//Validate validates the configs are set.
func (t *TwitterConfig) Validate() error {
	if t.Username == "" {
		return errors.New("TWITTER_USERNAME username is required")
	}

	if t.ConsumerKey == "" {
		return errors.New("CONSUMER_KEY username is required")
	}

	if t.ConsumerSecret == "" {
		return errors.New("CONSUMER_SECRET username is required")
	}

	if t.AccessToken == "" {
		return errors.New("ACCESS_TOKEN username is required")
	}

	if t.AccessSecret == "" {
		return errors.New("ACCESS_SECRET username is required")
	}

	return nil
}
