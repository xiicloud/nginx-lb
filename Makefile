.PHONY: build clean

build:
	(cd init && go build -ldflags '-w')
	docker build -t csphere/nginx-lb .

clean:
	rm -f init/init