package gateway

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/vontikov/stoa/pkg/pb"
)

func (s *server) QueueSize(ctx context.Context, v *pb.Name) (*pb.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_QUEUE_SIZE,
		Payload: &pb.ClusterCommand_Name{Name: v},
	}

	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.peer.R.Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.peer.R.Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}

	if r, ok := fut.Response().(*pb.Value); ok {
		return r, nil
	}
	panic(errIncorrectResponseType)
}

func (s *server) QueueClear(ctx context.Context, v *pb.Name) (*pb.Empty, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_QUEUE_CLEAR,
		Payload: &pb.ClusterCommand_Name{Name: v},
	}

	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.peer.R.Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.peer.R.Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	return &pb.Empty{}, nil
}

func (s *server) QueueOffer(ctx context.Context, v *pb.Value) (*pb.Result, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_QUEUE_OFFER,
		Payload: &pb.ClusterCommand_Value{Value: v},
	}
	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.peer.R.Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.peer.R.Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	return &pb.Result{Ok: true}, nil
}

func (s *server) QueuePoll(ctx context.Context, v *pb.Name) (*pb.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_QUEUE_POLL,
		Payload: &pb.ClusterCommand_Name{Name: v},
	}
	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.peer.R.Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.peer.R.Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	if r, ok := fut.Response().(*pb.Value); ok {
		return r, nil
	}
	panic(errIncorrectResponseType)
}

func (s *server) QueuePeek(ctx context.Context, v *pb.Name) (*pb.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}

	msg := pb.ClusterCommand{
		Command: pb.ClusterCommand_QUEUE_PEEK,
		Payload: &pb.ClusterCommand_Name{Name: v},
	}
	if err := processMetadata(ctx, &msg); err != nil {
		return nil, err
	}
	if msg.BarrierEnabled {
		fut := s.peer.R.Barrier(time.Duration(msg.BarrierMillis) * time.Millisecond)
		if err := fut.Error(); err != nil {
			return nil, err
		}
	}

	cmd, err := proto.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	fut := s.peer.R.Apply(cmd, 0)
	if err := fut.Error(); err != nil {
		return nil, wrapError(err)
	}
	if r, ok := fut.Response().(*pb.Value); ok {
		return r, nil
	}
	panic(errIncorrectResponseType)
}
