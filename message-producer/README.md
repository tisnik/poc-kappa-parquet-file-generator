# PoC: Parquet file generator based on Kappa architecture

## Message producer

Produce messages read and parsed from provided input file.

### Description

This utility reads and parses the input file that should contains message
contents split by newline character. Each message is represented as JSON (ie
the input file consists of several JSONs, one JSON per line). This JSON is
de-serialized and cluster ID is used to choose the `key` for message to be
written into Kafka. Message `value` contains the whole JSON.

### How to build the tool

```
go build
```

### Usage

```
Usage of message-producer:
  -broker string
        broker address (default "localhost:9092")
  -input string
        input data file
  -topic string
        topic name

```
