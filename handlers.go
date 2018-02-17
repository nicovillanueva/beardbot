package main

import (
	pers "github.com/nicovillanueva/beardbot/persistence"
	services "github.com/nicovillanueva/beardbot/services"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
    "fmt"
)

var s *pers.Storage

func registerHandlers(ctx services.SettingsContext, bot *tb.Bot) {
	bot.Handle("/ping", func(m *tb.Message) {
		bot.Reply(m, "Bot standing by")
	})

	bot.Handle("/start", func(m *tb.Message) {
		response := services.ReplyToQuery(ctx, "hi")  // Overwrite text so that it's handled as a greeting
        bot.Send(m.Sender, response)
	})

	bot.Handle(tb.OnLocation, func(m *tb.Message) {
		lat, long := fmt.Sprint(m.Location.Lat), fmt.Sprint(m.Location.Lng)
		log.Infof("Got location: %s:%s\n", lat, long)
		addrs := services.FindLocation(ctx, m.Location.Lat, m.Location.Lng)
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
		bot.Send(m.Sender, services.RecognizeImage(ctx, data))
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		response := services.ReplyToQuery(ctx, m.Text)
		bot.Send(m.Sender, response)
	})
}
