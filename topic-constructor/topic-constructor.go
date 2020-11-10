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

// see https://stackoverflow.com/questions/44094926/creating-kafka-topic-in-sarama
// for further Printfrmation

import (
	"flag"
	"log"

	"github.com/Shopify/sarama" // Sarama 1.22.0
)

// broker address we need to communicate to
const defaultBrokerAddress = "localhost:9092"

func createTopic(brokerAddress string, topicName string,
	numPartitions int32, replicationFactor int32) {
	brokerAddrs := []string{brokerAddress}

	// configuration handling
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	// try to create cluster admin object
	clusterAdmin, err := sarama.NewClusterAdmin(brokerAddrs, config)
	if err != nil {
		log.Fatal("Error: create cluster admin: ", err.Error())
	}
	log.Print("Cluster admin has been initialized")

	// cluster admin needs to be closed properly
	defer func() {
		err := clusterAdmin.Close()
		if err != nil {
			log.Fatal("Error: close cluster admin: ", err.Error())
		}
	}()

	// try to create a new topic
	err = clusterAdmin.CreateTopic(topicName,
		&sarama.TopicDetail{
			NumPartitions:     numPartitions,
			ReplicationFactor: int16(replicationFactor),
		}, false)
	if err != nil {
		log.Fatalf("Error: create topic '%s': %v", topicName, err.Error())
	}

	// everything's seems to be ok
	log.Printf("Topic '%s' has been created", topicName)
}

func main() {
	var topicName string
	var brokerAddress string
	var partitions int
	var replicationFactor int

	flag.StringVar(&topicName, "topic", "", "topic name")
	flag.StringVar(&brokerAddress, "broker", defaultBrokerAddress, "broker address")
	flag.IntVar(&partitions, "partitions", 1, "number of partitions")
	flag.IntVar(&replicationFactor, "replication-factor", 1, "replication factor")
	flag.Parse()

	if topicName == "" {
		log.Fatal("Topic name needs to be provided on CLI")
	}

	createTopic(brokerAddress, topicName, int32(partitions), int32(replicationFactor))
	// createTopic(brokerAddress, "test_partitions_1", 1, 1)
	// createTopic(brokerAddress, "test_partitions_2", 2, 1)
	// createTopic(brokerAddress, "test_partitions_4", 4, 1)
	// createTopic(brokerAddress, "test_partitions_8", 8, 1)
	// createTopic(brokerAddress, "test_partitions_16", 16, 1)
}
