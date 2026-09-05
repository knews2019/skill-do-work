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

# Finding 5: toolbox-shims-no-callers (REQ-551)
callerless_shims="$(
  for f in "$repo_root"/skills/do-work-toolbox/scripts/*.sh; do
    [ -e "$f" ] || continue
    b=$(basename "$f")
    n=$(rg -l --fixed-strings "$b" "$repo_root/skills" "$repo_root/tools" "$repo_root/suite" "$repo_root/README.md" "$repo_root/CLAUDE.md" "$repo_root/_dev/primes" --glob '!*CHANGELOG*' | grep -v "/$b$" | wc -l | tr -d ' ')
    [ "$n" -eq 0 ] && echo "$f"
  done
)"
if [ -n "$callerless_shims" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: caller-less toolbox shell shim found: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$callerless_shims"
fi

# Shipped shell delegating check (REQ-551 companion)
non_delegating="$(
  find "$repo_root/skills" "$repo_root/tools" -name '*.sh' | while read -r f; do
    rg -q do-work-cli "$f" || echo "NON-DELEGATING: $f"
  done
)"
if [ -n "$non_delegating" ]; then
  while IFS= read -r line; do
    [ -z "$line" ] && continue
    printf 'FAIL: %s\n' "$line" >&2
    failure_count=$((failure_count + 1))
  done <<< "$non_delegating"
fi

# Finding 8: dead-path-pointers-in-records (REQ-549)
# Topic-index `sources:` lists and the primes are read as live routing, so a token that
# matches no tracked file costs every reader a search. Dated records keep their citations
# as written and are not scanned: decisions/records/, decisions/audits/,
# decisions/imported-specs/ and decisions/log.md cite the tree as it was.
tracked_file_list="$(cd "$repo_root" && git ls-files)"

report_token_if_dead() {
  local candidate_token="$1"
  local citing_file="$2"
  local resolvable_token="$candidate_token"

  # Not a repo-relative file claim, so not checked: a glob pattern, a token rooted
  # somewhere the reader resolves from instead (leading "/"), or a bare `do-work/...`
  # token, which is the consuming project's queue state. Same carve-outs the shipped
  # reference contract draws — see _dev/primes/prime-action-files.md Cross-Referencing.
  case "$candidate_token" in
    '' | *' '* | *'*'* | /* | do-work/*) return 0 ;;
  esac
  while [ "${resolvable_token#../}" != "$resolvable_token" ]; do
    resolvable_token="${resolvable_token#../}"
  done
  # Substring match, the same reading `git ls-files | grep -F` gives: a token is live when
  # some tracked path ends with it. Matched in-shell because `grep -q` closes the pipe
  # early and `pipefail` would then read the writer's SIGPIPE as a miss.
  case "$tracked_file_list" in
    *"$resolvable_token"*) return 0 ;;
  esac
  printf '%s\t%s\n' "$candidate_token" "${citing_file#"$repo_root/"}"
}

dead_routing_tokens="$(
  for index_file in "$repo_root"/decisions/topics/*.md; do
    [ -e "$index_file" ] || continue
    awk '/^sources:/ { in_sources = 1; next }
         /^[a-z_]+:/ { in_sources = 0 }
         in_sources && /^  - / { sub(/^  - /, ""); print }' "$index_file" |
      while IFS= read -r source_entry; do
        report_token_if_dead "$source_entry" "$index_file"
      done
  done
  # A backticked token carrying a directory and a file extension is a path claim.
  for prime_file in "$repo_root"/_dev/primes/*.md; do
    [ -e "$prime_file" ] || continue
    grep -o '`[^`]*`' "$prime_file" | tr -d '`' | grep -E '/.*\.[a-z]{2,4}$' |
      while IFS= read -r prime_citation; do
        report_token_if_dead "$prime_citation" "$prime_file"
      done
  done
)"

if [ -n "$dead_routing_tokens" ]; then
  while IFS=$'\t' read -r dead_token citing_file; do
    [ -z "$dead_token" ] && continue
    printf 'FAIL: %s cites %s, which matches no tracked file\n' "$citing_file" "$dead_token" >&2
    failure_count=$((failure_count + 1))
  done <<< "$dead_routing_tokens"
fi

# Finding 2: cli-launcher-preamble-copied (REQ-553)
# The two exempt paths are named in full, not by trailing path shape: a third copy at
# skills/do-work-toolbox/tools/do-work-cli-preamble.sh is exactly the hand-rolled preamble
# this pins at zero, and a suffix filter would read it as the exempt file.
hand_rolled_preambles="$(
  rg -l --glob '*.sh' 'for cli_candidate in|^launcher_arguments=\(--format text\)$' \
    "$repo_root/skills" "$repo_root/tools" 2>/dev/null \
    | grep -vxF -e "$repo_root/tools/do-work-cli-preamble.sh" \
                -e "$repo_root/skills/do-work/tools/do-work-cli-preamble.sh"
)"
if [ -n "$hand_rolled_preambles" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$hand_rolled_preambles"
fi

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi

printf 'Audit lock-in regressions passed.\n'

