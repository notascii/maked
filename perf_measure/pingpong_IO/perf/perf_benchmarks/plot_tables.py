import matplotlib.pyplot as plt
from datetime import datetime

def parse_latency_file(latency_file):
    """Parse the latency file and extract the minimum latency based on message size."""
    min_latency = None  # Start with infinity to find the minimum
    min_message_size = float('inf')

    with open(latency_file, 'r') as file:
        for line in file:
            parts = line.split(":")
            if len(parts) == 2:
                message_size = int(parts[0].strip())
                latency = float(parts[1].strip()) * 1e3  # Convert seconds to microseconds
                if message_size < min_message_size:
                    min_latency = latency
                    min_message_size = message_size

    return min_latency


def parse_throughput_file(throughput_file):
    """Parse the throughput file and extract the maximum throughput based on message size."""
    max_throughput = None  # Start with negative infinity to find the maximum
    max_message_size = float('-inf')

    with open(throughput_file, 'r') as file:
        for line in file:
            parts = line.split(":")
            if len(parts) == 2:
                message_size = int(parts[0].strip())
                throughput = float(parts[1].strip()) / 1e3  # Convert MB/s to GB/s
                if message_size > max_message_size:
                    max_throughput = throughput
                    max_message_size = message_size

    return max_throughput 

def draw_table(data,outputfile):
    """Draw a table using matplotlib."""
    fig, ax = plt.subplots(figsize=(5, 2))
    ax.axis('tight')
    ax.axis('off')
    # Draw the table
    table = plt.table(
        cellText=data,
        colLabels=["Time Measured", "Latency (µs)", "Throughput (GB/s)"],
        loc="center",
        cellLoc="center"
    )
    table.auto_set_font_size(False)
    table.set_fontsize(10)
    table.auto_set_column_width(col=list(range(len(data[0]))))

    plt.savefig(outputfile, dpi=300, bbox_inches='tight')
    print(f"Table saved as {outputfile}")
    plt.close()

def main(latency_file, throughput_file):
    """Main function to process files and display the table."""
    # Extract required data
    latency = parse_latency_file(latency_file)
    throughput = parse_throughput_file(throughput_file)
    
    if latency is None or throughput is None:
        print("Failed to extract data. Please check the input files.")
        return
    
    # Prepare data for the table
    data = [[str(datetime.now()), f"{latency:.2f}", f"{throughput:.3f}"]]
    
    # Draw the table
    draw_table(data,"./perf/perf_benchmarks/tables/"+str(datetime.now())+"_perf_table.png")

if __name__ == "__main__":
    # Input files
    latency_file = "./perf/logs/latency.log"
    throughput_file = "./perf/logs/throughput.log"
    
    main(latency_file, throughput_file)
