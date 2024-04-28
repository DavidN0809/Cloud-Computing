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
  -d '{"username":"newuser","email":"newuser@example.com","password":"newuserpass"}'
```
### Get a User (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/users/get/<user_id> \
-H 'Authorization: Bearer <admin_token>' 
```
### Update a User (Admin only)
This operation should only succeed with admin privileges.
``` bash
curl -X PUT http://localhost:8000/users/update/<user_id> \
  -H "Content-Type: application/json" \
   -H 'Authorization: Bearer <admin_token>' \
  -d '{"username":"newuser_updated","email":"newuser_updated@example.com","password":"newuserpass_updated"}'
```

### Remove a User (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/users/remove/<user_id> \
-H 'Authorization: Bearer <admin_token>'

```

### List All Users (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/users/list \
      -H 'Authorization: Bearer <admin_token>' 
```

### Delete All Users (Testing only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/users/delete-all 
```

## CRUD Operations for Tasks

### Create a Parent Task
```bash
curl -X POST "http://localhost:8000/tasks/create" \
     -H "Content-Type: application/json" \
     -d '{
         "title": "Project Planning",
         "description": "Initial planning phase for the project.",
         "assigned_to": "<AssignedTo>",
         "status": "planned",
         "hours": 8,
         "start_date": "2024-04-01T00:00:00Z",
         "end_date": "2024-04-03T00:00:00Z"
     }'

```

### Create a Child Task
```
curl -X POST "http://localhost:8000/tasks/create" \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Example Child Task",
           "description": "This task is a child of another task.",
           "assigned_to": "<AssignedTo>",
           "status": "pending",
           "hours": 3,
           "start_date": "2024-06-01T09:00:00Z",
           "end_date": "2024-06-01T12:00:00Z",
           "parent_task": "<Parent Task ID>"
         }'

```
### Get a Task by task id
```bash
curl -X GET http://localhost:8000/tasks/get/<task_id>
```

### Get tasks by UserID
```bash
curl -X GET "http://localhost:8000/tasks/listByUser/<UserID>"
```

### Update a Parent Task
```bash
curl -X PUT "http://localhost:8000/tasks/update/<task_id>" \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Comprehensive Updated Title",
           "description": "Comprehensive updated description.",
           "assigned_to": "<User id>",
           "status": "done",
           "hours": 4.5,
           "start_date": "2024-06-02T09:00:00Z",
           "end_date": "2024-06-02T12:00:00Z"
         }'
```

### Update a Child/Convert to Child Task
#### Only requires task id field
```bash
curl -X PUT "http://localhost:8000/tasks/update/<task_id>" \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Comprehensive Updated Title",
           "description": "Comprehensive updated description.",
           "assigned_to": "<User id>",
           "status": "done",
           "hours": 4.5,
           "start_date": "2024-06-02T09:00:00Z",
           "end_date": "2024-06-02T12:00:00Z",
           "parent_task": "<parent id/new parent id>"
         }'
```

### Remove a Task (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/tasks/remove/<task_id> \
      -H 'Authorization: Bearer <admin_token>' 

```

### List All Tasks (Admin only)
```bash
curl -X GET http://localhost:8000/tasks/list \
      -H 'Authorization: Bearer <admin_token>' 
```

### Delete All Tasks (testing only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/tasks/removeAllTasks 
```

## CRUD Operations for Billing

### Create a Billing (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X POST http://localhost:8000/billings/create \
  -H "Content-Type: application/json" \
  -H 'Authorization: Bearer <admin_token>' \
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
curl -X GET http://localhost:8000/billings/get/<billing_id> \
      -H 'Authorization: Bearer <admin_token>' 
```

### Update a Billing (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X PUT http://localhost:8000/billings/update/<billing_id> \
  -H "Content-Type: application/json" \
  -H 'Authorization: Bearer <admin_token>'  \
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
     -H 'Authorization: Bearer <admin_token>' 
```

### List All Billings (Admin only)
This operation should only succeed with admin privileges.
```bash
curl -X GET http://localhost:8000/billings/list \
      -H 'Authorization: Bearer <admin_token>' 
```

### Delete All Billings (testing only)
This operation should only succeed with admin privileges.
```bash
curl -X DELETE http://localhost:8000/billings/removeAllBillings
```


