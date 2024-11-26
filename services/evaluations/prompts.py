from dataclasses import dataclass

import logging
import pandas as pd

from google.cloud import firestore
from google.cloud import bigquery

COLLECTION_NAME = "PromptTemplate"
GOLDENS = "goldens20241104"
ADVERSARIALS = "adversarial20241117"


@dataclass
class Template:
    model: str
    prompt: str
    date: int


def get_templates(*args, project_id: str, database_name: str) -> list[Template]:

    templates = []
    logger = logging.getLogger(__name__)

    client = firestore.Client(project=project_id, database=database_name)
    collection = client.collection(COLLECTION_NAME)

    for a in args:
        modelTemplates = (
            collection.document(a)
            .collection("Templates")
            .order_by("Created", direction=firestore.Query.ASCENDING)
            .stream()
        )
        results = [r for r in modelTemplates]

        if len(results) == 0:
            logger.warn(f"No templates found for {a}")
            continue

        template = results[0].to_dict()
        templates.append(
            Template(model=a, prompt=template["Prompt"], date=template["Created"])
        )

    return templates


def get_goldens(project_id: str, dataset_name: str) -> pd.DataFrame:
    bq_client = bigquery.Client(project_id)
    goldens_table_name = f"{project_id}.{dataset_name}.{GOLDENS}"
    sql = f"""
    SELECT prompt, reference
    FROM {goldens_table_name}
    """

    golden_dataset = bq_client.query_and_wait(sql).to_dataframe()
    return golden_dataset


def get_adversarials(project_id: str, dataset_name: str) -> pd.DataFrame:
    bq_client = bigquery.Client(project_id)
    adversarial_table_name = f"{project_id}.{dataset_name}.{ADVERSARIALS}"
    sql = f"""
    SELECT prompt, reference
    FROM {adversarial_table_name}
    """

    adversarial_dataset = bq_client.query_and_wait(sql).to_dataframe()
    return adversarial_dataset
