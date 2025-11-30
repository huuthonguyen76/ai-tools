# AI Tools - Streamlit UI

A modern, user-friendly web interface for AI-powered tools, built with Streamlit.

## 🚀 Features

- **Link Contextualizer**: Transform URLs into contextually rich, AI-enhanced links
- **Intuitive UI**: Clean, modern interface with sidebar navigation
- **Real-time Processing**: Instant feedback with loading states and error handling
- **Responsive Design**: Works seamlessly on different screen sizes

## 📋 Prerequisites

- Python 3.8 or higher
- Backend API server running (default: `http://localhost:8080`)

## 🛠️ Installation

1. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

2. **Configure backend URL (optional):**
   
   Set the `BACKEND_API_URL` environment variable if your backend is not running on `localhost:8080`:
   ```bash
   export BACKEND_API_URL="http://your-backend-url:port"
   ```

## 🎯 Usage

1. **Start the Streamlit app:**
   ```bash
   streamlit run app.py
   ```

2. **Open your browser:**
   
   The app will automatically open at `http://localhost:8501`

3. **Use the Link Contextualizer:**
   - Select "Link Contextualizer" from the sidebar (selected by default)
   - Enter a URL in the input field
   - Click "🚀 Contextualize Link"
   - View the results in the content blocks below

## 📁 Project Structure

```
ui/
├── app.py                      # Main application entry point
├── config.py                   # Configuration settings
├── requirements.txt            # Python dependencies
├── components/                 # UI components
│   ├── __init__.py
│   ├── sidebar.py             # Sidebar navigation
│   ├── link_tool.py           # Link contextualizer tool
│   └── content_blocks.py      # Results display
├── services/                   # Backend communication
│   ├── __init__.py
│   └── api_service.py         # API service layer
└── utils/                      # Utility functions
    ├── __init__.py
    └── helpers.py             # Helper functions
```

## 🔧 Configuration

Edit `config.py` to customize:

- `BACKEND_API_URL`: Backend API endpoint
- `API_TIMEOUT`: Request timeout in seconds
- `AVAILABLE_TOOLS`: List of available tools
- `APP_TITLE` and `APP_ICON`: Branding

## 🌐 Environment Variables

- `BACKEND_API_URL`: Override default backend URL (default: `http://localhost:8080`)
- `API_TIMEOUT`: Override default API timeout in seconds (default: `30`)

## 🧪 Testing

To test the application:

1. Ensure the backend API is running
2. Start the Streamlit app
3. Test with various URLs:
   - Valid URLs: `https://example.com`
   - Invalid URLs to test validation
   - Test error handling by stopping the backend

## 📝 API Integration

The UI communicates with the backend API:

- **Endpoint**: `GET /contextualize-link?link={url}`
- **Response Format**:
  ```json
  {
    "status_code": 200,
    "result": {
      "link": "https://example.com",
      "contextualized_link": "https://contextualized-link.com"
    },
    "error_msg": ""
  }
  ```

## 🎨 Customization

### Adding New Tools

1. Add tool name to `AVAILABLE_TOOLS` in `config.py`
2. Create new component in `components/`
3. Add routing logic in `app.py`

### Styling

Custom CSS can be added in `utils/helpers.py` in the `create_card_style()` function.

## 🐛 Troubleshooting

### Backend Connection Error
- **Issue**: "Cannot connect to backend API"
- **Solution**: Ensure the backend server is running and accessible at the configured URL

### Port Already in Use
- **Issue**: Port 8501 is already in use
- **Solution**: Use `streamlit run app.py --server.port 8502` to run on a different port

### Module Import Errors
- **Issue**: Import errors for components/services
- **Solution**: Ensure you're running the app from the `ui/` directory

## 📄 License

Part of the AI Tools project.

## 🤝 Contributing

Follow the project's contribution guidelines and workflow documented in `.cursor/dev_command/workflow.md`.

