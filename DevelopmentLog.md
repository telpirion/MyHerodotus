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
