import pytest

from unittest.mock import Mock, patch, MagicMock
from typing import List

from src.services.llm import OpenAIClient
from src.models.llm_model import LLMMessage


class TestOpenAIClient:
    """Test suite for OpenAIClient class."""

    @patch('src.services.llm.OpenAI')
    def test_with_valid_messages(self, mock_openai):
        """Test call_openai with valid messages."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance

        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = "Valid messages response"
        mock_client_instance.chat.completions.create.return_value = mock_response
        
        client = OpenAIClient()
        messages = [LLMMessage(role="user", content="Test message")]
        
        # Act
        result = client.call_openai(messages)
        
        # Assert
        assert result == "Valid messages response"
        mock_client_instance.chat.completions.create.assert_called_once_with(
            model="gpt-4.1-mini",
            messages=[{"role": "user", "content": "Test message"}]
        )

    @patch('src.services.llm.OpenAI')
    def test_call_openai_with_custom_model(self, mock_openai):
        """Test call_openai with a custom model parameter."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = "Custom model response"
        mock_client_instance.chat.completions.create.return_value = mock_response
        
        client = OpenAIClient()
        messages = [LLMMessage(role="user", content="Test message")]
        custom_model = "gpt-4"
        
        # Act
        result = client.call_openai(messages, model=custom_model)
        
        # Assert
        assert result == "Custom model response"
        mock_client_instance.chat.completions.create.assert_called_once_with(
            model=custom_model,
            messages=[{"role": "user", "content": "Test message"}]
        )

    @patch('src.services.llm.OpenAI')
    def test_call_openai_with_multiple_messages(self, mock_openai):
        """Test call_openai correctly handles multiple messages."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = "Response to conversation"
        mock_client_instance.chat.completions.create.return_value = mock_response
        
        client = OpenAIClient()
        messages = [
            LLMMessage(role="system", content="You are a helpful assistant."),
            LLMMessage(role="user", content="What is 2+2?"),
            LLMMessage(role="assistant", content="2+2 equals 4."),
            LLMMessage(role="user", content="What about 3+3?")
        ]
        
        # Act
        result = client.call_openai(messages)
        
        # Assert
        assert result == "Response to conversation"
        call_args = mock_client_instance.chat.completions.create.call_args
        assert len(call_args.kwargs['messages']) == 4

    @patch('src.services.llm.OpenAI')
    def test_call_openai_transforms_messages_correctly(self, mock_openai):
        """Test that LLMMessage objects are properly transformed to dictionaries."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = "Test"
        mock_client_instance.chat.completions.create.return_value = mock_response
        
        client = OpenAIClient()
        messages = [
            LLMMessage(role="user", content="Hello")
        ]
        
        # Act
        client.call_openai(messages)
        
        # Assert
        call_args = mock_client_instance.chat.completions.create.call_args
        messages_arg = call_args.kwargs['messages']
        assert isinstance(messages_arg, list)
        assert isinstance(messages_arg[0], dict)
        assert messages_arg[0] == {"role": "user", "content": "Hello"}

    @patch('src.services.llm.OpenAI')
    def test_call_openai_with_empty_response_content(self, mock_openai):
        """Test call_openai handles empty response content."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = ""
        mock_client_instance.chat.completions.create.return_value = mock_response
        
        client = OpenAIClient()
        messages = [LLMMessage(role="user", content="Test")]
        
        # Act
        result = client.call_openai(messages)
        
        # Assert
        assert result == ""

    @patch('src.services.llm.OpenAI')
    def test_call_openai_raises_exception_on_api_error(self, mock_openai):
        """Test that API errors are propagated correctly."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        mock_client_instance.chat.completions.create.side_effect = Exception("API Error")
        
        client = OpenAIClient()
        messages = [LLMMessage(role="user", content="Test")]
        
        # Act & Assert
        with pytest.raises(Exception) as exc_info:
            client.call_openai(messages)
        assert str(exc_info.value) == "API Error"

    @patch('src.services.llm.OpenAI')
    def test_call_openai_with_single_message(self, mock_openai):
        """Test call_openai with a single message."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_response = Mock()
        mock_response.choices = [Mock()]
        mock_response.choices[0].message.content = "Single message response"
        mock_client_instance.chat.completions.create.return_value = mock_response
        
        client = OpenAIClient()
        messages = [LLMMessage(role="user", content="Single message")]
        
        # Act
        result = client.call_openai(messages)
        
        # Assert
        assert result == "Single message response"
        call_args = mock_client_instance.chat.completions.create.call_args
        assert len(call_args.kwargs['messages']) == 1

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_with_valid_texts(self, mock_openai):
        """Test get_list_embedding with valid text list returns expected embeddings."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_embedding_1 = Mock()
        mock_embedding_1.embedding = [0.1, 0.2, 0.3]
        mock_embedding_2 = Mock()
        mock_embedding_2.embedding = [0.4, 0.5, 0.6]
        
        mock_response = Mock()
        mock_response.data = [mock_embedding_1, mock_embedding_2]
        mock_client_instance.embeddings.create.return_value = mock_response
        
        client = OpenAIClient()
        l_text = ["Hello world", "Test text"]
        
        # Act
        result = client.get_list_embedding(l_text)
        
        # Assert
        assert result == [[0.1, 0.2, 0.3], [0.4, 0.5, 0.6]]
        mock_client_instance.embeddings.create.assert_called_once_with(
            input=l_text,
            model="text-embedding-3-small"
        )

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_with_custom_model(self, mock_openai):
        """Test get_list_embedding with a custom embedding model."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_embedding = Mock()
        mock_embedding.embedding = [0.1, 0.2, 0.3, 0.4]
        
        mock_response = Mock()
        mock_response.data = [mock_embedding]
        mock_client_instance.embeddings.create.return_value = mock_response
        
        client = OpenAIClient()
        l_text = ["Custom model test"]
        custom_model = "text-embedding-3-large"
        
        # Act
        result = client.get_list_embedding(l_text, model=custom_model)
        
        # Assert
        assert result == [[0.1, 0.2, 0.3, 0.4]]
        mock_client_instance.embeddings.create.assert_called_once_with(
            input=l_text,
            model=custom_model
        )

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_with_single_text(self, mock_openai):
        """Test get_list_embedding with a single text."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_embedding = Mock()
        mock_embedding.embedding = [0.7, 0.8, 0.9]
        
        mock_response = Mock()
        mock_response.data = [mock_embedding]
        mock_client_instance.embeddings.create.return_value = mock_response
        
        client = OpenAIClient()
        l_text = ["Single text"]
        
        # Act
        result = client.get_list_embedding(l_text)
        
        # Assert
        assert result == [[0.7, 0.8, 0.9]]
        assert len(result) == 1

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_with_multiple_texts(self, mock_openai):
        """Test get_list_embedding correctly handles multiple texts."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance

        mock_embeddings = [Mock() for _ in range(5)]
        for i, mock_emb in enumerate(mock_embeddings):
            mock_emb.embedding = [float(i), float(i+1), float(i+2)]
        
        mock_response = Mock()
        mock_response.data = mock_embeddings
        mock_client_instance.embeddings.create.return_value = mock_response
        
        client = OpenAIClient()
        l_text = ["Text 1", "Text 2", "Text 3", "Text 4", "Text 5"]
        
        # Act
        result = client.get_list_embedding(l_text)
        
        # Assert
        assert len(result) == 5
        assert result[0] == [0.0, 1.0, 2.0]
        assert result[4] == [4.0, 5.0, 6.0]

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_with_empty_list(self, mock_openai):
        """Test get_list_embedding with an empty list."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_response = Mock()
        mock_response.data = []
        mock_client_instance.embeddings.create.return_value = mock_response
        
        client = OpenAIClient()
        l_text = []
        
        # Act
        result = client.get_list_embedding(l_text)
        
        # Assert
        assert result == []
        mock_client_instance.embeddings.create.assert_called_once_with(
            input=[],
            model="text-embedding-3-small"
        )

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_raises_exception_on_api_error(self, mock_openai):
        """Test that embedding API errors are propagated correctly."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        mock_client_instance.embeddings.create.side_effect = Exception("Embedding API Error")
        
        client = OpenAIClient()
        l_text = ["Test text"]
        
        # Act & Assert
        with pytest.raises(Exception) as exc_info:
            client.get_list_embedding(l_text)
        assert str(exc_info.value) == "Embedding API Error"

    @patch('src.services.llm.OpenAI')
    def test_get_list_embedding_extracts_embeddings_correctly(self, mock_openai):
        """Test that embeddings are correctly extracted from response data."""
        # Arrange
        mock_client_instance = MagicMock()
        mock_openai.return_value = mock_client_instance
        
        mock_data_1 = Mock()
        mock_data_1.embedding = [0.111, 0.222]
        mock_data_2 = Mock()
        mock_data_2.embedding = [0.333, 0.444]
        
        mock_response = Mock()
        mock_response.data = [mock_data_1, mock_data_2]
        mock_client_instance.embeddings.create.return_value = mock_response
        
        client = OpenAIClient()
        l_text = ["First", "Second"]
        
        # Act
        result = client.get_list_embedding(l_text)
        
        # Assert
        assert isinstance(result, list)
        assert len(result) == 2
        assert isinstance(result[0], list)
        assert isinstance(result[1], list)
        assert all(isinstance(val, float) for val in result[0])
        assert result == [[0.111, 0.222], [0.333, 0.444]]

