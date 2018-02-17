package services

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"os"
)

var prov *Provider
var created = false

// SettingsContext contains the path to a config file, and whether to destroy it or not after loading
type SettingsContext struct {
	ConfPath         *string
	DestroyAfterRead *bool
}

// Provider gives out config parameters
type Provider struct {
	CreatorID int
	URIs      urisProvider
	APIKeys   keysProvider
}

type keysProvider struct {
	TelegramBot   string
	GoogleAPI     string
	GiphyAPI      string
	DialogflowAPI string
}

type urisProvider struct {
	Mongo string
}

// GetProvider takes the path to the .toml file, and a boolean
// that indicates whether the config file should be deleted after being read
// Returns a pointer to a Provider
func GetProvider(ctx SettingsContext) *Provider {
	if created == true {
		return prov
	}
	var keysLocation = *ctx.ConfPath
	if _, err := toml.DecodeFile(keysLocation, &prov); err != nil {
		log.Fatalf("Could not read config file %s: %s", keysLocation, err)
	}
	log.Infof("Loaded settings file: %s", *ctx.ConfPath)
	if *ctx.DestroyAfterRead {
		log.Infof("Destroying %s", *ctx.ConfPath)
		os.Remove(*ctx.ConfPath)
	}
	created = true // TODO: Better, please (nil check; flags suck)
	return prov
}
