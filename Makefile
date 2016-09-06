all: clean build

clean:
	rm -f apollo:

build: syncproto genproto
	go build .

genproto:
	./proto/gen_go.sh

syncproto:
	cd proto && git pull origin master

install:
	git submodule update --init
	glide install
