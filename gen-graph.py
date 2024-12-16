import contextlib
import json
import sys
from pathlib import Path

import matplotlib.pyplot as plt

if __name__ == "__main__":
    if len(sys.argv) != 3:
        raise ValueError(
            "Wrong use of this file. Please use `python gen-graph.py -n <name>"
        )
    name = None
    for index, arg in enumerate(sys.argv):
        if arg == "-n":
            name = sys.argv[index + 1]
    if name is None:
        raise ValueError(
            "Wrong use of this file. Please use `python gen-makefile.py -n <name>"
        )

    data_path = f"without_nfs/server/json_storage/{name}"
    with contextlib.suppress(FileExistsError):
        path = Path(data_path).glob("**/*")
        files = [x for x in path if x.is_file()]
    without_nfs = {}
    for file in files:
        with open(file) as f:
            data = json.load(f)
            without_nfs[str(file).split("/")[-1].split(".")[0]] = data

    data_path = f"with_nfs/server/json_storage/{name}"
    with contextlib.suppress(FileExistsError):
        path = Path(data_path).glob("**/*")
        files = [x for x in path if x.is_file()]
    with_nfs = {}
    for file in files:
        with open(file) as f:
            data = json.load(f)
            with_nfs[str(file).split("/")[-1].split(".")[0]] = data

    fig = plt.figure()
    ax = fig.add_subplot(1, 1, 1)
    ax.plot(
        [int(key) for key in with_nfs],
        [value["makedDuration"] / 1_000_000 for value in with_nfs.values()],
        color="tab:blue",
        label="Maked NFS",
    )
    ax.plot(
        [int(key) for key in with_nfs],
        [value["makeDuration"] / 1_000_000 for value in with_nfs.values()],
        color="tab:red",
        label="Make",
    )
    ax.plot(
        [int(key) for key in without_nfs],
        [value["makedDuration"] / 1_000_000 for value in without_nfs.values()],
        color="tab:orange",
        label="Maked without NFS",
    )
    plt.xlabel("Number of nodes")
    plt.ylabel("Execution time (s)")
    plt.legend(loc="upper right")
    plt.title(f"Makefile execution times: {name}")
    plt.savefig("tmp.png")
