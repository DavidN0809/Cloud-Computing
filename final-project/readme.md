### curl commands for testing
#### user
```
# Create a user
curl -X POST -H "Content-Type: application/json" -d '{"username":"johndoe","email":"johndoe@example.com","password":"password"}' http://localhost:8001/users/create

# Get a user
curl -X GET http://localhost:8001/users/get/<user_id>

# Update a user
curl -X PUT -H "Content-Type: application/json" -d '{"username":"johndoe_updated","email":"johndoe_updated@example.com","password":"password_updated"}' http://localhost:8001/users/update/<user_id>

# Remove a user
curl -X DELETE http://localhost:8001/users/remove/<user_id>

# List all users
curl -X GET http://localhost:8001/users/list
```

#### tasks
```
# Create a task
curl -X POST -H "Content-Type: application/json" -d '{"title":"Task 1","description":"Task description","assigned_to":"<user_id>","status":"pending","hours":5}' http://localhost:8002/tasks/create

# Get a task
curl -X GET http://localhost:8002/tasks/get/<task_id>

# Update a task
curl -X PUT -H "Content-Type: application/json" -d '{"title":"Task 1 Updated","description":"Updated task description","assigned_to":"<user_id>","status":"in progress","hours":8}' http://localhost:8002/tasks/update/<task_id>

# Remove a task
curl -X DELETE http://localhost:8002/tasks/remove/<task_id>

# List all tasks
curl -X GET http://localhost:8002/tasks/list
```

#### billing
```
# Create a billing
curl -X POST -H "Content-Type: application/json" -d '{"user_id":"<user_id>","task_id":"<task_id>","hours":5,"amount":100}' http://localhost:8003/billings/create

# Get a billing
curl -X GET http://localhost:8003/billings/get/<billing_id>

# Update a billing
curl -X PUT -H "Content-Type: application/json" -d '{"user_id":"<user_id>","task_id":"<task_id>","hours":8,"amount":150}' http://localhost:8003/billings/update/<billing_id>

# Remove a billing
curl -X DELETE http://localhost:8003/billings/remove/<billing_id>

# List all billings
curl -X GET http://localhost:8003/billings/list
```

### docker

#### kill all running containers
```
 sudo docker stop $(sudo docker ps -aq)
```
