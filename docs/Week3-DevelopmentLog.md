# Week 3 Development Log

+ Unhappy: 👎
+ Anxious: 😬
+ Curious: 🤔
+ Happy: 👍

## Objectives

For this week's activities, we must do the following:

- [ ] Retrieve stored conversations with the user.
- [ ] Set the user's past conversations as context with Vertex (Gemini).
- [ ] Set the user's past conversations as context in template (Gemma).
- [ ] Find data about how to fine tune Gemini, Gemma, or Gemma 2.
- [ ] Fine tune a Gemini, Gemma, or Gemma 2 model using the OpenAssistant
      guanaco dataset on HuggingFace.
- [ ] Save the fine-tuned model to Model Registry, Model Garden, or HuggingFace.
- [ ] Deploy the model to an endpoint.
- [ ] Integrate the model into the web app.


## Retrieving user context for models

+ 👍 Setting the context cache for Gemini seems really simple. The API really helps with
  doing so.
+ I'm not sure whether the AI Platform API has a field for storing context. This might
  need to be included a la grounding.
+ 👎 Caching with Gen AI fails with because the cache isn't large enough?!

```sh
2024/10/25 19:56:49 error: 
Couldn't store conversation context: CreateCachedContent: rpc error: code = InvalidArgument desc = The cached content is of 2235 tokens. The minimum token count to start caching is 32768.
```

  - 👎 This is from a stored set of 22 back-and-forth query/responses from a real user (me) and the model!!!
  - 🤔 For users that need to provide context to the model, but DON'T have the 32768 token count to start caching, what is the recommended approach?
    RAG? Should our docs maybe recommend an approach?
  - 🤔 How would I know if I have the recommended token count without trial & error? If my ~22 response/replies is only 2235 tokens, assuming that each
    response/reply is about 110 tokens, then a conversation (from my app) needs to have a total of 330 response/replies from a user before caching is
    available ... 
  - 🤔 I wonder ... is it assumed that the system instructions are included in that token amount?

Sources:

+ https://go.dev/play/p/4rLkXhW570p
+ 👎 https://cloud.google.com/vertex-ai/generative-ai/docs/context-cache/context-cache-create
+ 👎 https://cloud.google.com/vertex-ai/generative-ai/docs/context-cache/context-cache-use
+ https://pkg.go.dev/errors#As 
+ https://stackoverflow.com/questions/54156119/range-over-string-slice-in-golang-template 

## Tuning a model

Sources:

+ https://huggingface.co/datasets/timdettmers/openassistant-guanaco