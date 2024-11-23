from typing import List

import os
import logging
import traceback
from datetime import datetime

import pandas as pd
import pandas_gbq

import bigframes.pandas as bpd
from google.cloud import bigquery

import vertexai
from vertexai import generative_models
from vertexai.generative_models import GenerativeModel
from vertexai.evaluation import EvalTask, Rouge, PointwiseMetric, PointwiseMetricPromptTemplate, MetricPromptTemplateExamples

from prompts import get_templates, get_goldens, get_adversarials
from metrics import get_metrics


def main():
    logger = logging.getLogger(__name__)
    project_id = os.getenv('PROJECT_ID')
    dataset_name = os.getenv('DATASET_NAME')

    if not project_id:
        logger.error('No project ID')
        return
    elif not dataset_name:
        logger.error('No dataset name')
        return

    logger.info(f"Project ID: {project_id}")
    logger.info(f"Dataset Name: {dataset_name}")
    location = "us-west1"
    tb = "exit status 0"

    vertexai.init(project=project_id, location=location)
    try:
        templates = get_templates("Gemini", "Gemma", project_id=project_id, database_name="l200")
        golden_dataset = get_goldens(project_id="erschmid-test-291318", dataset_name=dataset_name)
        adversarial_dataset = get_adversarials(project_id="erschmid-test-291318", dataset_name=dataset_name)

        metrics = get_metrics()
        timestamp = datetime.utcnow()
        timestamp_str = timestamp.strftime("%Y_%m_%d_%H_%M")
        
        tuned_model_endpoint = "1926929312049528832"
        tuned_model_name = f"projects/{project_id}/locations/{location}/endpoints/{tuned_model_endpoint}"
        
        gemma_model_endpoint = "3122353538139684864"
        gemma_model_name = f"projects/{project_id}/locations/{location}/endpoints/{gemma_model_endpoint}"
        
        models = [
            ("gemini-1.5-flash-001", "gemini_1_5_flash_001"),
            (tuned_model_name, "tuned_gemini"),
            #(gemma_model_name, "gemma"), # Raises "Template error: template not found"
        ]
        for m in models:
            model_id, model_name = m

            logger.info(f"{model_name} goldens eval started")
            results_df = run_eval(model_id=model_id, eval_dataset=golden_dataset, metrics=metrics)
            table_name = f"{project_id}.{dataset_name}.{model_name}_goldens_{timestamp_str}"
            store_results(results_df, table_name, project_id)
            logger.info(f"{model_name} goldens results written to log")

            logger.info(f"{model_name} adversarials eval started")

            # Relax safety settings
            safety_settings = [
                generative_models.SafetySetting(
                    category=generative_models.HarmCategory.HARM_CATEGORY_SEXUALLY_EXPLICIT,
                    threshold=generative_models.HarmBlockThreshold.BLOCK_ONLY_HIGH,
                ),
                generative_models.SafetySetting(
                    category=generative_models.HarmCategory.HARM_CATEGORY_DANGEROUS_CONTENT,
                    threshold=generative_models.HarmBlockThreshold.BLOCK_ONLY_HIGH,
                ),
            ]

            adversarials_df = run_eval(model_id=model_id, eval_dataset=adversarial_dataset, metrics=metrics, safety_settings=safety_settings)
            table_name = f"{project_id}.{dataset_name}.{model_name}_adversarials_{timestamp_str}"
            store_results(results_df, table_name, project_id)
            logger.info(f"{model_name} adversarials results written to log")

    except Exception as e:
        logger.error(e)
        tb = traceback.format_exc()
    finally:
        logger.error(tb)


def run_eval(model_id: str, eval_dataset: pd.DataFrame, metrics: List[any], safety_settings: List[any] = None) -> pd.DataFrame:
    candidate_model = GenerativeModel(model_id, safety_settings=safety_settings)
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
