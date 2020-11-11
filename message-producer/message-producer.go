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
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama" // Sarama 1.22.0
	"log"
	"os"
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

func New(brokerAddress string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, nil)
	if err != nil {
		log.Fatal("unable to create a new Kafka producer")
		return nil, err
	}

	return producer, nil
}

func produceMessage(producer sarama.SyncProducer, topic string, message string, key byte) (int32, int64, error) {
	bytes := []byte(message)
	keyBytes := []byte{key}

	producerMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(keyBytes),
		Value: sarama.ByteEncoder(bytes),
	}

	partition, offset, err := producer.SendMessage(producerMsg)
	if err != nil {
		log.Fatal("failed to produce message to Kafka")
	} else {
		log.Printf("message sent to partition %d at offset %d\n", partition, offset)
	}
	return partition, offset, err
}

func keyFromClusterID(clusterID string) byte {
	key := clusterID[0]
	if key >= '0' && key <= '9' {
		return key - '0'
	}
	if key >= 'a' && key <= 'f' {
		key -= 'a'
		return key + 10
	}
	return 0xff
}

func produceMessagesFromJSONs(producer sarama.SyncProducer, topic string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// we need a buffer with larger capacity as lines are very long
	const maxCapacity = 512 * 1024
	buffer := make([]byte, maxCapacity)
	scanner.Buffer(buffer, maxCapacity)

	for scanner.Scan() {
		text := scanner.Text()
		bytes := []byte(text)
		var report Report

		err := json.Unmarshal(bytes, &report)
		if err != nil {
			fmt.Println("err", err)
		}
		clusterID := report.Metadata.ClusterID
		key := clusterID[0]

		log.Printf("producing message for cluster %s using key %c", clusterID, key)
		produceMessage(producer, topic, text, key)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	producer, err := New(brokerAddress)
	if err != nil {
		log.Fatal("New producer", err)
	}
	log.Printf("Producer: %v\n", producer)
	defer producer.Close()

	produceMessagesFromJSONs(producer, "test_partitions_4", "../data/10.txt")
}
