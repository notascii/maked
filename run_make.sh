#!/bin/bash

kadeploy3 -e ubuntu2204-nfs

# Ensure a MAKEFILE_DIRECTORY argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <MAKEFILE_DIRECTORY>"
  exit 1
fi

# Assign the first argument to MAKEFILE_DIRECTORY
MAKEFILE_DIRECTORY="$1"

# Define the output file for execution time
OUTPUT_FILE="${MAKEFILE_DIRECTORY}_make.txt"

# Check if the OAR_NODEFILE environment variable is set
if [ -z "$OAR_NODEFILE" ]; then
  echo "Error: The OAR_NODEFILE environment variable is not set."
  exit 1
fi

# Read the list of unique nodes from OAR_NODEFILE
NODES=($(sort -u "$OAR_NODEFILE"))

# Ensure there is at least one node in the list
if [ ${#NODES[@]} -eq 0 ]; then
  echo "Error: No nodes found in $OAR_NODEFILE."
  exit 1
fi

# Define the target node as the first node in the list
TARGET_NODE="${NODES[0]}"

# Define the local directory to copy
LOCAL_DIRECTORY="./maked/"

# Define the remote destination directory
REMOTE_DIRECTORY="~/maked/"

# Copy the local directory to the remote node, excluding the .git directory
rsync -av --exclude='.git' "$LOCAL_DIRECTORY" "root@$TARGET_NODE:$REMOTE_DIRECTORY"

# Install make
apt install make


# Execute the make command on the remote node using TakTuk and measure the execution time
{ time taktuk -s -l root -f <(echo "$TARGET_NODE") broadcast exec [ "cd ${REMOTE_DIRECTORY}makefiles/${MAKEFILE_DIRECTORY} && chmod +x * && make" ]; } 2> "$OUTPUT_FILE"

echo "Make command executed on $TARGET_NODE. Execution time recorded in $OUTPUT_FILE."