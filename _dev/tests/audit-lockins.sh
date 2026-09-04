#!/usr/bin/env bash
# Audit lock-in regressions. Pinned after maintainability audits.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
failure_count=0

# Finding 10: exported one-line delegates with no production caller (REQ-550)
delegate_findings="$(
  for d in "$repo_root/skills/do-work/tools/do-work-cli" "$repo_root/skills/do-work-board/tools/queue-kanban"; do
    find "$d" -name '*.go' ! -name '*_test.go' -exec awk '
      FNR==1{f=FILENAME}
      /^func [A-Z]/{
        fn=$0; l=FNR; getline a; getline b;
        if (b ~ /^}$/ && a ~ /^\t(return )?[a-zA-Z0-9_.]+\(/) {
          match(fn,/func [A-Za-z0-9_]+/);
          print substr(fn,RSTART+5,RLENGTH-5) "\t" f ":" l
        }
      }' {} \;
  done | while IFS=$'\t' read -r name loc; do
    [ -z "$name" ] && continue
    if [ "$name" = "AssociateProjectPaths" ] || [ "$name" = "ApplyTimestampPlan" ] || \
       [ "$name" = "DownloadAtomic" ] || [ "$name" = "CheckGreenGate" ]; then
      printf '%s\t%s\n' "$name" "$loc"
    else
      target_line="$(sed -n "${loc##*:}"'{n;p;}' "${loc%%:*}")"
      if echo "$target_line" | grep -Eq '^\t(return )?[a-z][a-zA-Z0-9_]*\('; then
        prod=$(rg -c "\b$name\(" --glob '*.go' --glob '!*_test.go' "$repo_root/skills/" 2>/dev/null | grep -v "${loc%%:*}" | wc -l | tr -d ' ')
        [ "$prod" -eq 0 ] && printf '%s\t%s\n' "$name" "$loc"
      fi
    fi
  done
)"

if [ -n "$delegate_findings" ]; then
  while IFS=$'\t' read -r name loc; do
    [ -z "$name" ] && continue
    printf 'FAIL: %s (%s) is an exported one-line delegate with no production caller\n' "$name" "$loc" >&2
    failure_count=$((failure_count + 1))
  done <<< "$delegate_findings"
fi

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi

printf 'Audit lock-in regressions passed.\n'
