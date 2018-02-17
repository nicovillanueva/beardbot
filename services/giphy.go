package services

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var initialized = false
var captions = []string{
	"Here you have",
	"Hope your day gets better",
	"Enjoy",
	// "Cats make the world a better place",
	// "Cat-powered internet ftw",
}

// GetGif receives a tag, and queries the random giphy endpoint to get a
// gif with the given tag
func GetGif(ctx SettingsContext, tag string) (*tb.Video, error) {
	if !initialized {
		rand.Seed(time.Now().Unix())
		initialized = true
		log.Info("Initialized random seed")
	}

	ph, err := makeRequest(ctx, tag)
	if err != nil {
		return nil, err
	}

	ph.Caption = captions[rand.Intn(len(captions))]

	return &ph, nil
}

func makeRequest(ctx SettingsContext, tag string) (tb.Video, error) {
	// Based on https://github.com/paddycarey/gophy
	var endpoint = "https://api.giphy.com/v1/gifs/random?"
	qs := &url.Values{}
	qs.Set("api_key", GetProvider(ctx).APIKeys.GiphyAPI)
	qs.Set("tag", tag)
	endpoint += qs.Encode()

	resp, err := http.Get(endpoint)
	if err != nil {
		log.Errorf("Could not query Giphy: %s", err)
		return tb.Video{}, err
	}

	var response giphyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Errorf("Error decoding Giphy response: %s", err)
		return tb.Video{}, err
	}

	log.Infof("Found cat gif: %s", response.Data.ImageMp4URL)
	return tb.Video{
		File: tb.File{
			FileURL: response.Data.ImageMp4URL,
		},
	}, nil
}

type giphyResponse struct {
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
