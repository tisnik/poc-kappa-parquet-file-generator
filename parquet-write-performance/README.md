# PoC: Parquet file generator based on Kappa architecture

## Write-performance

Simple Parquet writing benchmark.

### Description

This tool is able to write selected number of records into Parquet file and
measure the time to finish writing. Currently, only records with the structure
`Report` is written correctly.

### How to build the tool

```
go build
```

### Usage

```
Usage of write-performance:
  -n int
        number of records to be written into Parquet file (default 1000)

```

### Parquet file structure

Currently, the simplified structure is used. It is based on the following Go
struct:

```
type Report struct {
        Id                   int64  `parquet:"name=id, type=INT64"`
        Key                  string `parquet:"name=key, type=UTF8, encoding=PLAIN_DICTIONARY"`
        ClusterID            string `parquet:"name=cluster_id, type=UTF8, encoding=PLAIN"`
        Path                 string `parquet:"name=path, type=UTF8, encoding=PLAIN"`
        ExternalOrganization string `parquet:"name=external_organization, type=UTF8, encoding=PLAIN"`
        Report               string `parquet:"name=report, type=UTF8, encoding=PLAIN"`
}
```
