package gateway

import (
	"errors"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrIncorrectResponseType is the error returned when a collection method is
// called and Raft returned unexpected response type.
var ErrIncorrectResponseType = errors.New("incorrect response type")

// ErrNotLeader is the error returned in case of a Raft follower call.
var ErrNotLeader = status.Errorf(codes.FailedPrecondition, raft.ErrNotLeader.Error())
