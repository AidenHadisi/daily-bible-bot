package bible

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//go:embed topics.json
var Topics []byte

//BibleVerseGetter gets bible verses.
type BibleVerseGetter interface {
	//GetVerse gets the given verse.
	GetVerse(verse string) (string, error)

	//GetRandomVerse gets a random verse from api.
	GetRandomVerse() (string, error)
}

type bibleApiResponse struct {
	Text string `json:"text"`
}

type BibleApi struct {
	client *http.Client
	cache  Cache
	topics []string
	verses map[string][]string
}

func NewBibleApi(client *http.Client, cache Cache) *BibleApi {
	return &BibleApi{
		client: client,
		cache:  cache,
		topics: make([]string, 0),
		verses: make(map[string][]string),
	}
}

func (b *BibleApi) Init() error {
	err := json.Unmarshal(Topics, &b.verses)
	if err != nil {
		return err
	}

	for k := range b.verses {
		b.topics = append(b.topics, k)
	}

	return nil
}

//GetRandomVerse gets a random verse from API.
func (b *BibleApi) GetRandomVerse() (string, error) {
	randomTopic := b.topics[rand.Intn(len(b.topics))]
	randomVerse := b.verses[randomTopic][rand.Intn(len(b.verses[randomTopic]))]
	resp, err := b.GetVerse(randomVerse)
	if err != nil {
		return "", err
	}
	return resp, nil
}

//GetVerse gets the requested verses from the API.
func (b *BibleApi) GetVerse(verse string) (string, error) {
	cached, err := b.cache.Get(verse)
	if err == nil {
		return cached, nil
	}
	url := fmt.Sprintf("https://bible-api.com/%s?translation=kjv", verse)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api returned code %d", resp.StatusCode)
	}

	bibleApiResult := &bibleApiResponse{}
	err = json.NewDecoder(resp.Body).Decode(bibleApiResult)
	if err != nil {
		return "", err
	}
	text := fmt.Sprintf("\"%s\" - %s", strings.ReplaceAll(bibleApiResult.Text, "\n", " "), verse)

	b.cache.Set(verse, text, time.Hour*2)
	return text, nil
}
