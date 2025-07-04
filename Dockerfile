FROM golang:1.22-alpine

RUN mkdir build

# We create folder named build for our app.
WORKDIR /build

COPY go.mod go.sum ./

# Download dependencies
RUN go install github.com/air-verse/air@latest
RUN go mod download