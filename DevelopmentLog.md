# Friction Log

## Learning how to create a templated web server with Go

Good:

+ https://go.dev/play/p/f4HPQ_WKv9
+ https://stackoverflow.com/questions/19546896/golang-template-how-to-render-templates

Mediocre:

+ https://chenyitian.gitbooks.io/gin-tutorials/content/gin/8.html

<<Unhappy>> Bad:

+ https://stackoverflow.com/questions/66658503/how-to-render-html-template-in-gin-gonic-golang
+ https://gin-gonic.com/docs/examples/html-rendering/#custom-template-renderer

## Learning how to style the web app

https://fonts.google.com/icons
https://bulma.io/documentation/elements/button/

## Architecture

Presentation tier: Bulma, plain-old JS, (TODO) Firebase sign-in
Application tier: Go, Gin templates, Vertex AI model
Data tier: Firestore

## Firestore

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

## Building Go on Cloud Shell

<<Unhappy>>Cloud Shell ran out of room with my successive Go builds. I had to clean the cache.

```sh
go clean -cache
```