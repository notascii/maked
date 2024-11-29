import os
import requests

user = input(f"Grid'5000 username (default is {os.getlogin()}): ") or os.getlogin()
password = input("Grid'5000 password (leave blank on frontends): ")
g5k_auth = (user, password) if password else None

sites = requests.get("https://api.grid5000.fr/stable/sites", auth=g5k_auth).json()["items"]

print("Grid'5000 sites:")
for site in sites:

    site_id = site["uid"]
    print(site_id + ":")

    site_clusters = requests.get(
        f"https://api.grid5000.fr/stable/sites/{site_id}/clusters",
        auth=g5k_auth,
    ).json()["items"]

    for cluster in site_clusters:
        print("-", cluster["uid"])