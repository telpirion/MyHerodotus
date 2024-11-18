FROM python:3.12-slim-bookworm AS build

WORKDIR /

COPY evaluations.py ./evaluations.py
COPY metrics.py ./metrics.py
COPY prompts.py ./prompts.py
COPY requirements.txt ./requirements.txt

RUN pip install -r requirements.txt

CMD [ "python", "./evaluations.py"]
