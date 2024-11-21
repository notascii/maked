NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
IFS='+' read -r -a NODE_LIST <<< "$NODES"
HOST1=${NODE_LIST[0]}
HOST2=${NODE_LIST[1]}

echo "Reserved nodes: $HOST1, $HOST2"

echo " Install Go on all hosts"
for HOST in "${NODE_LIST[@]}"; do
    ssh $HOST "wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz && sudo-g5k tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz && echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc"
done
echo "Go installed"

install_python_libraries() {
    echo "Installing Python libraries..."
    requirements_file="./perf/perf_benchmarks/requirements.txt"
    if [[ -f "$requirements_file" ]]; then
        pip3 install -r "$requirements_file"
        echo "Python libraries installed successfully."
    else
        echo "Requirements file not found: $requirements_file"
        exit 1
    fi
}

install_python_libraries
