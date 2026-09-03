#!/usr/bin/env bash
# REQ-501: real ownership-transfer behavior across recover, select, and reclaim.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
fail_count=0
# REQ-501: the public sole-authority path must cross one clean ownership-transfer
# boundary before ordinary selection and claiming resume. Exercise the shipped
# launcher because package-local state tests cannot prove the action's command chain.
if ! (
  set -euo pipefail
  recovery_fixture="$(mktemp -d "${TMPDIR:-/tmp}/do-work-rwr-recovery.XXXXXX")"
  trap 'rm -rf -- "$recovery_fixture"' EXIT
  git -C "$recovery_fixture" init -q
  git -C "$recovery_fixture" config user.name 'Contract Fixture'
  git -C "$recovery_fixture" config user.email 'contract@example.invalid'
  mkdir -p "$recovery_fixture/do-work/queue"
  cat > "$recovery_fixture/do-work/queue/REQ-901.md" <<'EOF'
---
id: REQ-901
title: Recovery fixture
status: pending
created_at: 2026-09-02T00:00:00Z
---

# Recovery Fixture

User-authored body stays.
EOF
  cat > "$recovery_fixture/do-work/CHECKPOINT.md" <<'EOF'
# Session Checkpoint

## In Progress (interrupted)
EOF
  printf 'committed implementation base\n' > "$recovery_fixture/implementation.txt"
  git -C "$recovery_fixture" add do-work implementation.txt
  git -C "$recovery_fixture" commit -qm 'seed recovery fixture'

  "$repo_root/skills/do-work/tools/do-work-cli.sh" \
    --repo-root "$recovery_fixture" --format json claim REQ-901 \
    --request-path do-work/queue/REQ-901.md --provenance explicit-req \
    --writer 'foreign:/checkout' --at 2026-09-02T01:00:00Z >/dev/null
  printf 'uncommitted implementation survives\n' >> "$recovery_fixture/implementation.txt"

  "$repo_root/skills/do-work/tools/do-work-cli.sh" \
    --repo-root "$recovery_fixture" --format json recover-claim REQ-901 \
    --request-path do-work/working/REQ-901.md \
    --checkpoint-writer 'foreign:/checkout' --assume-sole-writer --commit \
    --at 2026-09-02T02:00:00Z > "$recovery_fixture/recovery.json"
  "$repo_root/skills/do-work/tools/do-work-cli.sh" \
    --repo-root "$recovery_fixture" --format json next REQ-901 \
    > "$recovery_fixture/selection.json"
  python3 - "$recovery_fixture/recovery.json" "$recovery_fixture/selection.json" <<'PY'
import json
import pathlib
import sys

recovery = json.loads(pathlib.Path(sys.argv[1]).read_text())
selection = json.loads(pathlib.Path(sys.argv[2]).read_text())
if recovery.get("outcome") != "success":
    raise SystemExit(f"recover-claim did not succeed: {recovery!r}")
selected = [record.get("request_id") for record in selection.get("selected", [])]
if "REQ-901" not in selected:
    raise SystemExit(f"recovered REQ is not selectable: {selection!r}")
PY
  "$repo_root/skills/do-work/tools/do-work-cli.sh" \
    --repo-root "$recovery_fixture" --format json claim REQ-901 \
    --request-path do-work/queue/REQ-901.md --provenance explicit-req \
    --writer 'current:/checkout' --at 2026-09-02T03:00:00Z \
    > "$recovery_fixture/claim.json"
  python3 - "$recovery_fixture/claim.json" <<'PY'
import json
import pathlib
import sys

claim = json.loads(pathlib.Path(sys.argv[1]).read_text())
if claim.get("outcome") != "success":
    raise SystemExit(f"fresh claim did not succeed: {claim!r}")
PY
  grep -Fq 'uncommitted implementation survives' "$recovery_fixture/implementation.txt"
  ! grep -Fq 'writer: foreign:/checkout' "$recovery_fixture/do-work/CHECKPOINT.md"
  grep -Fq 'writer: current:/checkout' "$recovery_fixture/do-work/CHECKPOINT.md"
)
then
  printf 'FAIL: run-with-recovery cannot transfer ownership through recovery, selection, and a fresh claim (REQ-501).\n' >&2
  fail_count=$((fail_count + 1))
fi
[ "$fail_count" -eq 0 ] || exit 1
printf 'request-state contract probes passed.\n'
