package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	// "log"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Settings struct {
	APIKeys   KeysProvider
	CreatorID int
}

type KeysProvider struct {
	TelegramBot string
	GoogleAPI   string
	GiphyAPI    string
}

var SettingsProvider Settings

func main() {
	settPtr := flag.String("settings", "./settings.toml", "Path to the TOML settings file to use")
	destrPtr := flag.Bool("destroy", false, "Destroy the settings file after loading?")
	flag.Parse()

	log.Infof("Starting...")
	setupSettings(settPtr, destrPtr)

	b := assembleBot()

	registerHandlers(b)

	go b.Start()
	log.Infof("Systems ready. Bot online.")
	notifyStartup(b)

	log.Infof("Preparing teardown.")
	prepareShutdown(b)
}

func assembleBot() *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:  SettingsProvider.APIKeys.TelegramBot,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panicln(err)
	}
	log.Infof("Bot created: %s.", SettingsProvider.APIKeys.TelegramBot)
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
	if SettingsProvider.CreatorID != 0 {
		b.Send(&tb.User{ID: SettingsProvider.CreatorID}, "Bot online")
	}
}

func notifyShutdown(b *tb.Bot) {
	if SettingsProvider.CreatorID != 0 {
		b.Send(&tb.User{ID: SettingsProvider.CreatorID}, "Bot offline")
	}
}

func setupSettings(settingsPath *string, destroyAfterRead *bool) {
	var keysLocation = *settingsPath
	if _, err := toml.DecodeFile(keysLocation, &SettingsProvider); err != nil {
		log.Fatalf("Could not read config file %s: %s", keysLocation, err)
	}
	log.Infof("Loaded settings file %s", *settingsPath)
	if *destroyAfterRead {
		log.Infof("Destroying %s", *settingsPath)
		os.Remove(*settingsPath)
	}
}
