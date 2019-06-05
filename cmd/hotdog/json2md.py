#!/usr/bin/env python3
import json

import argparse

parser = argparse.ArgumentParser()
parser.add_argument('-o', '--output-filepath', required=True)
parser.add_argument('-t', '--title', required=True)
parser.add_argument('-i', '--input-filepath', required=True)
parser.add_argument('args', nargs='*')  # any length of args is ok

args, extra_args = parser.parse_known_args()

filepath = args.input_filepath
file = open(filepath, 'r', encoding="UTF-8")
data = json.load(file)
file.close()

title = args.title
output_filepath = args.output_filepath
with open(output_filepath, 'w', encoding="UTF-8") as file:
    file.write("# {}\n".format(title))
    for i in range(0, len(data)):
        file.write("## {}\n".format(data[i]['Name']))
        for key, val in data[i].items():
            if 'Name' not in key:
                file.write("* {}:{}\n".format(key, val))
        file.write("----\n")
