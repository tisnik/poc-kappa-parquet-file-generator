# PoC: Parquet file generator based on Kappa architecture

## Parquet read

Tool to read reports from specified Parquet file. Report is displayed on
standard output.

### Description

This tool is able to read all records stored in selected Parquet file.
Currently, only records with the structure `Report` is read correctly. Name of
input Parquet file needs to be selected from command line.

### How to build the tool

```
go build
```

### Usage

```
Usage of parquet-reader:
  -input string
        input data file

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
