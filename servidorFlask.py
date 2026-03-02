from flask import Flask, request

app = Flask(__name__)


@app.route("/")
def prompt_placeholder():
    prompt = request.args.get("prompt")
    return prompt


if __name__ == "__main__":
    app.run(host="127.0.0.1", port=8090)
