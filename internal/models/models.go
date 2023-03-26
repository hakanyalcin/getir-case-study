package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBModel struct {
	DB *mongo.Client
}

type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *mongo.Client) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Record is the type for all records
type Record struct {
	Key       string    `json:"key" bson:"key"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Counts    int       `json:"counts" bson:"totalCount"`
}
type RecordPayload struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}
type RecordResponse struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []Record `json:"records"`
}
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

	col := m.DB.Database("getircase-study").Collection("records")
	cursor, err := col.Aggregate(context.TODO(), pipe)
	defer cursor.Close(context.TODO())
	if err != nil {
		return records, err
	}
	if err = cursor.All(context.TODO(), &records); err != nil {
		return records, err
	}

	return records, nil

}
