#!/bin/bash

SOURCE_FOLDER="../pingpong_IO/client/disk"
LOG_FILE="./logs/delay2.log"

> "$LOG_FILE"

# Loop through each file in the folder
for file in "$SOURCE_FOLDER"/*; do
    if [ -f "$file" ]; then
    
        file_size=$(stat --printf="%s" "$file")
        
        # Get the start time
        start_time=$(date +%s%N)
        
        # Transfer the file using scp
        rsync -av "$file" $1        
        # Get the end time
        end_time=$(date +%s%N)
        
        # Calculate latency in seconds
        latency=$(echo "scale=3; ($end_time - $start_time) / 1000" | bc)
        
        # Log size and latency
        echo "$file_size:$latency" >> "$LOG_FILE"
    fi
done

echo "Transfer completed. Log file saved as $LOG_FILE."

python_file="./plot_graphs.py"

python3 "$python_file" rsync
