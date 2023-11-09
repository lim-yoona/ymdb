package db

import (
	"bytes"
	"encoding/gob"
	"github.com/rosedblabs/wal"
	"github/lim-yoona/tinyKVStore/index"
	"github/lim-yoona/tinyKVStore/options"
)

type DB struct {
	index   index.Index
	wal     *wal.WAL
	closed  bool
	options options.Option
}

func NewDB() *DB {
	walDefault, _ := wal.Open(options.DefalutOption.WalOption)
	return &DB{
		index:   index.NewSkiplistIndex(),
		wal:     walDefault,
		closed:  false,
		options: options.DefalutOption,
	}
}
func (db *DB) Put(key string, value string) bool {
	// 序列化value
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encoder.Encode(value)
	// 首先放入WAL
	position, _ := db.wal.Write(buf.Bytes())
	// 放入index
	db.index.Put(key, position)
	return true
}
func (db *DB) Get(key string) string {
	// 在跳表中查找位置
	position := db.index.Get(key)
	if position == nil {
		return ""
	}
	// 从wal中读取
	read, _ := db.wal.Read(position)
	var result string
	decoder := gob.NewDecoder(bytes.NewBuffer(read))
	_ = decoder.Decode(&result)
	return result
}
func (db *DB) Delete(key string) (*wal.ChunkPosition, error) {
	// 从跳表中删除
	chunkPosition, _ := db.index.Delete(key)
	return chunkPosition, nil
}
