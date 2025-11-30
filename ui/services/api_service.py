import requests
import config

from typing import Dict, Optional
from urllib.parse import quote


def get_contextualized_link(client: str, link: str) -> str:
    url = f"{config.BACKEND_API_URL}/contextualize-link?link={link}"
    response = requests.get(url)

    d_result = response.json()
    if response.status_code != 200:
        raise Exception(d_result)

    result = d_result["result"]
    format_link = f"{config.BACKEND_API_URL}/redirect/{client}/{result['contextualized_link']}"
    return {
        "original_link": link,
        "contextualized_link": format_link
    }
