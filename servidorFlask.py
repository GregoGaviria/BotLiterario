from flask import Flask, request
from dotenv import load_dotenv
from llama_cloud import AsyncLlamaCloud

import os

load_dotenv()

# Get the API key from environment variables
Llama_api_key = os.getenv("LLAMA_API_KEY")
client = AsyncLlamaCloud(Llama_api_key)


app = Flask(__name__)


@app.route("/")
def prompt_placeholder():
    prompt = request.args.get("prompt")
    return prompt


if __name__ == "__main__":
    app.run(host="127.0.0.1", port=8090)




async def Llamacall(text: Document):

    file_obj = await client.files.create(file = text, purpose="parse")
    result = await client.parsing.parse(
        file_id=file_obj.id,
        tier="agentic",
        version="latest",
        expand=["markdown_full", "text_full"],
    )

    print("Full markdown:")
    print(result.markdown_full)

    print("Full text:")
    print(result.text_full)
