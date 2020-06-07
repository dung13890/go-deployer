# go-deployer

## Run with docker for development
```sh
docker-compose up -d
docker-compose exec truck bash
```

## Setup project for development
```sh
# Inside docker
dep ensure -update
```

## Setup and running for development

```sh
# file config.yml
go run main.go
```
