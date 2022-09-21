# Container brdige

## Chart and docker versioning

The following files to be updated for chart and docker tag versionsing:
- Makefile -> BUILD parameter
- charts/agent/Chart.yaml
- charts/client/Chart.yaml
- .github/workflows/agent-docker-image.yaml
- .github/workflows/client-docker-image.yaml

Currently docker images are tagged latest always. Once stabilized versioning can be added to workflow with same version as chart by modifying above 3 files.

## Start docker compose manual test

```
docker compose -f ./docker-compose_manual_test.yaml up -d
```

## Stop docker compose

```
docker compose -f ./docker-compose_manual_test.yaml down -v
```

## Docker registry configuration

### Docker container registry

References:

- https://docs.docker.com/registry/deploying/

```
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

### Notification settings

References:

- https://docs.docker.com/registry/configuration/#notifications
- https://docs.docker.com/registry/notifications/

```
notifications:
  events:
    includereferences: true
  endpoints:
    - name: alistener
      disabled: false
      url: http://agent:8090/localregistry/event
      timeout: 10s
      threshold: 10
      backoff: 1s
```

### Container registry with podman

References:

``` https://www.redhat.com/sysadmin/simple-container-registry ```

# Testing

## Using docker client

### pull ubuntu latest image

```  docker pull ubuntu:latest ```

### Tag with local registry

Local registry listening on 0.0.0.0:5001

So have to tag with localhost:5001 prefix path

``` docker tag ubuntu:latest localhost:5001/ubuntu:v1 ```

### Push the tagged image to local registry

``` docker push localhost:5001/ubuntu:v1 ```

## Using podman client

### pull ubuntu latest image

``` podman pull ubuntu:latest ```

### Tag with version 2

``` podman tag ubuntu:latest localhost:5001/ubuntu:v2 ```

### Push the tagged image to local registry

``` podman push localhost:5001/ubuntu:v2 --tls-verify=false ```

# Check the events in grafana dashboard

- Open grafana GUI using ```localhost:3000```

- Install clickhouse plugin

- Add data source for clickhouse

# local environment with cpu&mem

## Docker

Reference: https://docs.docker.com/config/containers/resource_constraints/

```
docker run -it --cpus="0.5" --memory=256m container-bridge-agent:0.1.1
docker run -it --cpus="0.5" --memory=256m container-bridge-client:0.1.1
```

## Docker compose

Reference: https://docs.docker.com/compose/compose-file/compose-file-v3/#resources

```
For example:
version: "3.9"
services:
  agent:
    image: container-bridge-agent:0.1.1
    deploy:
      resources:
        limits:
          cpus: '0.50'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 256M
```

# Example docker event payload

```
{
   "events": [
      {
         "id": "d539f3c5-5734-47f0-a2a8-d27de5f9edb1",
         "timestamp": "2022-09-21T19:00:49.658010356Z",
         "action": "push",
         "target": {
            "mediaType": "application/octet-stream",
            "size": 29136663,           
            "digest": "sha256:12f42424f10d587d02674c8a0dab1c08d3fd81ab6bac5b7f5e3799215c6c52e6",
            "length": 29136663,         
            "repository": "ubuntu",     
            "url": "http://localhost:5001/v2/ubuntu/blobs/sha256:12f42424f10d587d02674c8a0dab1c08d3fd81ab6bac5b7f5e3799215c6c52e6"
         },
         "request": {
            "id": "3784878a-5075-4ea5-a51b-5288ebc54c9c",
            "addr": "172.22.0.1:36222",
            "host": "localhost:5001",
            "method": "PUT", 
            "useragent": "containers/5.16.0 (github.com/containers/image)"
         },
         "actor": {},
         "source": {
            "addr": "5b458294602c:5000",
            "instanceID": "da4bdd5f-3f85-45c6-a37d-bc97a82fbd39"
         }
      }
   ]
}
```
