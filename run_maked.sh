#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <MAKEFILE_DIRECTORY>"
  exit 1
fi

# Check if the OAR_NODEFILE environment variable is set
if [ -z "$OAR_NODEFILE" ]; then
  echo "Error: The OAR_NODEFILE environment variable is not set."
  exit 1
fi

# Read the list of unique nodes from OAR_NODEFILE
NODES=($(sort -u "$OAR_NODEFILE"))

# Local directory to copy
LOCAL_DIRECTORY="./maked/"

# Remote destination directory
REMOTE_DIRECTORY="/tmp/maked/"

# Makefile directory
MAKEFILE_DIRECTORY="$1"


# Copy the directory and execute commands on each node
for node in "${NODES[@]}"; do
  echo "Processing node: $node"

  # Copy the directory to the node using rsync and exclude the .git directory
  rsync -av --exclude='.git' "$LOCAL_DIRECTORY" "$node:$REMOTE_DIRECTORY"

  echo "Node $node setup complete"
done

echo "All nodes are set up"

# Start server on the first node
SERVER_NODE="${NODES[0]}"
echo "Starting server on $SERVER_NODE"
ssh $SERVER_NODE "cd ${REMOTE_DIRECTORY}server && mkdir -p server_storage && chmod +x main && nohup ./main ${MAKEFILE_DIRECTORY} > server.log 2>&1 &" &
echo "Server started on $SERVER_NODE"

# Start clients on the remaining nodes
CLIENT_NODES=("${NODES[@]:1}")
echo "Starting clients"

# Calculate the number of client nodes
NUM_CLIENT_NODES=${#CLIENT_NODES[@]}

# Name the output file based on the Makefile directory and the number of nodes
OUTPUT_FILE="${MAKEFILE_DIRECTORY}_${NUM_CLIENT_NODES}_nodes.txt"

rm -rf "${OUTPUT_FILE}"

{ time taktuk -s -l root -f <(printf "%s\n" "${CLIENT_NODES[@]}") broadcast exec [ "cd ${REMOTE_DIRECTORY}client && mkdir -p client_storage && chmod +x client && ./client ${SERVER_NODE}:8090" ]; } 2> "$OUTPUT_FILE"

echo "Ending clients"

# Wait for all background SSH processes to complete
wait
