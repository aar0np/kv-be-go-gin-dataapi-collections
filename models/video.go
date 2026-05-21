package models

import (
	"time"

	astratypes "github.com/datastax/astra-db-go/datatypes"
)

type Video struct {
	Videoid              astratypes.UUID `json:"videoId"`
	Userid               astratypes.UUID `json:"userId"`
	Name                 string          `json:"title"`
	Description          string          `json:"description"`
	Location             string          `json:"location"`
	PreviewImageLocation string          `json:"thumbnailUrl"`
	ContentFeatures      [384]float32    `json:"contentFeatures"`
	AddedDate            time.Time       `json:"submittedAt"`
	Views                int             `json:"views"`
	Score                float32         `json:"averageRating"`
	YouTubeId            string          `json:"youtubeVideoId"`
	Tags                 []string        `json:"tags"`
	ContentRating        string          `json:"content_rating"`
	Language             string          `json:"lanugage"`
	Category             string          `json:"category"`
	LocationType         int             `json:"location_type"`
}

func NewVideo() *Video {
	return &Video{Views: 0}
}
