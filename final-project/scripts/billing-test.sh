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


# Billing CRUD Tests

# Create a billing (admin only)
echo "Creating a billing (admin only)..."
CREATE_BILLING_RESPONSE=$(curl -s -X POST "$BASE_URL/billings/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN" \
  -d '{
        "user_id": "'$ADMIN_USER_ID'",
        "task_id": "'$TASK_ID'",
        "hours": 5,
        "hourly_rate": 100,
        "amount": 500
      }')
echo $CREATE_BILLING_RESPONSE
print_test_result "Create Billing (Admin)" $?

BILLING_ID=$(echo $CREATE_BILLING_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Billing ID: $BILLING_ID"
echo

# Read the billing (admin only)
echo "Reading the billing (admin only)..."
READ_BILLING_RESPONSE=$(curl -s -X GET "$BASE_URL/billings/get/$BILLING_ID" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN")
echo $READ_BILLING_RESPONSE
print_test_result "Read Billing (Admin)" $?
echo

# Update the billing (admin only)
echo "Updating the billing (admin only)..."
UPDATE_BILLING_RESPONSE=$(curl -s -X PUT "$BASE_URL/billings/update/$BILLING_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN" \
  -d '{
        "hours": 8,
        "hourly_rate": 120,
        "amount": 960
      }')
echo $UPDATE_BILLING_RESPONSE
print_test_result "Update Billing (Admin)" $?
echo

# Delete the billing (admin only)
echo "Deleting the billing (admin only)..."
DELETE_BILLING_RESPONSE=$(curl -s -X DELETE "$BASE_URL/billings/remove/$BILLING_ID" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN")
echo $DELETE_BILLING_RESPONSE
print_test_result "Delete Billing (Admin)" $?
echo

# Summary
echo "Billing Test Summary:"
if [ "$TEST_FAILED" = true ]; then
  echo "❌ Some billing tests failed. Please check the output for more details."
else
  echo "✅ All billing tests passed successfully!"
fi
