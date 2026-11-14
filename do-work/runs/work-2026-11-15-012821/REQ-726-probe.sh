#!/usr/bin/env bash
# REQ-726 GREEN probe: the "no external review canvas" decision record exists, is accepted,
# names its sources, and is linked from the log, the master index and the topic index.
# Fails at base by design (no such record yet). Runs in whichever tree is the current directory.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

fail() { printf 'REQ-726 probe FAIL: %s\n' "$1" >&2; exit 1; }

record_matches=(docs/decisions/records/adr-*review-canvas*.md)
[ -e "${record_matches[0]}" ] || fail "no docs/decisions/records/adr-*review-canvas*.md record"
[ "${#record_matches[@]}" -eq 1 ] || fail "more than one review-canvas record: ${record_matches[*]}"
record_path="${record_matches[0]}"
record_slug="$(basename "$record_path" .md)"
record_text="$(cat "$record_path")"

grep -q '^status: accepted$' <<<"$record_text" || fail "$record_path has no 'status: accepted' line"
for required_fact in \
  '2026-11-14_1452_proposal-tinka-review-surface' \
  '2026-11-14_1547_proposal-tinka-review-surface' \
  'REQ-721' '0.316.0' '0.316.1' '0.317.0' '127.0.0.1' '403'; do
  grep -qF -- "$required_fact" <<<"$record_text" || fail "$record_path does not name $required_fact"
done

for linking_file in docs/decisions/log.md docs/decisions/_master_index.md \
  docs/decisions/topics/_index_pipeline-deliverables.md; do
  linking_text="$(cat "$linking_file")"
  grep -qF -- "[[$record_slug]]" <<<"$linking_text" || fail "$linking_file does not link [[$record_slug]]"
done

printf 'REQ-726 probe PASS (%s)\n' "$record_slug"
