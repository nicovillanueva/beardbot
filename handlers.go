package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"strings"
)

func registerHandlers(bot *tb.Bot) {
	bot.Handle("/ping", func(m *tb.Message) {
		bot.Reply(m, "Bot standing by")
	})

	bot.Handle(tb.OnLocation, func(m *tb.Message) {
		lat, long := fmt.Sprint(m.Location.Lat), fmt.Sprint(m.Location.Lng)
		log.Infof("Got location: %s:%s\n", lat, long)
		addrs := FindLocation(m.Location.Lat, m.Location.Lng)
		log.Infof("Found: %v", addrs)
		bot.Send(m.Sender, fmt.Sprintf("I found this: %v", addrs))
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		fileID := m.Photo.File.FileID
		target := "/tmp/" + fileID
		log.Infof("Got image ID: %s\n", fileID)
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
    text := strings.ToLower(m.Text)
	switch {
	case strings.Contains(text, "cat"):
		v, err := GetGif("cat")
		if err != nil {
			b.Send(m.Sender, "hey, sorry, but I couldn't get you any cats... :(")
		} else {
			b.Send(m.Sender, v)
		}
    case strings.Contains(text, "dog"):
        v, err := GetGif("dog")
		if err != nil {
			b.Send(m.Sender, "hey, sorry, but I couldn't get you any dogs... you can always try with cats though ;)")
		} else {
			b.Send(m.Sender, v)
		}
	case strings.Contains(text, "online"):
		b.Send(m.Sender, "bot online")
	}
}
