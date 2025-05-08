package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	opt, err := nats.NkeyOptionFromSeed("../user_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	nc, err := nats.Connect(nats.DefaultURL, opt, nats.Name("Orders Consumer"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	stream, err := js.Stream(ctx, "orders")
	if err != nil {
		log.Fatal(err)
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    "order_processor",
		Durable: "order_processor",
	})
	if err != nil {
		log.Fatal(err)
	}

	cctx, err := cons.Consume(func(msg jetstream.Msg) {
		log.Printf("Received Message: %s \n", string(msg.Subject()))
		msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cctx.Stop()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
}
