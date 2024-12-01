import requests
import time

class Grid5000API:
    def __init__(self, user, password, site):
        self.user = user
        self.password = password
        self.auth = (self.user, self.password)
        self.site = site
        self.base_url = f"https://api.grid5000.fr/stable/sites/{site}"

    def submit_deployment_job(self, nodes, script_path):
        jobs_url = f"{self.base_url}/jobs/"
        job_data = {
            "resources": f"nodes={nodes}",
            "types": ["deploy"],
            "command": f"bash {script_path}",
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
                if state != old_state:
                    print(f"Current job state: {state}")
                    old_state = state
                if state in ['terminated', 'error', 'killed']:
                    return state
            else:
                print(f"Failed to retrieve job status: {response.status_code}")
                print("Error:", response.text)
                exit(1)
            time.sleep(10)  # Wait for 10 seconds before checking again

if __name__ == "__main__":
    login = "aabdelaz"
    password = "SCLK6yDs!m74tQG"
    site = "rennes"
    nodes = 4
    script_path = "./run_maked.sh"

    g5k = Grid5000API(login, password, site)
    job_id = g5k.submit_deployment_job(nodes, script_path)
    job_state = g5k.wait_for_job_completion(job_id)
    if job_state == 'terminated':
        print("Deployment completed successfully.")
    else:
        print("Job did not terminate successfully.")
