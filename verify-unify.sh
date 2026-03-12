#!/usr/bin/env bash
set -euo pipefail

RED='\033[1;31m'
GREEN='\033[1;32m'
NC='\033[0m'

# Ensure we are inside a git repository
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo -e "${RED}UNIFY FAILED: Not a git repository.${NC}" >&2
  exit 1
fi

# Capture only newly-added lines (starting with "+") from staged + unstaged diff
diff_output=$(git diff HEAD 2>/dev/null || git diff 2>/dev/null || true)

if [ -z "$diff_output" ]; then
  echo -e "${GREEN}UNIFY PASSED: Diff is clean.${NC}"
  exit 0
fi

# Grep added lines for common debug/leftover statements
matches=$(echo "$diff_output" \
  | grep '^+' \
  | grep -v '^+++' \
  | grep -En 'console\.log|print\(|TODO: AI|FIXME:' || true)

if [ -n "$matches" ]; then
  echo -e "${RED}UNIFY FAILED: Debug statements found. Please remove them.${NC}" >&2
  echo "$matches" >&2
  exit 1
fi

echo -e "${GREEN}UNIFY PASSED: Diff is clean.${NC}"
exit 0
