/*
Copyright © 2020 Pavel Tisnovsky

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

// defaultBrokerAddress represents address of broker running locally.
// Usually we need to communicate with this broker.
const defaultBrokerAddress = "localhost:9092"

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

// New constructor creates new instance of interface to Kafka consumer
func New(brokerAddress string) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{brokerAddress}, nil)

	// check if Kafka consumer has been initialized properly
	if err != nil {
		log.Fatal("error creating consumer", err)
		return nil, err
	}
	log.Print("consumer has been created")

	return consumer, nil
}

// startConsumer function initializes the consumer and starts processing
// messages
func startConsumer(brokerAddress string, topicName string, partition int) {
	// construct Kafka consumer
	consumer, err := New(brokerAddress)
	if err != nil {
		log.Fatal("New consumer", err)
	}

	// Kafka consumer needs to be closed properly
	defer func() {
		// try to close Kafka consumer and check if the operation was
		// successful
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	// construct Kafka consumer for selected topic and partition
	partitionConsumer, err := consumer.ConsumePartition(topicName, int32(partition), sarama.OffsetOldest)
	if err != nil {
		log.Fatal("New partition consumer", err)
	}

	// partition consumer needs to be closed properly
	defer func() {
		// try to close Kafka partition consumer and check if the
		// operation was successful
		if err := partitionConsumer.Close(); err != nil {
			panic(err)
		}
	}()

	consumed := consumeMessages(partitionConsumer)

	fmt.Printf("Consumed: %d\n", consumed)
}

// consumeMessages function consumes messages from Kafka and process them
// accordingly
func consumeMessages(partitionConsumer sarama.PartitionConsumer) int {
	consumed := 0
	for {
		msg := <-partitionConsumer.Messages()
		fmt.Printf("Consumed message with key '%s' at offset %d\n", msg.Key, msg.Offset)
		consumed++
	}
	return consumed
}

func main() {
	const noTopic = ""
	const noPartition = -1

	// filled via command line arguments
	var topicName string
	var brokerAddress string
	var partition int

	// find and parse all command line arguments
	flag.StringVar(&topicName, "topic", noTopic, "topic name")
	flag.StringVar(&brokerAddress, "broker", defaultBrokerAddress, "broker address")
	flag.IntVar(&partition, "partition", noPartition, "partition to consume")
	flag.Parse()

	// check if topic name has been specified on command line
	if topicName == noTopic {
		log.Fatal("Topic name needs to be provided on CLI")
	}

	// check if partition has been specified on command line
	if partition == noPartition {
		log.Fatal("Partition needs to be provided on CLI")
	}

	startConsumer(brokerAddress, topicName, partition)
}
