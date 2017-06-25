FROM golang:1.8.1-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN go get github.com/Shopify/sarama

ENV SOURCES /go/src/github.com/lreimer/advanced-cloud-native-go/Communication/Kafka/
COPY . ${SOURCES}

RUN cd ${SOURCES}subscriber/ && CGO_ENABLED=0 go build

ENV BROKER_ADDR localhost:9092

WORKDIR ${SOURCES}subscriber/
CMD ${SOURCES}subscriber/subscriber
