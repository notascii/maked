#!/bin/bash

# Directory to store the files
output_dir="./client/disk"

# Create the directory if it doesn't exist
mkdir -p "$output_dir"

# Initial size in MB
size_mb=100

# Loop to create 10 files
for i in {1..10}; do
    # Calculate the file size in bytes
    size_bytes=$((size_mb * 1024 * 1024))

    # File name
    filename="$output_dir/file_${size_mb}MB.bin"

    # Create the file with the specified size
    dd if=/dev/zero of="$filename" bs=1M count="$size_mb" status=none

    echo "Created file: $filename with size: ${size_mb} MB"

    # Increment size by 100 MB for the next file
    size_mb=$((size_mb + 100))
done

echo "All files created in directory: $output_dir"
