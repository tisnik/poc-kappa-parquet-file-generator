#!/bin/bash

./parquet-generator --db-name d0 > 0.txt 2>&1 &
./parquet-generator --db-name d1 > 1.txt 2>&1 &
./parquet-generator --db-name d2 > 2.txt 2>&1 &
./parquet-generator --db-name d3 > 3.txt 2>&1 &
./parquet-generator --db-name d4 > 4.txt 2>&1 &
./parquet-generator --db-name d5 > 5.txt 2>&1 &
./parquet-generator --db-name d6 > 6.txt 2>&1 &
./parquet-generator --db-name d7 > 7.txt 2>&1 &
./parquet-generator --db-name d8 > 8.txt 2>&1 &
./parquet-generator --db-name d9 > 9.txt 2>&1 &
./parquet-generator --db-name da > a.txt 2>&1 &
./parquet-generator --db-name db > b.txt 2>&1 &
./parquet-generator --db-name dc > c.txt 2>&1 &
./parquet-generator --db-name dd > d.txt 2>&1 &
./parquet-generator --db-name de > e.txt 2>&1 &
./parquet-generator --db-name df > f.txt 2>&1 &
