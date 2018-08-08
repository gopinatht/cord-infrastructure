appname = cord-infra-sidecar
app_output = ./$(appname)
docker_repo = gopinatht
docker_tag = latest
application_directory = ./src
docker_directory = .

all: push

image:
	docker build -t $(docker_repo)/$(appname):$(docker_tag) $(docker_directory)

push: image
	docker push $(docker_repo)/$(appname):$(docker_tag)
