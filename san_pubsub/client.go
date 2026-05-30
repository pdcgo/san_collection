package san_pubsub

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPubsubEmulator(ctx context.Context, projectID string) (c *pubsub.Client, err error) {
	conn, err := grpc.NewClient("localhost:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pubsub.NewClient(ctx, projectID, option.WithGRPCConn(conn))

}

func NewPubSubClient() (c *pubsub.Client, err error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT not set")
	}
	return pubsub.NewClient(context.Background(), projectID)
}

func NewPubSubClientWithContext(ctx context.Context) (c *pubsub.Client, err error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT not set")
	}
	return pubsub.NewClient(ctx, projectID)
}
