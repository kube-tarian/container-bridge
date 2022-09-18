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

- 