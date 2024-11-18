from typing import List, Mapping, Set

import os

import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim

from google.cloud import storage


CONTEXT_SIZE = 2
EMBEDDING_DIM = 10


class NGramLanguageModeler(nn.Module):
    def __init__(self, vocab_size, embedding_dim, context_size):
        super(NGramLanguageModeler, self).__init__()
        self.embeddings = nn.Embedding(vocab_size, embedding_dim)
        self.linear1 = nn.Linear(context_size * embedding_dim, 128)
        self.linear2 = nn.Linear(128, vocab_size)

    def forward(self, inputs):
        embeds = self.embeddings(inputs).view((1, -1))
        out = F.relu(self.linear1(embeds))
        out = self.linear2(out)
        log_probs = F.log_softmax(out, dim=1)
        return log_probs


def create_embeddings(ngrams: List[any],
                      word_to_ix: Mapping[str, any],
                      vocab: Set[any],
                      epochs: int = 2) -> NGramLanguageModeler:
    losses = []
    loss_function = nn.NLLLoss()
    model = NGramLanguageModeler(len(vocab), EMBEDDING_DIM, CONTEXT_SIZE)
    optimizer = optim.SGD(model.parameters(), lr=0.001)

    for epoch in range(epochs):
        total_loss = 0
        for context, target in ngrams:

            # Step 1. Prepare the inputs to be passed to the model (i.e, turn the words
            # into integer indices and wrap them in tensors)
            context_idxs = torch.tensor([word_to_ix[w] for w in context], dtype=torch.long)

            # Step 2. Recall that torch *accumulates* gradients. Before passing in a
            # new instance, you need to zero out the gradients from the old
            # instance
            model.zero_grad()

            # Step 3. Run the forward pass, getting log probabilities over next
            # words
            log_probs = model(context_idxs)

            # Step 4. Compute your loss function. (Again, Torch wants the target
            # word wrapped in a tensor)
            loss = loss_function(log_probs, torch.tensor([word_to_ix[target]], dtype=torch.long))

            # Step 5. Do the backward pass and update the gradient
            loss.backward()
            optimizer.step()

            # Get the Python number from a 1-element Tensor by calling tensor.item()
            total_loss += loss.item()
        losses.append(total_loss)

        print(f"Epoch {epoch} completed")
    return model 

def create_ngrams(herodotus_arr: str):
    ngrams = [
        (
            [herodotus_arr[i - j - 1] for j in range(CONTEXT_SIZE)],
            herodotus_arr[i]
        )
        for i in range(CONTEXT_SIZE, len(herodotus_arr))
    ]
    return ngrams

def get_herodotus(char_count: int = 0, front_matter: int = 285) -> str:
    herodotus_text = ""
    with open("history.mb.txt", "r") as f:
        herodotus_text = f.read()
    
    if char_count:
        return herodotus_text[front_matter:char_count]

    return herodotus_text[front_matter]


def save_model_to_storage(project_id: str, bucket_name: str, model_path: str):
    client = storage.Client(project=project_id)
    bucket = client.bucket(bucket_name)
    blob = bucket.blob(model_path)
    blob.upload_from_filename(model_path)


def main():
    project_id = os.getenv("PROJECT_ID")
    output_path = os.getenv("OUTPUT_PATH")
    bucket = os.getenv("BUCKET_NAME")
    torch.manual_seed(1)

    herodotus_text = get_herodotus(char_count=1985)
    herodotus_arr = herodotus_text.split()
    ngrams = create_ngrams(herodotus_arr)

    vocab = set(herodotus_arr)
    word_to_ix = {word: i for i, word in enumerate(vocab)}
    model = create_embeddings(ngrams, word_to_ix, vocab)
    torch.save(model.state_dict(), output_path)
    save_model_to_storage(
        project_id=project_id,
        bucket_name=bucket,
        model_path=output_path)



if __name__ == "__main__":
    print("start embeddings creation")
    main()
    print("embeddings created")