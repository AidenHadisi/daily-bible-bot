package twitter

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AidenHadisi/BibleBot/bible"
	"github.com/AidenHadisi/BibleBot/cron"
	"github.com/AidenHadisi/BibleBot/image"
	"github.com/dghubble/go-twitter/twitter"
)

//TwitterBot defines MyDailyBibleBot structure
type TwitterBot struct {
	api    *TwitterApi
	bible  *bible.BibleApi
	cron   cron.Cron
	config *TwitterConfig
}

func NewTwitterBot(api *TwitterApi, b *bible.BibleApi, cron cron.Cron, config *TwitterConfig) *TwitterBot {
	return &TwitterBot{
		api:    api,
		bible:  b,
		cron:   cron,
		config: config,
	}
}

func (b *TwitterBot) Init() error {
	//init the api client
	err := b.bible.Init()
	if err != nil {
		return err
	}

	//start listening to twitter
	c, err := b.api.ListenToMentions()
	if err != nil {
		return err
	}
	go b.handleMessages(c)

	//start the cron
	err = b.cron.CreateJob("0 */5 * * *", b.randomPost)
	if err != nil {
		return err
	}

	err = b.cron.StartCrons()
	if err != nil {
		return err
	}

	return nil
}

func (b *TwitterBot) handleMessages(messages <-chan any) {
	for message := range messages {
		if msg, ok := message.(*twitter.Tweet); ok {
			b.messageHandler(msg)
		}
	}
}

func (b *TwitterBot) messageHandler(tweet *twitter.Tweet) {
	if tweet.User.ScreenName == b.config.Username {
		return
	}

	verseRequest := bible.NewParser()
	err := verseRequest.Parse(tweet.Text)
	if err != nil {
		return
	}

	go b.reply(tweet, verseRequest)
}

func (b *TwitterBot) reply(tweet *twitter.Tweet, parsed *bible.Parser) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("reply recovered -- %v", err)
		}
	}()

	text, err := b.bible.GetVerse(parsed.GetPath())
	if err != nil {
		log.Println(err)
		return
	}

	if parsed.HasImage() {
		imageProcessor := image.NewImageProcessor(&http.Client{Timeout: time.Minute})
		by, err := imageProcessor.Process(parsed.Img, text, parsed.Size)
		if err != nil {
			message := fmt.Sprintf("@%s %s", tweet.User.ScreenName, "Sorry we weren't able to process that image.")
			b.api.Tweet(message, tweet.ID, nil)
			return
		}
		message := fmt.Sprintf("@%s %s", tweet.User.ScreenName, "")
		b.api.Tweet(message, tweet.ID, [][]byte{by})
	} else {
		message := fmt.Sprintf("@%s %s", tweet.User.ScreenName, text)
		if len([]rune(message)) > 280 {
			message = fmt.Sprintf("@%s %s", tweet.User.ScreenName, "Sorry the requested verse is too long for a single tweet.")
		}
		b.api.Tweet(message, tweet.ID, nil)
	}

}

func (b *TwitterBot) randomPost() {
	resp, err := b.bible.GetRandomVerse()
	if err != nil {
		return
	}
	imageProcessor := image.NewImageProcessor(&http.Client{Timeout: time.Minute})
	image, err := imageProcessor.Process("https://picsum.photos/1200/625", resp, 40)
	if err != nil {
		return
	}
	b.api.Tweet("", 0, [][]byte{image})

}

func (b *TwitterBot) Shutdown() {
	b.api.Stop()
	b.cron.StopCrons()
}
