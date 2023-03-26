package driver

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestConnectDB(t *testing.T) {
	dsn := "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getircase-study?retryWrites=true"
	client, err := ConnectDB(dsn)
	if err != nil {
		t.Errorf("Error connecting to database: %v", err)
	}
	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		t.Errorf("Error pinging database: %v", err)
	}
}
