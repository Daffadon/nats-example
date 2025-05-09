package service

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func JSSetup(ctx context.Context, connectionName, streamName, streamDesc string, subjects []string) (jetstream.JetStream, *nats.Conn) {
	opt, err := nats.NkeyOptionFromSeed("../user_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	nc, err := nats.Connect("nats://127.0.0.1:4224",
		opt, nats.Name(connectionName),
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

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        streamName,
		Description: streamDesc,
		Subjects:    subjects,
		MaxBytes:    1024 * 1024 * 1024,
	})
	if err != nil {
		log.Fatal(err)
	}
	return js, nc
}
