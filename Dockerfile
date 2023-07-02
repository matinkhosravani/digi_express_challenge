FROM golang:1.20.5-alpine

RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . /app
ENV APP_HOME /app
# Set the module cache directory
ENV GOMODCACHE /go/pkg/mod

WORKDIR "$APP_HOME"
RUN $GOPATH/bin/swag init -g ./cmd/main.go
RUN go build -v ./cmd/main.go
ENTRYPOINT "./main"
