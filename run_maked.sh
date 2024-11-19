#!/bin/bash
# oarsub -I -l host=10,walltime=1:45 -t deploy

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

# Start server on the first node and clients on the remaining nodes
for i in "${!NODES[@]}"; do
  node="${NODES[$i]}"
  if [ "$i" -eq 0 ]; then
    # First node: start the server
    echo "Starting server on $node"
    ssh root@$node "cd ${REMOTE_DIRECTORY}server && nohup go run . > server.log 2>&1 &" &
    echo "Server started on $node"
  else
    # Other nodes: start the client
    echo "Trying to connect client on $node to server ${NODES[0]}:8090"
    ssh root@$node "cd ${REMOTE_DIRECTORY}client && nohup go run client.go ${NODES[0]}:8090 > client.log 2>&1 &" &
    echo "Client started on $node"
  fi
done

# Wait for all background SSH processes to complete
wait
