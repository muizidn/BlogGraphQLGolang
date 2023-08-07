# syntax=docker/dockerfile:1

FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /blog-graphql-golang

CMD ["/blog-graphql-golang"]

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...