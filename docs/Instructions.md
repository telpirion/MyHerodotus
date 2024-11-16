# Instructions

This document contains the directions for running, building, and deploying the
MyHerodotus web app.

## Set environment variables

The local development environment must have the following variables:

+ PROJECT_ID
+ ENDPOINT_ID : pointing to the Gemma model deployed endpoint
+ COLLECTION_NAME : the Firestore Collection that contains user documents (conversations)
+ LOGGER_NAME : the name to use for Cloud Logging
+ CONFIGURATION_NAME : the name for this configuration, e.g. "HerodotusDev," "HerodotusStaging"
+ TUNED_MODEL_ENDPOINT_ID : pointing to the tuned Gemini model endpoint

The site must also include a file, `appInit.js`, that contains the Firebase client
configuration.

## Build and run the Docker image

To build and run the Docker image in the local development environment, you must 
run the following commands:

```sh
$ docker build . -t myherodotus -f Dockerfile --build-arg BUILD_VER=HerodotusStaging 
$ docker run -e PROJECT_ID=$PROJECT_ID -it --rm -p 8080:8080 --name myherodotus-running myherodotus 
```

To inspect the env vars of a container while its running, run the following command.

```sh
$ docker ps # to get the ID of the running container
$ docker inspect --format='{{.Config.Env}}' $CONTAINER_ID
```

## Upload a new Docker image to Artifact Registry

To tag and upload a new Docker image to Artifact Registry, run the
following commands. Be sure to set the `PROJECT_ID` and `SEMVER` environment
variables.

```sh
$ docker build . -t myherodotus -f Dockerfile --build-arg BUILD_VER=Herodotus
$ docker tag myherodotus us-west1-docker.pkg.dev/${PROJECT_ID}/my-herodotus/base-image:${SEMVER}
$ docker push us-west1-docker.pkg.dev/${PROJECT_ID}/my-herodotus/base-image:${SEMVER}
```

## Update protobuf files

This project uses [protocol buffer][protobuf] files to define types shared between the
different services and surfaces of the application. These files must be updated and the output
files regenerated when new fields are added.

**NOTE**: The tool registry used by `buf` has a rate limit of 10 unauthenticated requests
per hour :/. If this quota is reached, use [protoc directly instead][#protoc] 

1. Install the [`buf` CLI][buf].

    ```sh
    $ GO111MODULE=on go install github.com/bufbuild/buf/cmd/buf@v1.47.2
    ```

1. Check the installation.

    ```sh
    buf --version
    ```

1. Build the protos. From the `protos/` directory, run the following command.

    ```sh
    buf build
    ```

1. Generate the protos. From the `protos/` directory, run the following command.

    ```sh
    buf generate
    ```

### Generate with protoc-gen-go {#protoc}

1. 

[buf]:      https://buf.build/docs/tutorials/getting-started-with-buf-cli
[protobuf]: https://protobuf.dev/getting-started/gotutorial/