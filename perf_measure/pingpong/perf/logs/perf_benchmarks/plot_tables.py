import matplotlib.pyplot as plt
from datetime import datetime
def parse_latency_file(latency_file):
    """Parse the latency file and extract the latency for messageSize of 1."""
    with open(latency_file, 'r') as file:
        for line in file:
            parts = line.split(":")
            if len(parts) == 2:
                message_size = int(parts[0].strip())
                if message_size == 1:
                    latency = float(parts[1].strip()) * 1e3  # Convert seconds to microseconds
                    return latency
    return None

def parse_throughput_file(throughput_file):
    """Parse the throughput file and extract the last row."""
    with open(throughput_file, 'r') as file:
        lines = file.readlines()
        if lines:
            last_line = lines[-1]
            parts = last_line.split(":")
            if len(parts) == 2:
                throughput = float(parts[1].strip()) / 1e3  # Convert MB/s to GB/s
                return throughput
    return None

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
