package services

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var cli *maps.Client

func initClient(ctx SettingsContext) error {
	var err error
	cli, err = maps.NewClient(maps.WithAPIKey(GetProvider(ctx).APIKeys.GoogleAPI)) // TODO: APIKEY
	if err != nil {
		log.Errorf("fatal error building location client: %s", err)
		return err
	}
	return nil
}

// FindLocation Get a Telegram location, check GMaps what's there
func FindLocation(ctx SettingsContext, lat, long float32) []string {
	if cli == nil {
		err := initClient(ctx)
		if err != nil {
			return []string{"sorry, could not query location :("}
		}
	}
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: float64(lat),
			Lng: float64(long),
		},
	}
	result, err := cli.ReverseGeocode(context.Background(), req)
	if err != nil {
		log.Errorf("fatal error when rev geocoding: %s", err)
		return []string{"Could not get the geocoding, sorry"}
	}
	addresses := make([]string, len(result))
	for _, addr := range result {
		addresses = append(addresses, addr.FormattedAddress)
	}
	return addresses
}
