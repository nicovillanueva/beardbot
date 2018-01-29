package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
    "strings"
)

type ApiKeys struct {
	TelegramBot string
	GoogleAPI   string
	CreatorID   int
	GiphyAPI    string
}

var KeysProvider ApiKeys

func main() {
	log.Println("Starting...")
	var keysLocation string = "./keys.toml"
	if _, err := toml.DecodeFile(keysLocation, &KeysProvider); err != nil {
		log.Panicf("Could not read config file %s: %s", keysLocation, err)
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  KeysProvider.TelegramBot,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	} else {
		log.Printf("Bot created: %s. Registering handlers.", KeysProvider.TelegramBot)
	}

	registerHandlers(b)

	log.Println("Systems ready. Bot online.")
	go b.Start()
	notifyStartup(b)

	log.Println("Preparing teardown.")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	prepareShutdown(b, c)
}

func registerHandlers(bot *tb.Bot) {
	bot.Handle("/ping", func(m *tb.Message) {
		bot.Reply(m, "Bot standing by")
	})

	bot.Handle("cats!", func(m *tb.Message) {
		bot.Send(m.Sender, GetCatGif())
	})

	bot.Handle(tb.OnLocation, func(m *tb.Message) {
		lat, long := fmt.Sprint(m.Location.Lat), fmt.Sprint(m.Location.Lng)
		log.Printf("Got location: %s:%s\n", lat, long)
		addrs := FindLocation(m.Location.Lat, m.Location.Lng)
		log.Printf("Found: %v", addrs)
		bot.Send(m.Sender, fmt.Sprintf("I found this: %v", addrs))
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		fileID := m.Photo.File.FileID
		target := "/tmp/" + fileID
		log.Printf("Got image ID: %s\n", fileID)
		bot.Download(m.Photo.MediaFile(), target)
		data, err := ioutil.ReadFile(target)
		if err != nil {
			log.Fatalf("Error loading file! %s", err)
		}
		bot.Send(m.Sender, RecognizeImage(data))
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
        handleMessage(m, bot)
    })
}

func prepareShutdown(b *tb.Bot, c chan os.Signal) {
	<-c
	log.Println("Stopping bot...")
	notifyShutdown(b)
	b.Stop()
	log.Println("Done.")
	os.Exit(0)
}

func notifyStartup(b *tb.Bot) {
	if KeysProvider.CreatorID != 0 {
		b.Send(&tb.User{ID: KeysProvider.CreatorID}, "Bot online")
	}
}

func notifyShutdown(b *tb.Bot) {
	if KeysProvider.CreatorID != 0 {
		b.Send(&tb.User{ID: KeysProvider.CreatorID}, "Bot offline")
	}
}

// Match content of the message to do whatever
func handleMessage(m *tb.Message, b *tb.Bot) {
    switch {
    case strings.Contains(m.Text, "cat"):
        b.Send(m.Sender, GetCatGif())
    case strings.Contains(m.Text, "online"):
        b.Send(m.Sender, "bot online")
    }
}
