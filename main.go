package main

import (
	"github.com/BurntSushi/toml"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
    "flag"
)

type Settings struct {
	ApiKeys   KeysProvider
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

	log.Println("Starting...")
	setupSettings(settPtr, destrPtr)

	b := assembleBot()

	registerHandlers(b)

	go b.Start()
    log.Println("Systems ready. Bot online.")
	notifyStartup(b)

	log.Println("Preparing teardown.")
	prepareShutdown(b)
}

func assembleBot() *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token:  SettingsProvider.ApiKeys.TelegramBot,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panicln(err)
	}
    log.Printf("Bot created: %s.", SettingsProvider.ApiKeys.TelegramBot)
	return b
}

func prepareShutdown(b *tb.Bot) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Stopping bot...")
	notifyShutdown(b)
	b.Stop()
	log.Println("Done.")
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
		log.Panicf("Could not read config file %s: %s", keysLocation, err)
	}
    log.Printf("Loaded settings file %s", *settingsPath)
    if *destroyAfterRead {
        log.Printf("Destroying %s", *settingsPath)
        os.Remove(*settingsPath)
    }
}
