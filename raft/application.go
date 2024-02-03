package raft

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/rafterrors"
	"github.com/hashicorp/raft"
	"github.com/lim-yoona/tcpack"
	ser "github.com/lim-yoona/ymdb/interact/server"
	pb "github.com/lim-yoona/ymdb/proto"
)

type dataTracker struct {
	server *ser.Server
}

var _ raft.FSM = &dataTracker{}

func NewDataTracker(server2 *ser.Server) *dataTracker {
	return &dataTracker{
		server: server2,
	}
}

func (f *dataTracker) Apply(l *raft.Log) interface{} {
	tmpBytes := make([]byte, len(l.Data))
	copy(tmpBytes, l.Data)
	headDate := make([]byte, 8)
	copy(headDate, tmpBytes)
	tmpBytes = tmpBytes[8:]
	buffer := bytes.NewReader(headDate)
	msg := tcpack.NewMessage(0, 0, nil)
	if err := binary.Read(buffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil
	}
	if err := binary.Read(buffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil
	}
	if msg.GetDataLen() > 0 {
		msg.Data = make([]byte, msg.GetDataLen())
		copy(msg.Data, tmpBytes)
	}
	f.server.MsQueue <- msg
	return nil
}

func (f *dataTracker) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}

func (f *dataTracker) Restore(r io.ReadCloser) error {
	return nil
}

type snapshot struct{}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	return sink.Close()
}

func (s *snapshot) Release() {}

type rpcInterface struct {
	pb.UnimplementedExampleServer
	dataTracker *dataTracker
	raft        *raft.Raft
}

func (r rpcInterface) SendMsg(ctx context.Context, req *pb.SendMsgRequest) (*pb.SendMsgResponse, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.LittleEndian, req.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, req.Id); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, req.Data); err != nil {
		return nil, err
	}
	f := r.raft.Apply(buffer.Bytes(), time.Second)
	if err := f.Error(); err != nil {
		return nil, rafterrors.MarkRetriable(err)
	}
	responseMsg := <-r.dataTracker.server.ReQueue
	return &pb.SendMsgResponse{
		DataLen: responseMsg.DataLen,
		Id:      responseMsg.Id,
		Data:    responseMsg.Data,
	}, nil
}
