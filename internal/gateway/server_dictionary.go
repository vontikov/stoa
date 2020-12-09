package gateway

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vontikov/stoa/pkg/pb"
)

func (s *server) DictionarySize(ctx context.Context, v *pb.Name) (*pb.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_SIZE,
		Payload: &pb.ClusterCommand_Name{Name: v},
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
	if r, ok := fut.Response().(*pb.Value); ok {
		return r, nil
	}
	return &pb.Value{}, nil
}

func (s *server) DictionaryClear(ctx context.Context, v *pb.Name) (*pb.Empty, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_CLEAR,
		Payload: &pb.ClusterCommand_Name{Name: v},
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
	return &pb.Empty{}, nil
}

func (s *server) DictionaryPut(ctx context.Context, v *pb.KeyValue) (*pb.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_PUT,
		Payload: &pb.ClusterCommand_KeyValue{KeyValue: v},
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
	if r, ok := fut.Response().(*pb.Value); ok {
		return r, nil
	}
	panic(ErrIncorrectResponseType)
}

func (s *server) DictionaryPutIfAbsent(ctx context.Context, v *pb.KeyValue) (*pb.Result, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_PUT_IF_ABSENT,
		Payload: &pb.ClusterCommand_KeyValue{KeyValue: v},
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

func (s *server) DictionaryGet(ctx context.Context, v *pb.Key) (*pb.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_GET,
		Payload: &pb.ClusterCommand_Key{Key: v},
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
	if r, ok := fut.Response().(*pb.Value); ok {
		return r, nil
	}
	panic(ErrIncorrectResponseType)
}

func (s *server) DictionaryRemove(ctx context.Context, v *pb.Key) (*pb.Result, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_REMOVE,
		Payload: &pb.ClusterCommand_Key{Key: v},
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
	return &pb.Result{Ok: true}, nil
}

func (s *server) DictionaryRange(v *pb.Name, stream pb.Stoa_DictionaryRangeServer) error {
	if err := v.Validate(); err != nil {
		return err
	}

	if !s.cluster.IsLeader() {
		return ErrNotLeader
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_DICTIONARY_RANGE,
		Payload: &pb.ClusterCommand_Name{Name: v},
	}

	if err := processMetadata(stream.Context(), &msg); err != nil {
		return err
	}

	if msg.BarrierEnabled {
		fut := s.cluster.Raft().Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return err
		}
	}

	// bypass Raft
	ch, ok := s.cluster.Apply(&msg).(chan *pb.KeyValue)
	if !ok {
		panic(ErrIncorrectResponseType)
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		case kv := <-ch:
			if kv == nil {
				return nil
			}
			if err := stream.Send(kv); err != nil {
				return err
			}
		}
	}
}
