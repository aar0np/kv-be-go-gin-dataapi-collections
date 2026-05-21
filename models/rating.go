package models

import (
	"time"

	astratypes "github.com/datastax/astra-db-go/datatypes"
)

type Rating struct {
	Videoid    astratypes.UUID `json:"videoid"`
	Userid     astratypes.UUID `json:"userid"`
	Rating     string          `json:"rating"`
	RatingDate time.Time       `json:"rating_date"`
	Score      float32         `json:"score"`
}
