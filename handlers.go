package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"strings"
)

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

// Match content of the message to do whatever
func handleMessage(m *tb.Message, b *tb.Bot) {
	switch {
	case strings.Contains(strings.ToLower(m.Text), "cat"):
		b.Send(m.Sender, GetCatGif())
	case strings.Contains(m.Text, "online"):
		b.Send(m.Sender, "bot online")
	}
}
