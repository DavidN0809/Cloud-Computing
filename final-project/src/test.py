import requests

# Base URL for API
base_url = "http://localhost:8000"

# Headers for JSON content type
headers_json = {'Content-Type': 'application/json'}

# User credentials for testing
user_credentials = {
    'regular': {'username': 'regular_user', 'password': 'regular_pass'},
    'admin': {'username': 'admin_user', 'password': 'admin_pass'}
}

# Tokens to store session tokens
tokens = {}

def register_user(role):
    """Register users."""
    user = user_credentials[role]
    data = {
        "username": user['username'],
        "email": f"{user['username']}@example.com",
        "password": user['password'],
        "role": role
    }
    response = requests.post(f"{base_url}/auth/register", headers=headers_json, json=data)
    try:
        response_data = response.json()
        print(f"Register {role}: ", response_data)
    except requests.exceptions.JSONDecodeError:
        print(f"Failed to register {role}. Status Code: {response.status_code}, Response Text: {response.text}")


import json

def login_user(role):
    """Login users and store their tokens."""
    user = user_credentials[role]
    data = {
        "username": user['username'],
        "password": user['password']
    }
    response = requests.post(f"{base_url}/auth/login", headers=headers_json, json=data)
    print(f"Raw Response Text for {role}: {response.text}")  # Debug raw response

    if response.status_code == 200:
        try:
            # Split the response text to isolate the first JSON object
            first_part = response.text.split('\n')[0]
            response_data = json.loads(first_part)  # Parse the isolated first JSON object
            token = response_data.get('token')
            if token:
                tokens[role] = token
                print(f"Login {role} successful: Token stored")
            else:
                print(f"Token not found in the response for {role}.")
        except json.JSONDecodeError as e:
            print(f"Failed to decode JSON for {role}. Error: {e}")
    else:
        print(f"Login failed for {role}. Status Code: {response.status_code}, Response Text: {response.text}")




def create_user():
    """Create a new user by an admin."""
    data = {"username": "newuser", "email": "newuser@example.com", "password": "newuserpass"}
    headers = headers_json.copy()
    headers['Authorization'] = f"Bearer {tokens['admin']}"
    response = requests.post(f"{base_url}/users/create", headers=headers, json=data)
    print("Create User: ", response.json())

def list_users():
    """List all users by an admin."""
    headers = {'Authorization': f"Bearer {tokens['admin']}"}
    response = requests.get(f"{base_url}/users/list", headers=headers)
    print("List Users: ", response.json())

# Register users
register_user('regular')
register_user('admin')

# Login users
login_user('regular')
login_user('admin')

# Testing creating a user (requires admin)
create_user()

# Testing listing users (requires admin)
list_users()
