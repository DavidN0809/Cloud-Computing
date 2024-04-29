#!/bin/bash

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

# Clear the database
echo "Clearing the database..."
./clear-db.sh

# Run user tests
echo "Running user tests..."
./user-test.sh
USER_TEST_STATUS=$?
print_test_result "User Tests" $USER_TEST_STATUS

# Run task tests
echo "Running task tests..."
./task-test.sh
TASK_TEST_STATUS=$?
print_test_result "Task Tests" $TASK_TEST_STATUS

# Run billing tests
echo "Running billing tests..."
./billing-test.sh
BILLING_TEST_STATUS=$?
print_test_result "Billing Tests" $BILLING_TEST_STATUS

# Run task relationship tests
echo "Running task relationship tests..."
./task-relationship-test.sh
TASK_RELATIONSHIP_TEST_STATUS=$?
print_test_result "Task Relationship Tests" $TASK_RELATIONSHIP_TEST_STATUS

# Summary
echo "Test Summary:"
if [ "$TEST_FAILED" = true ]; then
  echo "❌ Some tests failed. Please check the output for more details."
else
  echo "✅ All tests passed successfully!"
fi
