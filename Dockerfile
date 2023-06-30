FROM golang:1.20.5-alpine

RUN go install github.com/githubnemo/CompileDaemon@latest
#RUN go install github.com/swaggo/swag/cmd/swag@latest

ENV APP_HOME /app
# Set the module cache directory
ENV GOMODCACHE /go/pkg/mod

WORKDIR "$APP_HOME"

#swag init -g ./cmd/main.go &&
ENTRYPOINT CompileDaemon -build="go build ./cmd/main.go" -command="./main"
