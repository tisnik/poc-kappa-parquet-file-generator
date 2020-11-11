package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

// broker address we need to communicate to
const brokerAddress = "localhost:9092"

// Metadata represents nested data structure containing report metadata
type Metadata struct {
	ClusterID            string `json:"cluster_id"`
	ExternalOrganization string `json:"external_organization"`
}

// Report represents the overall structure of report data (messages)
type Report struct {
	Path     string   `json:"path"`
	Metadata Metadata `json:"metadata"`
}

func New(brokerAddress string) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{brokerAddress}, nil)
	if err != nil {
		log.Fatal("error creating consumer", err)
		return nil, err
	}
	log.Print("consumer has been created")

	return consumer, nil
}

func main() {
	consumer, err := New(brokerAddress)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("test_partitions_4", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			panic(err)
		}
	}()

	consumed := 0
	for {
		msg := <-partitionConsumer.Messages()
		fmt.Printf("Consumed message with key '%s' at offset %d\n", msg.Key, msg.Offset)
		consumed++
	}

	fmt.Printf("Consumed: %d\n", consumed)
}
