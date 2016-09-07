all: clean build

clean:
	rm -f apollo:

build: genproto
	go build .

genproto:
	./proto/gen_go.sh

syncproto:
	cd proto && git pull origin master

install:
	git submodule update --init
	glide install
