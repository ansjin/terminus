FROM golang:latest as builder
# Copy the local package files to the container's workspace.
ENV GOBIN /go/bin
RUN export GOBIN=$GOPATH/bin
RUN export PATH=$GOPATH:$GOBIN:$PATH

RUN  apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    curl https://glide.sh/get | sh && \
    apt-get install -y git

RUN mkdir -p /go/src/TERMINUS
WORKDIR /go/src/TERMINUS

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
    go get github.com/hpcloud/tail/...

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

# STEP 2 build a small image
# start from scratch
FROM cloudd8/terminus-kops

# Copy our static executable
COPY --from=builder /go/src/TERMINUS/main .
COPY --from=builder /go/src/TERMINUS/assets ./assets/

# Document that the service listens on port 8080.
EXPOSE 8081
# Run the outyet command by default when the container starts.
ENTRYPOINT ["./main"]