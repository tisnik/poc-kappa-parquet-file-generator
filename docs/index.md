---
layout: default
---
# Description

A proof of concept of Parquet files generator. Data to be used are consumed
from selected Kafka topic and PostgreSQL is used as a cache for such data.

## Architecture

### A classic Kappa architecture utilizing one Kafka partition

![kappa_single_partition.png](kappa_single_partition.png)

### A classic Kappa architecture utilizing multiple Kafka partitions

![kappa_multiple_partitions.png](kappa_multiple_partitions.png)

### PostgreSQL used as a cache, not scaled up

![db-writers-single-generator.png](db-writers-single-generator.png)


### PostgreSQL used as a cache, scaled up

![db-writers-multiple-generators.png](db-writers-multiple-generators.png)

## Benchmark results

### Cummulative time for one writer to PostgreSQL database

![cummulative-time-1-writer.png](cummulative-time-1-writer.png)

### Cummulative time for 16 writers to PostgreSQL database

![cummulative-time-16-writers.png](cummulative-time-16-writers.png)
