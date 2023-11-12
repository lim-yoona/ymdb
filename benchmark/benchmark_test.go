package benchmark

import (
	"encoding/json"
	"fmt"
	"github.com/lim-yoona/tcpack"
	"github.com/lim-yoona/ymdb/config"
	"github.com/lim-yoona/ymdb/interact/server"
	"github.com/lim-yoona/ymdb/options"
	"github.com/lim-yoona/ymdb/route"
	"github.com/lim-yoona/ymdb/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var ymDB *server.Server

func BenchmarkPutGet(b *testing.B) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.PanicLevel)
	config.GetConfig()
	os.Remove(config.DefaultConfig.Wal.Restore.Path + "\\000000001.SEG")
	os.Remove(config.DefaultConfig.Wal.Store.Path + "\\000000001.SEG")
	options.DefaultOption()
	ymDB = server.NewServer()
	go route.RouterStart(ymDB)
	go ymDB.Start()
	b.Run("put", benchmarkPut)
	b.Run("get", benchmarkGet)
}

func benchmarkPut(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testKey := fmt.Sprintf("ymDB-testKey-%09d", i)
		putMsg := &util.Put{
			Key:   testKey,
			Value: fmt.Sprintf("%d", i),
		}
		marshal, _ := json.Marshal(putMsg)
		message := tcpack.NewMessage(util.PUTID, uint32(len(marshal)), marshal)
		ymDB.MsQueue <- message
	}
}

func benchmarkGet(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	go func() {
		for {
			result := <-ymDB.ReQueue
			assert.NotNil(b, result)
		}
	}()
	for n := 0; n < b.N; n++ {
		testKey := fmt.Sprintf("ymDB-testKey-%09d", n)
		getMsg := &util.Other{
			Data: testKey,
		}
		marshal, _ := json.Marshal(getMsg)
		message := tcpack.NewMessage(util.PUTID, uint32(len(marshal)), marshal)
		ymDB.MsQueue <- message
	}
}
