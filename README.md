# PoC: Parquet file generator based on Kappa architecture

[![Go Report Card](https://goreportcard.com/badge/github.com/tisnik/poc-kappa-parquet-file-generator)](https://goreportcard.com/report/github.com/tisnik/poc-kappa-parquet-file-generator)

Parquet file generator and pre-worker

## Tools available in this repository

1. [Topic constructor](topic-constructor/README.md)
1. [Topic cleaner](topic-cleaner/README.md)
1. [DB writer](db-writer/README.md)
1. [DB reader](db-reader/README.md)
1. [Parquet generator](parquet-generator/README.md)

## Database configuration

It is described in [this file](database/README.md).

## Notes

### Kafkacat basic usage in client mode

```
kafkacat -b localhost:9092 -t test_partitions_16 -C -K\\t
kafkacat -b localhost:9092 -t test_partitions_16 -C -f 'Key: %k\nValue: %s\n'
```
