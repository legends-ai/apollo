all: clean build

clean:
	rm apollo || :

build: genproto
	go build .

genproto:
	./proto/gen_go.sh

syncproto:
	cd proto && git pull origin master
