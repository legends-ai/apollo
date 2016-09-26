all: clean build

clean:
	rm -f apollo:

build: genproto
	go build .

syncbuild: syncproto genproto
	go build .

genproto:
	./proto/gen_go.sh

syncproto:
	cd proto && git pull origin master

init:
	git submodule update --init

install: init genproto
	glide install

docker-build:
	docker build -t apollo .

docker-push:
	docker tag apollo:latest 096202052535.dkr.ecr.us-west-2.amazonaws.com/apollo:latest
	docker push 096202052535.dkr.ecr.us-west-2.amazonaws.com/apollo:latest

docker-build-dev:
	docker build -t apollo:dev .

docker-push-dev:
	docker tag apollo:dev 096202052535.dkr.ecr.us-west-2.amazonaws.com/apollo:dev
	docker push 096202052535.dkr.ecr.us-west-2.amazonaws.com/apollo:dev

cassandratunnel:
	ssh -fNL 9042:node-0.cassandra.mesos:9042 centos@52.42.186.11
