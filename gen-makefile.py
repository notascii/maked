import contextlib
import sys
from pathlib import Path

if __name__ == "__main__":
    if len(sys.argv) != 9:
        raise ValueError(
            "Wrong use of this file. Please use `python gen-makefile.py -n <name> -c <nb_commands> -s <sleep_time> -z <file_size>"
        )
    name = None
    commands = None
    sleep = None
    size = None
    for index, arg in enumerate(sys.argv):
        if arg == "-n":
            name = sys.argv[index + 1]
        elif arg == "-c":
            commands = [int(x) for x in sys.argv[index + 1].split(",")]
        elif arg == "-s":
            sleep = float(sys.argv[index + 1])
        elif arg == "-z":
            size = int(sys.argv[index + 1])

    if sleep is None or commands is None or name is None or size is None:
        raise ValueError(
            "Wrong use of this file. Please use `python gen-makefile.py -n <name> -b <nb_branches> -s <sleep_time> -z <file_size>"
        )
    makefile_path = f"makefiles/{name}-c-{commands}-s-{sleep}-z-{size}"
    with contextlib.suppress(FileExistsError):
        Path.mkdir(makefile_path)

    root_commands = []
    for index, c in enumerate(commands):
        for length in range(c):
            if index == 0:
                root_commands.append(f"command-{length}")

    makefile = f"""all: {" ".join(root_commands)}
"""
    leaf_commands = []
    for n0 in range(commands[0]):
        makefile += f"""
command-{n0}: {" ".join([f"command-{n0}-{nb}" for nb in range(commands[1])])}
\tsleep {sleep} && echo "{"#" * size}" > command-{n0}
"""
        leaf_commands.append(f"command-{n0}")
        for n1 in range(commands[1]):
            makefile += f"""
command-{n0}-{n1}: {" ".join([f"command-{n0}-{n1}-{nb}" for nb in range(commands[2])])}
\tsleep {sleep} && echo "{"#" * size}" > command-{n0}-{n1}
"""
            leaf_commands.append(f"command-{n0}-{n1}")
            for n2 in range(commands[2]):
                makefile += f"""
command-{n0}-{n1}-{n2}:
\tsleep {sleep} && echo "{"#" * size}" > command-{n0}-{n1}-{n2}
"""
                leaf_commands.append(f"command-{n0}-{n1}-{n2}")

    makefile += f"""
clean: 
\t{" && ".join([f"rm {c}" for c in leaf_commands])}
"""

    with open(Path(makefile_path) / "Makefile", "w") as f:
        f.write(makefile)
    with open(Path(makefile_path) / ".gitignore", "w") as f:
        f.write("command-*")
