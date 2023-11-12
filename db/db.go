package db

import (
	"bytes"
	"encoding/gob"
	"github.com/lim-yoona/tcpack"
	"github.com/lim-yoona/ymdb/index"
	"github.com/lim-yoona/ymdb/options"
	"github.com/rosedblabs/wal"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

type DB struct {
	index        index.Index
	wal          *wal.WAL
	restoreWAL   *wal.WAL
	closed       bool
	options      options.Option
	RestoreQueue chan *tcpack.Message
	isRestore    bool
}

func NewDB() *DB {
	// 处理，如果已经存在WAL文件，则直接加载以初始化db
	resFileName := "\\000000001.SEG"
	restoreWALFilePath := options.DefalutOption.RestoreWALOption.DirPath + resFileName
	_, err2 := os.Stat(restoreWALFilePath)
	var restore bool
	restore = false
	// 如果文件已存在
	if !os.IsNotExist(err2) {
		restore = true
		log.Info().Msgf("[DB] >>> Find the WAL file is exist, restore from %s", restoreWALFilePath)
		// 删除用于存储的WAL
		walFilePath := options.DefalutOption.WalOption.DirPath + resFileName
		os.Remove(walFilePath)
	}
	// 创建或者打开wal文件
	walDefault, err := wal.Open(options.DefalutOption.WalOption)
	log.Info().Msgf("[DB] >>> storeWAL path is %s", options.DefalutOption.WalOption.DirPath)
	if err != nil {
		log.Panic().Err(err).Msg("[DB] >>> Create storeWAL failed")
	}
	restoreWALDefault, err := wal.Open(options.DefalutOption.RestoreWALOption)
	log.Info().Msgf("[DB] >>> restoreWAL path is %s", options.DefalutOption.RestoreWALOption.DirPath)
	if err != nil {
		log.Panic().Err(err).Msg("[DB] >>> Create restoreWAL failed")
	}
	return &DB{
		index:        index.NewSkiplistIndex(),
		wal:          walDefault,
		restoreWAL:   restoreWALDefault,
		closed:       false,
		options:      options.DefalutOption,
		RestoreQueue: make(chan *tcpack.Message, 100),
		isRestore:    restore,
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

func (db *DB) WriteRestoreWAL(data tcpack.Imessage) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encoder.Encode(data)
	_, err := db.restoreWAL.Write(buf.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("[DB] >>> Write restoreWAL failed")
	}
	return err
}

func (db *DB) Restore() {
	if db.isRestore {
		log.Info().Msg("[DB] >>> ymDB restore...")
		reader := db.restoreWAL.NewReader()
		for {
			val, _, err := reader.Next()
			if err == io.EOF {
				break
			}
			var entry tcpack.Message
			decoder := gob.NewDecoder(bytes.NewBuffer(val))
			_ = decoder.Decode(&entry)
			db.RestoreQueue <- &entry
		}
	}
	close(db.RestoreQueue)
}
