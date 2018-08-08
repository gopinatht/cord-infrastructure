# build stage
FROM golang:alpine AS build-env
RUN apk add --update --no-cache alpine-sdk bash python
WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN /root/librdkafka/configure
RUN make
RUN make install
ADD . /src
WORKDIR /src
RUN cd /src && go get -d ./... && go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o cord-infra-sidecar ./src

# final stage
FROM ubuntu:16.04
COPY --from=build-env /src/cord-infra-sidecar /
ENTRYPOINT ["/cord-infra-sidecar"]
