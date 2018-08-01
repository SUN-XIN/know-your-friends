package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/SUN-XIN/know-your-friends/types"
	"github.com/Shopify/sarama"
)

const (
	KAFKA_BROKER = "localhost:9092"
	KAFKA_TOPIC  = "test_zenly"
)

func main() {
	config := sarama.NewConfig()
	//  config.Producer.RequiredAcks = sarama.WaitForAll
	//  config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	config.Version = sarama.MaxVersion

	brokers := []string{KAFKA_BROKER}
	topic := KAFKA_TOPIC

	p, err := sarama.NewSyncProducer(brokers, config)
	defer p.Close()
	if err != nil {
		log.Printf("Failed NewSyncProducer: %+v", err)
		return
	}

	sess := types.SessionDetail{
		UserID1: "testuser2",
		UserID2: "testuser3",

		StartDate: 1532678400,
		EndDate:   1532678800,

		Lat: 48.847016,
		Lng: 2.355808,
	}
	var msg sarama.ProducerMessage

	b, err := json.Marshal(sess)
	if err != nil {
		log.Printf("Failed Marshal: %+v", err)
		return
	}
	msg = sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err := p.SendMessage(&msg); err != nil {
		log.Printf("Failed SendMessage: %+v", err)
		return
	}
	log.Printf("Send msg ok")
}
