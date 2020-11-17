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
import numpy as np
import matplotlib.pyplot as plt

# Check if command line argument is specified (it is mandatory).
if len(sys.argv) < 2:
    print("Usage:")
    print("  db-writer.py input_file.csv")
    print("Example:")
    print("  db-writer.py db-writer-0-0.csv")
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
    durations = [float(row[0]) for row in csv_reader]
    print(durations)

# Linear regression
x = np.arange(0, len(durations))
coef = np.polyfit(x, durations, 1)
poly1d_fn = np.poly1d(coef)

# Create new histogram graph
plt.plot(x, durations, "b", poly1d_fn(np.arange(0, len(durations))), 'y--')

# Title of a graph
plt.title("Duration of writing records into database")

# Add a label to x-axis
plt.xlabel("Record #")

# Add a label to y-axis
plt.ylabel("Cummulative time [ms]")

# Set the plot layout
plt.tight_layout()

# And save the plot into raster format and vector format as well
plt.savefig("db-writer-cummulative-time.png")
plt.savefig("db-writer-cummulative-time.svg")

# Try to show the plot on screen
plt.show()
