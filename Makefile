.PHONY: run dev binary setup glide start-mysql stop-mysql test update
SHELL := /bin/bash

all: run

run: binary
	scripts/run.sh

dev: stop-mysql start-mysql
	scripts/dev.sh

binary:
	GOARCH=amd64 GOOS=linux go build -i -o smartdev

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

start-mysql:
	docker run --name smartdevdb \
		-e MYSQL_ROOT_PASSWORD=blibb \
		-e MYSQL_DATABASE=smartdevdb \
		-e MYSQL_USER=blubb \
		-e MYSQL_PASSWORD=blabb \
		-p "3306:3306" \
		-d mariadb:10
	sleep 10
	scripts/db_setup.sh

stop-mysql:
	docker kill smartdevdb || true
	docker rm -f smartdevdb || true

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)

update:
	git checkout master
	git fetch --all
	git merge upstream/master
	git push
