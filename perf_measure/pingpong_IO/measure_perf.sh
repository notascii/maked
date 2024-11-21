#!/bin/bash


# Main function to orchestrate all tasks
main() {
    NODES=$(oarstat -u $USER -f | grep assigned_hostnames | sed 's/assigned_hostnames = //g')
    IFS='+' read -r -a NODE_LIST <<< "$NODES"

    HOST1=${NODE_LIST[0]}
    HOST2=${NODE_LIST[1]}
    ssh $HOST1 "cd ~/pingpong_IO && ./run_server.sh &"
    ssh $HOST2  "cd ~/pingpong_IO && ./run_client.sh $HOST1 &"

    scp $HOST2:~/pingpong_IO/perf/logs/*.log ~/perf/local_perf/
    scp $HOST2:~/pingpong_IO/perf/perf_benchmarks/graphs/*.png ~/perf/perf_plots/
}

# Execute the main function
main
