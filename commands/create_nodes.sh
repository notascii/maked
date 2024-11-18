#!/bin/bash
oarsub -l host=2 -I
NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
IFS='+' read -r -a NODE_LIST <<< "$NODES"
echo "HERE1"
HOST1=${NODE_LIST[0]}
HOST2=${NODE_LIST[1]}

echo "Reserved nodes: $HOST1, $HOST2"
