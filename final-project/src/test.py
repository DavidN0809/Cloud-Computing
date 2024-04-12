import json
import requests

def register_user(role):
    """Register users with error reporting for non-successful outcomes."""
    user = user_credentials[role]
    data = {
        "username": user['username'],
        "email": f"{user['username']}@example.com",
        "password": user['password'],
        "role": role
    }
    response = requests.post(f"{base_url}/auth/register", headers=headers_json, json=data)
    if response.status_code not in (200, 201):
        print(f"Registration failed for {role}. Status Code: {response.status_code}, Response Text: {response.text}")

def login_user(role):
    """Login users and store their tokens with error reporting for non-successful outcomes."""
    user = user_credentials[role]
    data = {
        "username": user['username'],
        "password": user['password']
    }
    response = requests.post(f"{base_url}/auth/login", headers=headers_json, json=data)
    if response.status_code == 200 or response.status_code == 201:
        try:
            first_part = response.text.split('\n')[0]
            response_data = json.loads(first_part)
            token = response_data.get('token')
            if token:
                tokens[role] = token
            else:
                print(f"Token not found in the response for {role}.")
        except json.JSONDecodeError as e:
            print(f"Failed to decode JSON for {role}. Error: {e}")
    else:
        print(f"Login failed for {role}. Status Code: {response.status_code}, Response Text: {response.text}")

def create_user():
    """Create a new user by an admin with error reporting for non-successful outcomes."""
    data = {"username": "newuser", "email": "newuser@example.com", "password": "newuserpass"}
    headers = headers_json.copy()
    headers['Authorization'] = f"Bearer {tokens.get('admin', '')}"
    response = requests.post(f"{base_url}/users/create", headers=headers, json=data)
    if response.status_code not in (200, 201):
        print("Failed to create user. Status Code: {}, Response Text: {}".format(response.status_code, response.text))

def list_users():
    """List all users by an admin with error reporting for non-successful outcomes."""
    headers = {'Authorization': f"Bearer {tokens.get('admin', '')}"}
    response = requests.get(f"{base_url}/users/list", headers=headers)
    if response.status_code not in (200, 201):
        print("Failed to list users. Status Code: {}, Response Text: {}".format(response.status_code, response.text))

# User credentials and base configuration
base_url = "http://localhost:8000"
headers_json = {'Content-Type': 'application/json'}
user_credentials = {
    'regular': {'username': 'regular_user', 'password': 'regular_pass'},
    'admin': {'username': 'admin_user', 'password': 'admin_pass'}
}
tokens = {}

# Execution of the registration and login tests
register_user('regular')
register_user('admin')

login_user('regular')
login_user('admin')

# Admin operations
create_user()
list_users()
