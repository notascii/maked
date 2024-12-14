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

# Remote work file
REMOTE_DIRECTORY_WORK_NO_NFS="/tmp/maked/without_nfs/" 

# Makefile directory
MAKEFILE_DIRECTORY="$1"

# Define the number of nodes to test on for each run
NODE_COUNTS=(2 3 4 5 6 7 8 9 10 11)

# Before running each test, ensure that all nodes have the necessary files
echo "Copying directory to all nodes..."
for node in "${NODES[@]}"; do
  echo "Copying to $node"
  rsync -av --exclude='.git' "$LOCAL_DIRECTORY" "$node:$REMOTE_DIRECTORY"
done
echo "All nodes are set up"

# Iterate over each node count you want to test
for COUNT in "${NODE_COUNTS[@]}"; do
  echo "==== Running test with $COUNT nodes ===="
  
  # Select the first $COUNT nodes from NODES
  SELECTED_NODES=("${NODES[@]:0:$COUNT}")

  # The first node is the server
  SERVER_NODE="${SELECTED_NODES[0]}"

  # Adjust the number of client nodes (COUNT - 1 because the first is the server)
  CLIENT_NODE_COUNT=$((COUNT - 1))

  # If we have more than 1 node, the rest are clients
  if [ $COUNT -gt 1 ]; then
    CLIENT_NODES=("${SELECTED_NODES[@]:1}")
  else
    CLIENT_NODES=()  # If we have only one node, no clients
  fi

  # Clean the storage directories on all selected nodes
  echo "Cleaning storage directories on all selected nodes..."
  for node in "${SELECTED_NODES[@]}"; do
    taktuk -s -f <(printf "%s\n" "$node") broadcast exec [ "rm -rf ${REMOTE_DIRECTORY_WORK_NO_NFS}client/client_storage/* ${REMOTE_DIRECTORY_WORK_NO_NFS}server/server_storage/*" ]
  done

  echo "Storage directories cleaned."

  # Start server on the first node
  echo "Starting server on $SERVER_NODE"
  taktuk -s -f <(printf "%s\n" "$SERVER_NODE") broadcast exec [ "export GOROOT=\$HOME/golang/go && export PATH=\$GOROOT/bin:\$PATH && cd ${REMOTE_DIRECTORY_WORK_NO_NFS}server && mkdir -p server_storage && chmod +x main && nohup go run . ${MAKEFILE_DIRECTORY} >> ~/maked/without_nfs/server/server_${CLIENT_NODE_COUNT}_clients.log 2>&1 &" ]
  echo "Server started on $SERVER_NODE"
  
  # Allow some time for the server to initialize
  sleep 5
  
  # Start clients on the remaining nodes
  echo "Starting $CLIENT_NODE_COUNT clients"

  # Name the output file based on the Makefile directory and the number of client nodes
  OUTPUT_FILE="${MAKEFILE_DIRECTORY}_${CLIENT_NODE_COUNT}_clients.txt"

  rm -f "${OUTPUT_FILE}"

  if [ $CLIENT_NODE_COUNT -gt 0 ]; then
    # Run client processes
    { time taktuk -s -f <(printf "%s\n" "${CLIENT_NODES[@]}") broadcast exec [ "export GOROOT=\$HOME/golang/go && export PATH=\$GOROOT/bin:\$PATH && cd ${REMOTE_DIRECTORY_WORK_NO_NFS}client && mkdir -p client_storage && go run client.go ${SERVER_NODE}:8090" ]; } 2> "$OUTPUT_FILE"
    echo "Clients finished for $CLIENT_NODE_COUNT clients"
  else
    echo "No clients to run for single-node test."
  fi

  # Ensure all background processes complete
  wait

done

# After all tests are done, we run the Python script once at the end
cd ./maked
python graph_generator.py
