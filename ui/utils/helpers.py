"""
Helper utility functions for the AI Tools application.
"""

import re
from datetime import datetime
from typing import Optional
import streamlit as st


def validate_url(url: str) -> bool:
    """
    Validate if a string is a valid URL.
    
    Args:
        url: String to validate as URL
        
    Returns:
        bool: True if valid URL, False otherwise
    """
    if not url or not isinstance(url, str):
        return False
    
    # Basic URL regex pattern
    url_pattern = re.compile(
        r'^https?://'  # http:// or https://
        r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+[A-Z]{2,6}\.?|'  # domain...
        r'localhost|'  # localhost...
        r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})'  # ...or ip
        r'(?::\d+)?'  # optional port
        r'(?:/?|[/?]\S+)$', re.IGNORECASE)
    
    return bool(url_pattern.match(url))


def format_timestamp(timestamp: Optional[str]) -> str:
    """
    Format timestamp string to human-readable format.
    
    Args:
        timestamp: ISO format timestamp string
        
    Returns:
        str: Formatted timestamp or 'N/A' if invalid
    """
    if not timestamp:
        return "N/A"
    
    try:
        dt = datetime.fromisoformat(timestamp.replace('Z', '+00:00'))
        return dt.strftime("%Y-%m-%d %H:%M:%S")
    except (ValueError, AttributeError):
        return "N/A"


def copy_to_clipboard_button(text: str, button_label: str = "📋 Copy") -> None:
    """
    Create a button that copies text to clipboard.
    
    Args:
        text: Text to copy to clipboard
        button_label: Label for the button
    """
    if st.button(button_label, key=f"copy_{hash(text)}"):
        st.code(text, language=None)
        st.success("✅ Copied to clipboard! (Use Ctrl+C / Cmd+C to copy from the box above)")


def display_error(message: str) -> None:
    """
    Display error message in a styled container.
    
    Args:
        message: Error message to display
    """
    st.error(f"❌ {message}")


def display_success(message: str) -> None:
    """
    Display success message in a styled container.
    
    Args:
        message: Success message to display
    """
    st.success(f"✅ {message}")


def display_info(message: str) -> None:
    """
    Display info message in a styled container.
    
    Args:
        message: Info message to display
    """
    st.info(f"ℹ️ {message}")


def create_card_style() -> str:
    """
    Return CSS for card-style containers.
    
    Returns:
        str: CSS styling for cards
    """
    return """
    <style>
    .card {
        padding: 1.5rem;
        border-radius: 0.5rem;
        background-color: #f0f2f6;
        margin: 1rem 0;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .card-title {
        font-size: 1.2rem;
        font-weight: bold;
        margin-bottom: 0.5rem;
        color: #262730;
    }
    .card-content {
        font-size: 1rem;
        color: #31333F;
    }
    .link-display {
        background-color: white;
        padding: 0.75rem;
        border-radius: 0.25rem;
        border-left: 4px solid #ff4b4b;
        margin: 0.5rem 0;
        word-break: break-all;
    }
    </style>
    """


def initialize_session_state(key: str, default_value) -> None:
    """
    Initialize session state key if it doesn't exist.
    
    Args:
        key: Session state key
        default_value: Default value to set
    """
    if key not in st.session_state:
        st.session_state[key] = default_value

