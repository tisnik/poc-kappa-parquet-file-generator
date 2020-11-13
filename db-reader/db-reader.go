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

// This tool reads all results stored in PostgreSQL database. It is possible to
// use command line flags to select the database and its options (user, etc.).
// Additionally it is possible to choose which column(s) to display.
package main

import (
	"flag"
	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL database driver
)

// default values for database connector
const (
	defaultDatabaseHostname = "localhost"
	defaultDatabasePort     = 5432
	defaultDatabaseName     = "d0"
	defaultDatabaseUser     = "postgres"
	defaultDatabasePassword = "postgres"
)

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

// readAndDisplayAllRecords function tries to read all records from database.
// Records are to be read sequentially by its unique ID. Content of recors is
// displayed on standard output.
func readAndDisplayAllRecords(storage *sql.DB, showID bool, showKey bool,
	showClusterID bool, showPath bool, showExternalOrganization bool,
	showReport bool) {

	const query = "SELECT id, key, cluster_id, path, external_organization, report FROM reports ORDER BY id"

	// try to select all records
	rows, err := storage.Query(query)
	if err != nil {
		log.Fatal("storage.Query", err)
	}
	defer rows.Close()

	// process all records, one by one
	for rows.Next() {
		var id int
		var key string
		var cluster_id string
		var path string
		var external_organization string
		var report string

		// try to read one record from database
		err = rows.Scan(&id, &key, &cluster_id, &path, &external_organization, &report)

		// skip errors
		if err != nil {
			log.Println("reading/scanning record", err)
			continue
		}

		// print the record to standard output
		if showID {
			fmt.Print(id, "\t")
		}
		if showKey {
			fmt.Print(key, "\t")
		}
		if showClusterID {
			fmt.Print(cluster_id, "\t")
		}
		if showPath {
			fmt.Print(path, "\t")
		}
		if showExternalOrganization {
			fmt.Print(external_organization, "\t")
		}
		if showReport {
			fmt.Print(report)
		}
		fmt.Println()
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal("process record(s)", err)
	}
}

func main() {
	// filled via command line arguments
	var databaseHost string
	var databasePort int
	var databaseUser string
	var databasePassword string
	var databaseName string

	var showID bool
	var showKey bool
	var showClusterID bool
	var showPath bool
	var showExternalOrganization bool
	var showReport bool

	// find and parse all command line arguments
	flag.StringVar(&databaseHost, "db-host", defaultDatabaseHostname, "database host name")
	flag.IntVar(&databasePort, "db-port", defaultDatabasePort, "database port")
	flag.StringVar(&databaseName, "db-name", defaultDatabaseName, "database name")
	flag.StringVar(&databaseUser, "db-user", defaultDatabaseUser, "database user")
	flag.StringVar(&databasePassword, "db-password", defaultDatabasePassword, "database password for given user")

	flag.BoolVar(&showID, "show-id", true, "ID of the record")
	flag.BoolVar(&showKey, "show-key", true, "Key of the record")
	flag.BoolVar(&showClusterID, "show-clusterid", true, "ClusterID")
	flag.BoolVar(&showPath, "show-path", true, "path")
	flag.BoolVar(&showExternalOrganization, "show-org", true, "external organization")
	flag.BoolVar(&showReport, "show-report", false, "the whole report")
	flag.Parse()

	// try to initialize the storage
	storage, err := initStorage(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
		log.Fatal("Storage init", err)
	}

	// storage needs to be closed properly
	defer storage.Close()

	// read and display all records
	readAndDisplayAllRecords(storage, showID, showKey, showClusterID,
		showPath, showExternalOrganization, showReport)
}
