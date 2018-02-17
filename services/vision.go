package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi/transport"
	vision "google.golang.org/api/vision/v1"
	"net/http"
	// youtube "google.golang.org/api/youtube/v3"
)

// type VisionResult struct {
// 	Descr      string
// 	Mid        string
// 	score      float32
// 	topicality float32
// }

func RecognizeImage(ctx SettingsContext, data []byte) string {
	enc := base64.StdEncoding.EncodeToString(data)
	img := &vision.Image{Content: enc}

	feature := &vision.Feature{
		Type:       "LABEL_DETECTION",
		MaxResults: 3,
	}
	req := &vision.AnnotateImageRequest{
		Image:    img,
		Features: []*vision.Feature{feature},
	}
	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: GetProvider(ctx).APIKeys.GoogleAPI}, // TODO: APIKEY
	}

	svc, err := vision.New(client)
	if err != nil {
		log.Errorf("Could not create vision client! %s", err)
		return "Could not contact google to do the recognizing thing :("
	}
	log.Infof("Calling Google Vision")
	res, err := svc.Images.Annotate(batch).Do()
	if err != nil {
		m := "Could not recognize images!"
		log.Errorf("%s %s\n", m, err)
		return m
	}
	body, _ := json.Marshal(res.Responses[0].LabelAnnotations)
	log.Infof("GVision says: %s", string(body))

	annotations := res.Responses[0].LabelAnnotations

	var s string
	switch len(annotations) {
	case 3:
		s += fmt.Sprintf(". Though it miiiight be a '%s' (%f)", annotations[2].Description, annotations[2].Score)
		fallthrough
	case 2:
		s = fmt.Sprintf(", but it could also be a '%s' (%f)", annotations[1].Description, annotations[1].Score) + s
		fallthrough
	case 1:
		s = fmt.Sprintf("Looks a lot like a '%s' (%f)", annotations[0].Description, annotations[0].Score) + s
	}

	return s
}

// var playlists = []string{
// 	"PLXD4mnw6H4dNyauYF7A2qHB5xLOCG64Bp",
// 	"PLOy0j9AvlVZPto6IkjKfpu0Scx--7PGTC",
// }
//
// func GetRandomVideo() {
// 	client := &http.Client{
// 		Transport: &transport.APIKey{Key: SettingsProvider.Keys.GoogleAPI},
// 	}
//
// 	svc, err := youtube.New(client)
// 	p := youtube.Playlist{
// 		Id: playlists[0],
// 	}
// 	playlistSvc := youtube.NewPlaylistItemsService(svc)
// }
