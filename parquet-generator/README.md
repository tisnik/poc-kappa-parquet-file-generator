# PoC: Parquet file generator based on Kappa architecture

## Parquet generator

Tool to read reports from database and create Parquet files from those reports.

### Description

### How to build the tool

```
go build
```

### Usage

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
