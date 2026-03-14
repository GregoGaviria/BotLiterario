from flask import Flask, request
from dotenv import load_dotenv
from llama_index.core import VectorStoreIndex, Document
import os

load_dotenv()
app = Flask(__name__)

# Variables globales simples para prueba local
index = None
chat_engine = None


@app.route("/")
def prompt_placeholder():
    prompt = request.args.get("prompt")
    return prompt or "No prompt received"


async def Llamacall(texto: str):
    global index, chat_engine

    documents = [Document(text=texto)]
    index = VectorStoreIndex.from_documents(documents)
    chat_engine = index.as_chat_engine()


async def Llamatalk(text: str):
    global chat_engine

    if chat_engine is None:
        return "No hay chat Engine, primero hay que mandar un archivo"

    response = chat_engine.chat(text)
    return str(response)


if __name__ == "__main__":
    app.run(host="127.0.0.1", port=8090)