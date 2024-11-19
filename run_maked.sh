#!/bin/bash
# oarsub -I -l host=2,walltime=1:45 -t deploy

# kadeploy3 -e ubuntu2204-nfs

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

# Name of the Go installation script
INSTALL_GO_SCRIPT="install_go.sh"

# Copy the directory and execute commands on each node
for i in "${!NODES[@]}"; do
  node="${NODES[$i]}"
  echo "Processing node: $node"

  # Copy the directory to the node
  scp -r "$LOCAL_DIRECTORY" "root@$node:$REMOTE_DIRECTORY"

  # Set execute permissions and run the Go installation script
  ssh root@$node "chmod +x ${REMOTE_DIRECTORY}${INSTALL_GO_SCRIPT} && ${REMOTE_DIRECTORY}${INSTALL_GO_SCRIPT}"

  # Determine the command to run based on node index
  if [ "$i" -eq 0 ]; then
    # First node: start the server
    echo "Starting server on $node"
    ssh root@$node "cd ${REMOTE_DIRECTORY}server && nohup go run . > server.log 2>&1 &" &
    echo "Server started on $node"
  else
    # Other nodes: start the client
    echo "Starting client on $node"
    ssh root@$node "cd ${REMOTE_DIRECTORY}client && nohup go run client.go ${NODES[0]}:8090 > client.log 2>&1 &" &
    echo "Client started on $node"
  fi

  echo "Node $node setup complete"
done

# Wait for all background SSH processes to complete
wait
