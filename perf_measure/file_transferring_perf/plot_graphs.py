import matplotlib.pyplot as plt
import datetime
import sys
# Function to read data from a file
def read_data(filename,metric=""):
    x = []
    y = []
    with open(filename, 'r') as file:
        for line in file:
            parts = line.split(":")
            if len(parts) == 2:
                x.append(float(parts[0].strip()) / pow(2,20))
                y.append(float(parts[1].strip()) / 1000) # Turn latency to micro seconds


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
    message_size_scp, latency_scp = read_data("./logs/delay.log")
    message_size_rs, latency_rs = read_data("./logs/delay2.log")
    # Turn message size into MB
    # Read throughput data
       
    # Plot latency
    if sys.argv[1] == "scp":
        plot_data(
            message_size_scp, latency_scp,
            title="Time delay vs Message Size using scp",
            xlabel="Message size (MB)",
            ylabel="time delay (us)",
            output_file="./graphs/"+str(datetime.datetime.now())+"_scp_delay%msg_size.png"
        )
    else:
        plot_data(
            message_size_rs, latency_rs,
            title="Time delay vs Message Size using rsync",
            xlabel="Message size (MB)",
            ylabel="time delay (us)",
            output_file="./graphs/"+str(datetime.datetime.now())+"_rsync_delay%msg_size.png"
        )
