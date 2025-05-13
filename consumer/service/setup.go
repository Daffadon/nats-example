package service

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func JSSetup(connectionName, streamName, consumerName string) (jetstream.Consumer, context.Context, *nats.Conn, context.CancelFunc) {
	// opt, err := nats.NkeyOptionFromSeed("../config/nkeys/user_key.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	nc, err := nats.Connect("nats://127.0.0.1:4221",
		nats.UserCredentials("../config/creds/admin.creds"),
		nats.Name(connectionName),
		nats.Timeout(10*time.Second),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		log.Fatal(err)
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:      consumerName,
		Durable:   consumerName,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cons, ctx, nc, cancel
}
