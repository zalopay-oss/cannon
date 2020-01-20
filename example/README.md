# Example

## Run with docker

### Requirement

You must install:

- [docker](https://docs.docker.com/install/)
- docker-compose

Start server first

```bash
$ cd example
$ docker-compose -f docker-compose-server.yaml up
```

Start cannon

```bash
$ cd example
$ docker-compose -f docker-compose-client.yaml up
```

Access [localhost:8089](http://localhost:8089) to view report
