"""
Sidebar component for AI Tools navigation.
"""

import streamlit as st
import config
from utils.helpers import initialize_session_state


def render_sidebar() -> str:
    """
    Render the sidebar with AI tools navigation.
    
    Returns:
        str: Selected tool name
    """
    # Initialize session state for selected tool
    initialize_session_state(
        config.SESSION_KEYS["selected_tool"],
        config.DEFAULT_TOOL
    )
    
    with st.sidebar:
        # App branding
        st.title(f"{config.APP_ICON} {config.APP_TITLE}")
        st.markdown("---")
        
        # Tool selection
        st.subheader("🛠️ Available Tools")
        
        selected_tool = st.radio(
            label="Select a tool",
            options=config.AVAILABLE_TOOLS,
            index=config.AVAILABLE_TOOLS.index(
                st.session_state[config.SESSION_KEYS["selected_tool"]]
            ),
            label_visibility="collapsed"
        )
        
        # Update session state
        st.session_state[config.SESSION_KEYS["selected_tool"]] = selected_tool
        
        st.markdown("---")
        
        # Additional sidebar info
        with st.expander("ℹ️ About"):
            st.markdown("""
            **AI Tools Platform**
            
            A collection of AI-powered tools to help with various tasks:
            
            - **Link Contextualizer**: Transform URLs into contextually rich links
            - More tools coming soon!
            """)
        
        # Backend status indicator
        st.markdown("---")
        st.caption(f"Backend: `{config.BACKEND_API_URL}`")
    
    return selected_tool

