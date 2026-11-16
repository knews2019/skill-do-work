#!/usr/bin/env python3
"""add-lesson.py <repo-root> <satellite-rel-path> <family-slug> <applies-phrase> <bullet-file>: add one lesson bullet at the top of a satellite and refresh its lessons-index row."""
import sys,os,re
root,sat,slug,phrase,bf=sys.argv[1:6]
sp=os.path.join(root,sat); s=open(sp).read(); bullet=open(bf).read().strip()
assert bullet.startswith("- [family: %s] "%slug)
i=s.index("\n- [family: "); s=s[:i]+"\n"+bullet+"\n"+s[i:]
open(sp,"w").write(s)
tokens=(len(s.encode())+3)//4
ip=os.path.join(root,"do-work/lessons-index.md"); t=open(ip).read().split("\n")
for n,line in enumerate(t):
    if line.startswith("| `%s` |"%sat):
        cells=line.split(" | ")
        assert len(cells)==5, len(cells)
        fam=set(cells[2].strip("`").split(", ")); fam.add(slug)
        cells[2]="`"+", ".join(sorted(fam))+"`"
        if phrase and phrase not in cells[1]: cells[1]=cells[1]+", or "+phrase if False else cells[1].rstrip()+", "+phrase
        cells[3]=str(tokens)
        t[n]=" | ".join(cells); break
else: raise SystemExit("row not found")
open(ip,"w").write("\n".join(t))
print("tokens",tokens)
