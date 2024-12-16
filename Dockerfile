FROM golang:1.23.2-alpine3.20

RUN mkdir build

# We create folder named build for our app.
WORKDIR /build

COPY go.mod go.sum ./

# Download dependencies
RUN go install github.com/air-verse/air@latest
RUN go mod download