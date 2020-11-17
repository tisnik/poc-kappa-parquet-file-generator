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

// This tool is able to read all records stored in selected Parquet file and
// measure the time to finish reading.  Currently, only records with the
// structure `Report` is read correctly.
package main

import (
	"log"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
)

// Report represents one record stored in Parquet file
type Report struct {
	ID                   int64  `parquet:"name=id, type=INT64"`
	Key                  string `parquet:"name=key, type=UTF8, encoding=PLAIN_DICTIONARY"`
	ClusterID            string `parquet:"name=cluster_id, type=UTF8, encoding=PLAIN"`
	Path                 string `parquet:"name=path, type=UTF8, encoding=PLAIN"`
	ExternalOrganization string `parquet:"name=external_organization, type=UTF8, encoding=PLAIN"`
	Report               string `parquet:"name=report, type=UTF8, encoding=PLAIN"`
}

// closeReader tries to close the given Parquet file reader
func closeReader(reader source.ParquetFile) {
	err := reader.Close()
	if err != nil {
		log.Println("close reader:", err)
	}
}

// readParquetFile function reads all records from Parquet file
func readParquetFile(fileName string) {
	t1 := time.Now()

	const parallelNumber = 1

	// construct the file reader and try to open the Parquet file for
	// reading
	fileReader, err := local.NewLocalFileReader(fileName)

	if err != nil {
		log.Fatal("Can't open file", err)
		return
	}

	// fileReader needs to be closed properly
	defer closeReader(fileReader)

	// initializa Parquet file reader
	parquetReader, err := reader.NewParquetReader(fileReader, new(Report),
		parallelNumber)

	if err != nil {
		log.Fatal("Can't create parquet reader", err)
		return
	}

	// parquetReader needs to be stopped
	defer parquetReader.ReadStop()

	readRecords(parquetReader)

	// compute and print duration
	t2 := time.Now()
	since := time.Since(t1)
	log.Println("Start time: ", t1)
	log.Println("End time:   ", t2)
	log.Println("Duration:   ", since)
}

func readRecords(parquetReader *reader.ParquetReader) {
	recordCount := int(parquetReader.GetNumRows())
	log.Println("Records to read", recordCount)

	record := make([]Report, 1)
	records := 0

	// try to read and display all records
	for i := 0; i < recordCount; i++ {
		// try to read record
		err := parquetReader.Read(&record)
		if err != nil {
			log.Println("Read error", err)
			continue
		} else {
			records++
		}
	}
	log.Println("Read", records, "records")
}

func main() {
	readParquetFile("records_compression_none.parquet")
	readParquetFile("records_compression_snappy.parquet")
	readParquetFile("records_compression_gzip.parquet")
}
