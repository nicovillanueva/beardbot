package main

import (
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"log"
)

func buildClient() *maps.Client {
	c, err := maps.NewClient(maps.WithAPIKey(KeysProvider.GoogleAPI))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return c
}

func FindLocation(lat, long float32) []string {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: float64(lat),
			Lng: float64(long),
		},
	}
	cli := buildClient()
	result, err := cli.ReverseGeocode(context.Background(), req)
	if err != nil {
		log.Fatalf("fatal error when rev geocoding: %s", err)
	}
	addresses := make([]string, len(result))
	for _, addr := range result {
		addresses = append(addresses, addr.FormattedAddress)
	}
	return addresses
}
