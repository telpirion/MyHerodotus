# Services

The MyHerodotus app has several microservices running behind the scenes to assist
with data collection, evaluation, and model tuning.

## Data collection

The [Firestore-to-BigQuery](../services/firestore-to-bigquery/) service updates a
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


### Sources

+ https://cloud.google.com/functions/docs/tutorials/storage
+ https://cloud.google.com/functions/docs/calling/eventarc
+ https://cloud.google.com/eventarc/docs/reference/supported-events#cloud-firestore
+ https://cloud.google.com/bigquery/docs/loading-data-cloud-firestore#python
+ https://cloud.google.com/firestore/docs/manage-data/export-import#gcloud
