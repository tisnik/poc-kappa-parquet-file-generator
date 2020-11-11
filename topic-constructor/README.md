# PoC: Parquet file generator based on Kappa architecture

## Topic constructor

Simple utility to create (construct) new topic.

### Description

This utility can be used to create new topic in Kafka broker. Broker address
and topic name needs to be specified on command line. Additionally it is
possible to specify number of partitions and replication factor for the new
topic.

### How to build the tool

```
go build
```

### Usage

```
Usage of topic-constructor:
  -broker string
        broker address (default "localhost:9092")
  -partitions int
        number of partitions (default 1)
  -replication-factor int
        replication factor (default 1)
  -topic string
        topic name

```

Please note that topic name is required.
