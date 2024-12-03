import os
import re
import matplotlib.pyplot as plt
from collections import defaultdict

# Function to extract real time from a file
def extract_real_time(file_path):
    with open(file_path, 'r') as file:
        for line in file:
            if line.startswith("real"):
                # Extract and convert real time to seconds
                time_parts = line.split()[1].replace(",", ".").split("m")
                minutes = float(time_parts[0])  # Extract minutes
                seconds = float(time_parts[1].replace("s", ""))  # Remove 's' and convert to float
                return minutes * 60 + seconds
    return None  # Return None if no real time is found

# Main function
def main():
    current_directory = os.getcwd()
    pattern = re.compile(r'(.+?)_(\d+)_nodes\.txt')
    data = defaultdict(list)  # Structure: {MAKEFILE_DIRECTORY: [(NUM_CLIENT_NODES, real_time)]}

    # Scan files in the current directory
    for file_name in os.listdir(current_directory):
        match = pattern.match(file_name)
        if match:
            makefile_directory, num_nodes = match.groups()
            num_nodes = int(num_nodes)
            file_path = os.path.join(current_directory, file_name)
            real_time = extract_real_time(file_path)
            if real_time is not None:
                data[makefile_directory].append((num_nodes, real_time))

    # Generate plots for each MAKEFILE_DIRECTORY
    for makefile_directory, values in data.items():
        # Sort values by the number of nodes
        values.sort(key=lambda x: x[0])
        nodes, real_times = zip(*values)

        # Create the plot
        plt.figure()
        plt.plot(nodes, real_times, marker='o', linestyle='-', label='Real Time')
        plt.title(f"Execution Time for {makefile_directory}")
        plt.xlabel("Number of Nodes")
        plt.ylabel("Execution Time (seconds)")
        plt.grid(True)
        plt.legend()
        plt.savefig(f"{makefile_directory}_execution_time.png")
        plt.show()

if __name__ == "__main__":
    main()
