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

# Test duplicate user
# Call signup twice
echo "Calling signup twice..."
SIGNUP_RESPONSE1=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "regular_user",
        "email": "regular@example.com",
        "password": "regular_pass",
        "role": "regular"
      }')
echo $SIGNUP_RESPONSE1

SIGNUP_RESPONSE2=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{
        "username": "regular_user",
        "email": "regular@example.com",
        "password": "regular_pass",
        "role": "regular"
      }')
echo $SIGNUP_RESPONSE2

# Check if the second response contains the error message
if echo "$SIGNUP_RESPONSE2" | grep -q "User with the same username already exists"; then
  echo "User with the same username already exists"
  print_test_result "Signup Error Message" 0
else
  echo "Error message not found"
  print_test_result "Signup Error Message" 1
fi

export ADMIN_USER_TOKEN=$(echo $ADMIN_USER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
export ADMIN_USER_ID=$(echo $ADMIN_USER_LOGIN | grep -o '"id":"[^"]*' | cut -d'"' -f4)

echo "Admin User Token: $ADMIN_USER_TOKEN"
echo "Admin User ID: $ADMIN_USER_ID"
echo

# User CRUD Tests

# Create a user (admin only)
echo "Creating a user (admin only)..."
CREATE_USER_RESPONSE=$(curl -s -X POST "$BASE_URL/users/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN" \
  -d '{
        "username": "new_user",
        "email": "new@example.com",
        "password": "new_pass",
        "role": "regular"
      }')
echo $CREATE_USER_RESPONSE
print_test_result "Create User (Admin)" $?

NEW_USER_ID=$(echo $CREATE_USER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "New User ID: $NEW_USER_ID"
echo

# Read the user (admin only)
echo "Reading the user (admin only)..."
READ_USER_RESPONSE=$(curl -s -X GET "$BASE_URL/users/get/$NEW_USER_ID" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN")
echo $READ_USER_RESPONSE
print_test_result "Read User (Admin)" $?
echo

# Update the user (admin only)
echo "Updating the user (admin only)..."
UPDATE_USER_RESPONSE=$(curl -s -X PUT "$BASE_URL/users/update/$NEW_USER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN" \
  -d '{
        "username": "updated_user",
        "email": "updated@example.com",
        "password": "updated_pass"
      }')
echo $UPDATE_USER_RESPONSE
print_test_result "Update User (Admin)" $?
echo

# Delete the user (admin only)
echo "Deleting the user (admin only)..."
DELETE_USER_RESPONSE=$(curl -s -X DELETE "$BASE_URL/users/remove/$NEW_USER_ID" \
  -H "Authorization: Bearer $ADMIN_USER_TOKEN")
echo $DELETE_USER_RESPONSE
print_test_result "Delete User (Admin)" $?
echo

# Summary
echo "User Test Summary:"
if [ "$TEST_FAILED" = true ]; then
  echo "❌ Some user tests failed. Please check the output for more details."
else
  echo "✅ All user tests passed successfully!"
fi
