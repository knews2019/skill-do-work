#!/usr/bin/env python3
"""pau.py <req-path> <plan-text> <apply-text> <unify-text>: tick the three P-A-U boxes with the builder's notes."""
import sys,re
p,plan,apply_,unify=sys.argv[1:5]
s=open(p).read()
for tag,text in (("PLAN",plan),("APPLY",apply_),("UNIFY",unify)):
    pat=re.compile(r"^- \[ \] \*\*\[%s\]:\*\* .*$"%tag, re.M)
    assert pat.search(s), tag
    s=pat.sub(lambda m: "- [x] **[%s]:** %s"%(tag,text), s, count=1)
open(p,"w").write(s)
