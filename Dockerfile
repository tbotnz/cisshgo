FROM golang:1.15.0-buster
ADD . /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod download
CMD go run cissh.go