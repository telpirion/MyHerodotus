# Friction Log

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

## Architecture

Presentation tier: Bulma, plain-old JS, (TODO) Firebase sign-in
Application tier: Go, Gin templates, Vertex AI model
Data tier: Firestore

## Integrating Firestore

<<Unhappy>> Updating array fields in Firestore is hard. Too hard. Go doesn't support the `arrayUnion` operation in Firestore :/

Switched to SubCollection of document

## Deploying Gemma model from Model Garden

<<Unhappy>> Try 1: tried deployment from Pantheon. I don't think it worked ... :/

The activity bar in Pantheon says that something is happening, but that's the only indication I have that a Gemma model is being deployed.

I'm going to use Gemini 1.5 Flash to continue prototyping.

<<Anxious>> I used the [official Gemma2 prediction sample](https://github.com/GoogleCloudPlatform/golang-samples/pull/4395/files) but
the output is awful. I think we need to change the temperature, top-p, and top-k settings.

<<Unhappy>> Gemma2 responses are ... awful. Even with changing the top-p and temperature settings. It looks like the parameters setting
is not required; I will remove it.

Overall,  we need to do a better job demonstrating how to prompt models and decode their responses, especially if we're connecting
directly to an Endpoint.

I had to phutz around with regexes and string parsing to extract the answer from the model.

Sources: 

+ https://pkg.go.dev/regexp#Regexp.FindAll
+ 👎 https://github.com/GoogleCloudPlatform/golang-samples/blob/4806377f298263add9ece7569940fd6b8cfde488/vertexai/gemma2/gemma2_predict_tpu.go
+ 👎 https://cloud.google.com/vertex-ai/generative-ai/docs/open-models/use-gemma
+ 👎 https://ai.google.dev/gemma/docs/gemma_chat

## Building Go on Cloud Shell


<<Unhappy>>Cloud Shell ran out of room with my successive Go builds. I had to clean the cache.

```sh
go clean -cache
```

Sources:

+ https://stackoverflow.com/questions/66564830/google-cloud-shell-no-space-left-on-device-even-though-disk-not-full

## Integrating Firebase auth

<<Unhappy>> Setting up my Authentication section ... To set up Google Auth, I need to provide a Web client ID and Web client secret.
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