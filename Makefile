all: compile docker

compile:
	CGO_ENABLED=0 go build -o bibletagapi

docker:
	sudo docker build --force-rm=true -t bibletagapi/bibletagapi .