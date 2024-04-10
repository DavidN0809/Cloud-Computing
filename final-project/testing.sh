#!/bin/bash
echo "must have jq installed..."
sudo apt install jq -y

# Direct output to both terminal and output.txt
exec > >(tee output.txt) 2>&1

# Function to test and log curl commands
# Arguments:
#   1 - Description of the command being tested
#   2 - The curl command to execute
#   3 - Expected condition to evaluate success
#   4 - Variable name to indicate failure of the section
function test_command {
    echo "Testing: $1"
    local cmd=$2
    local response=$(eval "$cmd")
    echo "Response: $response"
    if ! eval $3; then
        echo "Failure in $1"
        echo "Executed Command: $cmd"
        # Mark the section as failed
        eval "$4=1"
    fi
}

# Initialize failure indicators
user_fail=0
task_fail=0
billing_fail=0

# Docker commands for setup
echo "Stopping all running containers..."
sudo docker stop $(sudo docker ps -aq)

#echo "Starting up services with docker-compose..."
docker-compose up -d

#echo "Rebuilding services (if code was changed)..."
#docker-compose up --build -d

# Allow services to initialize
echo "Waiting for services to be fully started..."
sleep 10

# User CRUD operations
echo "----- User Tests -----"
test_command "Creating a new user" \
    "curl -s -X POST -H 'Content-Type: application/json' -d '{\"username\":\"johndoe\",\"email\":\"johndoe@example.com\",\"password\":\"password\"}' http://localhost:8000/users/create" \
    '[[ $(echo $response | jq -r ".id") != null ]]' \
    "user_fail"
user_response=$response
user_id=$(echo $user_response | jq -r '.id')

test_command "Updating the created user" \
    "curl -s -X PUT -H 'Content-Type: application/json' -d '{\"username\":\"johndoe_updated\",\"email\":\"johndoe_updated@example.com\",\"password\":\"password_updated\"}' http://localhost:8000/users/update/$user_id" \
    '[[ $(echo $response | jq -r ".username") == "johndoe_updated" ]]' \
    "user_fail"

test_command "Listing all users" \
    "curl -s -X GET http://localhost:8000/users/list" \
    '[[ $(echo $response | jq -r ".[] | select(.id == \"$user_id\") | .username") == "johndoe_updated" ]]' \
    "user_fail"

# Task CRUD operations
echo "----- Task Tests -----"
test_command "Creating a new task" \
    "curl -s -X POST -H 'Content-Type: application/json' -d '{\"title\":\"Task 1\",\"description\":\"Task description\",\"assigned_to\":\"$user_id\",\"status\":\"pending\",\"hours\":5}' http://localhost:8000/tasks/create" \
    '[[ $(echo $response | jq -r ".id") != null ]]' \
    "task_fail"
task_response=$response
task_id=$(echo $task_response | jq -r '.id')

test_command "Updating the created task" \
    "curl -s -X PUT -H 'Content-Type: application/json' -d '{\"title\":\"Task 1 Updated\",\"description\":\"Updated task description\",\"assigned_to\":\"$user_id\",\"status\":\"in progress\",\"hours\":8}' http://localhost:8000/tasks/update/$task_id" \
    '[[ $(echo $response | jq -r ".status") == "in progress" ]]' \
    "task_fail"

test_command "Listing all tasks" \
    "curl -s -X GET http://localhost:8000/tasks/list" \
    '[[ $(echo $response | jq -r ".[] | select(.id == \"$task_id\") | .title") == "Task 1 Updated" ]]' \
    "task_fail"

# Billing CRUD operations
echo "----- Billing Tests -----"
test_command "Creating a new billing" \
    "curl -s -X POST -H 'Content-Type: application/json' -d '{\"user_id\":\"$user_id\",\"task_id\":\"$task_id\",\"hours\":5,\"amount\":100}' http://localhost:8000/billings/create" \
    '[[ $(echo $response | jq -r ".id") != null ]]' \
    "billing_fail"
billing_response=$response
billing_id=$(echo $billing_response | jq -r '.id')

test_command "Updating the created billing" \
    "curl -s -X PUT -H 'Content-Type: application/json' -d '{\"user_id\":\"$user_id\",\"task_id\":\"$task_id\",\"hours\":8,\"amount\":150}' http://localhost:8000/billings/update/$billing_id" \
    '[[ $(echo $response | jq -r ".amount") == 150 ]]' \
    "billing_fail"

test_command "Listing all billings" \
    "curl -s -X GET http://localhost:8000/billings/list" \
    '[[ $(echo $response | jq -r ".[] | select(.id == \"$billing_id\") | .hours") == 8 ]]' \
    "billing_fail"

# Removal Tests
echo "----- Removal Tests -----"
test_command "Removing the created user" \
    "curl -s -X DELETE http://localhost:8000/users/remove/$user_id" \
    '[[ $(echo $response | jq -r ".message") == "User removed" ]]' \
    "user_fail"

test_command "Removing the created task" \
    "curl -s -X DELETE http://localhost:8000/tasks/remove/$task_id" \
    '[[ $(echo $response | jq -r ".message") == "Task removed" ]]' \
    "task_fail"

test_command "Removing the created billing" \
    "curl -s -X DELETE http://localhost:8000/billings/remove/$billing_id" \
    '[[ $(echo $response | jq -r ".message") == "Billing removed" ]]' \
    "billing_fail"

# Check each section for failures and report
echo "----- Final Report -----"
if [ $user_fail -eq 1 ]; then
    echo "User operations had failures."
else
    echo "User CRUD passed."
fi

if [ $task_fail -eq 1 ]; then
    echo "Task operations had failures."
else
    echo "Task CRUD passed."
fi

if [ $billing_fail -eq 1 ]; then
    echo "Billing operations had failures."
else
    echo "Billing CRUD passed."
fi

echo "Script execution completed."
