package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/daffadon/learn-nats/consumer/service"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	connectionName := "IoT Consumer"
	streamName := "iot"
	consumerName := "iot-durable"

	log.Printf("I'm consuming stream %s at consumer %s", strings.ToUpper(streamName), strings.ToUpper(consumerName))
	cons, ctx, nc, cancel := service.JSSetup(connectionName, streamName, consumerName)
	defer nc.Close()
	defer cancel()

	conn := service.InitDB()
	_, err := cons.Consume(func(msg jetstream.Msg) {
		log.Printf("Received Message: %s \n", string(msg.Data()))
		go service.SaveToDB(ctx, conn, msg.Data())
		go msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
