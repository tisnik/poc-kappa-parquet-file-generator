package main

import (
	"log"
	"os"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
)

type Report struct {
	Id                   int64  `parquet:"name=id, type=INT64"`
	Key                  string `parquet:"name=key, type=UTF8, encoding=PLAIN_DICTIONARY"`
	ClusterID            string `parquet:"name=cluster_id, type=UTF8, encoding=PLAIN"`
	Path                 string `parquet:"name=path, type=UTF8, encoding=PLAIN"`
	ExternalOrganization string `parquet:"name=external_organization, type=UTF8, encoding=PLAIN"`
	Report               string `parquet:"name=report, type=UTF8, encoding=PLAIN"`
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

	///read
	fr, err := local.NewLocalFileReader("flat.parquet")
	if err != nil {
		log.Println("Can't open file")
		return
	}

	pr, err := reader.NewParquetReader(fr, new(Report), 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return
	}
	num = int(pr.GetNumRows())
	for i := 0; i < num/10; i++ {
		if i%2 == 0 {
			pr.SkipRows(10) //skip 10 rows
			continue
		}
		report := make([]Report, 10) //read 10 rows
		if err = pr.Read(&report); err != nil {
			log.Println("Read error", err)
		}
		log.Println(report)
	}

	pr.ReadStop()
	fr.Close()

}
