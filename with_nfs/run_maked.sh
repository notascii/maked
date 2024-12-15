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
REMOTE_DIRECTORY="~/maked/"

# Makefile directory
MAKEFILE_DIRECTORY="$1"

# Ensure rsync is used for each node to copy the directory
echo "Copying directory to all nodes..."
for node in "${NODES[@]}"; do
  echo "Copying to $node"
  rsync -av --exclude='.git' "$LOCAL_DIRECTORY" "$node:$REMOTE_DIRECTORY"
done

echo "All nodes are set up"

# Start server on the first node
SERVER_NODE="${NODES[0]}"
echo "Starting server on $SERVER_NODE"

# Run the server process in the background
taktuk -s -f <(printf "%s\n" "$SERVER_NODE") broadcast exec [ "export GOROOT=\$HOME/golang/go && export PATH=\$GOROOT/bin:\$PATH && cd ${REMOTE_DIRECTORY}with_nfs/server && mkdir -p server_storage && chmod +x main && nohup go run . ${MAKEFILE_DIRECTORY} > server.log 2>&1 &" ]

echo "Server started on $SERVER_NODE"

# Allow some time for the server to initialize
sleep 5

# Start clients on the remaining nodes
CLIENT_NODES=("${NODES[@]:1}")
echo "Starting clients"

# Calculate the number of client nodes
NUM_CLIENT_NODES=${#CLIENT_NODES[@]}

# Name the output file based on the Makefile directory and the number of nodes
OUTPUT_FILE="${MAKEFILE_DIRECTORY}_${NUM_CLIENT_NODES}_nodes.txt"

rm -rf "${OUTPUT_FILE}"

# Run client processes
{ time taktuk -s -f <(printf "%s\n" "${CLIENT_NODES[@]}") broadcast exec [ "export GOROOT=\$HOME/golang/go && export PATH=\$GOROOT/bin:\$PATH && cd ${REMOTE_DIRECTORY}with_nfs/client && mkdir -p client_storage && go run client.go ${SERVER_NODE}:8090" ]; } 2> "$OUTPUT_FILE"

echo "Ending clients"

# Ensure all background processes complete
wait