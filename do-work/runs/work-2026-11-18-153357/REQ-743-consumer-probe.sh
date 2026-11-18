#!/usr/bin/env bash
# Exercise a real first install or update in a disposable Git project.
set -euo pipefail

if [[ $# -lt 3 || $# -gt 4 ]]; then
  printf 'Usage: bash %s BASE CANDIDATE fresh|run-do-work-update|do-work-update [hand-made]\n' "$0" >&2
  exit 2
fi
base="$1"
candidate="$2"
mode="$3"
hand_made="${4:-}"
case "$mode" in fresh|run-do-work-update|do-work-update) ;; *) exit 2 ;; esac
[[ -z "$hand_made" || "$hand_made" == hand-made ]] || exit 2

root="$(git rev-parse --show-toplevel)"
cd "$root"
scratch="$(mktemp -d "$root/do-work/runs/work-2026-11-18-153357/.consumer.XXXXXX")"
server_pid=''
cleanup() {
  if [[ -n "$server_pid" ]]; then
    kill "$server_pid" 2>/dev/null || true
    wait "$server_pid" 2>/dev/null || true
  fi
  rm -rf -- "$scratch"
}
trap cleanup EXIT
git archive --format=tar.gz --prefix=suite/ -o "$scratch/base.tar.gz" "$base"
git archive --format=tar.gz --prefix=suite/ -o "$scratch/candidate.tar.gz" "$candidate"
mkdir "$scratch/base" "$scratch/candidate" "$scratch/project"
tar xzf "$scratch/base.tar.gz" -C "$scratch/base" --strip-components=1
tar xzf "$scratch/candidate.tar.gz" -C "$scratch/candidate"
base_version="$(cat "$scratch/base/VERSION")"
candidate_version="$(cat "$scratch/candidate/suite/VERSION")"

if [[ "$mode" != fresh ]]; then
  # Keep tree contents from the candidate. Advance fixture metadata only if needed
  # to exercise the written-update path before a release version exists.
  candidate_version="$(python3 - "$scratch/candidate/suite" "$base_version" <<'PY'
from pathlib import Path
import re
import sys

root = Path(sys.argv[1])
base = tuple(map(int, sys.argv[2].split(".")))
version = (root / "VERSION").read_text().strip()
if tuple(map(int, version.split("."))) <= base:
    version = f"{base[0]}.{base[1]}.{base[2] + 1}"
    for name in ("VERSION", "skills/do-work/VERSION"):
        (root / name).write_text(version + "\n")
    action = root / "skills/do-work/actions/version.md"
    text, count = re.subn(r"(?m)^\*\*Current version\*\*:[^\n]*$", f"**Current version**: {version}", action.read_text())
    assert count == 1
    action.write_text(text)
print(version)
PY
)"
  tar czf "$scratch/candidate.tar.gz" -C "$scratch/candidate" suite
fi

project="$scratch/project"
git -C "$project" init -q
git -C "$project" config user.name 'REQ-743 verification'
git -C "$project" config user.email 'req-743@example.invalid'
printf '%s\n' 'Consumer content stays intact.' > "$project/README.md"
page="$project/docs/report-gallery/index.html"
if [[ -n "$hand_made" ]]; then
  mkdir -p "$(dirname "$page")"
  printf '<!doctype html><title>Hand-made gallery</title>\n' > "$scratch/hand-made.html"
  cp "$scratch/hand-made.html" "$page"
  git -C "$project" add -- docs/report-gallery/index.html
fi
git -C "$project" add -- README.md
git -C "$project" commit -qm 'Seed clean consumer'
installer=skills/do-work/tools/install-do-work-suite.sh
log="$scratch/operation.log"
if [[ "$mode" == fresh ]]; then
  bash "$scratch/candidate/suite/$installer" --project-root "$project" --archive "$scratch/candidate.tar.gz" <<< y >"$log" 2>&1 || { cat "$log"; exit 1; }
else
  bash "$scratch/base/$installer" --project-root "$project" --archive "$scratch/base.tar.gz" <<< y >"$scratch/seed.log" 2>&1 || { cat "$scratch/seed.log"; exit 1; }
  if [[ -z "$hand_made" && -e "$page" ]]; then
    printf 'The baseline install already has a gallery. Choose the pre-feature base.\n' >&2
    exit 1
  fi
  git -C "$project" add -A
  git -C "$project" commit -qm 'Seed installed baseline'
  [[ -z "$(git -C "$project" status --porcelain)" ]]
  python3 -u - "$scratch/candidate.tar.gz" "$scratch/port" <<'PY_SERVER' >"$scratch/server.log" 2>&1 &
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path
import signal
import sys

archive = Path(sys.argv[1]).read_bytes()
signal.alarm(120)
class ArchiveHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Length", str(len(archive)))
        self.end_headers()
        self.wfile.write(archive)

with HTTPServer(("127.0.0.1", 0), ArchiveHandler) as server:
    Path(sys.argv[2]).write_text(str(server.server_port))
    server.serve_forever()
PY_SERVER
  server_pid=$!
  for _ in {1..50}; do
    [[ -s "$scratch/port" ]] && break
    sleep 0.1
  done
  [[ -s "$scratch/port" ]] || { cat "$scratch/server.log"; exit 1; }
  url="http://127.0.0.1:$(cat "$scratch/port")/candidate.tar.gz"
  (cd "$project" && just "$mode" --upstream-url "$url" <<< y) >"$log" 2>&1 || { cat "$log"; exit 1; }
  grep -Fq -- "Updated to v$candidate_version" "$log" || { cat "$log"; exit 1; }
fi

printf '%s exited 0 (base v%s, installed v%s)\n' "$mode" "$base_version" "$candidate_version"
[[ "$(cat "$project/.claude/skills/do-work/VERSION")" == "$candidate_version" ]]
git -C "$project" diff --exit-code HEAD -- README.md
if [[ -n "$hand_made" ]]; then
  cmp "$scratch/hand-made.html" "$page"
  diagnostic_count="$(awk '/docs\/report-gallery\/index\.html/{count++} END{print count+0}' "$log")"
  [[ "$diagnostic_count" == 1 ]] || { cat "$log"; exit 1; }
  printf 'Hand-made gallery preserved and reported exactly once\n'
else
  [[ -f "$page" ]] || { printf 'Expected gallery missing after successful %s\n' "$mode" >&2; exit 1; }
  for sample in two-measures-side-by-side verdict-with-checks parts-over-time rhythm-over-cycle; do
    grep -Fq -- "$sample" "$page"
    [[ -f "$project/.claude/skills/do-work-toolbox/docs/report-gallery/$sample/sample.md" ]]
    [[ ! -e "$project/docs/report-gallery/$sample" ]]
  done
  printf 'Generated gallery includes all four installed samples\n'
fi
