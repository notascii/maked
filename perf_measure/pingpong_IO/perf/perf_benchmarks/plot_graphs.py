import matplotlib.pyplot as plt
import datetime
# Function to read data from a file
def read_data(filename,metric=""):
    x = []
    y = []
    with open(filename, 'r') as file:
        for line in file:
            parts = line.split(":")
            if len(parts) == 2:
                x.append(float(parts[0].strip()) / pow(2,20)) # Turn message size to MB
                if metric == "latency":
                    y.append(float(parts[1].strip()) * 1000) # Turn latency to micro seconds
                else:
                    y.append(float(parts[1].strip())) # Keep throughput MB/s

    sorted_x,sorted_y=zip(*sorted(zip(x,y)))
    return list(sorted_x),list(sorted_y)

# Function to plot data
def plot_data(x, y, title, xlabel, ylabel, output_file):
    plt.figure(figsize=(10, 6))
    plt.plot(x, y, marker='o', linestyle='-', color='b')
    plt.title(title)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    plt.grid(True)
    plt.savefig(output_file)
    print(f"Plot saved as {output_file}")
    plt.close()

# Main script
if __name__ == "__main__":
    # Read latency data
    message_size, latency = read_data("./perf/logs/latency.log","latency")
    # Turn message size into MB
    # Read throughput data
    _ ,throughput =read_data("./perf/logs/throughput.log")
    
    # Plot latency
    plot_data(
        message_size, latency,
        title="Latency vs Message Size",
        xlabel="Message size (MB)",
        ylabel="Latency (us)",
        output_file="./perf/perf_benchmarks/graphs/"+str(datetime.datetime.now())+"_latency%msg_size.png"
    )

    # Plot throughput
    plot_data(
        message_size, throughput,
        title="Throughput vs Message Size",
        xlabel="Message size (MB)",
        ylabel="Throughput (MB/s)",
        output_file="./perf/perf_benchmarks/graphs/"+str(datetime.datetime.now())+"_throughput%msg_size.png"
    )
