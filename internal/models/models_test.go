package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetRecords(t *testing.T) {
	// initialize the test data
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	db := mongoClient.Database("getircase-study")
	col := db.Collection("records")

	payload := RecordPayload{
		StartDate: "2022-01-01",
		EndDate:   "2022-02-01",
		MinCount:  10,
		MaxCount:  20,
	}

	// insert sample records into the collection
	records := []interface{}{
		bson.M{
			"key":       "test-key-1",
			"createdAt": time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC),
			"counts":    []int{5, 15},
		},
		bson.M{
			"key":       "test-key-2",
			"createdAt": time.Date(2022, 1, 10, 0, 0, 0, 0, time.UTC),
			"counts":    []int{10, 20},
		},
		bson.M{
			"key":       "test-key-3",
			"createdAt": time.Date(2022, 1, 20, 0, 0, 0, 0, time.UTC),
			"counts":    []int{5, 25},
		},
	}
	_, err = col.InsertMany(context.Background(), records)
	require.NoError(t, err)

	// initialize the DBModel and call GetRecords
	model := &DBModel{DB: mongoClient}
	result, err := model.GetRecords(payload)
	require.NoError(t, err)

	// check if the result is as expected
	require.Equal(t, []Record{
		{Key: "test-key-1", CreatedAt: time.Date(2022, 1, 5, 0, 0, 0, 0, time.UTC)},
		{Key: "test-key-2", CreatedAt: time.Date(2022, 1, 10, 0, 0, 0, 0, time.UTC)},
	}, result)
}
