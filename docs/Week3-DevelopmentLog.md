# Week 3 Development Log

+ Unhappy: 👎
+ Anxious: 😬
+ Curious: 🤔
+ Happy: 👍

## Objectives

For this week's activities, we must do the following:

- [x] Retrieve stored conversations with the user.
- [x] Set the user's past conversations as context with Vertex (Gemini).
- [x] Set the user's past conversations as context in template (Gemma).
- [x] Find data about how to fine tune Gemini, Gemma, or Gemma 2.
- [ ] Fine tune a Gemini, Gemma, or Gemma 2 model using the OpenAssistant
      guanaco dataset on HuggingFace.
- [ ] Save the fine-tuned model to Model Registry, Model Garden, or HuggingFace.
- [ ] Deploy the model to an endpoint.
- [ ] Integrate the model into the web app.

Nice-to-haves:

- [ ] Create a toast UI element that informs the user when their response rating was received.

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

+ 🤔 Storing the context as a RAG part of the prompt seems to be working okay for Gemini. I wonder what would happen if I
  used Gemini outputs as context for the Gemma prompt? Do I need to filter out Gemma & Gemini context histories?
 
  - 👎 Oof, the model failed with a cryptic `rpc error: code = Internal desc = {"error":"Incomplete generation","error_type":"Incomplete generation"}`
    error. I honestly don't know how to debug that error...
  - Looking in the logs, I see that there were TOO many new tokens for the Gemma model:

    ```json
  {"timestamp":"2024-10-25T22:35:49.174225Z","level":"ERROR","message":"`inputs` tokens + `max_new_tokens` must be <= 2048. Given: 4375 `inputs` tokens and 100 `max_new_tokens`","target":"text_generation_router::infer","filename":"router/src/infer/mod.rs","line_number":102,"span":{"name":"generate_stream"},"spans":[{"name":"vertex_compatibility"},{"name":"generate"},{"name":"generate_stream"}]}
    ```
  
  - Looking into Go tokenizers ... it looks like `SentencePiece` is what Gemma uses, which has a C++ binary associated with it (?).
    * https://github.com/eliben/go-sentencepiece
    * go-sentence piece needs this file: https://github.com/google/gemma_pytorch/blob/main/tokenizer/tokenizer.model 
    * https://github.com/google/sentencepiece
  
  - Other open source tokenizers:
    * https://github.com/tiktoken-go/tokenizer (for OpenAI models)

  - 🤔 Maybe we need to record the number of tokens in each ConversationBit, and then only collect the first 2000.
    We might be able to use the Firestore aggregation filters to get this out-of-the-box.

Sources:

+ https://go.dev/play/p/4rLkXhW570p
+ 👎 https://cloud.google.com/vertex-ai/generative-ai/docs/context-cache/context-cache-create
+ 👎 https://cloud.google.com/vertex-ai/generative-ai/docs/context-cache/context-cache-use
+ https://pkg.go.dev/errors#As 
+ https://stackoverflow.com/questions/54156119/range-over-string-slice-in-golang-template
+ https://ai.google.dev/gemma/docs/model_card_2  

## Creating a UI toast

+ Going to use CSS animations to show and hide the toast notification, using keyframes.
+ Getting the timing just right is the tough part, making sure that the notification shows
  and then is hidden again.

Sources:

+ https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_animations/Using_CSS_animations
+ 👎 https://stackoverflow.com/questions/16670931/hide-scroll-bar-but-while-still-being-able-to-scroll


## Finding information about fine tuning

+ 😬 It looks like some of the versions of the Guanaco dataset aren't well supported. One version of the
  dataset said that there could be some inappropriate content in the dataset.
+ Fine tuning seemingly is only documented for Python. It _should_ be possible in other languages
  that the Vertex client library is available in.

Sources:

+ https://cloud.google.com/vertex-ai/generative-ai/docs/models/gemini-use-supervised-tuning#python
+ https://huggingface.co/datasets/timdettmers/openassistant-guanaco
+ https://guanaco-model.github.io