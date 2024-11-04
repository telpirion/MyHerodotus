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
- [x] Fine tune a Gemini, Gemma, or Gemma 2 model using the OpenAssistant
      guanaco dataset on HuggingFace.
- [x] Create a notebook for managing data science tasks
- [x] Save the fine-tuned model to Model Registry, Model Garden, or HuggingFace.
- [x] Deploy the model to an endpoint.
- [x] Integrate the model into the web app.

Nice-to-haves:

- [x] Create a toast UI element that informs the user when their response rating was received.
- [x] Refactor Docker file to accept env vars as argument.
- [x] Refactor frontend, backend into separate files.
- [x] Integrate token-counting into DB.


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

+ 👍 COUNT TOKENS! https://github.com/GoogleCloudPlatform/golang-samples/blob/main/vertexai/token-count/tokencount.go 
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
+ Going to try fine-tuning Gemini through the UI (Cloud Console).
+ 🤔 Why is it possible to upload a training set from my local machine but not the validation set?
+ 👎 Tuning the Gemini model, using the HuggingFace dataset, fails because the JSON file is in the
  wrong format. The error message doesn't give me enough details about _how_ the file is in the
  wrong format
+ 🤔 Is this a good dataset to use for fine-tuning? It has multiple languages ... which we don't have
  as a requirement for our model. Why not use a dataset that is specific towards travel?
+ It looks like we need to set up an ETL pipeline first before the dataset is ready for tuning.

Sources:

+ https://cloud.google.com/vertex-ai/generative-ai/docs/models/gemini-use-supervised-tuning#python
+ https://huggingface.co/datasets/timdettmers/openassistant-guanaco
+ https://guanaco-model.github.io
+ https://github.com/GoogleCloudPlatform/vertex-ai-samples/blob/main/notebooks/official/training/pytorch_gcs_data_training.ipynb
+ https://cloud.google.com/blog/topics/developers-practitioners/pytorch-google-cloud-how-train-and-tune-pytorch-models-vertex-ai
+ https://cloud.google.com/vertex-ai/generative-ai/docs/models/tune_gemini/text_tune
+ https://cloud.google.com/vertex-ai/generative-ai/docs/models/tune-models
+ https://pkg.go.dev/cloud.google.com/go/aiplatform@v1.68.0/apiv1/aiplatformpb#CreateCustomJobRequest 
+ Other datasets:
  - https://huggingface.co/datasets/osunlp/TravelPlanner
  - https://www.kaggle.com/datasets/niraliivaghani/chatbot-dataset
  - https://huggingface.co/datasets/Binaryy/reddit-travel-qa?row=1
  - https://huggingface.co/datasets/NLPC-UOM/Travel-Dataset-5000
  - https://huggingface.co/datasets/Binaryy/travel_sample_extended
  - https://huggingface.co/datasets/thari01/travel_data/viewer/default/train?p=1
  - https://huggingface.co/datasets/soniawmeyer/travel-conversations-finetuning?row=0

## Fine tune the Gemini model

+ It looks like I will need to transform the Guanaco dataset into a JSON structure that
  Vertex will accept. I can probably write some simple Go scripts to do this, but it
  will probably be easier to do in a Python notebook. (Also more true to real-life examples.)
+ I'm going to go with a Python notebook since it allows me to use the Vertex SDK with all its
  extra features. I'll need to spin up a new notebook instance in my project.
+ 👎👎 I keep getting a `Missing permissions: storage.objects.get` error message when I attempt to
  tune a Gemini model from a dataset file on GCS. I don't know which SA I need to grant the
  role to. This is a hard blocker.
  - I filed a bug against the docs for this: b/376106542
+ 👍 Looks like I've successfully trained a Gemini model!
  - Note that I had to reject multiple entries in the Guanaco dataset since it didn't work for my
    simple transform.

Sources:

+ https://cloud.google.com/vertex-ai/generative-ai/docs/models/gemini-use-supervised-tuning#python
+ https://huggingface.co/datasets/timdettmers/openassistant-guanaco
+ https://python-jsonschema.readthedocs.io/en/latest/validate/
+ https://huggingface.co/docs/datasets/en/installation
+ https://json-schema.org/learn/miscellaneous-examples 

## Fine tune the Gemma model

+ For the Gemma model, I'll need to create a custom training application and run it as
  a job on Vertex AI's GAPIC layer (I think).
+ 👍👍 The Kaggle model card has an "Open on Vertex AI" button!! Gonna press it!
  - So this just opens the Gemma model card in the Vertex Model Card, which is still pretty nifty.
  - I see that there is a button to allow me to tune the Gemma model from a JSONL file stored
    on GCS.
  - Downloading the recommended file, I see that the dataset format is PEFT--yet another differet
    JSON shape than either Guanaco or what Gemini requires.
+ 👍 I decided to open the fine tuning notebook for Gemma in Colab, which was remarkably easy to
  do from the UI. Hopefully there is an option to look at the dataset or at least understand
  the shape of the dataset.
+ According to this Colab, the dataset format for tuning Gemma IS the Guanaco dataset! Hmm.
+ Going to attempt to tune a Gemma 2 model using the linked Colab.

Sources:

+ https://github.com/GoogleCloudPlatform/vertex-ai-samples/blob/main/notebooks/community/model_garden/model_garden_gemma2_finetuning_on_vertex.ipynb 
+ https://huggingface.co/google/gemma-2-27b-it-pytorch
+ https://www.kaggle.com/models/google/gemma-2
+ https://huggingface.co/blog/peft 

## Call the tuned model

+ 👎 I have tried a couple of different ways to call my tuned model: the raw predict API, using the GenerativeModel() method in Go
  Neither of them work!
  - I think that the [Godoc for the `genai` package](https://pkg.go.dev/github.com/google/generative-ai-go/genai#Client.GenerativeModel)
    is either wrong or misleading. The structure `tunedModels/NAME` where NAME is replaced by the tuned model's name, doesn't
    work.
  - The `NAME` of a tuned model is ambiguous. I'll try again, this time with the fully-qualified resource name.
  - Setting the `NAME` to the resource name of the model fails: `tunedModels/projects/1025771077852/locations/us-west1/models/8135194174937890816@1`
  - Setting the `NAME` to the resource name of the endpoint fails: `tunedModels/projects/1025771077852/locations/us-west1/endpoints/1926929312049528832`
  - Setting the `NAME` to just the endpoint name PANICs because no candidates are returned: `projects/1025771077852/locations/us-west1/endpoints/1926929312049528832`
+ 👎👎 The Go libraries leave a lot to be desired.
  - `ListModels()` method for v1 version of API
  - Clearer Godoc
  - Better error messages
+ 😬 The tuned model produces worse responses than the OOTB model:
  - The tuned model often returns an empty response as its first response
  - The tuned model ignores history context provided in prompt

Sources:

+ 👎👎 https://pkg.go.dev/cloud.google.com/go/vertexai/genai#Client.GenerativeModel 
+ https://cloud.google.com/vertex-ai/docs/generative-ai/start/quickstarts/quickstart-multimodal