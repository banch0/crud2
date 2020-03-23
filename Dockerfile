FROM golang:latest

RUN mkdir -p $GOPATH/src/app

WORKDIR $GOPATH/src/app

COPY . $GOPATH/src/app

EXPOSE 8888

VOLUME $GOPATH/src

CMD ["go", "run", "./cmd/crud/main.go"]