"""
API Service layer for communicating with the backend API.
"""

import requests
from typing import Dict, Optional
from urllib.parse import quote
import config


class APIService:
    """Service class for handling all backend API communications."""
    
    def __init__(self, base_url: Optional[str] = None, timeout: Optional[int] = None):
        """
        Initialize API service.
        
        Args:
            base_url: Base URL of the backend API
            timeout: Request timeout in seconds
        """
        self.base_url = base_url or config.BACKEND_API_URL
        self.timeout = timeout or config.API_TIMEOUT
        self.session = requests.Session()
    
    def _make_request(self, method: str, endpoint: str, **kwargs) -> Dict:
        """
        Make HTTP request to the API.
        
        Args:
            method: HTTP method (GET, POST, etc.)
            endpoint: API endpoint path
            **kwargs: Additional arguments to pass to requests
            
        Returns:
            Dict: Parsed JSON response
            
        Raises:
            APIException: On API errors
        """
        url = f"{self.base_url}{endpoint}"
        
        try:
            response = self.session.request(
                method=method,
                url=url,
                timeout=self.timeout,
                **kwargs
            )
            response.raise_for_status()
            return response.json()
            
        except requests.exceptions.Timeout:
            raise APIException("Request timed out. Please try again.")
        except requests.exceptions.ConnectionError:
            raise APIException(
                f"Cannot connect to backend API at {self.base_url}. "
                "Please ensure the backend server is running."
            )
        except requests.exceptions.HTTPError as e:
            if response.status_code == 400:
                raise APIException("Bad request. Please check your input.")
            elif response.status_code == 404:
                raise APIException("Resource not found.")
            elif response.status_code == 500:
                error_msg = "Internal server error."
                try:
                    error_data = response.json()
                    if "error_msg" in error_data:
                        error_msg = error_data["error_msg"]
                except:
                    pass
                raise APIException(error_msg)
            else:
                raise APIException(f"HTTP error occurred: {str(e)}")
        except requests.exceptions.RequestException as e:
            raise APIException(f"Request failed: {str(e)}")
        except ValueError:
            raise APIException("Invalid JSON response from server.")
    
    def contextualize_link(self, link: str) -> Dict:
        """
        Contextualize a link using the backend API.
        
        Args:
            link: Original URL to contextualize
            
        Returns:
            Dict: API response containing:
                - result: Dict with 'link' and 'contextualized_link'
                - status_code: HTTP status code
                - error_msg: Error message if any
                
        Raises:
            APIException: On API errors
        """
        if not link:
            raise APIException("Link cannot be empty.")
        
        endpoint = config.API_ENDPOINTS["contextualize_link"]
        params = {"link": link}
        
        response = self._make_request("GET", endpoint, params=params)
        
        # Validate response structure
        if "result" not in response:
            raise APIException("Invalid response format from server.")
        
        return response
    
    def get_redirect_url(self, client: str, contextualized_link: str) -> str:
        """
        Construct redirect URL for a contextualized link.
        
        Args:
            client: Client identifier
            contextualized_link: Contextualized link to redirect
            
        Returns:
            str: Full redirect URL
        """
        endpoint = config.API_ENDPOINTS["redirect"].format(
            client=quote(client),
            link=quote(contextualized_link)
        )
        return f"{self.base_url}{endpoint}"
    
    def health_check(self) -> bool:
        """
        Check if the backend API is healthy.
        
        Returns:
            bool: True if API is healthy, False otherwise
        """
        try:
            response = self._make_request("GET", "/healthz")
            return response.get("message") == "OK"
        except:
            return False


class APIException(Exception):
    """Custom exception for API errors."""
    pass

