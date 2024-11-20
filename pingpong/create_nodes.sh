#!/bin/bash
reserve_nodes(){
    oarsub -l host=2 -I
    NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
    IFS='+' read -r -a NODE_LIST <<< "$NODES"
    HOST1=${NODE_LIST[0]}
    HOST2=${NODE_LIST[1]}
    ./install_dependencies
    echo "Reserved nodes: $HOST1, $HOST2"
}

reserve_nodes