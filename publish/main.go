package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/daffadon/learn-nats/publish/service"
	"github.com/daffadon/learn-nats/types"
)

func main() {

	connectionName := "IoT Publisher"
	streamName := "iot"
	streamDesc := "Stream for iot.>"
	subjects := []string{"iot.>"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	js, nc := service.JSSetup(ctx, connectionName, streamName, streamDesc, subjects)
	defer nc.Drain()

	iotIDs := []string{"sensor1", "sensor2", "sensor3", "sensor4", "sensor5"}
	temperatures := []float32{22.5, 25.0, 18.3, 30.1, 15.8}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping publisher.")
			return
		case <-ticker.C:
			id := rand.Intn(5)
			temp := rand.Intn(5)
			data := &types.IoTReq{
				IotID:       iotIDs[id],
				Temperature: temperatures[temp],
			}

			dataBytes, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshalling data: %v", err)
				continue
			}

			_, err = js.Publish(ctx, fmt.Sprintf("iot.%s", iotIDs[id]), []byte(dataBytes))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Publish message: iot_id: %s, temperature: %.1f \n", data.IotID, data.Temperature)
		}
	}
}
