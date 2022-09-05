# Container brdige

## Start docker compose

```
export GITLAB_HOME=./gitlab-repo-server
docker compose up -d
```

## Stop docker compose

```
docker compose down -v
```

## Docker registry configuration

### Image

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
