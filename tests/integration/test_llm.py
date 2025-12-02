from src.services.llm import OpenAIClient
from src.models.llm_model import LLMMessage


class TestOpenAIClient:
    def test_simple_response(self) -> None:
        openai_client = OpenAIClient()
        response = openai_client.call_openai([
            LLMMessage(role="user", content="Hello, how are you?")
        ])

        assert response is not None

    def test_simple_embedding(self) -> None:
        openai_client = OpenAIClient()
        response = openai_client.get_list_embedding([
            "Hello, how are you?",
            "I am fine, thank you!"
        ])

        assert len(response) == 2
        assert len(response[0]) == 1536
        assert len(response[1]) == 1536
