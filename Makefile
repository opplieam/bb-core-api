# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

tidy:
	go mod tidy

dev-up:
	minikube start --memory 5000 --cpus 2
	helm repo add bitnami-labs https://bitnami-labs.github.io/sealed-secrets/
	helm install bb-sealed-secrets bitnami-labs/sealed-secrets --version 2.16.2 --namespace kube-system
	helm upgrade --install ingress-nginx ingress-nginx \
  		--repo https://kubernetes.github.io/ingress-nginx \
  		--namespace ingress-nginx --create-namespace
	#kubectl apply -f ./k8s/secret/bitnami-sealed-secrets-v0.27.1.yaml

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

docker-push:
	docker push $(SERVICE_IMAGE)

docker-build-push: docker-build-prod docker-push

gen-prod-chart:
	rm -rf .genmanifest
	helm template $(SERVICE_NAME) chart -f chart/values.yaml \
		--set image=$(SERVICE_IMAGE) \
		--set imagePullPolicy=Always \
		--set env.webService=prod \
		--output-dir .genmanifest

helm-prod:
	helm upgrade --install -f ./chart/values.yaml \
	--set image=$(SERVICE_IMAGE) \
	--set imagePullPolicy=Always \
	--set env.webService=prod \
	bb-core-api ./chart

kus-dev:
	kubectl apply -k k8s/dev/
helm-dev:
	helm upgrade --install -f ./chart/values.yaml bb-core-api ./chart
dev-restart:
	kubectl rollout restart deployment $(DEPLOYMENT_NAME) --namespace=$(NAMESPACE)
dev-stop:
	kubectl delete -k k8s/dev/
dev-helm-stop:
	helm uninstall bb-core-api

dev-apply: tidy docker-build-dev kus-dev apply-secret-dev dev-restart

dev-apply-helm: tidy docker-build-dev helm-dev apply-secret-dev dev-restart

# ------------------------------------------------------------
# Seal secret
apply-seal-controller:
	helm upgrade --install bb-sealed-secrets bitnami-labs/sealed-secrets --version 2.16.2 --namespace kube-system
seal-fetch-cert-dev:
	kubeseal --controller-name bb-sealed-secrets --fetch-cert > ./k8s/secret/dev/publickey.pem
seal-fetch-cert-prod:
	kubeseal --controller-name bb-sealed-secrets --fetch-cert > ./k8s/secret/prod/publickey.pem
seal-secret-dev:
	kubeseal --controller-name bb-sealed-secrets --cert ./k8s/secret/dev/publickey.pem < ./k8s/secret/dev/encoded-secret-dev.yaml > ./k8s/secret/dev/sealed-env-dev.yaml
seal-secret-prod:
	kubeseal --controller-name bb-sealed-secrets --cert ./k8s/secret/prod/publickey.pem < ./k8s/secret/prod/encoded-secret-prod.yaml > ./k8s/secret/prod/sealed-env-prod.yaml
apply-seal-dev:
	kubectl apply -f ./k8s/secret/dev/sealed-env-dev.yaml
apply-seal-prod:
	kubectl apply -f ./k8s/secret/prod/sealed-env-prod.yaml

apply-secret-dev: seal-fetch-cert-dev seal-secret-dev apply-seal-dev
apply-secret-prod: seal-fetch-cert-prod seal-secret-prod apply-seal-prod

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

test-unit:
	go test ./... -short

# ------------------------------------------------------------
# Helper function
sleep-%:
	sleep $(@:sleep-%=%)