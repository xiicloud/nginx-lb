.PHONY: build clean

build: go-build
	docker build -t csphere/nginx-lb:1.11.1.3 .

alpine-1.13: go-build
	docker build -t csphere/nginx-lb:1.13-alpine.1 -f Dockerfile.1.13-alpine .

go-build:
	(cd init && go build -ldflags '-w')

clean:
	rm -f init/init
