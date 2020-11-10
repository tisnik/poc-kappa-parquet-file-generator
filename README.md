# poc-kappa-parquet-file-generator
Parquet file generator and pre-worker

## Tools available in this repository

1. [Topic cleaner](topic-cleaner/README.md)

## Notes

### Kafkacat basic usage in client mode

```
kafkacat -b localhost:9092 -t test_partitions_16 -C -K\\t
kafkacat -b localhost:9092 -t test_partitions_16 -C -f 'Key: %k\nValue: %s\n'
```
