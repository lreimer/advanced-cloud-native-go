FROM golang:1.8.1-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN go get -u github.com/micro/micro && \
    go get github.com/micro/protobuf/proto && \
    go get -u github.com/micro/protobuf/protoc-gen-go && \
    go get github.com/micro/go-plugins/wrapper/breaker/hystrix && \
    go get github.com/afex/hystrix-go/hystrix

ENV SOURCES /go/src/github.com/lreimer/advanced-cloud-native-go/Communication/Go-Micro/
COPY . ${SOURCES}

RUN cd ${SOURCES}client/ && CGO_ENABLED=0 go build

ENV CONSUL_HTTP_ADDR localhost:8500

WORKDIR ${SOURCES}client/
CMD ${SOURCES}client/client
