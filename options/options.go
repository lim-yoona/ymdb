package options

import "github.com/rosedblabs/wal"

type Option struct {
	WalOption wal.Options
}

var DefalutOption = Option{
	WalOption: struct {
		DirPath        string
		SegmentSize    int64
		SegmentFileExt string
		BlockCache     uint32
		Sync           bool
		BytesPerSync   uint32
	}{DirPath: "E:\\project-about\\tinyKV\\walDir", SegmentSize: wal.GB, SegmentFileExt: ".SEG", BlockCache: 32 * wal.KB * 10, Sync: false, BytesPerSync: 0},
}
