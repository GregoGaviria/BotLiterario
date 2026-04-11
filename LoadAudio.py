import argparse
import whisper

parser = argparse.ArgumentParser()
parser.add_argument("-f", "--filename")
args = parser.parse_args()
filename = args.filename

model = whisper.load_model("turbo")
result = model.transcribe(filename)

with open("transcript.txt", "w+") as f:
    f.write(result["text"])
