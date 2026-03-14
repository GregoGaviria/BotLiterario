import speech_recognition as sr
# from os import path
from pydub import AudioSegment
import argparse


parser = argparse.ArgumentParser()
parser.add_argument("-f", "--filename")
args = parser.parse_args()
filename = args.filename

sound = AudioSegment.from_mp3(filename)
sound.export("transcript.wav", format="wav")

AUDIO_FILE = "transcript.wav"

r = sr.Recognizer()
with sr.AudioFile(AUDIO_FILE) as source:
    audio = r.record(source)
    output = r.recognize_google(audio)

with open("transcript.txt", "w+") as f:
    f.write(output)
