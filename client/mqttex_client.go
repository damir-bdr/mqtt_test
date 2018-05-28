package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	rootTopic      = "go-mqtt/"
	mqttBrokerHost = "tcp://localhost:1883"
	period         = 500
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [text]\n", os.Args[0])
		os.Exit(2)
	}
	args := os.Args[1:]

	var wg sync.WaitGroup
	wg.Add(1)

	for _, sarg := range args {
		go func(subtopic string) {
			opts := MQTT.NewClientOptions().AddBroker(mqttBrokerHost)
			opts.SetClientID(uuid.NewV4().String())

			c := MQTT.NewClient(opts)
			if token := c.Connect(); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
			defer c.Disconnect(250)

			topic := rootTopic + subtopic
			for {
				time.Sleep(time.Duration(1000+rand.Intn(period)) * time.Millisecond)

				text := subtopic
				log.Printf("publish topic: %s msg: %s\n", topic, text)
				token := c.Publish(topic, 0, false, text)
				token.Wait()
			}
		}(sarg)
	}

	wg.Wait()
}
