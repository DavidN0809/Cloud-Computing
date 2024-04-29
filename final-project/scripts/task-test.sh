#!/bin/bash

# Base URL of the gateway
BASE_URL="http://localhost:8000"

# Function to print test results
print_test_result() {
  local test_name=$1
  local status=$2
  if [ $status -eq 0 ]; then
    echo "✅ $test_name: PASS"
  else
    echo "❌ $test_name: FAIL"
    TEST_FAILED=true
  fi
}

# Register a regular user
echo "Registering regular user..."
REGULAR_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "regular_user",
        "email": "regular@example.com",
        "password": "regular_pass",
        "role": "regular"
      }')
echo $REGULAR_USER_RESPONSE
print_test_result "Register Regular User" $?

# Login as the regular user and export the token and user ID
echo "Logging in as regular user..."
REGULAR_USER_LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "regular_user",
        "password": "regular_pass"
      }')
echo $REGULAR_USER_LOGIN
print_test_result "Login Regular User" $?

export REGULAR_USER_TOKEN=$(echo $REGULAR_USER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
export REGULAR_USER_ID=$(echo $REGULAR_USER_LOGIN | grep -o '"id":"[^"]*' | cut -d'"' -f4)

echo "Regular User Token: $REGULAR_USER_TOKEN"
echo "Regular User ID: $REGULAR_USER_ID"
echo

# Register an admin user
echo "Registering admin user..."
ADMIN_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "admin_user",
        "email": "admin@example.com",
        "password": "admin_pass",
        "role": "admin"
      }')
echo $ADMIN_USER_RESPONSE
print_test_result "Register Admin User" $?

# Login as the admin user and export the token and user ID
echo "Logging in as admin user..."
ADMIN_USER_LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "admin_user",
        "password": "admin_pass"
      }')
echo $ADMIN_USER_LOGIN
print_test_result "Login Admin User" $?

export ADMIN_USER_TOKEN=$(echo $ADMIN_USER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
export ADMIN_USER_ID=$(echo $ADMIN_USER_LOGIN | grep -o '"id":"[^"]*' | cut -d'"' -f4)

echo "Admin User Token: $ADMIN_USER_TOKEN"
echo "Admin User ID: $ADMIN_USER_ID"
echo

# Task CRUD Tests

# Create a task (regular user)
echo "Creating a task (regular user)..."
CREATE_TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "title": "Sample Task",
        "description": "This is a sample task.",
        "assigned_to": "'$REGULAR_USER_ID'",
        "status": "in_progress",
        "hours": 5,
        "start_date": "2024-06-01T00:00:00Z",
        "end_date": "2024-06-03T00:00:00Z"
      }')
echo $CREATE_TASK_RESPONSE
print_test_result "Create Task (Regular User)" $?

TASK_ID=$(echo $CREATE_TASK_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Task ID: $TASK_ID"
echo

# Read the task (regular user)
echo "Reading the task (regular user)..."
READ_TASK_RESPONSE=$(curl -s -X GET "$BASE_URL/tasks/get/$TASK_ID" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN")
echo $READ_TASK_RESPONSE
print_test_result "Read Task (Regular User)" $?
echo

# Update the task (regular user)
echo "Updating the task (regular user)..."
UPDATE_TASK_RESPONSE=$(curl -s -X PUT "$BASE_URL/tasks/update/$TASK_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "title": "Updated Sample Task",
        "description": "This is an updated sample task.",
        "status": "done"
      }')
echo $UPDATE_TASK_RESPONSE
print_test_result "Update Task (Regular User)" $?
echo

# Delete the task (admin only)
echo "Deleting the task (admin only)..."
DELETE_TASK_RESPONSE=$(curl -s -X DELETE "$BASE_URL/tasks/remove/$TASK_ID" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN")
echo $DELETE_TASK_RESPONSE
print_test_result "Delete Task (Admin)" $?
echo

# Summary
echo "Task Test Summary:"
if [ "$TEST_FAILED" = true ]; then
  echo "❌ Some task tests failed. Please check the output for more details."
else
  echo "✅ All task tests passed successfully!"
fi
