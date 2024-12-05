# Buy-Better Core API
## Table of contents
- [Overview](#overview)
- [Project structure](#project-structure)
- [Dependencies](#dependencies)
- [Developer Setup](#developer-setup)
- [Running in local cluster minikube](#running-in-local-cluster-minikube)
- [Useful Command/Makefile](#useful-commandmakefile)
- [Database Schema](#database-schema)

## Overview
Buy-Better Core API is the server that act like API gateway with additional logic. This service can 
authenticate with `Oauth2` and authorization before calling other service. It's similar to `Orchestration design pattern`

`NOTE: This project is for learning purpose and not fully complete yet`

## Project structure
```
├── .gen                # auto generated by jet-db
│   ├── buy-better-core
├── bin                 # go binary
├── cmd
│   ├── api             # main package for core api server
│   └── dbhelper        # a dev tool for database helper
├── deploy                    # helm chart
│   ├── bb-core-api     # project helm chart
│   ├── secrets         # bitnami secret
├── internal
│   ├── middleware      # api middleware
│   ├── store           # database logic
│   ├── utils           # global utilities
│   └── v1
│       ├── auth        # Oauth2 handler
│       ├── probe       # liveness & readiness handler
│       └── product     # product handler
│       └── dev
└── migrations          # migration sql
```

## Dependencies
#### Infrastructure
- docker / docker-compose
- minikube
- kubectl / kustomize
- helm
- kubeseal
#### Database tools
- CLI go [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- CLI [jet-db](https://github.com/go-jet/jet?tab=readme-ov-file#installation)

## Developer Setup
1. Create environment variable store in `.env` file at root directory
```
WEB_SERVICE_ENV="dev"
WEB_ADDR=":3030"
WEB_READ_TIMEOUT=5
WEB_WRITE_TIMEOUT=40
WEB_IDLE_TIMEOUT=120
WEB_SHUTDOWN_TIMEOUT=20

DB_DRIVER="postgres"
DB_DSN="postgresql://postgres:admin1234@localhost:5433/buy-better-core?sslmode=disable"
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_IDLE_TIME="15m"

TOKEN_ENCODED="1c0021bc344fa16c72fc522c53bfe9f77a2a597507374e56e3a275759c4c1562"

SESSION_SECRET="your goth session secret"
GOOGLE_CLIENT_ID="your google clientid"
GOOGLE_CLIENT_SECRET="your google secret"
GOOGLE_CALLBACK="http://127.0.0.1:3030/v1/auth/google/callback"

PRODUCT_SERVICE_ADDR="localhost:3031"
```
> For `TOKEN_ENCODED`, you can random generate using [this](https://www.browserling.com/tools/random-hex) and use 64 digits

2. Create postgres environment variable in `.postgres.env` file at root directory. This will be used by `docker-compose`.

```
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="admin1234"
POSTGRES_DB="buy-better-core"
```
> postgres environment variables must be match with Makefile

3. Visit `Makefile` There are 4 important variables for local development. Feel free to edit.
```
DB_DSN
DB_NAME
DB_USERNAME
CONTAINER_NAME		
```

4. Start the development postgres db `make dev-db-up` this command does follow
    * docker-compose with postgres image
    * sleep for 3 seconds
    * migrate up
    * seed the fake data with `dbhelper` tool

5. `go run cmd/api` start the server with port `:3030`

## Running in local cluster minikube
1.  create `encoded-secret.yaml` under k8s/secret/dev
```
apiVersion: v1
kind: Secret
metadata:
  name: core-api-secret
  namespace: buy-better

type: Opaque
data:
  TOKEN_ENCODED: "your base64 encode"
  GOOGLE_CLIENT_ID: "your base64 encode"
  GOOGLE_CLIENT_SECRET: "your base64 encode"
  SESSION_SECRET: "your base64 encode"
```
2. `make dev-up-all`
   * starting postgres db -> migration -> seed fake data
   * starting minikube
   * apply `bitnami-sealed-secrets` controller
3. `minikbue tunnel` to expose load balancer
4. `make dev-apply`
   * go mod tidy
   * building an image with docker
   * kustomize apply resources
   * generate and apply bitnami secret
   * restart deployment (due to bitnami seal secret controller changing certificate everytime when starting a new cluster)

## Useful Command/Makefile

Please visit `Makefile` for the full command.
- `make jet-gen` generate a type safe from database. run this command everytime there is a change in database schema.
- `make dev-db-reset` restart the postgres container. run when you want to reset the database

## Database Schema
![db](https://github.com/opplieam/bb-core-api/blob/main/Buy-Better-Core.png?raw=true)