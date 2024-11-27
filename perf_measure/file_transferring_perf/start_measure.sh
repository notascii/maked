#!/bin/bash

ssh $1 "cd file_transferring_perf  && ./measure_scp_perf.sh $2" &

ssh $1 "cd file_transferring_perf  && ./measure_rsync_perf.sh $2" &
