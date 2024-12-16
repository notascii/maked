import contextlib
import json
import sys
from pathlib import Path

import matplotlib.pyplot as plt

# import pandas as pd
# import seaborn as sns

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
    without_nfs = []
    for file in files:
        with open(file) as f:
            data = json.load(f)
            without_nfs.append((str(file).split("/")[-1].split(".")[0], data))

    data_path = f"with_nfs/server/json_storage/{name}"
    with contextlib.suppress(FileExistsError):
        path = Path(data_path).glob("**/*")
        files = [x for x in path if x.is_file()]
    with_nfs = []
    for file in files:
        with open(file) as f:
            data = json.load(f)
            with_nfs.append((str(file).split("/")[-1].split(".")[0], data))

    with_nfs = sorted(with_nfs, key=lambda x: x[0])
    without_nfs = sorted(without_nfs, key=lambda x: x[0])

    fig = plt.figure()
    ax = fig.add_subplot(1, 1, 1)

    ax.plot(
        [int(elem[0]) for elem in with_nfs],
        [elem[1]["makeDuration"] / 1_000_000 for elem in with_nfs],
        color="tab:red",
        label="Make",
    )
    ax.plot(
        [int(elem[0]) for elem in with_nfs],
        [elem[1]["makedDuration"] / 1_000_000 for elem in with_nfs],
        color="tab:blue",
        label="Maked NFS",
    )
    ax.plot(
        [int(elem[0]) for elem in without_nfs],
        [elem[1]["makedDuration"] / 1_000_000 for elem in without_nfs],
        color="tab:orange",
        label="Maked without NFS",
    )
    plt.xlabel("Number of nodes")
    plt.ylabel("Execution time (s)")
    plt.legend(loc="upper right")
    plt.title(f"Makefile execution times: {name}")
    plt.xticks(
        list(
            set(
                [
                    *[int(elem[0]) for elem in with_nfs],
                    *[int(elem[0]) for elem in without_nfs],
                ]
            )
        )
    )
    plt.savefig(f"without_nfs/server/json_storage/{name}/speed-compare.png")
    plt.savefig(f"with_nfs/server/json_storage/{name}/speed-compare.png")
    plt.close()
    fig = plt.figure()
    ax = fig.add_subplot(1, 1, 1)

    ax.plot(
        [int(elem[0]) for elem in with_nfs],
        [elem[1]["makeDuration"] / elem[1]["makedDuration"] * 100 for elem in with_nfs],
        color="tab:blue",
        label="Maked NFS",
    )
    ax.plot(
        [int(elem[0]) for elem in without_nfs],
        [
            elem[1]["makeDuration"] / elem[1]["makedDuration"] * 100
            for elem in without_nfs
        ],
        color="tab:orange",
        label="Maked without NFS",
    )
    plt.xlabel("Number of nodes")
    plt.ylabel("Relative speed increase (%)")
    plt.legend(loc="upper right")
    plt.title(f"Makefile execution relative speed increase: {name}")
    plt.xticks(
        list(
            set(
                [
                    *[int(elem[0]) for elem in with_nfs],
                    *[int(elem[0]) for elem in without_nfs],
                ]
            )
        )
    )
    plt.savefig(f"{data_path}/speed-relative.png")

    # print(sns.load_dataset("tips"))
    # dataset = pd.DataFrame(
    #     [
    #         {
    #             "nodes": str(key),
    #         }
    #         for key in list(
    #             set(
    #                 [
    #                     *[int(key) for key in with_nfs],
    #                     *[int(key) for key in without_nfs],
    #                 ]
    #             )
    #         )
    #     ]
    # )
    # sns.violinplot(data={
    # })
