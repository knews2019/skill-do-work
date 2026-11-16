#!/usr/bin/env python3
"""set-estimate.py <req-path> <p50> <confidence> <basis>... : insert the frozen estimate block before claimed_at."""
import sys,datetime
p,p50,conf=sys.argv[1:4]; basis=sys.argv[4:]
s=open(p).read()
assert "\nestimate:" not in s.split("\n---\n",1)[0], "estimate already present"
now=datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
block="estimate:\n  p50_active_minutes: %s\n  confidence: %s\n  calculated_at: %s\n  basis:\n%s" % (p50,conf,now,"".join("    - %s\n"%b for b in basis))
assert "\nclaimed_at:" in s
s=s.replace("\nclaimed_at:","\n"+block+"claimed_at:",1)
open(p,"w").write(s)
