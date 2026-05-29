package repositories

import (
	"context"
	"fmt"
	"killrvideo/go-backend-astra-dataapi/models"
	"time"

	astradb "github.com/datastax/astra-db-go/astra"
	"github.com/datastax/astra-db-go/astra/filter"
	"github.com/datastax/astra-db-go/astra/options"
	"github.com/datastax/astra-db-go/astra/sort"
	"github.com/datastax/astra-db-go/astra/update"
)

type VideoDAL struct {
	DB  *astradb.Db
	Ctx context.Context
}

func NewVideoDAL(db *astradb.Db, ctx context.Context) *VideoDAL {
	return &VideoDAL{
		DB:  db,
		Ctx: ctx,
	}
}

func (r *VideoDAL) GetVideo(id string) (*models.Video, error) {
	collection := r.DB.Collection("videos")

	video := &models.Video{}

	err1 := collection.FindOne(r.Ctx, filter.Eq("videoid", id)).Decode(&video)
	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	if video == nil {
		return nil, fmt.Errorf("video not found")
	}

	return video, nil
}

func (r *VideoDAL) SaveVideo(video models.Video) {
	collection := r.DB.Collection("videos")

	collection.InsertOne(r.Ctx, video)
}

func (r *VideoDAL) UpdateYoutubeId(videoid string, youtubeId string) {
	collection := r.DB.Collection("videos")

	collection.UpdateOne(r.Ctx, filter.Eq("videoid", videoid),
		update.Coll().Set("youtube_id", youtubeId))
}

func (r *VideoDAL) UpdateVideoView(videoid string, views int) {
	collection := r.DB.Collection("videos")
	collection.UpdateOne(r.Ctx, filter.Eq("videoid", videoid),
		update.Coll().Set("views", views))
}

func (r *VideoDAL) GetLatestVideosToday(day time.Time, limit int) (*[]models.Video, error) {
	collection := r.DB.Collection("videos")

	cursor := collection.Find(filter.Gte("added_date", day),
		options.CollectionFind().
			SetSort(sort.Desc("added_date")).
			SetLimit(limit))
	defer cursor.Close()

	var videos []models.Video

	for cursor.Next(r.Ctx) {
		var video models.Video
		if err1 := cursor.Decode(&video); err1 == nil {
			// if a video errors out while decoding, just skip it
			videos = append(videos, video)
		}
	}

	return &videos, nil
}

func (r *VideoDAL) GetLatestVideos(limit int) (*[]models.Video, error) {
	collection := r.DB.Collection("videos")

	cursor := collection.Find(filter.Filter{},
		options.CollectionFind().
			SetSort(sort.Desc("added_date")).
			SetLimit(limit))
	defer cursor.Close()

	var videos []models.Video
	//var videos json.RawMessage

	for cursor.Next(r.Ctx) {
		var video models.Video
		if err1 := cursor.Decode(&video); err1 == nil {
			// if a video errors out while decoding, just skip it
			videos = append(videos, video)
		}
	}

	return &videos, nil
}

func (r *VideoDAL) GetVideosByVector(vector [384]float32, limit int) (*[]models.Video, error) {
	collection := r.DB.Collection("videos")

	cursor := collection.Find(filter.F{},
		options.CollectionFind().
			SetSort(sort.Vector(vector)).
			SetLimit(limit))
	defer cursor.Close()

	var videos []models.Video

	for cursor.Next(r.Ctx) {
		var video models.Video
		if err1 := cursor.Decode(&video); err1 == nil {
			// if a video errors out while decoding, just skip it
			videos = append(videos, video)
		}
	}

	return &videos, nil
}
