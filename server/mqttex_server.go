package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"sync"
	"time"
)

const (
	rootTopic      = "go-mqtt/"
	mqttBrokerHost = "tcp://localhost:1883"
)

func sha1hash(data []string) []byte {
	h := sha1.New()
	for _, wd := range data {
		h.Write([]byte(wd))
	}
	return h.Sum(nil)
}

func worker(queue <-chan string, args []string) {
	sample := sha1hash(args)

	ts := time.Now()

	s := make([]string, len(args))
	for elem := range queue {
		s = append(s[1:], elem)
		if bytes.Equal(sha1hash(s), sample) {
			now := time.Now()
			diff := now.Sub(ts)
			ts = now
			fmt.Printf("%v <- %s [Found! Elapsed time: %v]\n", s, elem, diff)
		} else {
			fmt.Printf("%v <- %s\n", s, elem)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [text]\n", os.Args[0])
		os.Exit(2)
	}
	args := os.Args[1:]

	queue := make(chan string, len(args))
	defer close(queue)
	go worker(queue, args)

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

			msgRcvd := func(client MQTT.Client, message MQTT.Message) {
				queue <- fmt.Sprintf("%s", message.Payload())
			}

			if token := c.Subscribe(topic, 0, msgRcvd); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
				os.Exit(1)
			}
			log.Printf("Subscribed on topic: %s\n", topic)

			defer func() {
				if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
					log.Fatal(token.Error())
					os.Exit(1)
				}

			}()

			log.Printf("Waiting on topic: %s\n", topic)
			done := make(chan bool)
			<-done
		}(sarg)
	}

	wg.Wait()
}
