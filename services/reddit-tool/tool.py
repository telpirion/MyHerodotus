from typing import List, Mapping

import json
import os
import praw

import vertexai
from vertexai.preview import reasoning_engines
from google.cloud import secretmanager

from secret import PROJECT_ID

project_id = staging_bucket = ""

SUBREDDIT_NAME = "travel"
LIMIT = 10
LOCATION = "us-central1"
MODEL = "gemini-1.5-pro"

def get_secrets() -> Mapping[str, str]:
    secret_name = f"projects/{PROJECT_ID}/secrets/reddit-api-key/versions/1"
    secret_client = secretmanager.SecretManagerServiceClient()
    secret = secret_client.access_secret_version(name=secret_name)
    reddit_key_json = json.loads(secret.payload.data)
    return reddit_key_json

def get_posts(query: str, credentials: Mapping[str, str]) -> List[Mapping[str, str]]:
    reddit = praw.Reddit(
        client_id=credentials["client_id"],
        client_secret=credentials["secret"],
        user_agent=credentials["user_agent"],
    )
    posts = reddit.subreddit(SUBREDDIT_NAME).search(query, time_filter="year")
    reddit_messages = []
    for p in posts:
        comments_list = p.comments.list()
        if len(comments_list) == 0:
            continue
        reddit_messages.append({
            "title": p.title,
            "message": p.selftext,
            "top_comment": comments_list[0].body
        })
    return reddit_messages

def get_reddit_reviews(query: str) -> List[Mapping[str, str]]:
    """Gets a list of place reviews from Reddit.

    Arguments:
        query: the user query, specifically a request for travel information about a destination

    Returns:
        A list of dictionaries, where dictionaries have a 'message' and a 'comment'
    """    
    reddit_key_json = get_secrets()
    messages = get_posts(query, credentials=reddit_key_json)
    return messages

def deploy():
    project_id = os.environ["PROJECT_ID"]
    staging_bucket = os.environ["BUCKET"]
    vertexai.init(project=project_id, location=LOCATION, staging_bucket=staging_bucket)

    system_instruction = """
You are a helpful AI travel assistant. The user wants to hear Reddit reviews about a
specific location. You are going to use the get_reddit_reviews tool to get Reddit posts
about the specific location that the user wants to know about.
"""

    agent = reasoning_engines.LangchainAgent(
        system_instruction=system_instruction,
        model=MODEL,
        model_kwargs={"temperature": 0.6}, # Try to avoid "I can't help you" answers
        tools=[
            get_reddit_reviews,
        ],
    )

    remote_agent = reasoning_engines.ReasoningEngine.create(
        agent,
        requirements=[
            "google-cloud-aiplatform[langchain,reasoningengine]",
            "cloudpickle==3.0.0",
            "pydantic==2.7.4",
            "google-cloud-secret-manager",
            "praw",
        ],
        extra_packages=[
            "./secret.py"
        ]
    )

    # Test remote
    response = remote_agent.query(
        input="""I want to take a trip to Crete. Where should I stay? What sites should I go see?"""
    )
    output = response["output"]
    print(output)


if __name__ == "__main__":
    deploy()