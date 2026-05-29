package models

import (
	"time"
)

type Video struct {
	Videoid              string       `json:"videoid"`
	Userid               string       `json:"userid"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Location             string       `json:"location"`
	PreviewImageLocation string       `json:"preview_image_location"`
	ContentFeatures      [384]float32 `json:"$vector"`
	AddedDate            time.Time    `json:"added_date"`
	Views                int          `json:"views"`
	Score                float32      `json:"averageRating"`
	YouTubeId            string       `json:"youtubeVideoId"`
	Tags                 []string     `json:"tags"`
	ContentRating        string       `json:"content_rating"`
	Language             string       `json:"lanugage"`
	Category             string       `json:"category"`
	LocationType         int          `json:"location_type"`
}

func NewVideo() *Video {
	return &Video{Views: 0}
}
