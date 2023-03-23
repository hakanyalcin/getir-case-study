package cache

import (
	"errors"
	"getir-case-study/internal/models"
)

type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CachedEnty struct {
	Entry
}

type LocalCache struct {
	Entries map[string]CachedEnty
}

func NewLocalCache() *LocalCache {
	lc := &LocalCache{
		Entries: make(map[string]CachedEnty),
	}
	return lc

}

func (lc *LocalCache) SetEntry(e models.CachePayload) (Entry, error) {
	entry := Entry{
		Key:   e.Key,
		Value: e.Value,
	}
	lc.Entries[entry.Key] = CachedEnty{
		Entry: entry,
	}
	return entry, nil
}

var (
	errUserNotInCache = errors.New("cache missing: The key isn't in cache")
)

func (lc *LocalCache) GetEntry(key string) (Entry, error) {

	cu, ok := lc.Entries[key]
	if !ok {
		return Entry{}, errUserNotInCache
	}

	return cu.Entry, nil
}
