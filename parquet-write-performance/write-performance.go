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
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
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

func generateColor() string {
	var colors []string = []string{
		"black",
		"blue",
		"red",
		"magenta",
		"green",
		"cyan",
		"yellow",
		"white",
	}
	return colors[rand.Int()%len(colors)]
}

func writeRecords(pw *writer.ParquetWriter, n int) {
	// create report structure to be stored in Parquet file
	record := Report{}

	for i := 0; i < n; i++ {
		clusterID := faker.UUIDHyphenated()
		record.Id = int64(i)
		record.Key = clusterID[0:1]
		record.ClusterID = clusterID
		record.Path = faker.Name()
		record.ExternalOrganization = strconv.Itoa(int(rand.Int31n(1000)))
		report := ""
		for j := 0; j < 500; j++ {
			report += faker.Paragraph()
		}
		record.Report = report

		// write the record structure into Parquet file
		err := pw.Write(record)
		if err != nil {
			log.Println("Write into Parquet error", err)
		}
	}
}

// stopWrite function stop writing into Parquet file
func stopWrite(pw *writer.ParquetWriter) {
	err := pw.WriteStop()

	// most write errors are caught at this time
	if err != nil {
		log.Println("WriteStop error", err)
	}
}

func createAndWriteIntoParquetFile(filename string, records int, compression parquet.CompressionCodec) {
	t1 := time.Now()

	w, err := os.Create(filename)
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	defer w.Close()

	// initialize Parquet file writer
	pw, err := writer.NewParquetWriterFromWriter(w, new(Report), 1)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = compression

	defer stopWrite(pw)

	writeRecords(pw, records)

	log.Println("Write Finished")

	// compute and print duration
	t2 := time.Now()
	since := time.Since(t1)
	log.Println("Start time: ", t1)
	log.Println("End time:   ", t2)
	log.Println("Duration:   ", since)
}

func main() {
	const defaultNumberOfRecords = 1000

	// filled via command line argument
	var numberOfRecrods int

	// find and parse all command line arguments
	flag.IntVar(&numberOfRecrods, "n", defaultNumberOfRecords, "number of records to be written into Parquet file")
	flag.Parse()

	createAndWriteIntoParquetFile("records_compression_none.parquet", numberOfRecrods, parquet.CompressionCodec_UNCOMPRESSED)
	createAndWriteIntoParquetFile("records_compression_snappy.parquet", numberOfRecrods, parquet.CompressionCodec_SNAPPY)
	createAndWriteIntoParquetFile("records_compression_gzip.parquet", numberOfRecrods, parquet.CompressionCodec_GZIP)
}
