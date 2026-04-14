from flask import Flask, request, jsonify
from openai import OpenAI
from dotenv import load_dotenv
import os
import uuid

load_dotenv()

app = Flask(__name__)

client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))
MODEL = os.getenv("OPENAI_MODEL", "gpt-5.4-nano")

# Memoria simple en el servidor
# session_id -> { transcript: str, last_response_id: str | None }
sesiones = {}

MASTER_INSTRUCTIONS = """
Eres un asistente que analiza transcripts de llamadas de call center.
Responde únicamente con base en el transcript proporcionado.
No inventes información.
No uses conocimiento externo.
Si la respuesta no está claramente en el transcript, responde exactamente:
No encuentro esa información en el transcript.
Responde en español, de forma clara y breve.
"""


@app.route("/", methods=["GET"])
def home():
    return jsonify({
        "ok": True,
        "message": "Servidor Flask funcionando correctamente."
    })


@app.route("/set-context", methods=["POST"])
def set_context():
    try:
        data = request.get_json(silent=True)

        transcript = ""
        if data and "transcript" in data:
            transcript = data.get("transcript", "").strip()
        else:
            transcript = request.data.decode("utf-8").strip()

        if not transcript:
            return jsonify({
                "ok": False,
                "error": "Debes enviar el transcript."
            }), 400

        session_id = str(uuid.uuid4())

        sesiones[session_id] = {
            "transcript": transcript,
            "last_response_id": None
        }

        return jsonify({
            "ok": True,
            "message": "Contexto cargado correctamente.",
            "session_id": session_id,
            "characters": len(transcript)
        })

    except Exception as error:
        return jsonify({
            "ok": False,
            "error": str(error)
        }), 500


@app.route("/ask", methods=["POST"])
def ask():
    try:
        data = request.get_json(silent=True)

        if not data:
            return jsonify({
                "ok": False,
                "error": "Debes enviar JSON."
            }), 400

        session_id = data.get("session_id", "").strip()
        question = data.get("question", "").strip()

        if not session_id:
            return jsonify({
                "ok": False,
                "error": "El campo 'session_id' es obligatorio."
            }), 400

        if not question:
            return jsonify({
                "ok": False,
                "error": "El campo 'question' es obligatorio."
            }), 400

        if session_id not in sesiones:
            return jsonify({
                "ok": False,
                "error": "No existe una sesión con ese session_id."
            }), 404

        transcript = sesiones[session_id]["transcript"]
        last_response_id = sesiones[session_id]["last_response_id"]

        input_text = f"""
TRANSCRIPT:
{transcript}

PREGUNTA:
{question}
"""

        # Primera pregunta de la sesión
        if last_response_id is None:
            response = client.responses.create(
                model=MODEL,
                instructions=MASTER_INSTRUCTIONS,
                input=input_text
            )
        else:
            # Preguntas siguientes, manteniendo continuidad
            response = client.responses.create(
                model=MODEL,
                instructions=MASTER_INSTRUCTIONS,
                previous_response_id=last_response_id,
                input=input_text
            )

        sesiones[session_id]["last_response_id"] = response.id

        return jsonify({
            "ok": True,
            "model": MODEL,
            "session_id": session_id,
            "answer": response.output_text
        })

    except Exception as error:
        return jsonify({
            "ok": False,
            "error": str(error)
        }), 500


@app.route("/clear-context", methods=["POST"])
def clear_context():
    try:
        data = request.get_json(silent=True)

        if not data:
            return jsonify({
                "ok": False,
                "error": "Debes enviar JSON."
            }), 400

        session_id = data.get("session_id", "").strip()

        if not session_id:
            return jsonify({
                "ok": False,
                "error": "El campo 'session_id' es obligatorio."
            }), 400

        if session_id in sesiones:
            del sesiones[session_id]

        return jsonify({
            "ok": True,
            "message": "Contexto eliminado correctamente."
        })

    except Exception as error:
        return jsonify({
            "ok": False,
            "error": str(error)
        }), 500


if __name__ == "__main__":
    app.run(debug=True, host="127.0.0.1", port=5000)