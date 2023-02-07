#!/bin/bash

# Use jq to extract the values of the "size" field
values=$(cat ./out/result.json | jq -r '.[]'.size)

# Sum the values
sum=0
for value in $values; do
  sum=$(($sum + $value))
done

# Convert the sum from bytes to gigabytes
factor=1000/1000/1000

# Print the result with 4 decimal places and the unit "gb"
echo $sum /$factor | bc -l | awk '{printf "%.4f", $0} END {print " GB"}'