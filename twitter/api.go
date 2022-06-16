package twitter

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

//TwitterApi is a facade for 3rd party twitter api client
type TwitterApi struct {
	client     *twitter.Client
	config     *TwitterConfig
	streams    []*twitter.Stream
	streamLock sync.Mutex
}

func NewTwitterApi(cfg *TwitterConfig) *TwitterApi {
	config := oauth1.NewConfig(cfg.ConsumerKey, cfg.ConsumerSecret)
	token := oauth1.NewToken(cfg.AccessToken, cfg.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return &TwitterApi{
		config:  cfg,
		client:  twitter.NewClient(httpClient),
		streams: make([]*twitter.Stream, 0),
	}
}

//ListenToMentions starts listening to a twitter mentions of given username.
//It returns a channel that can be listened on.
func (t *TwitterApi) ListenToMentions() (<-chan interface{}, error) {
	_, _, err := t.client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("credentials invalid -- %v", err)
	}

	t.streamLock.Lock()
	defer t.streamLock.Unlock()
	params := &twitter.StreamFilterParams{
		Track:         []string{fmt.Sprintf("@%s", t.config.Username)},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := t.client.Streams.Filter(params)
	if err != nil {
		return nil, err
	}

	t.streams = append(t.streams, stream)

	return stream.Messages, nil
}

//Tweet tweets a response at a user
func (t *TwitterApi) Tweet(text string, inReplyTo int64, images [][]byte) error {
	var mediaIds []int64
	var err error
	if len(images) > 0 {
		mediaIds, err = t.uploadImagesToTwitter(images)
		if err != nil {
			return err
		}
	}
	_, _, err = t.client.Statuses.Update(text, &twitter.StatusUpdateParams{
		MediaIds:          mediaIds,
		InReplyToStatusID: inReplyTo,
	})

	return err
}

func (t *TwitterApi) uploadImagesToTwitter(images [][]byte) ([]int64, error) {
	var mediaIds []int64
	for _, img := range images {
		res, _, err := t.client.Media.Upload(img, "png")
		if err != nil || res.MediaID <= 0 {
			return nil, errors.New("failed to upload images to twitter")
		}

		mediaIds = append(mediaIds, res.MediaID)
	}

	return mediaIds, nil
}

func (t *TwitterApi) Stop() {
	t.streamLock.Lock()
	defer t.streamLock.Unlock()
	for _, s := range t.streams {
		s.Stop()
	}
}
