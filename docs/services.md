# Services: Instructions

The MyHerodotus app has several microservices running behind the scenes to assist
with data collection, evaluation, model tuning, and creating embeddings.

## Data collection

The [Firestore-to-BigQuery](../services/data-collection/) service updates a
BigQuery table with user data and responses from the MyHerodotus app. The service
is triggered by a specific event: when a document is updated in the Firestore database.
This event is used for data collection because it only occurs when a user has rated
a response provided by the app.

All data collected has had PII removed from it, specifically first name, last name,
age, and email addresses. This list of deidentified info types is configurable in
the app.

The following code shows the equivalent gcloud command for exporting from Firestore.

```sh
$ gcloud firestore export gs://myherodotus --database=l200 --collection-ids=HerodotusDev,Conversations
```

### Deploy the service to Cloud functions

To deploy the `data-collection` function to Cloud Run, run the following command from the
`data-collection/` directory. Be sure to set the project ID using `gcloud config set project`.

**IMPORTANT**: Make sure that `$PROJECT_ID` and `$DATASET_NAME` env vars are set before deploying
the function!

```sh
$ gcloud functions deploy data-collection \
  --gen2 \
  --runtime=go121 \
  --region="us-west1" \
  --trigger-location="us-west1" \
  --source=. \
  --entry-point=CollectData \
  --set-env-vars PROJECT_ID=${PROJECT_ID},DATASET_NAME=${DATASET_NAME},BUILD_VER=Herodotus \
  --trigger-event-filters="type=google.cloud.firestore.document.v1.updated" \
  --trigger-event-filters="database=l200" \
  --trigger-event-filters-path-pattern=document='Herodotus/{userId}/Conversations/{conversationId}'
```

### Sources

+ https://cloud.google.com/functions/docs/calling/cloud-firestore
+ https://cloud.google.com/functions/docs/tutorials/storage
+ https://cloud.google.com/functions/docs/calling/eventarc
+ https://cloud.google.com/eventarc/docs/reference/supported-events#cloud-firestore
+ https://cloud.google.com/bigquery/docs/loading-data-cloud-firestore#python
+ https://cloud.google.com/firestore/docs/manage-data/export-import#gcloud


## Evaluations

The [evaluations](../services/evaluations/) microservice runs as a job on [Cloud Run][jobs].
The microservice, written in Python (to make use of data science libraries), is packaged as a
container image, uploaded to Artifact Registry, and then executed as a job.

### Run the job locally

From the root of the evaluations microservice, run the following commands to build and run
the evaluation microservice. Be sure to set the `PROJECT_ID` and `DATASET_NAME` environment variables
before running this service.

```sh
$ docker build . -t evaluations -f Dockerfile
$ docker run -e PROJECT_ID=$PROJECT_ID -e DATASET_NAME=$DATASET_NAME -it --rm --name evaluations-running evaluations 
```

## Embeddings

The [embeddings](../services/embeddings/) microservice uses [PyTorch][pytorch] to
generate custom embeddings from a specified text, in this case the _Histories_ by
Herodotus. To build this

This service can be run either as a Vertex AI job or as a Cloud Run job.

**TIP**: To list the available PyTorch learning containers, visit the
[list of PyTorch containers][pytorch-containers].


### Run the job locally

**NOTE**: The size required for building this Docker image will likely cause Cloud
Shell to run out of space. Instead of running locally, you might need to send
the request to build this image to [Cloud Build][build].

1. Set the following environment variables.

  + `PROJECT_ID`
  + `BUCKET_NAME`
  + `OUTPUT_PATH` (optional)

1. From the root of the embeddings microservice, run the following commands.

    ```sh
    $ ./build.sh
    $ docker run -e PROJECT_ID=$PROJECT_ID -e BUCKET_NAME=$BUCKET_NAME -it --rm --name embeddings-running embeddings 
    ```

### Build the Docker container on Cloud Build

1. Set the following environment variables.

  + `PROJECT_ID`
  + `SEMVER`

From the root of the embeddings microservice, run the following command.

```sh
$ gcloud builds submit --region=us-west1 --config cloudbuild.yaml
```

### Run the job on Cloud Run

Run the following command, making sure that you have the `PROJECT_ID` and `SEMVER` environment variables set.

```sh
$ gcloud run jobs create embeddings \
  --region us-west1 \
  --image us-west1-docker.pkg.dev/${PROJECT_ID}/my-herodotus/embeddings:${SEMVER}
$ gcloud run jobs execute embeddings --region us-west1
```

### Sources

+ https://cloud.google.com/run/docs/create-jobs#client-libraries
+ https://cloud.google.com/build/docs/build-push-docker-image#build_an_image_using_a_build_config_file
+ https://cloud.google.com/build/docs/build-config-file-schema
+ https://cloud.google.com/deep-learning-containers/docs/choosing-container#pytorch
+ https://cloud.google.com/vertex-ai/docs/workbench/user-managed/custom-container#make_sure_your_custom_container_is_ready
+ https://pytorch.org/tutorials/beginner/saving_loading_models.html

[build]: https://cloud.google.com/build/docs/build-push-docker-image
[jobs]: https://cloud.google.com/run/docs/create-jobs
[pytorch]: https://pytorch.org/
[pytorch-containers]: https://cloud.google.com/deep-learning-containers/docs/choosing-container#pytorch