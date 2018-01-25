appname = cord-infra-initializer
app_output = ./$(appname)
docker_repo = gopinatht
docker_tag = latest
application_directory = .
docker_directory = .

all: push

build:
	GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o $(app_output) $(application_directory)

image: build
	docker build -t $(docker_repo)/$(appname):$(docker_tag) $(docker_directory)
	rm -r $(app_output)

push: image
	docker push $(docker_repo)/$(appname):$(docker_tag)
