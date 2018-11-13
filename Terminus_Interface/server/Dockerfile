FROM golang:latest as builder
# Copy the local package files to the container's workspace.
ENV GOBIN /go/bin
RUN export GOBIN=$GOPATH/bin
RUN export PATH=$GOPATH:$GOBIN:$PATH

RUN  apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    curl https://glide.sh/get | sh && \
    apt-get install -y git

RUN mkdir -p /go/src/TERMINUSINTF
WORKDIR /go/src/TERMINUSINTF

#COPY glide.yaml ./

#RUN glide install

RUN go get gopkg.in/mgo.v2 && \
    go get github.com/aws/aws-sdk-go && \
    git clone https://github.com/golang/crypto.git $GOPATH/src/golang.org/x/crypto && \
    go get github.com/gorilla/mux && \
    go get github.com/swaggo/swag/cmd/swag &&\
    go get github.com/logmatic/logmatic-go &&\
    go get github.com/bshuster-repo/logrus-logstash-hook &&\
    go get github.com/sirupsen/logrus &&\
    go get github.com/influxdata/influxdb/client/v2 &&\
    go get github.com/hpcloud/tail/...

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

# STEP 2 build a small image
# start from scratch
FROM ubuntu:16.04
RUN sed -i -e 's/us.archive.ubuntu.com/archive.ubuntu.com/g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y curl && apt-get install apt-transport-https
RUN curl -sL https://repos.influxdata.com/influxdb.key | apt-key add -
RUN /bin/bash -c 'source /etc/lsb-release && echo "deb https://repos.influxdata.com/${DISTRIB_ID,,} ${DISTRIB_CODENAME} stable"' | tee /etc/apt/sources.list.d/influxdb.list
RUN apt-get update && apt-get install -y influxdb
# Copy our static executable
COPY --from=builder /go/src/TERMINUSINTF/main .
COPY --from=builder /go/src/TERMINUSINTF/assets ./assets/
COPY --from=builder /go/src/TERMINUSINTF/docs ./docs/
COPY --from=builder /go/src/TERMINUSINTF/templates ./templates/

# Document that the service listens on port 8082.
EXPOSE 8082
# Run the outyet command by default when the container starts.
ENTRYPOINT ["./main"]