from pydantic import BaseModel, field_validator

class LLMMessage(BaseModel):
    role: str
    content: str

    @field_validator('role')
    def validate_role(cls, v):
        if v not in ['system', 'user', 'assistant']:
            raise ValueError('Invalid role')
        return v
