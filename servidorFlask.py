from flask import Flask, request, jsonify
from openai import OpenAI
from dotenv import load_dotenv
import os

load_dotenv()

app = Flask(__name__)

client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))
MODEL = os.getenv("OPENAI_MODEL", "gpt-5.4")


@app.route("/", methods=["GET"])
def home():
    return jsonify({
        "ok": True,
        "message": "Servidor Flask funcionando correctamente."
    })


@app.route("/ask", methods=["POST"])
def ask_openai():
    try:
        data = request.get_json()

        if not data:
            return jsonify({
                "ok": False,
                "error": "No se recibió JSON en la petición."
            }), 400

        prompt = data.get("prompt", "").strip()

        if not prompt:
            return jsonify({
                "ok": False,
                "error": "El campo 'prompt' es obligatorio."
            }), 400

        response = client.chat.completions.create(
            model=MODEL,
            messages=[
                {"role": "system", "content": "Responde de forma clara, breve y precisa."},
                {"role": "user", "content": prompt}
            ]
        )

        return jsonify({
            "ok": True,
            "model": MODEL,
            "answer": response.choices[0].message.content
        })

    except Exception as error:
        return jsonify({
            "ok": False,
            "error": str(error)
        }), 500


if __name__ == "__main__":
    app.run(debug=True, host="127.0.0.1", port=5000)