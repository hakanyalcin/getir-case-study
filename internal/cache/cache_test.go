package cache

import (
	"getir-case-study/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetEntry(t *testing.T) {
	data := models.CachePayload{
		Key:   "getir-key",
		Value: "getir-value",
	}

	cache := NewLocalCache()
	_, err := cache.SetEntry(data)
	if err != nil {
		return
	}

	require.Equal(t, data.Value, cache.Entries["getir-key"].Entry.Value)
	require.Equal(t, data.Key, cache.Entries["getir-key"].Entry.Key)
}

func TestGetEntry(t *testing.T) {
	data := models.CachePayload{
		Key:   "getir-key",
		Value: "getir-value",
	}
	cache := NewLocalCache()
	_, err := cache.SetEntry(data)
	if err != nil {
		return
	}
	actual, _ := cache.GetEntry("getir-key")
	require.Equal(t, data.Key, actual.Key)
	require.Equal(t, data.Value, actual.Value)

}
