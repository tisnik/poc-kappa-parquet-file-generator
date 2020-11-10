# PoC: Parquet file generator based on Kappa architecture

## Topic constructor

Simple utility to create (construct) new topic.

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
