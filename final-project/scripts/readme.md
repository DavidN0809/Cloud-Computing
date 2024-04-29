# Cloud Computing Final Project - Test Scripts

This directory contains a set of test scripts for the Cloud Computing Final Project. The scripts are designed to test the functionality of the user, task, billing, and task relationship endpoints.

## Prerequisites

Before running the test scripts, ensure that you have the following:

- Bash shell (for running the scripts)
- curl command-line tool (for making HTTP requests)
- The gateway service is running and accessible at [http://localhost:8000](http://localhost:8000) (or update the `BASE_URL` variable in the scripts if running on a different URL/port)

## Test Scripts

The following test scripts are available:

- `user-test.sh`: Tests the user-related endpoints (registration, login, user CRUD operations).
- `task-test.sh`: Tests the task-related endpoints (task CRUD operations).
- `billing-test.sh`: Tests the billing-related endpoints (billing CRUD operations).
- `task-relationship-test.sh`: Tests the task relationship endpoints (creating parent and child tasks, updating task status, and checking invoice generation).
- `clear-db.sh`: Script to clear the database before running the tests.
- `test-all.sh`: Script to run all the test scripts in sequence.

## Running the Tests

To run the tests, follow these steps:

1. Make sure you have the necessary prerequisites mentioned above.
2. Open a terminal and navigate to the directory containing the test scripts.
3. Make the test scripts executable by running the following command:
    ```bash
    chmod +x *.sh
    ```
4. Run the `test-all.sh` script to execute all the test scripts in sequence:
    ```bash
    ./testl-all.sh
    ```
    This script will clear the database, run the user tests, task tests, billing tests, and task relationship tests, and provide a summary of the test results.
5. Alternatively, you can run each test script individually:
    ```bash
    ./user-test.sh
    ./task-test.sh
    ./billing-test.sh
    ./task-relationship-test.sh
    ```
6. Review the output of the test scripts to see the results of each test case. The script will indicate whether each test passed or failed.

## Test Script Details

Here's a brief overview of each test script:

- `user-test.sh`:
  - Registers an admin user and a regular user
  - Tests user CRUD operations (create, read, update, delete) as an admin user

- `task-test.sh`:
  - Registers a regular user and an admin user
  - Tests task CRUD operations (create, read, update) as a regular user
  - Tests deleting a task as an admin user

- `billing-test.sh`:
  - Registers an admin user
  - Tests billing CRUD operations (create, read, update, delete) as an admin user

- `task-relationship-test.sh`:
  - Registers a regular user
  - Tests creating a parent task and a child task
  - Tests updating the child task and parent task to "done" status
  - Checks if the child task and parent task invoices are generated
  - Tests creating and updating a regular task to "done" status
  - Checks if the regular task invoice is generated

- `clear-db.sh`:
  - Script to clear the database before running the tests
  - Ensure that you have implemented the necessary functionality to clear the database

- `test-all.sh`:
  - Script to run all the test scripts in sequence
  - Clears the database before running the tests
  - Provides a summary of the test results
