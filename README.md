# PoC: Parquet file generator based on Kappa architecture

[![Go Report Card](https://goreportcard.com/badge/github.com/tisnik/poc-kappa-parquet-file-generator)](https://goreportcard.com/report/github.com/tisnik/poc-kappa-parquet-file-generator)
[![GitHub Pages](https://img.shields.io/badge/%20-GitHub%20Pages-informational)](https://tisnik.github.io/poc-kappa-parquet-file-generator/)

Parquet file generator and pre-worker

## Tools available in this repository

1. [Topic constructor](topic-constructor/README.md)
1. [Topic cleaner](topic-cleaner/README.md)
1. [DB writer](db-writer/README.md)
1. [DB reader](db-reader/README.md)
1. [Message producer](message-producer/README.md)
1. [Parquet generator](parquet-generator/README.md)
1. [Parquet reader](parquet-reader/README.md)
1. [Parquet read performance (benchmark)](parquet-read-performance/README.md)
1. [Parquet write performance (benchmark)](parquet-write-performance/README.md)

## Makefile

```
Usage: make <OPTIONS> ... <TARGETS>

Available targets are:

create-db                            Create database and initialize table(s)
drop-db                              Drop database completely
topic-cleaner/topic-cleaner          Build topic-cleaner tool
topic-constructor/topic-constructor  Build topic-constructor tool
message-producer/message-producer    Build message producer tool
db-reader/db-reader                  Build DB-reader tool
db-writer/db-writer                  Build DB-writer tool
parquet-generator/parquet-generator  Build parquet-generator tool
parquet-reader/parquet-reader        Build parquet-reader tool
parquet-read-performance/read-performance Build parquet-read-performance benchmark
parquet-write-performance/write-performance Build parquet-write-performance benchmark
clean                                Perform cleanup
fmt                                  Run Go formatter over all source files
linter                               Run Go linter over all source files
errcheck                             Run Go error checker over all source files
vet                                  Run go vet. Report likely mistakes in source code
gocyclo                              Run gocyclo
help                                 Show this help screen
```

## Database configuration

It is described in [this file](database/README.md).

## Notes

### Kafkacat basic usage in client mode

```
kafkacat -b localhost:9092 -t test_partitions_16 -C -K\\t
kafkacat -b localhost:9092 -t test_partitions_16 -C -f 'Key: %k\nValue: %s\n'
```
