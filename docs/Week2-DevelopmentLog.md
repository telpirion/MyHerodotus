# Week 2 Development Log

<p><span style="background-color: red;">Unhappy</span></p>
<p><span style="background-color: yellow;">Anxious</span></p>
<p><span style="background-color: rgb(100, 150, 200);">Curious</span></p>
<p><span style="background-color: rgb(100, 200, 150);">Happy</span></p>

## Architecture

TODO(telpirion): Insert architectural diagram

+ **Presentation tier**: Bulma, plain-old JS, Firebase auth
+ **Application tier**: Go, Gin templates, Vertex AI model, Cloud Logging
+ **Data tier**: Firestore

## Objectives

For this week's activities, we need to accomplish the following:

- [x] Extend app to connect to a Gemini endpoint instead of a Gemma endpoint.

   _I've integrated Gemini 1.5 Flash into the app._

- [x] Add version control for prompts. 

   _I have moved prompts out into their own, version-controlled files_

- [ ] Select version of prompt based upon model being called.
- [ ] Experiment with prompting strategies to evaluate outcomes.
- [ ] Provide feedback mechanism for users to rate responses.
- [ ] Track user feedback across sessions and users.
- [ ] Tag user feedback with model type, endpoint ID, prompt

## Moving prompts to their own templates

Sources:

+ https://gobyexample.com/text-templates

## Refactoring UI to allow user ratings, model selection toggle, progress bar

<span style="background-color: rgb(100, 150, 200);">Curious</span>
My piecemeal approach to frontend development for this project makes me appreciate
frontend frameworks more. It is difficult to keep track of which HTML templates
go with what JS scripts.

<span style="background-color: red;">Unhappy</span>
The material symbols docs need to
EMPHASIZE that the order of icons must be alphabetically sorted!

<span style="background-color: red;">Unhappy</span>
CONTINUOUS running out of memory errors with Cloud Shell. I wish I could just set
a checkbox somewhere that cleans up the cache every time I close the window!

```sh
du -hs $(ls -A)
go clean --cache
go clean -modcache
```

<span style="background-color: rgb(100, 200, 150);">Happy</span>
Being able to debug server-side AND client-side code in the same interface is AWESOME.

Sources:

+ https://cloud.google.com/shell/docs/quotas-limits
+ https://bulma.io/documentation/form/select/
+ https://developer.mozilla.org/en-US/docs/Web/API/Window/sessionStorage
+ 👎 (CONFUSING) https://developers.google.com/fonts/docs/material_symbols 

