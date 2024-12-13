#!/bin/bash

ssh $1 "cd file_transferring_perf  && ./measure_scp_perf.sh $2" 
