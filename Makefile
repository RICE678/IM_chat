APP=myapp
VERSION=v1
SERVER=root@47.110.92.199
IMAGE=$(APP):$(VERSION)

build:
	docker build -t $(IMAGE) .


upload: build
	docker save $(IMAGE) | gzip | ssh $(SERVER) "docker load"


deploy: upload
	ssh $(SERVER) "docker stop $(APP) || true && \
	               docker rm $(APP) || true && \
	               docker run -d --name $(APP) \
	               --restart always \
	               -p 8082:8081 $(IMAGE)"

release: deploy