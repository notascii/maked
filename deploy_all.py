import os
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

    def submit_deployment_job(self, nodes, script_path, makefile_directory):
        jobs_url = f"{self.base_url}/jobs/"
        job_data = {
            "resources": f"nodes={nodes}",
            "command": f"bash {script_path} {makefile_directory}",
            "name": "DeployNode"
        }
        response = requests.post(jobs_url, json=job_data, auth=self.auth)
        if response.status_code == 201:
            job = response.json()
            print(f"Job submitted: ID {job['uid']}")
            return job['uid']
        else:
            print(f"Job submission failed: {response.status_code}")
            print("Error:", response.text)
            exit(1)

    def wait_for_job_completion(self, job_id):
        job_url = f"{self.base_url}/jobs/{job_id}"
        old_state = ""
        while True:
            response = requests.get(job_url, auth=self.auth)
            if response.status_code == 200:
                job_info = response.json()
                state = job_info['state']
                if state != old_state:
                    print(f"Current job state: {state}")
                    old_state = state
                if state in ['terminated', 'error', 'killed']:
                    return state
            else:
                print(f"Failed to retrieve job status: {response.status_code}")
                print("Error:", response.text)
                exit(1)

if __name__ == "__main__":
    login = input("Enter login: ")
    password = getpass.getpass()
    site = os.getenv("GRID5000_SITE", "rennes")
    script_path = "./maked/run_maked.sh"
    script_init_path = "./maked/run_make.sh"    
    directory_list = ["matrix", "premier_tiny", "premier"]
    list_int = [2, 3, 4, 6, 8, 11]

    for directory in directory_list:
        print(f"Deployment for classic Make")
        start = time.time()
        g5k = Grid5000API(login, password, site)
        job_id = g5k.submit_deployment_job(1, script_init_path, directory)
        job_state = g5k.wait_for_job_completion(job_id)
        end = time.time()
        if job_state == 'terminated':
            print(f"Deployment initial completed successfully in {end - start:.2f} seconds.")
        else:
            print("Job did not terminate successfully.")
        for number_of_nodes in list_int:
            start = time.time()
            print(f"Deployment with {number_of_nodes-1} clients")
            g5k = Grid5000API(login, password, site)
            job_id = g5k.submit_deployment_job(number_of_nodes, script_path, directory)
            job_state = g5k.wait_for_job_completion(job_id)
            end = time.time()
            if job_state == 'terminated':
                print(f"Deployment completed successfully in {end - start:.2f} seconds.")
            else:
                print("Job did not terminate successfully.")
