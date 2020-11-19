package gateway

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vontikov/stoa/pkg/pb"
)

func (s *server) MutexTryLock(ctx context.Context, v *pb.Id) (*pb.Result, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_MUTEX_TRY_LOCK,
		Payload: &pb.ClusterCommand_Id{Id: v},
	}

	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.cluster.Raft().Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.cluster.Raft().Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	if r, ok := fut.Response().(*pb.Result); ok {
		return r, nil
	}
	panic(ErrIncorrectResponseType)
}

func (s *server) MutexUnlock(ctx context.Context, v *pb.Id) (*pb.Result, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_MUTEX_UNLOCK,
		Payload: &pb.ClusterCommand_Id{Id: v},
	}

	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.cluster.Raft().Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.cluster.Raft().Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	if r, ok := fut.Response().(*pb.Result); ok {
		return r, nil
	}
	panic(ErrIncorrectResponseType)
}
