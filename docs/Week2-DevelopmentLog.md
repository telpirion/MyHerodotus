# Week 2 Development Log

## Architecture

TODO(telpirion): Insert architectural diagram

+ **Presentation tier**: Bulma, plain-old JS, Firebase auth
+ **Application tier**: Go, Gin templates, Vertex AI model, Cloud Logging
+ **Data tier**: Firestore

## Objectives

For this week's activities, we need to accomplish the following:

- [x] Extend app to connect to a Gemini endpoint instead of a Gemma endpoint.

   _I've integrated Gemini 1.5 Flash into the app._

- [ ] Add version control for prompts. 
- [ ] Select version of prompt based upon model being called.
- [ ] Experiment with prompting strategies to evaluate outcomes.
- [ ] Provide feedback mechanism for users to rate responses.
- [ ] Track user feedback across sessions and users.
- [ ] Tag user feedback with model type, endpoint ID, prompt