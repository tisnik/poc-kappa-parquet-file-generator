# PoC: Parquet file generator based on Kappa architecture

## DB writer

Tool to consume messages from Kafka topic and store them in relational
database, PostgreSQL in the default variant.

### Description

This tool consumes messages from selected Kafka topic and partition. Such
messages are unmarshalled, parsed, and stored into PostgreSQL database. It is
possible to use command line flags to select the database, Kafka topic, and
Kafka partition.

### How to build the tool

```
go build
```

### Usage

```
Usage of /tmp/ramdisk/go-build323414172/b001/exe/db-writer:
  -broker string
        broker address (default "localhost:9092")
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
  -partition int
        partition to consume (default -1)
  -topic string
        topic name

```

### PostgreSQL table structure

Currently, the simplified SQL table structure is used:

```
CREATE TABLE reports (
    id                    SERIAL PRIMARY KEY,
    key                   CHAR(1) NOT NULL,
    cluster_id            CHAR(36) NOT NULL,
    path                  VARCHAR(200) NOT NULL,
    external_organization VARCHAR(20) NOT NULL,
    report                TEXT NOT NULL
);
```
