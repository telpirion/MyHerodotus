FROM golang:1.23 AS build

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
COPY prompts ./server/templates
COPY server/favicon.ico ./server/favicon.ico
COPY server/generated ./server/generated
COPY server/ai ./server/ai
COPY server/databases ./server/databases
COPY server/*.go ./server

COPY server/go.mod server/go.sum ./server/
WORKDIR /server
RUN go mod download

RUN go build -o main .

# Set the entry point command to run the built binary
CMD ["./main"]

EXPOSE 8080