# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

tidy:
	go mod tidy

dev-up:
	minikube start
	kubectl apply -f ./k8s/secret/bitnami-sealed-secrets-v0.27.1.yaml

dev-down:
	minikube delete

dev-up-all: dev-db-up dev-up
dev-down-all: dev-down dev-db-down

DB_DSN				:= "postgresql://postgres:admin1234@localhost:5433/buy-better-core?sslmode=disable"
DB_NAME				:= "buy-better-core"
DB_USERNAME			:= "postgres"
CONTAINER_NAME		:= "pg-core-dev-db"

BASE_IMAGE_NAME 	:= opplieam
SERVICE_NAME    	:= bb-core-api
VERSION         	:= "0.0.1-$(shell git rev-parse --short HEAD)"
VERSION_DEV         := "cluster-dev"
SERVICE_IMAGE   	:= $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
SERVICE_IMAGE_DEV   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION_DEV)

DEPLOYMENT_NAME		:= core-api-deployment
SECRET_NAME			:= core-api-secret
NAMESPACE			:= buy-better

# ------------------------------------------------------------
# Deploy in local cluster

docker-build-dev:
	@eval $$(minikube docker-env);\
	docker build \
		-t $(SERVICE_IMAGE_DEV) \
    	--build-arg BUILD_REF=$(VERSION_DEV) \
    	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    	-f dev.Dockerfile \
    	.
docker-build-prod:
	docker build \
		-t $(SERVICE_IMAGE) \
    	--build-arg BUILD_REF=$(VERSION) \
    	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    	.

kus-dev:
	kubectl apply -k k8s/dev/
dev-restart:
	kubectl rollout restart deployment $(DEPLOYMENT_NAME) --namespace=$(NAMESPACE)
dev-stop:
	kubectl delete -k k8s/dev/

dev-apply: tidy docker-build-dev kus-dev apply-secret dev-restart

# ------------------------------------------------------------
# Seal secret
apply-seal-controller:
	kubectl apply -f ./k8s/secret/bitnami-sealed-secrets-v0.27.1.yaml
seal-fetch-cert:
	kubeseal --fetch-cert > ./k8s/secret/dev/publickey.pem
seal-secret:
	kubeseal --cert ./k8s/secret/dev/publickey.pem < ./k8s/secret/dev/encoded-secret.yaml > ./k8s/secret/dev/sealed-env-dev.yaml
apply-seal:
	kubectl apply -f ./k8s/secret/dev/sealed-env-dev.yaml

apply-secret: seal-fetch-cert seal-secret apply-seal

# ------------------------------------------------------------
# DB
docker-compose-up:
	docker compose up -d
docker-compose-down:
	docker compose down

migrate-up:
	migrate -path=./migrations \
	-database=$(DB_DSN) \
	up

migrate-down:
	migrate -path=./migrations \
    -database=$(DB_DSN) \
    down

dev-db-seed: dbhelper-build dbhelper-seed-all

dbhelper-build:
	go build -o ./bin/dbhelper ./cmd/dbhelper
dbhelper-seed-all:
	./bin/dbhelper

dev-db-up: docker-compose-up sleep-3 migrate-up dev-db-seed
dev-db-down: docker-compose-down
dev-db-reset: dev-db-down sleep-1 dev-db-up

jet-gen:
	jet -dsn=$(DB_DSN) -path=./.gen


# ------------------------------------------------------------
# Helper function
sleep-%:
	sleep $(@:sleep-%=%)