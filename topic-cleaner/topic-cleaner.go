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
	"log"

	"github.com/Shopify/sarama" // Sarama 1.22.0
)

// defaultBrokerAddress represents address of broker running locally.
// Usually we need to communicate with this broker.
const defaultBrokerAddress = "localhost:9092"

// deleteTopic function deletes existing topic from Kafka
func deleteTopic(brokerAddress string, topicName string) {
	brokerAddresses := []string{brokerAddress}

	// configuration handling
	// please note the versioning, it needs to be specified explicitly
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	// try to create and initialize cluster admin object
	// (it will reach out Kafka broker, so it checks the connection as well)
	clusterAdmin, err := sarama.NewClusterAdmin(brokerAddresses, config)

	// check if cluster admin has been initialized successfully
	if err != nil {
		log.Fatal("Error: create cluster admin: ", err.Error())
	}

	// everything's seems to be ok
	log.Print("Cluster admin has been initialized")

	// cluster admin needs to be closed properly
	defer func() {
		// try to close cluster admin
		err := clusterAdmin.Close()

		// check if cluster admin has been closed successfully
		if err != nil {
			log.Fatal("Error: close cluster admin: ", err.Error())
		}
	}()

	// try to delete a new topic via cluster admin
	err = clusterAdmin.DeleteTopic(topicName)
	if err != nil {
		log.Fatalf("Error: delete topic '%s': %v", topicName, err.Error())
	}

	// everything's seems to be ok -> topic has been deleted
	log.Printf("Topic '%s' has been deleted", topicName)
}

func main() {
	const noTopic = ""

	// filled via command line arguments
	var topicName string
	var brokerAddress string

	// find and parse all command line arguments
	flag.StringVar(&topicName, "topic", noTopic, "topic name")
	flag.StringVar(&brokerAddress, "broker", defaultBrokerAddress, "broker address")
	flag.Parse()

	// check if topic name has been specified on command line
	if topicName == noTopic {
		log.Fatal("Topic name needs to be provided on CLI")
	}

	deleteTopic(brokerAddress, topicName)

	// Examples how to delete existing topics:
	// deleteTopic(brokerAddress, "test_partitions_1")
	// deleteTopic(brokerAddress, "test_partitions_2")
	// deleteTopic(brokerAddress, "test_partitions_4")
	// deleteTopic(brokerAddress, "test_partitions_8")
	// deleteTopic(brokerAddress, "test_partitions_16")
}
