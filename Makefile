.PHONY: build clean

build:
	(cd init && go build -ldflags '-w')
	docker build -t csphere/nginx-lb:1.11.1.1 .

clean:
	rm -f init/init
