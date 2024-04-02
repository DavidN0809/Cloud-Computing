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
### docker

#### kill all running containers
```
 sudo docker stop $(sudo docker ps -aq)
```
