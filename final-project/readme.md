# API Testing Guide

This document outlines the steps to test the user management system's functionalities, including registration, login, CRUD operations for users and tasks, and billing management.

## Setup
Install docker
``` 
sudo curl -fsSL https://get.docker.com -o get-docker.sh
```
```
sudo sh get-docker.sh 
```

Before you begin testing, ensure your local server is running:
```
docker compose up -d
```
To rebuild after making changes to the code, use:
```
docker compose up --build -d
```


## User Registration and Login
### Register a Regular User
```
curl -X POST http://localhost:8000/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "regular_user",
        "email": "regular@example.com",
        "password": "regular_pass",
        "role": "regular"
      }'
```
### Login as Regular User
```
curl -X POST http://localhost:8000/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "regular_user",
        "password": "regular_pass"
      }'
```

### Register a Admin User
```
curl -X POST http://localhost:8000/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "admin_user",
        "email": "admin@example.com",
        "password": "admin_pass",
        "role": "admin"
      }'

```
### Login as Admin User
```
curl -X POST http://localhost:8000/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "admin_user",
        "password": "admin_pass"
      }'
```

Note: Be sure to update the placeholder `<admin_token>` with the actual admin JWT token obtained after logging in as an admin. Similarly, replace `<user_id>`, `<task_id>`, and `<billing_id>` with actual IDs as you proceed with the tests. The commands assuming the API is listening on `localhost` and port `8000`. Adjust the port if your services are running on different ports.

## CRUD Operations for Users
### Create a User
```bash
curl -X POST http://localhost:8000/users/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{"username":"newuser","email":"newuser@example.com","password":"newuserpass"}'
```
### Get a User (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/users/get/<user_id>
```
### Update a User (Admin only)
This operation should only succeed with admin privileges.
``` bash
curl -X PUT http://localhost:8000/users/update/<user_id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{"username":"newuser_updated","email":"newuser_updated@example.com","password":"newuserpass_updated"}'
```

### Remove a User (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/users/remove/<user_id> \
  -H "Authorization: Bearer <admin_token>"
```

### List All Users (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/users/list \
  -H "Authorization: Bearer <admin_token>"
```

### Delete All Users (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/users/delete-all \
  -H "Authorization: Bearer <admin_token>"
```

## CRUD Operations for Tasks

### Create a Task (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X POST http://localhost:8000/tasks/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
        "title": "New Task",
        "description": "Task description",
        "assigned_to": "<user_id>",
        "status": "pending",
        "hours": 5
      }'
```

### Get a Task
```bash
curl -X GET http://localhost:8000/tasks/get/<task_id>
```

### Update a Task (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X PUT http://localhost:8000/tasks/update/<task_id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
        "title": "Task 1 Updated",
        "description": "Updated task description",
        "assigned_to": "<user_id>",
        "status": "in progress",
        "hours": 8
      }'
```

### Remove a Task (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/tasks/remove/<task_id> \
  -H "Authorization: Bearer <admin_token>"
```

### List All Tasks
```bash
curl -X GET http://localhost:8000/tasks/list
```

### Delete All Tasks (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/tasks/removeAllTasks \
  -H "Authorization: Bearer <admin_token>"
```

## CRUD Operations for Billing

### Create a Billing (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X POST http://localhost:8000/billings/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
        "user_id": "<user_id>",
        "task_id": "<task_id>",
        "hours": 5,
        "amount": 100
      }'
```

### Get a Billing  (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/billings/get/<billing_id>
```

### Update a Billing (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X PUT http://localhost:8000/billings/update/<billing_id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
        "user_id": "<user_id>",
        "task_id": "<task_id>",
        "hours": 8,
        "amount": 150
      }'
```

### Remove a Billing (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/billings/remove/<billing_id> \
  -H "Authorization: Bearer <admin_token>"
```

### List All Billings (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/billings/list \
  -H "Authorization: Bearer <admin_token>"
```

### Delete All Billings (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/billings/removeAllBillings \
  -H "Authorization: Bearer <admin_token>"
```

