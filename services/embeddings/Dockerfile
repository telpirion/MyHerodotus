FROM us-docker.pkg.dev/deeplearning-platform-release/gcr.io/pytorch-2-0-gpu-v20240922-debian-11

ARG OUTPUT_PATH=/artifacts
ENV OUTPUT_PATH=$OUTPUT_PATH

WORKDIR /

COPY embeddings.py ./embeddings.py
COPY requirements.txt ./requirements.txt
COPY history.mb.txt ./history.mb.txt

RUN mkdir $OUTPUT_PATH

RUN pip install -r requirements.txt

CMD [ "python", "./embeddings.py"]
