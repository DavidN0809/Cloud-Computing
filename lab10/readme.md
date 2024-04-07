####how to build compose
```
sudo docker-compose build
sudo docker-compose up
```


####setting up registry
```
microk8s kubectl get svc -n container-registry
```
```
docker tag my-webserver localhost:32000/webserver
docker push localhost:32000/webserver
```

####edit /etc/docker/daemon.json
```
{
  "insecure-registries" : ["localhost:32000"]
}
```
####more registry finishing
```
sudo systemctl restart docker
sudo docker build -t my-webserver .
sudo docker tag my-webserver:latest localhost:32000/webserver:latest
sudo docker push localhost:32000/webserver:latest
```

####adding dns
```
sudo microk8s enable dns
```
####deploying files
```
sudo microk8s kubectl apply -f webserver-deployment.yaml
sudo microk8s kubectl apply -f webserver-service.yaml
sudo microk8s kubectl apply -f mongodb-deployment.yaml
sudo microk8s kubectl apply -f mongodb-service.yaml
```
####checking deployments and pods wait till all are 1/1
```
sudo microk8s kubectl get deployments
sudo microk8s kubectl get pods
```
####shows what port is open
```
sudo microk8s kubectl get svc webserver
```
####testing webserver
```
curl http://localhost:30080/list
```
