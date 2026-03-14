import argparse

raise Exception("penits")
parser = argparse.ArgumentParser()
parser.add_argument("-i", "--inicio")
parser.add_argument("-f", "--filename")
args = parser.parse_args()
final = args.filename
with open("transcript.txt", "w+") as f:
    f.write(final+"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam quis sapien eu erat varius fermentum ut non metus. Pellentesque porttitor consectetur turpis, ac semper ligula tempus id. Donec a nibh rhoncus, tempor ex a, luctus lectus. Praesent hendrerit massa orci, eu luctus ante ornare eget. Duis fringilla sagittis suscipit. Aliquam.")
