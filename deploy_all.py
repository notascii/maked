import os
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
    login = input("Enter login: ")
    password = getpass.getpass()
    site = os.getenv("GRID5000_SITE", "rennes")
    script_path = "./maked/run_maked.sh"
    script_init_path = "./maked/run_make.sh"    
    directory_list = ["premier_tiny", "matrix", "premier"]
    for directory in directory_list:
        print(f"#################### {directory} #########################")
        print(f"Deployment of make")
        g5k = Grid5000API(login, password, site)
        job_id = g5k.submit_deployment_job(1, script_init_path, directory, f"make_{directory}")
        print(f"Deployment with clients")
        g5k = Grid5000API(login, password, site)
        job_id = g5k.submit_deployment_job(11, script_path, directory, f"maked_{directory}")