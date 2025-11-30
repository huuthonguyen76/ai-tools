"""
Configuration settings for the AI Tools Streamlit application.
"""

import os

# Backend API Configuration
BACKEND_API_URL = os.getenv("BACKEND_API_URL", "http://localhost:8080")
API_TIMEOUT = int(os.getenv("API_TIMEOUT", "30"))  # seconds

# Layout Configuration
SIDEBAR_WIDTH_PERCENT = 30
MAIN_CONTENT_WIDTH_PERCENT = 70

# Available AI Tools
AVAILABLE_TOOLS = [
    "Link Contextualizer",
    # Add more tools here as they become available
]

# Default Tool
DEFAULT_TOOL = "Link Contextualizer"

# App Configuration
APP_TITLE = "AI Tools"
APP_ICON = "🤖"
PAGE_LAYOUT = "wide"

# API Endpoints
API_ENDPOINTS = {
    "contextualize_link": "/contextualize-link",
    "redirect": "/redirect/{client}/{link}",
}

# Session State Keys
SESSION_KEYS = {
    "selected_tool": "selected_tool",
    "link_data": "link_data",
    "error_message": "error_message",
    "loading": "loading",
}

