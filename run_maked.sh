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

echo "" >~/.ssh/known_hosts

# Read the list of unique nodes from OAR_NODEFILE
NODES=($(sort -u "$OAR_NODEFILE"))

# Directories and variables
LOCAL_DIRECTORY="./maked/"
REMOTE_DIRECTORY="/tmp/maked/" # Used for the without_nfs scenario
MAKEFILE_DIRECTORY="$1"

# Delete files in specified directories
echo "Cleaning up specified directories..."
rm -rf ./maked/without_nfs/server/json_storage/"$MAKEFILE_DIRECTORY"/* ./maked/without_nfs/server/*.log
rm -rf ./maked/with_nfs/server/json_storage/"$MAKEFILE_DIRECTORY"/* ./maked/with_nfs/server/*.log
rm -rf ./make/with_nfs/commun_storage/*
echo "Directories cleaned."

# Dynamically define the number of nodes to test on for each run (2 to total number of nodes)
NODE_COUNTS=($(seq 2 ${#NODES[@]}))

# Sync the local directory to all nodes for the without_nfs scenario
echo "Copying directory to all nodes for without_nfs scenario..."
for node in "${NODES[@]}"; do
  echo "Copying to $node"
  rsync -av --exclude='.git' "$LOCAL_DIRECTORY" "$node:$REMOTE_DIRECTORY"
done
echo "All nodes are set up for the without_nfs scenario."

# Function to run tests for a given scenario (without_nfs or with_nfs)
run_tests_for_directory() {
  LOCAL_TEST_DIRECTORY="$1" # "without_nfs" or "with_nfs"

  if [ "$LOCAL_TEST_DIRECTORY" = "without_nfs" ]; then
    TEST_WORK_DIR="${REMOTE_DIRECTORY}without_nfs/"
  else
    TEST_WORK_DIR="./maked/with_nfs/"
  fi

  echo "=== Running tests for ${LOCAL_TEST_DIRECTORY} ==="

  FIRST_RUN=1  # Initialize the first run flag for the loop

  for COUNT in "${NODE_COUNTS[@]}"; do
    echo "==== Running test with $COUNT nodes for ${LOCAL_TEST_DIRECTORY} ===="

    SELECTED_NODES=("${NODES[@]:0:$COUNT}")
    SERVER_NODE="${SELECTED_NODES[0]}"
    CLIENT_NODE_COUNT=$((COUNT - 1))

    if [ $COUNT -gt 1 ]; then
      CLIENT_NODES=("${SELECTED_NODES[@]:1}")
    else
      CLIENT_NODES=()
    fi

    echo "Cleaning storage directories on all selected nodes..."
    for node in "${SELECTED_NODES[@]}"; do
      taktuk -s -f <(printf "%s\n" "$node") broadcast exec [ "rm -rf ${TEST_WORK_DIR}client/client_storage/* ${TEST_WORK_DIR}server/server_storage/*" ]
    done
    echo "Storage directories cleaned."

    echo "Starting server on $SERVER_NODE in $LOCAL_TEST_DIRECTORY"
    taktuk -s -f <(printf "%s\n" "$SERVER_NODE") broadcast exec [ "export GOROOT=\$HOME/golang/go && export PATH=\$GOROOT/bin:\$PATH && cd ${TEST_WORK_DIR}server && mkdir -p server_storage ~/maked/${LOCAL_TEST_DIRECTORY}/server/json_storage/${MAKEFILE_DIRECTORY} && chmod +x main && nohup go run . ${MAKEFILE_DIRECTORY} $FIRST_RUN > ~/maked/${LOCAL_TEST_DIRECTORY}/server/json_storage/${MAKEFILE_DIRECTORY}/server_${CLIENT_NODE_COUNT}_clients.log 2>&1 &" ]
    echo "Server started on $SERVER_NODE"

    # Allow some time for the server to initialize
    sleep 5

    echo "Starting $CLIENT_NODE_COUNT clients"
    OUTPUT_FILE="${MAKEFILE_DIRECTORY}_${CLIENT_NODE_COUNT}_clients_${LOCAL_TEST_DIRECTORY}.txt"
    rm -f "${OUTPUT_FILE}"

    if [ $CLIENT_NODE_COUNT -gt 0 ]; then
      { taktuk -s -f <(printf "%s\n" "${CLIENT_NODES[@]}") broadcast exec [ "export GOROOT=\$HOME/golang/go && export PATH=\$GOROOT/bin:\$PATH && cd ${TEST_WORK_DIR}client && mkdir -p client_storage && go run client.go ${SERVER_NODE}:8090" ]; }
      echo "Clients finished for $CLIENT_NODE_COUNT clients in $LOCAL_TEST_DIRECTORY"
    else
      echo "No clients to run for single-node test in $LOCAL_TEST_DIRECTORY."
    fi

    # Ensure all background processes complete
    wait
    python maked/gen-graph.py -n $MAKEFILE_DIRECTORY

    # Update FIRST_RUN after the first iteration
    FIRST_RUN=0

  done

  echo "=== Finished tests for ${LOCAL_TEST_DIRECTORY} ==="
}

run_tests_for_directory "without_nfs"
run_tests_for_directory "with_nfs"

cd ./maked
python graph_generator.py
