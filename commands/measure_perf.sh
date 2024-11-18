#!/bin/bash

NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
IFS='+' read -r -a NODE_LIST <<< "$NODES"

HOST1=${NODE_LIST[0]}
HOST2=${NODE_LIST[1]}



# Execute server and client programs
ssh $HOST2 "cd ~/pingpong && go build -o server ./server.go && ./server "

ssh $HOST1 "cd ~/pingpong && go build -o  client ./client.go && ./client $HOST2"

# Wait for execution to finish (add sleep if needed)

# Transfer performance logs back to the local machine
scp $HOST1:~/pingpong/perf/logs/*.log ~/local_perf/

echo "Performance logs have been transferred to ~/local_perf/"
exit 0

