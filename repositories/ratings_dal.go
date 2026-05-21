package repositories

import (
	"context"
	"killrvideo/go-backend-astra-dataapi/models"
	"strconv"

	astradb "github.com/datastax/astra-db-go"
	astratypes "github.com/datastax/astra-db-go/datatypes"
	"github.com/datastax/astra-db-go/filter"
)

type RatingsDAL struct {
	DB  *astradb.Db
	Ctx context.Context
}

func NewRatingsDAL(db *astradb.Db, ctx context.Context) *RatingsDAL {
	return &RatingsDAL{
		DB:  db,
		Ctx: ctx,
	}
}

func (r *RatingsDAL) GetSingleRating(videoid astratypes.UUID) (*models.Rating, error) {
	collection := r.DB.Collection("video_ratings")

	cursor := collection.Find(filter.Eq("videoid", videoid))
	defer cursor.Close()

	var returnRating models.Rating
	ratingCount := 0
	ratingTotal := 0

	for cursor.Next(r.Ctx) {
		var rating models.Rating
		if err := cursor.Decode(&rating); err != nil {
			return nil, err
		}

		ratingCount++
		ratingTotal += parseRating(rating.Rating)
	}

	if err2 := cursor.Err(); err2 != nil {
		return nil, err2
	}

	return &returnRating, nil
}

func parseRating(rating string) int {

	ratingInt, err := strconv.Atoi(rating)
	if err != nil {
		return 0
	}
	return ratingInt
}
