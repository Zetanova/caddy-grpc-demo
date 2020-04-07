FROM golang:1.14-alpine3.11

RUN apk add --no-cache \
    git \
    ca-certificates

WORKDIR $GOPATH/src

RUN go get -u google.golang.org/grpc
	
WORKDIR $GOPATH/src/google.golang.org/grpc/examples/helloworld

RUN go install greeter_server/main.go \
 && go install greeter_server/main.go

EXPOSE 50051

#CMD ["go", "run", "greeter_server/main.go"]
#CMD ["go", "run", "greeter_client/main.go"]
CMD ["go", "version"]
