# Week 1 Development Log

## Architecture

TODO(telpirion): Insert architectural diagram

+ **Presentation tier**: Bulma, plain-old JS, Firebase auth
+ **Application tier**: Go, Gin templates, Vertex AI model, Cloud Logging
+ **Data tier**: Firestore

## Objectives

For this week's activities, we need to accomplish the following:

- [x] Deploy a Gemma2 model to an endpoint.

   _I've chosen to deploy to Vertex._

- [x] Build a sample chat app that connects to the model. 

  _I built this with Go, Gin, Bulma (FE)._

- [x]  Deploy the sample chat app.

  _The app has been deployed [here](https://myherodotus-1025771077852.us-west1.run.app/)._

- [x]  Instrument the application for Cloud Observability (Logging).

  _I have instrumented the application for Cloud Logging._

- [ ]  Instrument the application for Cloud Observability (Monitoring).

- [x] Persist model interactions into a Database. 

  _I have integrated Firestore into the app._

- [x]  Identify data that needs to be persisted to make response history useful.

  _I have integrated Firebase auth into the app. This asks users to sign in so that their interactions are
  stored keyed into the user's email. I may want to separately store query & responses from the models to
  track their accuracy._

## Tracking future upgrades to app

- [ ] Provide feedback mechanism for users to rate responses.
- [ ] Track user feedback across sessions and users.
- [ ] Tag user feedback with model type, endpoint ID, prompt

## Learning how to create a templated web server with Go

Sources:

+ 👍 https://go.dev/play/p/f4HPQ_WKv9
+ 👍 https://stackoverflow.com/questions/19546896/golang-template-how-to-render-templates
+ https://chenyitian.gitbooks.io/gin-tutorials/content/gin/8.html
+ 👎  https://stackoverflow.com/questions/66658503/how-to-render-html-template-in-gin-gonic-golang
+ 👎  https://gin-gonic.com/docs/examples/html-rendering/#custom-template-renderer

## Learning how to style the web app

Sources:

+ 👍 https://fonts.google.com/icons
+ 👍 https://bulma.io/documentation/elements/button/

## Integrating Firestore

<span style="background-color: red;">Unhappy</span> Updating array fields in Firestore is hard. Too hard. Go doesn't support the `arrayUnion` operation in Firestore :/

Switched to SubCollection of document

## Deploying Gemma model from Model Garden

<span style="background-color: red;">Unhappy</span> Try 1: tried deployment from Pantheon. I don't think it worked ... :/

The activity bar in Pantheon says that something is happening, but that's the only indication I have that a Gemma model is being deployed.

I'm going to use Gemini 1.5 Flash to continue prototyping.

<span style="background-color: yellow;">Anxious</span> I used the [official Gemma2 prediction sample](https://github.com/GoogleCloudPlatform/golang-samples/pull/4395/files) but
the output is awful. I think we need to change the temperature, top-p, and top-k settings.

<span style="background-color: red;">Unhappy</span> Gemma2 responses are ... awful. Even with changing the top-p and temperature settings. It looks like the parameters setting
is not required; I will remove it.

Even with removing the temperature settings, the responses are garbage. If I were releasing this publicly, I would rely
on a more accurate model like Gemini. The Gemma2 model may require better tuning for good results.

Overall,  we need to do a better job demonstrating how to prompt models and decode their responses, especially if we're connecting
directly to an Endpoint.

I had to phutz around with regexes and string parsing to extract the answer from the model.

Sources: 

+ https://pkg.go.dev/regexp#Regexp.FindAll
+ 👎 https://github.com/GoogleCloudPlatform/golang-samples/blob/4806377f298263add9ece7569940fd6b8cfde488/vertexai/gemma2/gemma2_predict_tpu.go
+ 👎 https://cloud.google.com/vertex-ai/generative-ai/docs/open-models/use-gemma
+ 👎 https://ai.google.dev/gemma/docs/gemma_chat

## Building Go on Cloud Shell

<span style="background-color: red;">Unhappy</span>Cloud Shell ran out of room with my successive Go builds. I had to clean the cache.

```sh
go clean -cache
```

Sources:

+ https://stackoverflow.com/questions/66564830/google-cloud-shell-no-space-left-on-device-even-though-disk-not-full

## Integrating Firebase auth

<span style="background-color: red;">Unhappy</span> Setting up my Authentication section ... To set up Google Auth, I need to provide a Web client ID and Web client secret.
I'm sure that I can find it, but I have to go hunting around for it

After finally hacking something that works (not documented), I can't get the redirect to work.

Sources:
+ 👎 (OLD?) https://github.com/firebase/firebaseui-web/blob/master/README.md
+ https://cloud.google.com/docs/authentication/use-cases#app-users
+ https://firebase.google.com/docs/auth/web/redirect-best-practices
+ https://firebase.google.com/docs/auth/where-to-start
+ https://firebase.google.com/docs/auth/web/firebaseui
+ https://firebase.google.com/docs/auth/web/start
+ 👎 (OLD?) https://firebase.google.com/docs/auth/web/firebaseui#email_address_and_password
+ https://firebase.google.com/docs/web/setup#add-sdk-and-initialize
+ https://firebase.google.com/docs/auth/web/google-signin

## Deploying

<span style="background-color: red;">Unhappy</span>I'm a novice with Dockerfiles. I wanted to make sure that all my Go and client files are deployed to the image.
(They are.) But then I need to make sure that the application knows how to load them (it doesn't). How do I test
this?

I decided to build & run the Dockerfile locally (in Cloud Shell) to check that it works correct. The deployment
process to Cloud Run/Build is sooooo long and I want to make sure that I get it right.

<span style="background-color: red;">Unhappy</span>It's very frustrating that I can't install the Docker client on my work laptop. It would
be so much more convenient for me to build & run locally so that I can inspect the built artifacts and logs :/

To build & run from project root:

```sh
docker build . -t myherodotus -f Dockerfile

docker run -it --rm -p 8080:8080 --name myherodotus-running myherodotus
```

To kill the Docker container that is running:

```sh
# Show all the running containers
docker ps -a

docker kill [CONTAINER_NAME]
```

To create the Docker repository in Artifact Registry:

```sh
gcloud artifacts repositories create my-herodotus --repository-format=docker \
    --location=us-west1 --description="Docker repository" \
    --project=${PROJECT_ID}
```

To upload the Docker image to Artifact Registry (where `vSEMVER` is the tag in format "v.NN.NN.NN")

```sh
docker tag myherodotus \
us-west1-docker.pkg.dev/${PROJECT_ID}/my-herodotus/base-image:${SEMVER}

docker push us-west1-docker.pkg.dev/${PROJECT_ID}/my-herodotus/base-image:${SEMVER}
```

<span style="background-color: rgb(100, 200, 150);">Happy</span> Deploying a new version of my web app from Artifact Registry was shockingly intuitive.

Sources:

+ https://github.com/telpirion/telpirion_com/blob/main/README.md
+ https://medium.com/@manzurulhoque/deploying-a-golang-web-app-to-google-cloud-run-a-step-by-step-guide-619e6bb1836e
+ https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service
+ https://phoenixnap.com/kb/docker-environment-variables

## Integrating Cloud Observability (Logging)

The basic quickstart is in Python only :/.

<span style="background-color: yellow;">Anxious</span> The version of the tutorial in Go says "use standard logging," but it shows how to use the
cloud.google.com/go/logging library (not the standard `log` package).

<span style="background-color: rgb(100, 150, 200);">Curious</span> It seems like reinitializing the LoggingClient each time I need to log a message is
a bit repetitive. I wonder if there is a better pattern for this?

Sources:

+ https://cloud.google.com/logging/docs/setup/go
+ https://cloud.google.com/logging/docs/write-query-log-entries-python

# Hardening authentication (client- & server-side)

Realized over the weekend that users can navigate directly to /home without authenticating.

+ [x] Add a check (on the server) that the user's email is propogated. 
+ [x] Add a check (in the JS) that the user has signed-in

TODO(telpirion): ensure that user emails are stored in a obfuscated manner

Sources:

+ https://gin-gonic.com/docs/examples/redirects/
+ https://firebase.google.com/docs/auth/web/manage-users 

## Integrating Cloud Monitoring (Open Telemetry)