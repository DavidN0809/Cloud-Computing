import requests

BASE_URL = "http://localhost:8000"

def register_user(username, email, password, role):
    url = f"{BASE_URL}/auth/register"
    payload = {
        "username": username,
        "email": email,
        "password": password,
        "role": role
    }
    response = requests.post(url, json=payload)
    assert response.status_code == 201
    return response.json()

def login(username, password):
    url = f"{BASE_URL}/auth/login"
    payload = {
        "username": username,
        "password": password
    }
    response = requests.post(url, json=payload)
    assert response.status_code == 200
    return response.json()['token']

def create_task(admin_token, task_data):
    url = f"{BASE_URL}/tasks/create"
    headers = {"Authorization": f"Bearer {admin_token}"}
    response = requests.post(url, headers=headers, json=task_data)
    assert response.status_code == 201
    return response.json()['id']

def assign_task(admin_token, user_id, task_id):
    # Assuming you have an endpoint for assigning a task
    url = f"{BASE_URL}/tasks/assign/{task_id}"
    headers = {"Authorization": f"Bearer {admin_token}"}
    payload = {"assigned_to": user_id}
    response = requests.put(url, headers=headers, json=payload)
    assert response.status_code == 200

def update_task(admin_token, task_id, task_data):
    url = f"{BASE_URL}/tasks/update/{task_id}"
    headers = {"Authorization": f"Bearer {admin_token}"}
    response = requests.put(url, headers=headers, json=task_data)
    assert response.status_code == 200

def view_task(task_id):
    url = f"{BASE_URL}/tasks/get/{task_id}"
    response = requests.get(url)
    assert response.status_code == 200
    return response.json()

def remove_task(admin_token, task_id):
    url = f"{BASE_URL}/tasks/remove/{task_id}"
    headers = {"Authorization": f"Bearer {admin_token}"}
    response = requests.delete(url, headers=headers)
    assert response.status_code == 200

def remove_user(admin_token, user_id):
    url = f"{BASE_URL}/users/remove/{user_id}"
    headers = {"Authorization": f"Bearer {admin_token}"}
    response = requests.delete(url, headers=headers)
    assert response.status_code == 200

def billing_operations(admin_token, billing_data):
    url = f"{BASE_URL}/billings/create"
    headers = {"Authorization": f"Bearer {admin_token}"}
    response = requests.post(url, headers=headers, json=billing_data)
    assert response.status_code == 201
    return response.json()['id']

def run_test():
    # Register users
    regular_user = register_user("regular_user", "regular@example.com", "password", "user")
    admin_user = register_user("admin_user", "admin@example.com", "password", "admin")
    
    # Login users
    user_token = login("regular_user", "password")
    admin_token = login("admin_user", "password")
    
    # Perform admin actions: create, update, remove tasks and users
    task_id = create_task(admin_token, {"title": "New Task", "description": "Do something"})
    
    # Assign task to regular user
    assign_task(admin_token, regular_user['id'], task_id)
    
    # Update task
    update_task(admin_token, task_id, {"title": "Updated Task", "description": "Do something else"})
    
    # View task
    view_task(task_id)
    
    # Perform billing operations for the task
    billing_id = billing_operations(admin_token, {"task_id": task_id, "user_id": regular_user['id'], "amount": 100})
    
    # Remove task
    remove_task(admin_token, task_id)
    
    # Remove user
    remove_user(admin_token, regular_user['id'])

if __name__ == "__main__":
    run_test()
