package index

import "github.com/rosedblabs/wal"

type Index interface {
	Put(key string, position *wal.ChunkPosition) bool
	Get(key string) *wal.ChunkPosition
	Delete(key string) (*wal.ChunkPosition, bool)
}
