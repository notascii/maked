#!/bin/bash

# Used when API doesn't work

if [ -z "$1" ]; then
  echo "Usage: $0 <RUN_MAKED_SCRIPT>"
  exit 1
fi

RUN_MAKED_SCRIPT="$1"

# Ensure the script to be executed exists
if [ ! -f "$RUN_MAKED_SCRIPT" ]; then
  echo "Error: Script $RUN_MAKED_SCRIPT not found."
  exit 1
fi

# Array representing the number of nodes to reserve (1 to 10)
NODE_COUNTS=(1 2 3 4 5 6 7 8 9 10)

OUTPUT_DIR="output_logs"
mkdir -p "$OUTPUT_DIR"

# Iterate over each node count
for NODE_COUNT in "${NODE_COUNTS[@]}"; do
  echo "Reserving $NODE_COUNT node(s) and launching script..."

  # Reserve nodes and directly launch the script
  oarsub -l nodes=$NODE_COUNT,walltime=01:00:00 -S "$RUN_MAKED_SCRIPT" > "$OUTPUT_DIR/reservation_${NODE_COUNT}.log" 2>&1

  if [ $? -ne 0 ]; then
    echo "Error reserving nodes or executing script for $NODE_COUNT nodes. Skipping."
    continue
  fi

  echo "Execution completed for $NODE_COUNT node(s). Logs saved in $OUTPUT_DIR."
done

echo "All node reservations and executions completed."