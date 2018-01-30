package main

import (
	"encoding/json"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var initialized bool = false
var captions = []string{
	"Here you have",
	"Hope your day gets better",
	"Enjoy",
	"Cats make the world a better place",
	"Cat-powered internet ftw",
}

func GetCatGif() *tb.Video {
	if !initialized {
		rand.Seed(time.Now().Unix())
		log.Println("Initialized random seed")
	}

	ph := makeRequest()
	ph.Caption = captions[rand.Intn(len(captions))]

	return &ph
}

func makeRequest() tb.Video {
	// Based on https://github.com/paddycarey/gophy
	var endpoint string = "https://api.giphy.com/v1/gifs/random?"
	qs := &url.Values{}
	qs.Set("api_key", SettingsProvider.ApiKeys.GiphyAPI)
	qs.Set("tag", "cats")
	endpoint += qs.Encode()

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalf("Could not query Giphy: %s", err)
	}

	var response GiphyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("Error decoding Giphy response: %s", err)
	}

	log.Printf("Found cat gif: %s", response.Data.ImageMp4URL)
	return tb.Video{
		File: tb.File{
			FileURL: response.Data.ImageMp4URL,
		},
	}
}

type GiphyResponse struct {
	Data struct {
		Caption                      string `json:"caption"`
		FixedHeightDownsampledHeight string `json:"fixed_height_downsampled_height"`
		FixedHeightDownsampledURL    string `json:"fixed_height_downsampled_url"`
		FixedHeightDownsampledWidth  string `json:"fixed_height_downsampled_width"`
		FixedHeightSmallHeight       string `json:"fixed_height_small_height"`
		FixedHeightSmallStillURL     string `json:"fixed_height_small_still_url"`
		FixedHeightSmallURL          string `json:"fixed_height_small_url"`
		FixedHeightSmallWidth        string `json:"fixed_height_small_width"`
		FixedWidthDownsampledHeight  string `json:"fixed_width_downsampled_height"`
		FixedWidthDownsampledURL     string `json:"fixed_width_downsampled_url"`
		FixedWidthDownsampledWidth   string `json:"fixed_width_downsampled_width"`
		FixedWidthSmallHeight        string `json:"fixed_width_small_height"`
		FixedWidthSmallStillURL      string `json:"fixed_width_small_still_url"`
		FixedWidthSmallURL           string `json:"fixed_width_small_url"`
		FixedWidthSmallWidth         string `json:"fixed_width_small_width"`
		ID                           string `json:"id"`
		ImageFrames                  string `json:"image_frames"`
		ImageHeight                  string `json:"image_height"`
		ImageMp4URL                  string `json:"image_mp4_url"`
		ImageOriginalURL             string `json:"image_original_url"`
		ImageURL                     string `json:"image_url"`
		ImageWidth                   string `json:"image_width"`
		Type                         string `json:"type"`
		URL                          string `json:"url"`
		Username                     string `json:"username"`
	} `json:"data"`
	Meta struct {
		Msg        string `json:"msg"`
		ResponseID string `json:"response_id"`
		Status     int    `json:"status"`
	} `json:"meta"`
}
