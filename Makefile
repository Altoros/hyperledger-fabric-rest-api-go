.PHONY: all dev clean build up down run restart

all: clean build up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
up:
	@echo "Start environment ..."
	@cd test/fixtures && docker-compose up --force-recreate -d
	@echo "Environment up"

down:
	@echo "Stop environment ..."
	@cd test/fixtures && docker-compose down
	@echo "Environment down"

restart: clean up

##### RUN
run:
	@echo "Start app ..."
	@./heroes-service

##### CLEAN
clean: down
	@echo "Clean up ..."
	@docker rm -f -v `docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
