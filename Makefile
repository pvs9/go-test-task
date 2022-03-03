build:
	./docker/build.sh

.PHONY: vendor
vendor:
	docker run --rm -e GO111MODULE=on -v "$(PWD)":/go/src/go-test-task.io vpdev/golang-dev sh -c 'go mod tidy -compat=1.17 && go mod vendor'

up:
	@make build
	@make vendor
	docker-compose -f docker-compose.yml up -d
logs:
	docker-compose logs -f
rm:
	docker-compose rm  -sfv
start:
	docker-compose start
stop:
	docker-compose stop
rest:
	@make rm
	@make up
	@make logs

bash:
	docker-compose exec test-task bash

aws-cli:
	docker run --network=app-network-test  vpdev/awscli --endpoint-url="http://localstack:4566" $(filter-out $@,$(MAKECMDGOALS))

help:
	@echo "make container=vpdev_test-task_1 bash \t\t\t: exec a bash shell in the specific container"
