import os
import sys
import time
from getpass import getpass

import requests
import time
import getpass

class Grid5000API:
    def __init__(self, user, password, site):
        self.user = user
        self.password = password
        self.auth = (self.user, self.password)
        self.site = site
        self.base_url = f"https://api.grid5000.fr/stable/sites/{site}"

    def submit_deployment_job(self, nodes, script_path, makefile_directory, name):
        jobs_url = f"{self.base_url}/jobs/"
        job_data = {
            "resources": f"nodes={nodes}",
            "command": f"bash {script_path} {makefile_directory}",
            "name": name
        }
        response = requests.post(jobs_url, json=job_data, auth=self.auth)
        if response.status_code == 201:
            job = response.json()
            print(f"Job submitted: ID {job['uid']}")
            return job["uid"]
        else:
            print(f"Job submission failed: {response.status_code}")
            print("Error:", response.text)
            exit(1)


if __name__ == "__main__":
    # Ensure the script is run with the required argument
    if len(sys.argv) != 2:
        print("Usage: python launch.py <number_of_nodes>")
        sys.exit(1)

    # Get the number of nodes from the argument
    try:
        node_count = int(sys.argv[1])
        if node_count < 1:
            raise ValueError
    except ValueError:
        print("Error: <number_of_nodes> must be a positive integer.")
        sys.exit(1)
    login = input("Enter login: ")
    password = getpass.getpass()
    site = os.getenv("GRID5000_SITE", "rennes")
    script_path = "./maked/run_maked.sh"
    directory_list = ["premier_tiny"]
    for directory in directory_list:
        print(f"#################### {directory} #########################")
        print(f"Deployment of maked")
        g5k = Grid5000API(login, password, site)
        job_id = g5k.submit_deployment_job(node_count, script_path, directory, f"maked_{directory}")