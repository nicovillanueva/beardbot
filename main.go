package main

import (
	services "github.com/nicovillanueva/beardbot/services"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var prov *services.Provider

func main() {
	settPtr := flag.String("settings", "./settings.toml", "Path to the TOML settings file to use")
	destrPtr := flag.Bool("destroy", false, "Destroy the settings file after loading?")
	flag.Parse()

	log.Infof("Starting...")
	var ctx = services.SettingsContext{
		ConfPath:         settPtr,
		DestroyAfterRead: destrPtr,
	}
	prov = services.GetProvider(ctx)

	b := assembleBot()

	registerHandlers(ctx, b)

	go b.Start()
	log.Infof("Systems ready. Bot online.")
	notifyStartup(b)

	log.Infof("Preparing teardown.")
	prepareShutdown(b)
}

func assembleBot() *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:  prov.APIKeys.TelegramBot,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panicln(err)
	}
	log.Infof("Bot created: %s", prov.APIKeys.TelegramBot)
	return b
}

func prepareShutdown(b *tb.Bot) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Infof("Stopping bot...")
	notifyShutdown(b)
	b.Stop()
	log.Infof("Done.")
	os.Exit(0)
}

func notifyStartup(b *tb.Bot) {
	if prov.CreatorID != 0 {
		b.Send(&tb.User{ID: prov.CreatorID}, "[status] Bot online")
	}
}

func notifyShutdown(b *tb.Bot) {
	if prov.CreatorID != 0 {
		b.Send(&tb.User{ID: prov.CreatorID}, "[status] Bot offline")
	}
}
