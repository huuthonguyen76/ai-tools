import os

from dotenv import load_dotenv

load_dotenv()

# Backend API Configuration
BACKEND_API_URL = os.getenv("DOMAIN_API")
API_TIMEOUT = 30

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

# Session State Keys
SESSION_KEYS = {
    "selected_tool": "selected_tool",
    "link_data": "link_data",
    "error_message": "error_message",
    "loading": "loading",
}

