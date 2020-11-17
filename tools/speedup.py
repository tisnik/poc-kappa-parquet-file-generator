#!/usr/bin/env python3

# Copyright © 2020 Pavel Tisnovsky
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import sys
import csv
import matplotlib.pyplot as plt

# Check if command line argument is specified (it is mandatory).
if len(sys.argv) < 2:
    print("Usage:")
    print("  speedup.py input_file.csv")
    print("Example:")
    print("  speedup.py overall.csv")
    sys.exit(1)

# First command line argument should contain name of input CSV.
input_csv = sys.argv[1]

# Try to open the CSV file specified.
with open(input_csv) as csv_input:
    # And open this file as CSV
    csv_reader = csv.reader(csv_input)

    # Skip header
    next(csv_reader, None)
    rows = 0

    # Read all rows from the provided CSV file
    # and read just first digits from cluster IDs
    data = [(int(row[0]), float(row[1]), float(row[2]), float(row[3]), float(row[4]), float(row[5])) for row in csv_reader]
    print(data)

# Linear regression
topics = [item[0] for item in data]
producing = [item[1] for item in data]
db_write = [item[2] for item in data]
parquet_none = [item[3] for item in data]
parquet_snappy = [item[4] for item in data]
parquet_gzip = [item[5] for item in data]

# Create new histogram graph
plt.plot(topics, producing, label='write messages to Kafka')
plt.plot(topics, db_write, label='read from Kafka, write to DB')
plt.plot(topics, parquet_none, label='create Parquet, None compression')
plt.plot(topics, parquet_snappy, label='create Parquet, Snappy compression')
plt.plot(topics, parquet_gzip, "b", label='create Parquet, GZIP compression')

# Title of a graph
plt.title("Duration of all operations based on number of topics")

# Add a label to x-axis
plt.xlabel("Topics")

# Add a label to y-axis
plt.ylabel("Time [s]")

plt.legend(loc="upper right")

# Set the plot layout
plt.tight_layout()

# And save the plot into raster format and vector format as well
plt.savefig("speedup.png")
plt.savefig("speedup.svg")

# Try to show the plot on screen
plt.show()
