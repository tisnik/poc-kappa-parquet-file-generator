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

// This tool reads all records from provided SQL database and writes all the
// records into Parquet file. It is possible to use command line flags to
// select the database and its options (user, etc.). Additionally it is
// possible to choose the output file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL database driver

	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

// default values for database connector
const (
	defaultDatabaseHostname = "localhost"
	defaultDatabasePort     = 5432
	defaultDatabaseName     = "d0"
	defaultDatabaseUser     = "postgres"
	defaultDatabasePassword = "postgres"
)

const defaultOutputFile = "flat.parquet"

// Report represents one record stored in Parquet file
type Report struct {
	Id                   int64  `parquet:"name=id, type=INT64"`
	Key                  string `parquet:"name=key, type=UTF8, encoding=PLAIN_DICTIONARY"`
	ClusterID            string `parquet:"name=cluster_id, type=UTF8, encoding=PLAIN"`
	Path                 string `parquet:"name=path, type=UTF8, encoding=PLAIN"`
	ExternalOrganization string `parquet:"name=external_organization, type=UTF8, encoding=PLAIN"`
	Report               string `parquet:"name=report, type=UTF8, encoding=PLAIN"`
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

// readAndExportAllRecords function read all reports from SQL database and
// write them into Parquet file.
func readAndExportAllRecords(storage *sql.DB, pw *writer.ParquetWriter) {
	const query = "SELECT id, key, cluster_id, path, external_organization, report FROM reports ORDER BY id"

	// some statistic about processing
	readRecords := 0
	readErrors := 0
	writtenRecords := 0
	writeErrors := 0

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
		var clusterID string
		var path string
		var externalOrganization string
		var report string

		// try to read one record from database
		err = rows.Scan(&id, &key, &clusterID, &path, &externalOrganization, &report)

		// skip errors
		if err != nil {
			log.Println("reading/scanning record", err)
			readErrors++
			continue
		}

		readRecords++

		// create report structure to be stored in Parquet file
		reportRecord := Report{
			Id:                   int64(id),
			Key:                  key,
			ClusterID:            clusterID,
			Path:                 path,
			ExternalOrganization: externalOrganization,
			Report:               report,
		}

		// write the record structure into Parquet file
		err := pw.Write(reportRecord)
		if err != nil {
			log.Println("Write into Parquet error", err)
			writeErrors++
			continue
		}
		writtenRecords++
	}

	fmt.Println("read records:    ", readRecords)
	fmt.Println("read errors:     ", readErrors)
	fmt.Println("written records: ", writtenRecords)
	fmt.Println("write errors:    ", writeErrors)

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal("process record(s)", err)
	}

}

// stopWrite function stop writing into Parquet file
func stopWrite(pw *writer.ParquetWriter) {
	err := pw.WriteStop()
	if err != nil {
		log.Println("WriteStop error", err)
	}
}

func main() {
	// filled via command line arguments
	var databaseHost string
	var databasePort int
	var databaseUser string
	var databasePassword string
	var databaseName string

	var outputFile string

	// find and parse all command line arguments
	flag.StringVar(&databaseHost, "db-host", defaultDatabaseHostname, "database host name")
	flag.IntVar(&databasePort, "db-port", defaultDatabasePort, "database port")
	flag.StringVar(&databaseName, "db-name", defaultDatabaseName, "database name")
	flag.StringVar(&databaseUser, "db-user", defaultDatabaseUser, "database user")
	flag.StringVar(&databasePassword, "db-password", defaultDatabasePassword, "database password for given user")
	flag.StringVar(&outputFile, "output", defaultOutputFile, "output file (Parquet)")
	flag.Parse()

	t1 := time.Now()

	// try to initialize the storage
	storage, err := initStorage(databaseHost, databasePort, databaseUser, databasePassword, databaseName)
	if err != nil {
		log.Fatal("Storage init", err)
	}

	// storage needs to be closed properly
	defer storage.Close()

	w, err := os.Create("flat.parquet")
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	defer w.Close()

	// initialize Parquet file writer
	pw, err := writer.NewParquetWriterFromWriter(w, new(Report), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	defer stopWrite(pw)

	readAndExportAllRecords(storage, pw)

	log.Println("Write Finished")

	// compute and print duration
	t2 := time.Now()
	since := time.Since(t1)
	log.Println("Start time: ", t1)
	log.Println("End time:   ", t2)
	log.Println("Duration:   ", since)
}
