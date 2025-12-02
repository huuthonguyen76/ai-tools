from openai import OpenAI
from typing import List

from src.models.llm_model import LLMMessage


class OpenAIClient:
    def __init__(self):
        self.client = OpenAI()

    def call_openai(self, l_message: List[LLMMessage], model="gpt-4.1-mini") -> str:
        l_message = [message.model_dump() for message in l_message]
        response = self.client.chat.completions.create(
            model=model,
            messages=l_message
        )
        return response.choices[0].message.content

    def get_list_embedding(self, l_text: List[str], model="text-embedding-3-small") -> List[List[float]]:
        response = self.client.embeddings.create(
            input=l_text,
            model=model
        )
        return [data.embedding for data in response.data]
