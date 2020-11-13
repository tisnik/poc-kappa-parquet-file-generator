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
	"fmt"
	"log"
	"os"

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

func readAndDisplayAllRecords(storage *sql.DB) {
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
		fmt.Println()
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal("process record(s)", err)
	}
}

func main() {
	var err error
	w, err := os.Create("flat.parquet")
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	pw, err := writer.NewParquetWriterFromWriter(w, new(Report), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY
	num := 100
	for i := 0; i < num; i++ {
		stu := Report{
			Id:                   int64(i),
			Key:                  "a",
			ClusterID:            "StudentName",
			Path:                 "Path",
			ExternalOrganization: "eorg",
			Report:               "report",
		}
		if err = pw.Write(stu); err != nil {
			log.Println("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
		return
	}
	log.Println("Write Finished")
	w.Close()
}
