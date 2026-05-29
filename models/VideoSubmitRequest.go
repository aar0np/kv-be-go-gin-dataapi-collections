package models

import astratypes "github.com/datastax/astra-db-go/astra/datatypes"

type VideoSubmitRequest struct {
	YouTubeUrl  string          `json:"youtubeUrl"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	UserId      astratypes.UUID `json:"userid"`
}
