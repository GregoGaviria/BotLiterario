import os
import webbrowser
from threading import Thread

# flaskThread = Thread(target=lambda: os.system("python servidorFlask.py"))
goThread = Thread(target=lambda: os.system("go run ."))
# flaskThread.start()
goThread.start()
webbrowser.open("localhost:8000")