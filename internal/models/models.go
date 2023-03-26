package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBModel is the type for database connection values
type DBModel struct {
	DB *mongo.Database
}

// Record is the type for record on the db
type Record struct {
	Key       string    `json:"key" bson:"key"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Counts    int       `json:"counts" bson:"totalCount"`
}

// RecordPayload is the type for request body
type RecordPayload struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

// RecordResponse is the type for response body
type RecordResponse struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []Record `json:"records"`
}

// CachePayload is the type for request body
type CachePayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *DBModel) GetRecords(payload RecordPayload) ([]Record, error) {
	records := []Record{}

	pipe := []bson.M{
		{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gt": payload.StartDate,
					"$lt": payload.EndDate,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"key":        1,
				"createdAt":  1,
				"totalCount": bson.M{"$sum": "$counts"},
			},
		},
		{
			"$match": bson.M{
				"totalCount": bson.M{
					"$gt": payload.MinCount,
					"$lt": payload.MaxCount,
				},
			},
		},
	}

	col := m.DB.Collection("records")
	cursor, err := col.Aggregate(context.TODO(), pipe)

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(cursor, context.TODO())
	if err != nil {
		return records, err
	}
	if err = cursor.All(context.TODO(), &records); err != nil {
		return records, err
	}

	return records, nil
}
