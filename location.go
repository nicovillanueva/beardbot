package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func buildClient() (*maps.Client, error) {
	c, err := maps.NewClient(maps.WithAPIKey(SettingsProvider.APIKeys.GoogleAPI))
	if err != nil {
		log.Errorf("fatal error: %s", err)
        return nil, err
	}
	return c, nil
}

func FindLocation(lat, long float32) []string {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: float64(lat),
			Lng: float64(long),
		},
	}
	cli, err := buildClient()
    if err != nil {
        return []string{"Could not query google maps :("}
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
