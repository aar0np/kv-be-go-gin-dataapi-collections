package repositories

import (
	"context"
	"killrvideo/go-backend-astra-dataapi/models"

	"github.com/datastax/astra-db-go/astra"
	"github.com/datastax/astra-db-go/astra/filter"
	"github.com/datastax/astra-db-go/astra/options"
)

type CommentDAL struct {
	DB  *astra.Db
	Ctx context.Context
}

func NewCommentDAL(db *astra.Db, ctx context.Context) *CommentDAL {
	return &CommentDAL{
		DB:  db,
		Ctx: ctx,
	}
}

func (c *CommentDAL) GetCommentsByVideoId(videoid string, pageSize int) (*[]models.Comment, error) {
	collection := c.DB.Collection("comments")

	cursor := collection.Find(filter.Eq("videoid", videoid), options.CollectionFind().SetLimit(pageSize))
	defer cursor.Close()

	var comments []models.Comment
	if err1 := cursor.DecodeAll(c.Ctx, &comments); err1 != nil {
		return nil, err1
	}
	return &comments, nil
}

func (c *CommentDAL) SaveComment(comment models.Comment) {
	collection := c.DB.Collection("comments")

	collection.InsertOne(c.Ctx, comment)
}
