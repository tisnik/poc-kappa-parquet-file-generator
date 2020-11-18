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

// This utility reads and parses the input file that should contains message
// contents split by newline character. Each message is represented as JSON (ie
// the input file consists of several JSONs, one JSON per line). This JSON is
// de-serialized and cluster ID is used to choose the `key` for message to be
// written into Kafka. Message `value` contains the whole JSON.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama" // Sarama 1.22.0
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

// New constructor creates new instance of interface to Kafka broker
func New(brokerAddress string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, nil)

	// check if Kafka producer has been initialized properly
	if err != nil {
		log.Fatal("unable to create a new Kafka producer")
		return nil, err
	}

	return producer, nil
}

// produceMessage function produce (send) one message with specified key to
// selected Kafka topic
func produceMessage(producer sarama.SyncProducer, topic string, message string, key byte) (int32, int64, error) {
	// message key needs to be represented as raw bytes
	keyBytes := []byte{key}

	// message value needs to be represented as raw bytes
	bytes := []byte(message)

	// try to produce (end) the message into specified topic
	producerMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(keyBytes),
		Value: sarama.ByteEncoder(bytes),
	}

	// check if send has been successful and get additional information
	// where the message was stored
	partition, offset, err := producer.SendMessage(producerMsg)
	if err != nil {
		log.Fatal("failed to produce message to Kafka")
	} else {
		log.Printf("message sent to partition %d at offset %d\n", partition, offset)
	}
	return partition, offset, err
}

// keyFromClusterID constructs a key (byte value 0x00-0x0f) from cluster ID.
// Key is generated from the first character of cluster ID.
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

// closeFile function tries to close the Parquet file
func closeFile(file *os.File) {
	err := file.Close()

	// check for any error during close operation
	if err != nil {
		log.Fatal("Can't close the Parquet file")
	}
}

// produceMessagesFromJSONs read provided input file line by line, parse JSON
// values from each line, creates Kafka message and send the message to Kafka
// broker into selected topic.
func produceMessagesFromJSONs(producer sarama.SyncProducer, topic string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile(file)

	scanner := bufio.NewScanner(file)

	// we need a buffer with larger capacity as lines are very long
	// (larger that the default max capacity 64kB)
	const maxCapacity = 1024 * 1024
	buffer := make([]byte, maxCapacity)
	scanner.Buffer(buffer, maxCapacity)

	// process the file line by line
	for scanner.Scan() {
		text := scanner.Text()
		bytes := []byte(text)
		var report Report

		err := json.Unmarshal(bytes, &report)
		if err != nil {
			fmt.Println("err", err)
		}

		// generate message key from cluster ID
		clusterID := report.Metadata.ClusterID
		key := clusterID[0]

		log.Printf("producing message for cluster %s using key %c", clusterID, key)
		produceMessage(producer, topic, text, key)
	}

	// check if file processing was successful
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// closeProducer function tries to close Kafka producer
func closeProducer(producer sarama.SyncProducer) {
	err := producer.Close()
	if err != nil {
		log.Fatal("producer.Close()", err)
	}
}

func main() {
	const noTopic = ""
	const noFile = ""

	// filled via command line arguments
	var topicName string
	var brokerAddress string
	var fileName string

	// find and parse all command line arguments
	flag.StringVar(&topicName, "topic", noTopic, "topic name")
	flag.StringVar(&brokerAddress, "broker", defaultBrokerAddress, "broker address")
	flag.StringVar(&fileName, "input", noFile, "input data file")
	flag.Parse()

	// check if topic name has been specified on command line
	if topicName == noTopic {
		log.Fatal("Topic name needs to be provided on CLI")
	}

	// check if topic name has been specified on command line
	if fileName == noFile {
		log.Fatal("Input file name needs to be provided on CLI")
	}

	t1 := time.Now()

	// construct new interface to Kafka broker
	producer, err := New(brokerAddress)
	if err != nil {
		log.Fatal("New producer", err)
	}
	log.Printf("Producer has been initialized: %v\n", producer)
	defer closeProducer(producer)

	// and start producing messages to it
	produceMessagesFromJSONs(producer, topicName, fileName)

	log.Println("Producing messages finished")

	// compute and print duration
	t2 := time.Now()
	since := time.Since(t1)
	log.Println("Start time: ", t1)
	log.Println("End time:   ", t2)
	log.Println("Duration:   ", since)
}
