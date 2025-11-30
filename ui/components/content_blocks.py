"""
Content Blocks Component for displaying results.
"""

import streamlit as st
from typing import Dict, Optional
from datetime import datetime
from utils.helpers import create_card_style


def render_content_blocks(data: Optional[Dict]) -> None:
    """
    Render content blocks to display link contextualization results.
    
    Args:
        data: Dictionary containing API response with 'result' key
    """
    # Apply custom CSS
    st.markdown(create_card_style(), unsafe_allow_html=True)
    
    if not data:
        # Show placeholder when no data
        st.info("👆 Enter a URL above to see the contextualized results here.")
        return

    # Extract result data
    result = data
    original_link = result.get("original_link", "N/A")
    contextualized_link = result.get("contextualized_link", "N/A")

    st.markdown("---")
    st.subheader("📊 Results")
    
    # Create three columns for layout
    col1, col2 = st.columns([1, 1])
    
    # Block 1: Original Link Display
    with col1:
        with st.container():
            st.markdown("##### 🔗 Original Link")
            st.markdown(
                f'<div class="link-display">{original_link}</div>',
                unsafe_allow_html=True
            )
            if original_link != "N/A":
                st.markdown(f"[Open Original Link ↗]({original_link})")
    
    # Block 2: Contextualized Link Display
    with col2:
        with st.container():
            st.markdown("##### ✨ Contextualized Link")
            st.markdown(
                f'<div class="link-display">{contextualized_link}</div>',
                unsafe_allow_html=True
            )
            if contextualized_link != "N/A":
                # Copy to clipboard functionality
                if st.button("📋 Copy Contextualized Link", key="copy_contextualized"):
                    st.code(contextualized_link, language=None)
                    st.success("✅ Link ready to copy! Use Ctrl+C / Cmd+C")
    
    # Additional details in expander
    with st.expander("🔍 View Full API Response"):
        st.json(data)
    
    # Action buttons
    st.markdown("---")
    col_action1, col_action2, col_action3 = st.columns([1, 1, 2])
    
    with col_action1:
        if st.button("🔄 Process Another Link", use_container_width=True):
            st.session_state.clear()
            st.rerun()
    
    with col_action2:
        if st.button("💾 Download Results", use_container_width=True):
            # Create downloadable JSON
            import json
            json_str = json.dumps(data, indent=2)
            st.download_button(
                label="⬇️ Download JSON",
                data=json_str,
                file_name=f"contextualized_link_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json",
                mime="application/json",
                use_container_width=True
            )


def render_empty_state() -> None:
    """
    Render empty state when no tool is selected or no data available.
    """
    st.info("👈 Select a tool from the sidebar to get started!")
    
    st.markdown("""
    ### Welcome to AI Tools! 🎉
    
    This platform provides various AI-powered tools to enhance your workflow.
    
    **Available Tools:**
    - 🔗 **Link Contextualizer**: Transform URLs into contextually rich links
    
    Get started by selecting a tool from the sidebar!
    """)

