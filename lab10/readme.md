```sudo docker-compose build```
```sudo docker-compose up```



```
microk8s kubectl get svc -n container-registry
```
```
docker tag my-webserver localhost:32000/webserver
docker push localhost:32000/webserver
```

edit /etc/docker/daemon.json
```
{
  "insecure-registries" : ["localhost:32000"]
}
```
```
sudo systemctl restart docker
sudo docker build -t my-webserver .
sudo docker tag my-webserver:latest localhost:32000/webserver:latest
sudo docker push localhost:32000/webserver:latest
```
```
sudo microk8s enable dns
```

```
sudo microk8s kubectl apply -f webserver-deployment.yaml
sudo microk8s kubectl apply -f webserver-service.yaml
sudo microk8s kubectl apply -f mongodb-deployment.yaml
sudo microk8s kubectl apply -f mongodb-service.yaml
```

```
sudo microk8s kubectl get deployments
sudo microk8s kubectl get pods
```
```
sudo microk8s kubectl get svc webserver
```

```
curl http://localhost:30080/list
```
