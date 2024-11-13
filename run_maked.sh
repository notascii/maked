#!/bin/bash

# Should be used once you are placed in g5k grenoble and created host

# Environment deployment

kadeploy3 debian11-min

# Copy of files

for node in $(uniq "$OAR_NODEFILE"); do   
    scp ./texte.txt root@"$node":/home/text.txt
    echo "$node"
done