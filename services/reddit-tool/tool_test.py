import logging
import os

import vertexai
from vertexai.preview import reasoning_engines

from tool import get_reddit_reviews

LOCATION = "us-west1"
MODEL = "gemini-1.5-pro"
REASONING_ENGINE_ID = "1823623597550206976"
LOGGER = logging.getLogger()
INPUT = """I want to take a trip to Crete. Where should I stay? I want to see
ancient ruins. What are the best archaeological sites to see?"""


def test_create_agent_local():
    project_id = os.environ["PROJECT_ID"]
    staging_bucket = os.environ["BUCKET"]
    vertexai.init(project=project_id, location=LOCATION,
                  staging_bucket=staging_bucket)
    system_instruction = """
You are a helpful AI travel assistant. The user wants to hear Reddit reviews
about a specific location. You are going to use the get_reddit_reviews tool to
get Reddit posts about the specific location that the user wants to know about.
"""
    agent = reasoning_engines.LangchainAgent(
        model=MODEL,
        model_kwargs={"temperature": 0.5},
        tools=[
            get_reddit_reviews,
        ],
        system_instruction=system_instruction,
    )

    response = agent.query(
        input=INPUT
    )
    output = response["output"]
    LOGGER.info(output)
    assert output != ""


def test_query_agent_remote():
    project_number = os.environ["PROJECT_NUMBER"]
    agent_name = f'projects/{project_number}/locations/us-central1/reasoningEngines/{REASONING_ENGINE_ID}'
    reasoning_engine = vertexai.preview.reasoning_engines.ReasoningEngine(
        agent_name)
    response = reasoning_engine.query(
        input="""I want to take a trip to Crete. Where should I stay? I want
        to see ancient ruins. What are the best archaeological sites to see?"""
    )
    output = response['output']

    LOGGER.info(output)

    assert output != ""
