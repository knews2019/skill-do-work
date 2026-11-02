#!/usr/bin/env python3
"""Persist an executed estimate record as the frozen `estimate:` frontmatter block.
usage: put-estimate.py <request-path> <advance-json-file> <calculated-at>"""
import json, sys
req_path, record_file, calculated_at = sys.argv[1:4]
record = json.load(open(record_file))
lines = next(g for g in record["advance"]["gate_records"] if g["gate_id"] == "estimate-p50")["output_lines"]
text = open(req_path).read()
assert text.startswith("---\n")
end = text.index("\n---\n", 4)
front, rest = text[:end + 1], text[end + 1:]
assert "\nestimate:" not in "\n" + front, "estimate block already present"
block = ["estimate:"]
for line in lines:
    if line.startswith("- "):
        block.append("    " + line)
    elif line == "basis:":
        block.append("  calculated_at: " + calculated_at)
        block.append("  basis:")
    else:
        block.append("  " + line)
open(req_path, "w").write(front + "\n".join(block) + "\n" + rest)
