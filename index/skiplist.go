package index

import (
	"github.com/huandu/skiplist"
	"github.com/rosedblabs/wal"
	"sync"
)

type skiplistIndex struct {
	skiplist *skiplist.SkipList
	mu       sync.Mutex
}

func NewSkiplistIndex() *skiplistIndex {
	return &skiplistIndex{
		skiplist: skiplist.New(skiplist.String),
	}
}

func (sl *skiplistIndex) Put(key string, position *wal.ChunkPosition) bool {
	defer sl.mu.Unlock()
	sl.mu.Lock()
	sl.skiplist.Set(key, position)
	return true
}
func (sl *skiplistIndex) Get(key string) *wal.ChunkPosition {
	result := sl.skiplist.Get(key)
	if result == nil {
		return nil
	}
	position := result.Value.(*wal.ChunkPosition)
	return position
}
func (sl *skiplistIndex) Delete(key string) (*wal.ChunkPosition, bool) {
	defer sl.mu.Unlock()
	sl.mu.Lock()
	elem := sl.skiplist.Remove(key)
	if elem != nil {
		position := elem.Value.(*wal.ChunkPosition)
		return position, true
	} else {
		return nil, false
	}
}
