#!/bin/bash
# oarsub -I -l host=3,walltime=1:45 -t deploy

kadeploy3 -e ubuntu2204-nfs

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

# Copy the directory and execute commands on each node
for node in "${NODES[@]}"; do
  echo "Processing node: $node"

  # Copy the directory to the node using rsync and exclude the .git directory
  rsync -av --exclude='.git' "$LOCAL_DIRECTORY" "root@$node:$REMOTE_DIRECTORY"

  # Install Go on the node
  ssh root@$node "snap install go --classic"

  echo "Node $node setup complete"
done

echo "All nodes are set up"


# Start the server on the first node and clients on the remaining nodes
SERVER_NODE="${NODES[0]}"
CLIENT_NODES=("${NODES[@]:1}")

# Start the server on the first node
taktuk -f <(printf "%s\n" "$SERVER_NODE") broadcast exec [ "cd ${REMOTE_DIRECTORY}server && mkdir -p server_storage && nohup go run . > server.log 2>&1 &" ]

# Start clients on the remaining nodes
taktuk -f <(printf "%s\n" "${CLIENT_NODES[@]}") broadcast exec [ "cd ${REMOTE_DIRECTORY}client && mkdir -p client_storage && nohup go run client.go ${SERVER_NODE}:8090 > client.log 2>&1 &" ]

echo "All nodes are set up and services are running."

# Wait for all background SSH processes to complete
wait
