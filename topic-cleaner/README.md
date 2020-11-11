# PoC: Parquet file generator based on Kappa architecture

## Topic cleaner

Simple utility to clean (delete) selected topic.

### Description

This utility can be used to clean (i.e. delete) selected topic from Kafka
broker. Broker address and topic name needs to be specified on command line.

### How to build the tool

```
go build
```

### Usage

```
Usage of topic-cleaner:
  -broker string
        broker address (default "localhost:9092")
  -topic string
        topic name
```

Please note that topic name is required.
