FROM golang:1.23 AS build

# Avoid dynamic linking of libc, since we are using a different deployment image
# that might have a different version of libc.
ENV CGO_ENABLED=0
ENV ENDPOINT_ID=3122353538139684864
ENV COLLECTION_NAME=HerodotusStaging

# TODO(telpirion): Delete this before deploying
ENV PROJECT_ID=erschmid-test-291318

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY js ./js
COPY templates ./templates

RUN go build -o main .

# Set the entry point command to run the built binary
CMD ["./main"]

EXPOSE 8080