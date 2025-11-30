"""
AI Tools - Main Streamlit Application

A collection of AI-powered tools for various tasks.
"""

import streamlit as st
import config
from components.sidebar import render_sidebar
from components.link_tool import render_link_tool
from components.content_blocks import render_content_blocks, render_empty_state


def configure_page():
    """Configure Streamlit page settings."""
    st.set_page_config(
        page_title=config.APP_TITLE,
        page_icon=config.APP_ICON,
        layout=config.PAGE_LAYOUT,
        initial_sidebar_state="expanded"
    )


def main():
    """Main application entry point."""
    # Configure page
    configure_page()
    
    # Render sidebar and get selected tool
    selected_tool = render_sidebar()
    
    # Main content area
    if selected_tool == "Link Contextualizer":
        # Render link tool
        data = render_link_tool()
        
        # Render content blocks with results
        render_content_blocks(data)
    
    else:
        # Placeholder for other tools
        st.header(f"🛠️ {selected_tool}")
        st.info(f"The {selected_tool} is coming soon! Stay tuned.")
        render_empty_state()
    
    # Footer
    st.markdown("---")
    st.markdown(
        """
        <div style='text-align: center; color: #666; padding: 2rem 0;'>
            <p>AI Tools Platform | Powered by Streamlit & AI</p>
        </div>
        """,
        unsafe_allow_html=True
    )


if __name__ == "__main__":
    main()

