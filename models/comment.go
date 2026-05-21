package models

import (
	"time"

	astratypes "github.com/datastax/astra-db-go/datatypes"
)

type Comment struct {
	Commentid      astratypes.UUID `json:"commentid"`
	Videoid        astratypes.UUID `json:"videoid"`
	Userid         astratypes.UUID `json:"userid"`
	UserName       string          `json:"user_name"`
	CommentText    string          `json:"comment"`
	SentimentScore float32         `json:"sentiment_score"`
	Timestamp      time.Time       `json:"timestamp"`
}

func NewComment() *Comment {
	return &Comment{SentimentScore: 0.0}
}
