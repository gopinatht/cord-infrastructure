# build stage
# FROM golang:alpine AS build-env
# RUN apk add --update --no-cache alpine-sdk bash python
# ADD . /src
# WORKDIR /src
# RUN cd /src && go get -d ./... && GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o cord-infra-initializer ./src

# final stage
FROM ubuntu:16.04
#COPY --from=build-env /src/cord-infra-initializer /
COPY ./cord-infra-initializer /cord-infra-initializer
ENTRYPOINT ["/cord-infra-initializer"]
