import contextlib
import json
import sys
from pathlib import Path
import matplotlib.pyplot as plt

if __name__ == "__main__":
    if len(sys.argv) != 3:
        raise ValueError(
            "Wrong use of this file. Please use `python gen-graph.py -n <name>`"
        )
    name = None
    for index, arg in enumerate(sys.argv):
        if arg == "-n":
            name = sys.argv[index + 1]
    if name is None:
        raise ValueError(
            "Wrong use of this file. Please use `python gen-makefile.py -n <name>`"
        )

    data_path = f"maked/without_nfs/server/json_storage/{name}"
    with contextlib.suppress(FileExistsError):
        path = Path(data_path).glob("**/*.json")
        files = [x for x in path if x.is_file()]
    without_nfs = []
    for file in files:
        with open(file) as f:
            data = json.load(f)
            without_nfs.append((str(file).split("/")[-1].split(".")[0], data))

    data_path = f"maked/with_nfs/server/json_storage/{name}"
    with contextlib.suppress(FileExistsError):
        path = Path(data_path).glob("**/*.json")
        files = [x for x in path if x.is_file()]
    with_nfs = []
    for file in files:
        with open(file) as f:
            data = json.load(f)
            with_nfs.append((str(file).split("/")[-1].split(".")[0], data))

    with_nfs = sorted(with_nfs, key=lambda x: int(x[0]))
    without_nfs = sorted(without_nfs, key=lambda x: int(x[0]))

    # Extract x-values and ensure they are sorted
    x_values = sorted(set([int(elem[0]) for elem in with_nfs] + [int(elem[0]) for elem in without_nfs]))

    # Apply a nicer style
    plt.style.use("seaborn-darkgrid")

    # First figure: Execution times
    fig, ax = plt.subplots(figsize=(10, 6))

    ax.plot(
        [int(elem[0]) for elem in without_nfs],
        [elem[1]["makeDuration"] / 1_000_000 for elem in without_nfs],
        color="tab:red",
        label="Make",
        marker='o'
    )
    ax.plot(
        [int(elem[0]) for elem in with_nfs],
        [elem[1]["makedDuration"] / 1_000_000 for elem in with_nfs],
        color="tab:blue",
        label="Maked (with NFS)",
        marker='o'
    )
    ax.plot(
        [int(elem[0]) for elem in without_nfs],
        [elem[1]["makedDuration"] / 1_000_000 for elem in without_nfs],
        color="tab:orange",
        label="Maked (without NFS)",
        marker='o'
    )

    ax.set_xlabel("Number of nodes", fontsize=12)
    ax.set_ylabel("Execution time (s)", fontsize=12)
    ax.set_title(f"Makefile Execution Times: {name}", fontsize=14)
    ax.set_xticks(x_values)
    ax.legend(fontsize=12)
    ax.grid(True)

    fig.tight_layout()
    plt.savefig(f"maked/without_nfs/server/json_storage/{name}/speed-compare.png")
    plt.close(fig)

    # Second figure: Relative speed increase
    fig, ax = plt.subplots(figsize=(10, 6))

    # Note: Ensure division by zero does not occur. If makedDuration = 0, handle gracefully.
    def relative_speed(data):
        return [
            (d["makeDuration"] / d["makedDuration"] * 100) if d["makedDuration"] != 0 else None
            for d in data
        ]

    ax.plot(
        [int(elem[0]) for elem in with_nfs],
        relative_speed([elem[1] for elem in with_nfs]),
        color="tab:blue",
        label="Maked (with NFS)",
        marker='o'
    )
    ax.plot(
        [int(elem[0]) for elem in without_nfs],
        relative_speed([elem[1] for elem in without_nfs]),
        color="tab:orange",
        label="Maked (without NFS)",
        marker='o'
    )

    ax.set_xlabel("Number of nodes", fontsize=12)
    ax.set_ylabel("Relative speed increase (%)", fontsize=12)
    ax.set_title(f"Makefile Execution Relative Speed Increase: {name}", fontsize=14)
    ax.set_xticks(x_values)
    ax.legend(fontsize=12)
    ax.grid(True)

    fig.tight_layout()
    plt.savefig(f"maked/without_nfs/server/json_storage/{name}/speed-relative.png")
    plt.close(fig)
