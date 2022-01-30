
```python
import requests

username = "username"
password = "password"

url = "http://127.0.0.1:3333/api"

# Create New User
create = requests.post(url + "/user/create", json={"username": username, "password": password}).text

# Authentication User
token = requests.post(url + "/user/auth", json={"username": username, "password": password}).json()["token"]

# Get User
getuser = requests.get(url + "/user", auth=(token, token)).json()

# Create New Wallet
wallet = requests.post(url + "/wallet/create", json={"label": "default"}, auth=(token, token)).json()

# Get list wallets
wallets = requests.get(url + "/wallets", auth=(token, token)).json()

# Get Wallet
get_wallet_default = requests.get(url + "/wallet/" + wallets[0]["wallet_id"], auth=(getuser["master_api_key"], wallets[0]["wallet_admin_key"])).json()

# Rename Wallet
rename_wallet = requests.put(url + "/wallet/" + wallets[0]["wallet_id"] + "/rename", json={"label": "bitcoin"}, auth=(token, token))

# Delete Wallet
delete_wallet = requests.delete(url + "/wallet/" + wallet["wallet_id"] + "/delete", auth=(token, token))

# Receive Bitcoin Address.
receive_wallet = requests.get(url + "/wallet/" + wallet["wallet_id"] + "/receive?network=bitcoin", auth=(getuser["master_api_key"], wallets[0]["wallet_admin_key"]))

```
