package services

import "testing"

var settingsPath = "settings_example.toml"
var destroyAfterRead = false

func TestSetupSettings(t *testing.T) {
	var p = GetProvider(SettingsContext{
		ConfPath:         &settingsPath,
		DestroyAfterRead: &destroyAfterRead,
	})

	t.Logf("Read: %+v", p)

	if p.CreatorID != 0 {
		t.Error("Could not read Creator ID")
	}

	if p.URIs.Mongo != "localhost:27017" {
		t.Error("Could not read URIs subkey")
	}

	if p.APIKeys.TelegramBot != "0" || p.APIKeys.GoogleAPI != "1" || p.APIKeys.GiphyAPI != "2" {
		t.Error("Could not read APIKeys subkey")
	}
}
