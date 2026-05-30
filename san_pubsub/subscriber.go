package san_pubsub

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateTemporarySub(ctx context.Context, client *pubsub.Client, expiration time.Duration, topic string, sub string) (subName string, close func() error, err error) {
	subName = fmt.Sprintf("projects/%s/subscriptions/%s", client.Project(), sub)
	topic = fmt.Sprintf("projects/%s/topics/%s", client.Project(), topic)

	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subName,
		Topic: topic,
		ExpirationPolicy: &pubsubpb.ExpirationPolicy{
			Ttl: durationpb.New(expiration),
		},
		AckDeadlineSeconds:       10,
		RetainAckedMessages:      true,
		MessageRetentionDuration: durationpb.New(expiration),
	})

	close = func() error {
		return DeleteTemporarySub(ctx, client, subName)
	}

	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return subName, close, nil
		}
	}

	return subName, close, err
}

func DeleteTemporarySub(ctx context.Context, client *pubsub.Client, subName string) error {
	return client.SubscriptionAdminClient.DeleteSubscription(ctx, &pubsubpb.DeleteSubscriptionRequest{
		Subscription: subName,
	})
}

func SeekTime(ctx context.Context, client *pubsub.Client, subName string, ts time.Time) error {
	_, err := client.SubscriptionAdminClient.Seek(ctx, &pubsubpb.SeekRequest{
		Subscription: subName,
		Target: &pubsubpb.SeekRequest_Time{
			Time: timestamppb.New(ts),
		},
	})

	return err
}
