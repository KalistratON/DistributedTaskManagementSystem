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
	docker-compose up $(SUBSTANANCE)

.PHONY: kafka
kafka:
	docker-compose up $(KAFKA)

.PHONY: ui
ui:
	docker-compose up $(UI)

.PHONY: services
services:
	docker-compose up $(SERVICES)

.PHONY: services-stop
services-stop:
	docker stop $(SERVICES)

.PHONY: all-stop
all-stop:
	docker stop $(SERVICES) $(UI) $(SUBSTANANCE) $(KAFKA)

download:
	echo download
