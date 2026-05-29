package models

import (
	"time"
)

type Comment struct {
	Commentid      string    `json:"commentid"`
	Videoid        string    `json:"videoid"`
	Userid         string    `json:"userid"`
	UserName       string    `json:"user_name"`
	CommentText    string    `json:"comment"`
	SentimentScore float32   `json:"sentiment_score"`
	Timestamp      time.Time `json:"timestamp"`
}

func NewComment() *Comment {
	return &Comment{SentimentScore: 0.0}
}
