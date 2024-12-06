from typing import List
from dataclasses import dataclass, field
import requests

from google.protobuf.json_format import ParseDict

from vertexai.generative_models import GenerativeModel
from google.cloud.aiplatform_v1.types.prediction_service import GenerateContentResponse


class HerodotusModel(GenerativeModel):
    base_url = "https://myherodotus-1025771077852.us-west1.run.app/predict"

    def __init__(self, modality):
        self.modality = modality

    @property
    def _model_name(self) -> str:
        return self.modality

    def generate_content(self, prompt: str):
        payload = {"message": prompt, "model": self.modality}
        resp = requests.post(self.base_url, json=payload, verify=False)
        resp_json = resp.json()
        response_payload = {
            "candidates": [
                {
                    "finish_reason": 1,
                    "content": {"parts": [{"text": resp_json["Message"]["Message"]}]},
                },
            ],
        }
        proto_ver = ParseDict(response_payload, GenerateContentResponse()._pb)
        return proto_ver
