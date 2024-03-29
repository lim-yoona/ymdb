package route

import (
	"encoding/json"

	"github.com/lim-yoona/tcpack"
	"github.com/lim-yoona/ymdb/db"
	"github.com/lim-yoona/ymdb/interact/server"
	"github.com/lim-yoona/ymdb/util"
	"github.com/rs/zerolog/log"
)

type Router struct {
	server *server.Server
	db     *db.DB
}

func NewRouter(server2 *server.Server) *Router {
	return &Router{
		db:     db.NewDB(),
		server: server2,
	}
}

func (r *Router) ReadMsQueue() {
	for {
		message := <-r.server.MsQueue
		go r.Handle(message, false)
	}
}

func (r *Router) Handle(imessage tcpack.Imessage, restore bool) {
	switch imessage.GetMsgId() {
	case util.PUTID:
		if !restore {
			err := r.writeWAL(imessage)
			if err != nil {
				log.Error().Err(err).Msg("[Route] >>> Route Handle WriteWAL failed")
			}
		}
		var putMsg util.Put
		json.Unmarshal(imessage.GetMsgData(), &putMsg)
		r.db.Put(putMsg.Key, putMsg.Value)
		log.Info().Msgf("[Route] >>> Put data: %s", putMsg)
		sendMsg := &util.Other{
			Data: "PUT status = ok",
		}
		marshal, _ := json.Marshal(sendMsg)
		msgSend := tcpack.NewMessage(util.RESPONSEID, uint32(len(marshal)), marshal)
		r.server.ReQueue <- msgSend
		break
	case util.GETID:
		var getMsg util.Other
		json.Unmarshal(imessage.GetMsgData(), &getMsg)
		get := r.db.Get(getMsg.Data)
		log.Info().Msgf("[Route] >>> Get data: %s", getMsg)
		sendMsg := &util.Other{
			Data: get,
		}
		marshal, _ := json.Marshal(sendMsg)
		msgSend := tcpack.NewMessage(util.RESPONSEID, uint32(len(marshal)), marshal)
		r.server.ReQueue <- msgSend
		break
	case util.DELETEID:
		if !restore {
			err := r.writeWAL(imessage)
			if err != nil {
				log.Error().Err(err).Msg("[Route] >>> Route Handle WriteWAL failed")
			}
		}
		var deleteMsg util.Other
		json.Unmarshal(imessage.GetMsgData(), &deleteMsg)
		r.db.Delete(deleteMsg.Data)
		log.Info().Msgf("[Route] >>> Delete data: %s", deleteMsg)
		sendMsg := &util.Other{
			Data: "DELETE status = ok",
		}
		marshal, _ := json.Marshal(sendMsg)
		msgSend := tcpack.NewMessage(util.RESPONSEID, uint32(len(marshal)), marshal)
		r.server.ReQueue <- msgSend
		break
	}
}
func (r *Router) writeWAL(imessage tcpack.Imessage) error {
	err := r.db.WriteRestoreWAL(imessage)
	if err != nil {
		log.Error().Err(err).Msg("[Route] >>> route WriteWAL failed")
	}
	return err
}

func (r *Router) Restore() {
	restoreCount := 0
	go r.db.Restore()
	for {
		restoreEntry := <-r.db.RestoreQueue
		if restoreEntry == nil {
			break
		}
		restoreCount++
		go r.Handle(restoreEntry, true)
	}
	log.Info().Msgf("[DB] >>> Restore finished! Total restored %s entries.", restoreCount)
}

func RouterStart(server2 *server.Server) {
	router := NewRouter(server2)
	router.Restore()
	router.ReadMsQueue()
}
