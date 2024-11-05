# Week 4 Development Log

+ Unhappy: 👎
+ Anxious: 😬
+ Curious: 🤔
+ Happy: 👍

## Objectives

For this week's activities, we must do the following:

- [ ] Set up an evaluation pipeline to compare Gemini, Gemma, and/or tuned model.
- [ ] Export evaluation to a "table" (BQ?).
- [ ] Set up a rapid evaluation pipeline to see the specific performance of a model.

Nice to haves:

- [ ] Limit context passed to Gemma model based upon token count
- [ ] Train Gemma model on Guanaco dataset
- [ ] Upgrade ALL the things to Genkit

## Set up an evaluation pipeline

+ 😬🤔 The evaluation overview in the docs reads more like marketing copy than a 
  technical overview.
+ In the quickstart, are "Fluency" and "Entertaining" both metrics that are pre-defined?
  Or can these metrics be any arbitrary measurement provided that they have a explanation?
+ 👎👎 Quickstart doesn't give an example of a correct experiment name!!
  [here](https://cloud.google.com/vertex-ai/generative-ai/docs/models/evaluation-quickstart#import_libraries).
  - Filed a bug: b/376756582
+ It would be nice to have some examples of commonly-used metrics for Gen AI models -- plus how to 
  potentially define them when setting up the result.
+ Additional notes collected in notebook.

Sources:

+ https://cloud.google.com/vertex-ai/generative-ai/docs/models/evaluation-overview