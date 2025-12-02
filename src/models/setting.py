from pydantic_settings import BaseSettings


class SettingModel(BaseSettings):
    OPEN_AI_TOKEN: str = ''

    API_ENV: str = ''
    ELASTICSEARCH_URL: str = ''
    ELASTICSEARCH_USERNAME: str = ''
    ELASTICSEARCH_PASSWORD: str = ''
    ELASTICSEARCH_INDEX: str = 'k_bot-kb-index'
    
    APIFY_TOKEN: str = ''
    APIFY_ACTOR_URL: str = ''
    
    DEFAULT_PROMPT: str = ''
    GOOGLE_APPLICATION_CREDENTIALS: str = ''

    SLACK_HEALTH_CHECK_CHANNEL: str = ''
    
    AZURE_GPT4_ENDPOINT: str = ''
    AZURE_GPT4_TOKEN: str = ''
    
    AZURE_GPT35_ENDPOINT: str = ''
    AZURE_GPT35_TOKEN: str = ''
    
    AZURE_EMBEDDING_ENDPOINT: str = ''
    AZURE_EMBEDDING_TOKEN: str = ''
    
    IS_AZURE: int = 0

    PROMETHEUS_MULTIPROC_DIR: str = './temp'

    DIFY_AGENT_TOKEN: str = ''
    DIFY_DATASET_TOKEN: str = ''
    DIFY_DATASET_ID: str = ''

    IMAGE_STORAGE_BUCKET_NAME: str = "multi-modal-image"

    JINA_READER_TOKEN: str = ''
