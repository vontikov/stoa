package gateway

import (
	"testing"

	"context"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/vontikov/stoa/pkg/client"
	"github.com/vontikov/stoa/pkg/pb"
)

func TestProcessMetadata(t *testing.T) {
	// Make time-based calls testable
	client.TimeNowInMillis = func() int64 { return 0 }

	t.Run("Should enable barrier",
		func(t *testing.T) {
			md := client.MetadataFromCallOptions(client.WithBarrier())
			ctx := metadata.NewIncomingContext(context.Background(), md)

			msg := pb.ClusterCommand{}
			err := processMetadata(ctx, &msg)

			assert.Nil(t, err)
			assert.True(t, msg.BarrierEnabled)
			assert.Equal(t, int64(0), msg.BarrierMillis)
		})

	t.Run("Should enable waiting barrier",
		func(t *testing.T) {
			md := client.MetadataFromCallOptions(client.WithWaitingBarrier(5 * time.Second))
			ctx := metadata.NewIncomingContext(context.Background(), md)

			msg := pb.ClusterCommand{}
			err := processMetadata(ctx, &msg)
			assert.Nil(t, err)

			assert.True(t, msg.BarrierEnabled)
			assert.Equal(t, int64(5000), msg.BarrierMillis)
		})

	t.Run("Should enable message TTL",
		func(t *testing.T) {
			md := client.MetadataFromCallOptions(client.WithTTL(42 * time.Second))
			ctx := metadata.NewIncomingContext(context.Background(), md)

			msg := pb.ClusterCommand{}
			err := processMetadata(ctx, &msg)
			assert.Nil(t, err)

			assert.True(t, msg.TtlEnabled)
			assert.Equal(t, int64(42000), msg.TtlMillis)
		})
}
