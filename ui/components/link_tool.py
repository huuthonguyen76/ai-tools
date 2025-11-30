"""
Link Contextualizer Tool Component.
"""

import streamlit as st
from typing import Optional, Dict
import config
from utils.helpers import (
    validate_url,
    initialize_session_state,
    display_error,
    display_success
)
from services.api_service import get_contextualized_link


def render_link_tool() -> Optional[Dict]:
    """
    Render the link contextualization tool interface.
    
    Returns:
        Optional[Dict]: API response data if successful, None otherwise
    """
    # Initialize session state
    initialize_session_state(config.SESSION_KEYS["link_data"], None)
    initialize_session_state(config.SESSION_KEYS["loading"], False)
    initialize_session_state(config.SESSION_KEYS["error_message"], None)
    
    st.header("🔗 Link Contextualizer")
    st.markdown("""
    Transform your URLs into contextually rich, AI-enhanced links that provide better 
    engagement and understanding.
    """)
    
    # Create the input form
    with st.form("link_form", clear_on_submit=False):
        st.subheader("Enter URL")
        
        link_input = st.text_input(
            label="Original Link",
            placeholder="https://example.com/article",
            help="Enter the URL you want to contextualize",
            key="link_input"
        )
        
        col1, col2, col3 = st.columns([1, 1, 2])
        with col1:
            submit_button = st.form_submit_button(
                "🚀 Contextualize Link",
                use_container_width=True,
                type="primary"
            )
        with col2:
            clear_button = st.form_submit_button(
                "🗑️ Clear",
                use_container_width=True
            )
    
    # Handle clear button
    if clear_button:
        st.session_state[config.SESSION_KEYS["link_data"]] = None
        st.session_state[config.SESSION_KEYS["error_message"]] = None
        st.rerun()

    # Handle form submission
    if submit_button:
        # Clear previous error
        st.session_state[config.SESSION_KEYS["error_message"]] = None
        
        # Validate input
        if not link_input:
            display_error("Please enter a URL")
            return None
        
        if not validate_url(link_input):
            display_error("Please enter a valid URL (must start with http:// or https://)")
            return None
        
        # Call API
        try:
            with st.spinner("🤖 AI is contextualizing your link..."):
                try:
                    d_response = get_contextualized_link("tho", link_input)

                    # Store in session state
                    st.session_state[config.SESSION_KEYS["link_data"]] = d_response                

                    display_success(
                        f"Link successfully contextualized! "
                        f"Check the results below."
                    )
                    return d_response
                except:
                    display_error("Unexpected response format from API")
                    return None

        except Exception as e:
            error_msg = f"Unexpected error: {str(e)}"
            st.session_state[config.SESSION_KEYS["error_message"]] = error_msg
            display_error(error_msg)
            return None
    
    # Return existing data from session state if available
    return st.session_state.get(config.SESSION_KEYS["link_data"])

