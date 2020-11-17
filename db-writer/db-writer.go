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

// This tool consumes messages from selected Kafka topic and partition. Such
// messages are unmarshalled, parsed, and stored into PostgreSQL database. It
// is possible to use command line flags to select the database, Kafka topic,
// and Kafka partition.

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL database driver

	"github.com/Shopify/sarama"
)

// defaultBrokerAddress represents address of broker running locally.
// Usually we need to communicate with this broker.
const defaultBrokerAddress = "localhost:9092"

// default values for database connector
const (
	defaultDatabaseHostname = "localhost"
	defaultDatabasePort     = 5432
	defaultDatabaseName     = "d0"
	defaultDatabaseUser     = "postgres"
	defaultDatabasePassword = "postgres"
)

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
func startConsumer(storage *sql.DB, brokerAddress string, topicName string, partition int, csvWriter *csv.Writer) {
	// construct Kafka consumer
	consumer, err := New(brokerAddress)
	if err != nil {
		log.Fatal("New consumer", err)
	}

	// Kafka consumer needs to be closed properly
	defer closeConsumer(consumer)

	// construct Kafka consumer for selected topic and partition
	partitionConsumer, err := consumer.ConsumePartition(topicName, int32(partition), sarama.OffsetOldest)
	if err != nil {
		log.Fatal("New partition consumer", err)
	}

	// partition consumer needs to be closed properly
	defer closePartitionConsumer(partitionConsumer)

	consumed, errors := consumeMessages(storage, partitionConsumer, csvWriter)

	fmt.Printf("Consumed: %d\n", consumed)
	fmt.Printf("Errors:   %d\n", errors)
}

// closeConsumer function tries to close Kafka consumer and checks if the
// operation was successful
func closeConsumer(consumer sarama.Consumer) {
	err := consumer.Close()
	if err != nil {
		log.Fatal("closeConsumer", err)
	}
}

// closePartitionConsumer function tries to close Kafka partition consumer and
// check if the operation was successful
func closePartitionConsumer(partitionConsumer sarama.PartitionConsumer) {
	err := partitionConsumer.Close()
	if err != nil {
		log.Fatal("closePartitionConsumer", err)
	}
}

// consumeMessages function consumes messages from Kafka and process them
// accordingly
func consumeMessages(storage *sql.DB, partitionConsumer sarama.PartitionConsumer, csvWriter *csv.Writer) (int, int) {
	t1 := time.Now()
	consumed := 0
	errors := 0

	for {
		msg := <-partitionConsumer.Messages()
		// fmt.Printf("Consumed message with key '%s' at offset %d with length %d\n", msg.Key, msg.Offset, len(msg.Value))
		err := writeMessageIntoDatabase(storage, msg)
		if err != nil {
			log.Println(err)
			errors++
		}
		consumed++

		// compute and print duration
		// t2 := time.Now()
		since := time.Since(t1)
		// log.Println("Start time: ", t1)
		// log.Println("End time:   ", t2)
		log.Printf("Consumed: %d  Offset: %d  Duration: %v", consumed, msg.Offset, since)

		time := strconv.FormatFloat(float64(since)/1000.0/1000.0, 'f', 1, 64)
		csvWriter.Write([]string{time})
		csvWriter.Flush()
	}
	return consumed, errors
}

// writeMessageIntoDatabase function tries to write the message into database
func writeMessageIntoDatabase(storage *sql.DB, message *sarama.ConsumerMessage) error {
	var report Report

	err := json.Unmarshal(message.Value, &report)
	if err != nil {
		fmt.Println("Message unmarshalling", err)
		return err
	}

	// all info needed for insert new record into database
	key := message.Key
	clusterID := report.Metadata.ClusterID
	orgID := report.Metadata.ExternalOrganization
	path := report.Path
	value := string(message.Value)

	// prepare INSERT statement
	insertStatement := `INSERT INTO reports(key, cluster_id, external_organization, path, report) values ($1, $2, $3, $4, $5)`

	// run the INSERT statement
	_, err = storage.Exec(insertStatement, key, clusterID, orgID, path, value)
	if err != nil {
		fmt.Println("Insert into DB", err)
		return err
	}

	// everything's fine
	return nil
}

// initStorage iniciates connection to DB storage.
func initStorage(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	log.Println("Connecting to database: " + psqlconn)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal("Postgres driver initialization", err)
	}
	log.Println("Postgres driver initialization: OK")

	// check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping to Postgres", err)
	}
	log.Println("Connection to database: OK")

	return db, err
}

func main() {
	const noTopic = ""
	const noPartition = -1

	// filled via command line arguments
	var topicName string
	var brokerAddress string
	var partition int

	var databaseHost string
	var databasePort int
	var databaseUser string
	var databasePassword string
	var databaseName string

	// find and parse all command line arguments
	flag.StringVar(&topicName, "topic", noTopic, "topic name")
	flag.StringVar(&brokerAddress, "broker", defaultBrokerAddress, "broker address")
	flag.IntVar(&partition, "partition", noPartition, "partition to consume")
	flag.StringVar(&databaseHost, "db-host", defaultDatabaseHostname, "database host name")
	flag.IntVar(&databasePort, "db-port", defaultDatabasePort, "database port")
	flag.StringVar(&databaseName, "db-name", defaultDatabaseName, "database name")
	flag.StringVar(&databaseUser, "db-user", defaultDatabaseUser, "database user")
	flag.StringVar(&databasePassword, "db-password", defaultDatabasePassword, "database password for given user")
	flag.Parse()

	// check if topic name has been specified on command line
	if topicName == noTopic {
		log.Fatal("Topic name needs to be provided on CLI")
	}

	// check if partition has been specified on command line
	if partition == noPartition {
		log.Fatal("Partition needs to be provided on CLI")
	}

	// try to initialize the storage
	storage, err := initStorage(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
		log.Fatal("Storage init", err)
	}

	// storage needs to be closed properly
	defer storage.Close()

	csvFileName := fmt.Sprintf("db-writer-%d.csv", partition)
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		log.Fatal("Create CSV file", err)
	}

	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()
	csvWriter.Write([]string{"Cummulative time"})

	// start consuming messages and store them into opened storage
	startConsumer(storage, brokerAddress, topicName, partition, csvWriter)
}
