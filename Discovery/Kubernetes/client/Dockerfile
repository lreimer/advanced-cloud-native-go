FROM golang:1.8.1-alpine

ENV SOURCES /go/src/github.com/lreimer/advanced-cloud-native-go/Discovery/Kubernetes/
COPY . ${SOURCES}

RUN cd ${SOURCES}client/ && CGO_ENABLED=0 go build

WORKDIR ${SOURCES}client/
CMD ${SOURCES}client/client
