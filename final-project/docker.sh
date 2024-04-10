
# Docker commands for setup
echo "Stopping all running containers..."
sudo docker stop $(sudo docker ps -aq)

#echo "Starting up services with docker-compose..."
#docker-compose up -d

echo "Rebuilding services (if code was changed)..."
docker-compose up -d --build

