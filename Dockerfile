FROM golang:1.23 AS build

# Avoid dynamic linking of libc, since we are using a different deployment image
# that might have a different version of libc.
ARG BUILD_VER=HerodotusStaging

ENV CGO_ENABLED=0
ENV ENDPOINT_ID=3122353538139684864
ENV COLLECTION_NAME=$BUILD_VER
ENV LOGGER_NAME=$BUILD_VER
ENV CONFIGURATION_NAME=$BUILD_VER
ENV TUNED_MODEL_ENDPOINT_ID=1926929312049528832

WORKDIR /

COPY site/js ./site/js
COPY site/css ./site/css
COPY site/html ./site/html
COPY server/templates ./server/templates
COPY server/favicon.ico ./server/favicon.ico
COPY server/* ./server

COPY server/go.mod server/go.sum ./server/
WORKDIR /server
RUN go mod download

RUN go build -o main .

# Set the entry point command to run the built binary
CMD ["./main"]

EXPOSE 8080