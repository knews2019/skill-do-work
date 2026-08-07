#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fail_count=0

fail() {
  printf 'FAIL: %s\n' "$1" >&2
  fail_count=$((fail_count + 1))
}

require_file() {
  local relative_path="$1"

  if [ ! -f "$repo_root/$relative_path" ]; then
    fail "required staged-suite file is missing: $relative_path"
  fi
}

core_files=(
  skills/do-work/SKILL.md
  skills/do-work/next-steps.md
  skills/do-work/actions/capture.md
  skills/do-work/actions/capture-reference.md
  skills/do-work/actions/work.md
  skills/do-work/actions/work-reference.md
  skills/do-work/actions/verify-requests.md
  skills/do-work/actions/review-work.md
  skills/do-work/actions/kb-lessons-handoff.md
  skills/do-work/actions/clarify.md
  skills/do-work/actions/abandon.md
  skills/do-work/actions/cleanup.md
  skills/do-work/actions/commit.md
  skills/do-work/actions/roadmap.md
  skills/do-work/actions/forensics.md
  skills/do-work/actions/version.md
  skills/do-work/actions/help.md
  skills/do-work/actions/pipeline.md
  skills/do-work/actions/pipeline-reference.md
  skills/do-work/crew-members/coding-guardrails.md
  skills/do-work/crew-members/clear-questions.md
  skills/do-work/specs/bug-fix.md
  skills/do-work/specs/refactor.md
  skills/do-work/hooks/hooks.json
  skills/do-work/hooks/session-start.sh
  skills/do-work/hooks/pipeline-guard.sh
  skills/do-work/tools/checks/preflight.sh
  skills/do-work/tools/checks/qualify.sh
  skills/do-work/tools/checks/record-commit-hash.sh
  skills/do-work/tools/do-work-update.sh
  skills/do-work/tools/validate-suite-manifest.sh
  skills/do-work/tools/replace-text-section.sh
)

for core_file in "${core_files[@]}"; do
  require_file "$core_file"
done

if [ -d "$repo_root/skills" ]; then
  maintainer_citations="$(grep -rIEn '\[[^]]*\]\([^)]*(CLAUDE|AGENTS)\.md|(^|[[:space:]])(see|per|according to)[[:space:]]+`?(CLAUDE|AGENTS)\.md' "$repo_root/skills" || true)"
  if [ -n "$maintainer_citations" ]; then
    printf 'FAIL: staged runtime files must not cite repository maintainer instructions:\n%s\n' "$maintainer_citations" >&2
    fail_count=$((fail_count + 1))
  fi
fi

assert_core_sibling_reference() {
  local relative_file="$1"
  local sibling_name="$2"

  if [ -f "$repo_root/skills/do-work/$relative_file" ] \
    && ! grep -Fq -- "$sibling_name/" "$repo_root/skills/do-work/$relative_file"; then
    fail "core runtime must resolve $relative_file through sibling $sibling_name"
  fi
}

assert_core_sibling_reference actions/capture.md do-work-board
assert_core_sibling_reference actions/work.md do-work-board
assert_core_sibling_reference actions/forensics.md do-work-board
assert_core_sibling_reference actions/kb-lessons-handoff.md do-work-knowledge
assert_core_sibling_reference actions/pipeline.md do-work-toolbox

if [ -f "$repo_root/skills/do-work/SKILL.md" ]; then
  while IFS= read -r action_path; do
    if [ ! -f "$repo_root/skills/do-work/${action_path#./}" ]; then
      fail "core router dispatch target does not resolve: $action_path"
    fi
  done < <(grep -Eo '`\./actions/[A-Za-z0-9._/-]+\.md`' "$repo_root/skills/do-work/SKILL.md" | tr -d '`' | sort -u)
fi

if [ -d "$repo_root/skills/do-work" ]; then
  if ! python3 - "$repo_root/skills/do-work" "$repo_root/suite/modules.tsv" <<'PY'
import pathlib
import re
import sys

skill_root = pathlib.Path(sys.argv[1])
manifest_file = pathlib.Path(sys.argv[2])
declared_siblings = {
    pathlib.PurePosixPath(line.split("\t", 1)[0]).name
    for line in manifest_file.read_text().splitlines()[1:]
    if "\t" in line
}
reference_pattern = re.compile(
    r"(?<![A-Za-z0-9_./-])"
    r"(?P<path>(?:\.\./do-work-(?:board|knowledge|toolbox)/)?"
    r"(?:actions|tools|hooks|crew-members|specs|docs)/[A-Za-z0-9._/-]+)"
)
missing = []
for source in sorted(skill_root.rglob("*.md")):
    if source.name == "CHANGELOG.md":
        continue
    for line_number, line in enumerate(source.read_text(errors="replace").splitlines(), 1):
        for match in reference_pattern.finditer(line):
            reference = match.group("path").rstrip(".,:;)")
            if reference in {"docs/prime-bar.md", "docs/prime-foo.md"}:
                continue
            if reference.startswith("docs/design/"):
                continue
            if reference.startswith("../"):
                target = (skill_root / reference).resolve()
                sibling_root = (skill_root / reference.split("/", 2)[0] / reference.split("/", 2)[1]).resolve()
                if not sibling_root.exists():
                    if sibling_root.name in declared_siblings:
                        continue
                    missing.append(
                        f"{source.relative_to(skill_root)}:{line_number}: undeclared sibling {reference}"
                    )
                    continue
            else:
                target = skill_root / reference
            if not target.exists():
                missing.append(
                    f"{source.relative_to(skill_root)}:{line_number}: {reference}"
                )
if missing:
    print("unresolved staged core runtime references:", file=sys.stderr)
    print("\n".join(missing), file=sys.stderr)
    raise SystemExit(1)
PY
  then
    fail_count=$((fail_count + 1))
  fi
fi

if [ -f "$repo_root/SKILL.md" ]; then
  if ! grep -Fq '| board' "$repo_root/SKILL.md"; then
    fail 'active root router lost its board route before modular cutover'
  fi
  if [ ! -f "$repo_root/actions/board.md" ]; then
    fail 'active root board action was removed before modular cutover'
  fi
fi

if [ -f "$repo_root/skills/do-work/hooks/hooks.json" ]; then
  if ! python3 - "$repo_root/skills/do-work/hooks/hooks.json" "$repo_root/skills/do-work" <<'PY'
import json
import pathlib
import re
import sys

hooks_file = pathlib.Path(sys.argv[1])
skill_root = pathlib.Path(sys.argv[2])
data = json.loads(hooks_file.read_text())
missing = []

def strings(value):
    if isinstance(value, str):
        yield value
    elif isinstance(value, dict):
        for child in value.values():
            yield from strings(child)
    elif isinstance(value, list):
        for child in value:
            yield from strings(child)

for value in strings(data):
    for command in re.findall(r'\.claude/skills/do-work/([^" ]+)', value):
        target = skill_root / command
        if not target.is_file():
            missing.append(command)
if missing:
    print("unresolved core hook targets: " + ", ".join(sorted(set(missing))), file=sys.stderr)
    raise SystemExit(1)
PY
  then
    fail_count=$((fail_count + 1))
  fi
fi

if [ "$fail_count" -ne 0 ]; then
  exit 1
fi

printf 'staged skills contract: PASS\n'
