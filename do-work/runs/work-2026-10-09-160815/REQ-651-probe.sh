#!/bin/bash
# Focused checks for REQ-651: the ai-report bundle for REQ-632 exists and carries the two titled parts.
set -euo pipefail
bundle=$(ls -d ai-reports/*REQ-632*/ | tail -1)
test -f "$bundle/index.html"
grep -q "Deterministic git structure" "$bundle/index.html"
grep -q "Scaffold around agent behaviour" "$bundle/index.html"
