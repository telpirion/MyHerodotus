from typing import List

import functions_framework

import os
import logging
import traceback
from datetime import datetime

import pandas as pd
import pandas_gbq

import bigframes.pandas as bpd
from google.cloud import bigquery

import vertexai
from vertexai.generative_models import GenerativeModel
from vertexai.evaluation import EvalTask, Rouge, PointwiseMetric, PointwiseMetricPromptTemplate, MetricPromptTemplateExamples

@functions_framework.http
def entrypoint(request):
    """HTTP Cloud Function.
    Args:
        request (flask.Request): The request object.
        <https://flask.palletsprojects.com/en/1.1.x/api/#incoming-request-data>
    Returns:
        The response text, or any set of values that can be turned into a
        Response object using `make_response`
        <https://flask.palletsprojects.com/en/1.1.x/api/#flask.make_response>.
    """
    main()
    
def main():
    logger = logging.getLogger(__name__)
    project_id = os.getenv('PROJECT_ID')
    location = "us-west1"
    tb = "no error"

    try:
        golden_dataset = get_goldens()
        metrics = get_metrics()
        timestamp = datetime.utcnow()
        timestamp_str = timestamp.isoformat(timespec='hours')
        
        tuned_model_endpoint = "1926929312049528832"
        tuned_model_name = f"projects/{project_id}/locations/{location}/endpoints/{tuned_model_endpoint}"
        
        gemma_model_endpoint = "3122353538139684864"
        gemma_model_name = f"projects/{project_id}/locations/{location}/endpoints/{gemma_model_endpoint}"
        
        models = [
            ("gemini-1.5-flash-001", "gemini_1_5_flash_001"),
            (tuned_model_name, "tuned_gemini"),
            (gemma_model_name, "gemma"), # Raises "Template error: template not found"
        ]
        for m in models:
            model_id, model_name = m
            logger.info(f"{model_name} eval started")

            results_df = run_eval(model_id=model_id, eval_dataset=golden_dataset, metrics=metrics)
            table_name = f"{project_id}.myherodotus.{model_name}dd{timestamp_str}"

            store_results(results_df, table_name, project_id)
            logger.info(f"{model_name} results written to log")

    except Exception as e:
        logger.error(e)
        tb = traceback.format_exc()
    finally:
        logger.error(tb)
    
    
def get_goldens() -> pd.DataFrame:
    project_id = os.getenv('PROJECT_ID')
    bq_client = bigquery.Client(project_id)
    goldens_table_name = f"{project_id}.myherodotus.goldens20241104"
    sql = f"""
    SELECT prompt, reference
    FROM {goldens_table_name}
    """

    golden_dataset = bq_client.query_and_wait(sql).to_dataframe()
    return golden_dataset

def get_metrics() -> List[any]:
    # My set of metrics
    open_domain = '''
    In this conversation between a human and the AI, the AI is helpful and friendly, 
    and when it does not know the answer it says \"I don’t know\".\n
    '''

    closed_domain = '''
    The user wants to travel to a country to see historical landmarks and archaeological sites.
    The AI is a helpful travel guide. Please provide 3 to 5 destination suggestions.
    '''

    prompteng_metrics = PointwiseMetric(
        metric="prompteng_metrics",
        metric_prompt_template=PointwiseMetricPromptTemplate(
            criteria={
                "open_domain": open_domain,
                "closed_domain": closed_domain,
            },
            rating_rubric={
                "1": "The response performs well on both criteria.",
                "0.5": "The response performs well on one but not the other criteria.",
                "0": "The response falls short on both criteria",
            },
        ),
    )

    rouge = Rouge(rouge_type="rouge1")
    metrics = [
        prompteng_metrics,
        rouge,
        MetricPromptTemplateExamples.Pointwise.GROUNDEDNESS,
        MetricPromptTemplateExamples.Pointwise.COHERENCE,
    ]
    return metrics
    
def run_eval(model_id: str, eval_dataset: pd.DataFrame, metrics: List[any]) -> pd.DataFrame:
    candidate_model = GenerativeModel(model_id)
    pointwise_eval_task = EvalTask(
        dataset=eval_dataset,
        metrics=metrics,
    )
    pointwise_result = pointwise_eval_task.evaluate(
        model=candidate_model,
    )
    results = pointwise_result.metrics_table
    return results
        
def store_results(results_df: pd.DataFrame, table_name: str, project_id: str) -> bool:
    clean_results = cleanup_column_names(results_df)
    pandas_gbq.to_gbq(clean_results, table_name, project_id=project_id)

def cleanup_column_names(df: pd.DataFrame) -> pd.DataFrame:
    new_names = {}
    for series_name, _ in df.items():
        new_names[series_name] = series_name.replace("/", "_")

    return df.rename(columns=new_names)

if __name__ == "__main__":
    main()
