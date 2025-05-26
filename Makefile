REGISTRY_URL="docker.io"
KAFKA = zookeeper \
		kafka \
		kafka-ui

SUBSTANANCE = mongo \
			  redis \
			  postgres

UI = mongo-express \
	 redis-ui

SERVICES = task-service \
		   auth-service \
		   notification-service \
		   gateway-service


.PHONY: all
all: help

.PHONY: main
main:
	docker-compose -f ./docker-compose.yml up $(SUBSTANANCE)

.PHONY: kafka
kafka:
	docker-compose -f ./docker-compose.yml up $(KAFKA)

.PHONY: ui
ui:
	docker-compose -f ./docker-compose.yml up $(UI)

.PHONY: services
services:
	docker-compose -f ./docker-compose.yml up $(SERVICES)

.PHONY: services-stop
services-stop:
	docker stop $(SERVICES)

.PHONY: all-stop
all-stop:
	docker stop $(SERVICES) $(UI) $(SUBSTANANCE) $(KAFKA)

.PHONY: test-env
test-env: main

.PHONY: test
test:
	export MONGODB_CONNECTION_STRING=mongodb://root:password@localhost:27017; go test -v ./internal/repository/user;
	export MONGODB_CONNECTION_STRING=mongodb://root:password@localhost:27017; go test -bench=. ./internal/repository/user -benchmem -benchtime=100x

download:
	echo download
