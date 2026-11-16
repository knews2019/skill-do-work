#!/usr/bin/env python3
"""insert-section.py <req-path> <before-heading-or-END> <file>: insert a full section (with its ## heading) before a heading line, or at the end."""
import sys
p,before,f=sys.argv[1:4]
s=open(p).read(); body=open(f).read().rstrip("\n")+"\n"
if before=="END":
    s=s.rstrip("\n")+"\n\n"+body
else:
    marker="\n"+before+"\n"
    assert s.count(marker)==1, (before, s.count(marker))
    s=s.replace(marker,"\n"+body+"\n"+before+"\n",1)
open(p,"w").write(s)
