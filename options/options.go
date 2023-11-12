package options

import (
	"github.com/lim-yoona/ymdb/config"
	"github.com/rosedblabs/wal"
)

type Option struct {
	WalOption        wal.Options
	RestoreWALOption wal.Options
}

var DefalutOption Option

func DefaultOption() {
	DefalutOption.WalOption = wal.Options{
		DirPath:        config.DefaultConfig.Wal.Store.Path,
		SegmentSize:    wal.GB,
		SegmentFileExt: config.DefaultConfig.Wal.SegmentFileExt,
		BlockCache:     32 * wal.KB * 10,
		Sync:           false, BytesPerSync: 0,
	}
	DefalutOption.RestoreWALOption = wal.Options{
		DirPath:        config.DefaultConfig.Wal.Restore.Path,
		SegmentSize:    wal.GB,
		SegmentFileExt: config.DefaultConfig.Wal.SegmentFileExt,
		BlockCache:     32 * wal.KB * 10,
		Sync:           false, BytesPerSync: 0,
	}
}
