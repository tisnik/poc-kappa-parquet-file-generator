# PoC: Parquet file generator based on Kappa architecture

## Parquet generator

Tool to read reports from database and create Parquet files from those reports.

### Description

This tool reads all records from provided SQL database and writes all the
records into Parquet file. It is possible to use command line flags to select
the database and its options (user, etc.). Additionally it is possible to
choose the output file.

### How to build the tool

```
go build
```

### Usage

```
Usage of parquet-generator:
  -db-host string
        database host name (default "localhost")
  -db-name string
        database name (default "d0")
  -db-password string
        database password for given user (default "postgres")
  -db-port int
        database port (default 5432)
  -db-user string
        database user (default "postgres")
  -output string
        output file (Parquet) (default "flat.parquet")
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
