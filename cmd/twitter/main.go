package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AidenHadisi/BibleBot/bible"
	"github.com/AidenHadisi/BibleBot/cron"
	"github.com/AidenHadisi/BibleBot/twitter"
)

func main() {
	fmt.Println("hello")
	log.Println("Starting twitter bot ...")
	log.Println("twitter bot has start.")

	cfg := twitter.NewTwitterConfig()
	cfg.Load()

	err := cfg.Validate()
	if err != nil {
		log.Panicf("Invalid twitter configs -- %s", err)
	}

	api := twitter.NewTwitterApi(cfg)

	bible := bible.NewBibleApi(&http.Client{Timeout: time.Minute}, bible.NewMemoryCache())

	bot := twitter.NewTwitterBot(api, bible, cron.NewSimpleCron(), cfg)

	err = bot.Init()
	if err != nil {
		log.Panicf("failed to start twitter bot -- %s", err)
	}

	log.Println("twitter bot has start.")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	bot.Shutdown()
}
