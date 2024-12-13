#!/bin/bash

NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
IFS='+' read -r -a NODE_LIST <<< "$NODES"

HOST1=${NODE_LIST[0]}
HOST2=${NODE_LIST[1]}

./start_measure.sh $HOST1 $HOST2
