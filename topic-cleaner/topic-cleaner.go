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

// broker address we need to communicate to
const defaultBrokerAddress = "localhost:9092"

func deleteTopic(brokerAddress string, topicName string) {
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

	// try to delete a new topic
	err = clusterAdmin.DeleteTopic(topicName)
	if err != nil {
		log.Fatalf("Error: delete topic '%s': %v", topicName, err.Error())
	}

	// everything's seems to be ok
	log.Printf("Topic '%s' has been deleted", topicName)
}

func main() {
	var topicName string
	var brokerAddress string

	flag.StringVar(&topicName, "topic", "", "topic name")
	flag.StringVar(&brokerAddress, "broker", defaultBrokerAddress, "broker address")
	flag.Parse()

	if topicName == "" {
		log.Fatal("Topic name needs to be provided on CLI")
	}

	deleteTopic(brokerAddress, topicName)

	// deleteTopic(brokerAddress, "test_partitions_1")
	// deleteTopic(brokerAddress, "test_partitions_2")
	// deleteTopic(brokerAddress, "test_partitions_4")
	// deleteTopic(brokerAddress, "test_partitions_8")
	// deleteTopic(brokerAddress, "test_partitions_16")
}
