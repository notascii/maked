import requests
from getpass import getpass
import time

class Grid5000API:
    def __init__(self, user, password, site):
        self.user = user
        self.password = password
        self.auth = (self.user, self.password)
        self.site = site
        self.base_url = f"https://api.grid5000.fr/stable/sites/{site}"

    def submit_deployment_job(self, nodes, commands):
        jobs_url = f"{self.base_url}/jobs/"
        job_data = {
            "resources": f"nodes={nodes}",
            "types": ["deploy"],     # Specify that this job is for deployment
            "command": commands,  # Deployment command
            "name": "DeployUbuntuNFS"
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
                if (state != old_state):
                    print(f"Current job state: {state}")
                    old_state = state
                if state in ['terminated', 'error', 'killed']:
                    return state
            else:
                print(f"Failed to retrieve job status: {response.status_code}")
                print("Error:", response.text)
                exit(1)

def fuse_command(list_command):
    return " && ".join(list_command)

if __name__ == "__main__":
    login = "aabdelaz"
    password = ""
    g5k = Grid5000API(login, password, site="rennes")
    # First we deploy the nodes
    command1 = "kadeploy3 -f $OAR_NODEFILE -e ubuntu2204-nfs "
    # Copy ./maked to each node
    command2 = "taktuk -s -l root -f <(uniq $OAR_FILE_NODES) broadcast exec [ date ]"
    command3 = "taktuk -s -l root -f <(uniq $OAR_FILE_NODES) broadcast exec [ \"apt install git\" ]"

    # Copy
    
    # We fuse all commands inside one script 
    commands = fuse_command([command1, command2, command3])
    
    # We deploy nodes
    job_id = g5k.submit_deployment_job(nodes=4, commands=commands)
    # Check each states of completion
    job_state = g5k.wait_for_job_completion(job_id)
    if job_state == 'terminated':
        print("Deployment completed successfully.")
    else:
        print("Job did not terminate successfully.")
