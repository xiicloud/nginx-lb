.PHONY: build clean

IMG_DEBIAN := csphere/nginx-lb:1.11.1.4
IMG_ALPINE := csphere/nginx-lb:1.13-alpine.2
PUSH_REPO := index.csphere.cn
DEST_IMG_DEBIAN := $(PUSH_REPO)/$(IMG_DEBIAN)
DEST_IMG_ALPINE := $(PUSH_REPO)/$(IMG_ALPINE)

build: go-build
	docker build -t $(IMG_DEBIAN) .

alpine-1.13: go-build
	docker build -t $(IMG_ALPINE) -f Dockerfile.1.13-alpine .

go-build:
	(cd init && go build -ldflags '-w')

clean:
	rm -f init/init

push: build alpine-1.13
	docker tag $(IMG_DEBIAN) $(DEST_IMG_DEBIAN)
	docker tag $(IMG_ALPINE) $(DEST_IMG_ALPINE)
	docker push $(DEST_IMG_DEBIAN)
	docker push $(DEST_IMG_ALPINE)
	docker rmi $(DEST_IMG_DEBIAN)
	docker rmi $(DEST_IMG_ALPINE)
