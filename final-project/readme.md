### curl commands for testing
#### user
```
# Create a regular user
curl -X POST -H "Content-Type: application/json" -d '{"username":"john_doe","email":"john@example.com","password":"password123","role":"regular"}' http://localhost:8000/auth/register

# Create an admin user
curl -X POST -H "Content-Type: application/json" -d '{"username":"admin","email":"admin@example.com","password":"adminpass","role":"admin"}' http://localhost:8000/auth/register

# User login
curl -X POST -H "Content-Type: application/json" -d '{"email":"john@example.com","password":"password123"}' http://localhost:8000/auth/login

# Get a user
curl -X GET http://localhost:8000/users/{user_id}/

# Update a user
curl -X PUT -H "Content-Type: application/json" -d '{"username":"john_doe_updated","email":"john_updated@example.com","password":"newpassword"}' http://localhost:8000/users/{user_id}/

# Delete a user
curl -X DELETE http://localhost:8000/users/{user_id}/

# List all users
curl -X GET http://localhost:8000/users/
```
#### tasks
```
# Create a task
curl -X POST -H "Content-Type: application/json" -d '{"title":"Task 1","description":"Description of Task 1","assigned_to":"{user_id}","status":"open","hours":5}' http://localhost:8000/tasks/

# Get a task
curl -X GET http://localhost:8000/tasks/{task_id}/

# Update a task
curl -X PUT -H "Content-Type: application/json" -d '{"title":"Updated Task 1","description":"Updated description of Task 1","assigned_to":"{user_id}","status":"in_progress","hours":8}' http://localhost:8000/tasks/{task_id}/

# Delete a task
curl -X DELETE http://localhost:8000/tasks/{task_id}/

# List all tasks
curl -X GET http://localhost:8000/tasks/
```
#### billing
```
# Create a billing
curl -X POST -H "Content-Type: application/json" -d '{"user_id":"{user_id}","task_id":"{task_id}","hours":5,"amount":100}' http://localhost:8000/billings/

# Get a billing
curl -X GET http://localhost:8000/billings/{billing_id}/

# Update a billing
curl -X PUT -H "Content-Type: application/json" -d '{"user_id":"{user_id}","task_id":"{task_id}","hours":8,"amount":160}' http://localhost:8000/billings/{billing_id}/

# Delete a billing
curl -X DELETE http://localhost:8000/billings/{billing_id}/

# List all billings
curl -X GET http://localhost:8000/billings/
```

### docker

#### kill all running containers
```
 sudo docker stop $(sudo docker ps -aq)
```
