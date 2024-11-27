import os
import functions_framework

import vertexai
from vertexai.preview import reasoning_engines

from tool import get_reddit_reviews

LOCATION = "us-west1"
MODEL = "gemini-1.5-pro"


@functions_framework.http
def get_agent_request(request):
    """HTTP Cloud Function.
    Args:
        request (flask.Request): The request object.
        <https://flask.palletsprojects.com/en/1.1.x/api/#incoming-request-data>
    Returns:
        The response text, or any set of values that can be turned into a
        Response object using `make_response`
        <https://flask.palletsprojects.com/en/1.1.x/api/#flask.make_response>.
    """
    query = ""
    request_json = request.get_json(silent=True)

    if request_json and "query" in request_json:
        query = request_json["query"]

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
        system_instruction=system_instruction,
        model=MODEL,
        # Try to avoid "I can't help you" answers
        model_kwargs={"temperature": 0.6},
        tools=[
            get_reddit_reviews,
        ],
    )

    response = agent.query(
        input=query
    )
    output = response["output"]

    return {
        "response": output
    }
