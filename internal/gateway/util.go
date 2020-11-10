package gateway

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/raft"

	"github.com/vontikov/stoa/internal/common"
	"github.com/vontikov/stoa/pkg/pb"
)

func processMetadata(ctx context.Context, msg *pb.ClusterCommand) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	if kv, ok := md[common.MetaKeyBarrier]; ok {
		b, err := strconv.Atoi(kv[0])
		if err != nil {
			return err
		}
		msg.BarrierEnabled = true
		msg.BarrierMillis = int64(b)
	}

	if kv, ok := md[common.MetaKeyTTL]; ok {
		t, err := strconv.Atoi(kv[0])
		if err != nil {
			return err
		}
		msg.TtlEnabled = true
		msg.TtlMillis = int64(t)
	}

	return nil
}

func wrapError(err error) error {
	if errors.Is(err, raft.ErrNotLeader) {
		status.Errorf(codes.ResourceExhausted, "%s", err)
	}
	return err
}
