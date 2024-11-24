# Server: Instructions

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
