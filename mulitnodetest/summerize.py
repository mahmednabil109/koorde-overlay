#!/usr/bin/python

import re
import os
import pandas as pd
from pathlib import Path

files = Path('.').glob('*.log')
data = []

for f in files:
    if not re.match(r".+\d{4}\.log$", f.name):
        continue
    with f.open() as log:
        text = f.read_text()
        lkpm = len(re.findall("IN lookup", text))
        brom = len(re.findall("IN Broadcast", text))
        corr = len(re.findall("Correction ->", text))
        forw = len(re.findall("Forward ->", text))
        succ = len(re.findall("Successor \|\|", text))
        init = len(re.findall("init", text))
        data.append([corr, forw, succ, init, lkpm, brom])

df = pd.DataFrame(
    data,
    columns=['Correction', 'Forward', 'Successor',
             'init', 'Lookup Messages', 'Broadcast Messages']
)
df.to_csv(f'data{len(data)}.csv')
df.describe()
