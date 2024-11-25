from typing import List
from dataclasses import dataclass, field
import requests

from vertexai.generative_models import GenerativeModel

class Response:
    candidates: [list]

class HerodotusModel(GenerativeModel):
    base_url = "http://localhost:8080/predict"
    def generate_content(self, prompt: str):
        payload = {
            "message": prompt,
            "model": "gemini"
        }
        resp = requests.post(self.base_url, json=payload, verify=False)
        resp_json = resp.json()
        return resp_json["Message"]["Message"]
