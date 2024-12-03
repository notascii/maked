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
    efficiency_data = defaultdict(list)  # Structure: {prefix: [(num_nodes, efficiency_ratio)]}

    # Identify unique prefixes and process files
    prefixes = set()
    for file_name in os.listdir(current_directory):
        match = pattern.match(file_name)
        if match:
            prefixes.add(match.group(1))

    for prefix in prefixes:
        # Extract real time from the corresponding makefile
        makefile_path = os.path.join(current_directory, f"{prefix}_make.txt")
        makefile_time = extract_real_time(makefile_path)
        if makefile_time is None:
            print(f"Error: Could not find or parse '{prefix}_make.txt'. Skipping prefix.")
            continue

        # Process files matching the current prefix
        for file_name in os.listdir(current_directory):
            match = pattern.match(file_name)
            if match and match.group(1) == prefix:
                num_nodes = int(match.group(2))
                file_path = os.path.join(current_directory, file_name)
                real_time = extract_real_time(file_path)
                if real_time is not None:
                    efficiency_ratio = makefile_time / real_time
                    efficiency_data[prefix].append((num_nodes, efficiency_ratio))

    # Generate plots for each prefix
    for prefix, values in efficiency_data.items():
        # Sort data by number of nodes
        values.sort(key=lambda x: x[0])
        nodes, efficiency_ratios = zip(*values)

        # Create the efficiency ratio plot
        plt.figure()
        plt.plot(nodes, efficiency_ratios, marker='o', linestyle='-', label=f'{prefix} Efficiency Ratio')
        plt.title(f"Efficiency Ratio for {prefix} Files")
        plt.xlabel("Number of Nodes")
        plt.ylabel("Efficiency Ratio (Makefile / Nodes)")
        plt.grid(True)
        plt.legend()
        plt.savefig(f"{prefix}_efficiency_ratio.png")
        plt.show()

if __name__ == "__main__":
    main()
