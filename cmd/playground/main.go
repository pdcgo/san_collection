package main

// func main() {
// 	ctx := context.Background()
// 	client, err := san_pubsub.NewPubSubClientWithContext(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	subname, _, err := san_pubsub.CreateTemporarySub(ctx, client, time.Hour*24, "selling-topic", "tmp-sub")
// 	if err != nil {
// 		panic(err)
// 	}

// 	slog.Info("created", "subname", subname)

// 	slog.Info("seek on start 5 day")
// 	_, err = client.SubscriptionAdminClient.Seek(ctx, &pubsubpb.SeekRequest{
// 		Subscription: subname,
// 		Target: &pubsubpb.SeekRequest_Time{
// 			Time: timestamppb.New(time.Now().AddDate(0, 0, -5)),
// 		},
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	slog.Info("start listening", "subname", subname)

// 	var msgCount atomic.Int64

// 	err = client.Subscriber(subname).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {

// 		debugtool.LogJson(msg)
// 		msg.Ack()
// 		msgCount.Add(1)
// 		slog.Info("message count", "count", msgCount.Load())
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	// defer delete()

// }
