NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
IFS='+' read -r -a NODE_LIST <<< "$NODES"
echo "HERE1"
HOST1=${NODE_LIST[0]}
HOST2=${NODE_LIST[1]}

echo "Reserved nodes: $HOST1, $HOST2

# Install Go on all hosts
for HOST in "${NODE_LIST[@]}"; do
    ssh $HOST "wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz && sudo-g5k tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz && echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc"
done
echo "Go installed"
