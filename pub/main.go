package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	opt, err := nats.NkeyOptionFromSeed("../user_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	nc, err := nats.Connect("nats://127.0.0.1:4221", opt, nats.Name("Orders Publisher"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "orders",
		Description: "Messages for orders",
		Subjects: []string{
			"orders.>",
		},
		MaxBytes: 1024 * 1024 * 1024,
	})
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for {
		i++
		_, err = js.Publish(ctx, fmt.Sprintf("orders.%d", i), []byte("Hello orders"))
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("Publisher message %d \n", i)
	}

}
