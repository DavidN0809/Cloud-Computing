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

# Task Relationship Tests

# Create a parent task (regular user)
echo "Creating a parent task (regular user)..."
CREATE_PARENT_TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "title": "Parent Task",
        "description": "This is a parent task.",
        "assigned_to": "'$REGULAR_USER_ID'",
        "status": "in_progress",
        "hours": 10,
        "start_date": "2024-05-01T00:00:00Z",
        "end_date": "2024-05-03T00:00:00Z"
      }')
echo $CREATE_PARENT_TASK_RESPONSE
print_test_result "Create Parent Task (Regular User)" $?

PARENT_TASK_ID=$(echo $CREATE_PARENT_TASK_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Parent Task ID: $PARENT_TASK_ID"
echo

# Create a child task (regular user)
echo "Creating a child task (regular user)..."
CREATE_CHILD_TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "title": "Child Task",
        "description": "This is a child task.",
        "assigned_to": "'$REGULAR_USER_ID'",
        "status": "in_progress",
        "hours": 5,
        "start_date": "2024-05-01T09:00:00Z",
        "end_date": "2024-05-01T14:00:00Z",
        "parent_task": "'$PARENT_TASK_ID'"
      }')
echo $CREATE_CHILD_TASK_RESPONSE
print_test_result "Create Child Task (Regular User)" $?

CHILD_TASK_ID=$(echo $CREATE_CHILD_TASK_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Child Task ID: $CHILD_TASK_ID"
echo

# Update the child task to done (regular user)
echo "Updating the child task to done (regular user)..."
UPDATE_CHILD_TASK_RESPONSE=$(curl -s -X PUT "$BASE_URL/tasks/update/$CHILD_TASK_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "status": "done"
      }')
echo $UPDATE_CHILD_TASK_RESPONSE
print_test_result "Update Child Task to Done (Regular User)" $?
echo

# Check the child task invoice ID
echo "Checking the child task invoice ID..."
CHILD_TASK_DETAILS=$(curl -s -X GET "$BASE_URL/tasks/get/$CHILD_TASK_ID" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN")
echo $CHILD_TASK_DETAILS

CHILD_TASK_INVOICE_ID=$(echo $CHILD_TASK_DETAILS | grep -o '"invoice_id":"[^"]*' | cut -d'"' -f4)
echo "Child Task Invoice ID: $CHILD_TASK_INVOICE_ID"

if [ -n "$CHILD_TASK_INVOICE_ID" ] && [ "$CHILD_TASK_INVOICE_ID" != "000000000000000000000000" ]; then
  print_test_result "Child Task Invoice Generated" 0
else
  print_test_result "Child Task Invoice Generated" 1
fi
echo

# Update the parent task to done (regular user)
echo "Updating the parent task to done (regular user)..."
UPDATE_PARENT_TASK_RESPONSE=$(curl -s -X PUT "$BASE_URL/tasks/update/$PARENT_TASK_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "status": "done"
      }')
echo $UPDATE_PARENT_TASK_RESPONSE
print_test_result "Update Parent Task to Done (Regular User)" $?
echo

# Check the parent task invoice ID
echo "Checking the parent task invoice ID..."
PARENT_TASK_DETAILS=$(curl -s -X GET "$BASE_URL/tasks/get/$PARENT_TASK_ID" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN")
echo $PARENT_TASK_DETAILS

PARENT_TASK_INVOICE_ID=$(echo $PARENT_TASK_DETAILS | grep -o '"invoice_id":"[^"]*' | cut -d'"' -f4)
echo "Parent Task Invoice ID: $PARENT_TASK_INVOICE_ID"

if [ -n "$PARENT_TASK_INVOICE_ID" ] && [ "$PARENT_TASK_INVOICE_ID" != "000000000000000000000000" ]; then
  print_test_result "Parent Task Invoice Generated" 0
else
  print_test_result "Parent Task Invoice Generated" 1
fi
echo

# Create a regular task (regular user)
echo "Creating a regular task (regular user)..."
CREATE_REGULAR_TASK_RESPONSE=$(curl -s -X POST "$BASE_URL/tasks/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "title": "Regular Task",
        "description": "This is a regular task.",
        "assigned_to": "'$REGULAR_USER_ID'",
        "status": "in_progress",
        "hours": 8,
        "start_date": "2024-06-01T00:00:00Z",
        "end_date": "2024-06-02T00:00:00Z"
      }')
echo $CREATE_REGULAR_TASK_RESPONSE
print_test_result "Create Regular Task (Regular User)" $?

REGULAR_TASK_ID=$(echo $CREATE_REGULAR_TASK_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Regular Task ID: $REGULAR_TASK_ID"
echo

# Update the regular task to done (regular user)
echo "Updating the regular task to done (regular user)..."
UPDATE_REGULAR_TASK_RESPONSE=$(curl -s -X PUT "$BASE_URL/tasks/update/$REGULAR_TASK_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN" \
  -d '{
        "status": "done"
      }')
echo $UPDATE_REGULAR_TASK_RESPONSE
print_test_result "Update Regular Task to Done (Regular User)" $?
echo

# Check the regular task invoice ID
echo "Checking the regular task invoice ID..."
REGULAR_TASK_DETAILS=$(curl -s -X GET "$BASE_URL/tasks/get/$REGULAR_TASK_ID" \
  -H "Authorization: Bearer $REGULAR_USER_TOKEN")
echo $REGULAR_TASK_DETAILS

REGULAR_TASK_INVOICE_ID=$(echo $REGULAR_TASK_DETAILS | grep -o '"invoice_id":"[^"]*' | cut -d'"' -f4)
echo "Regular Task Invoice ID: $REGULAR_TASK_INVOICE_ID"

if [ -n "$REGULAR_TASK_INVOICE_ID" ] && [ "$REGULAR_TASK_INVOICE_ID" != "000000000000000000000000" ]; then
  print_test_result "Regular Task Invoice Generated" 0
else
  print_test_result "Regular Task Invoice Generated" 1
fi
echo

# Summary
echo "Task Relationship Test Summary:"
if [ "$TEST_FAILED" = true ]; then
  echo "❌ Some task relationship tests failed. Please check the output for more details."
else
  echo "✅ All task relationship tests passed successfully!"
fi
