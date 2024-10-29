FROM golang:1.23 AS build

# Avoid dynamic linking of libc, since we are using a different deployment image
# that might have a different version of libc.
ENV CGO_ENABLED=0
ENV ENDPOINT_ID=3122353538139684864
ENV COLLECTION_NAME=HerodotusStaging
ENV LOGGER_NAME=HerodotusStaging
ENV CONFIGURATION_NAME=Herodotus
ENV TUNED_MODEL_ENDPOINT_ID=1926929312049528832

# TODO(telpirion): Change this before checking in
#ENV PROJECT_ID=definitely-not-my-project

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY js ./js
COPY templates ./templates
COPY favicon.ico ./favicon.ico

RUN go build -o main .

# Set the entry point command to run the built binary
CMD ["./main"]

EXPOSE 8080