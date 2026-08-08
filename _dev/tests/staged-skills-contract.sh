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
  skills/do-work/actions/moved-command-shim.md
  skills/do-work/crew-members/coding-guardrails.md
  skills/do-work/crew-members/clear-questions.md
  skills/do-work/specs/bug-fix.md
  skills/do-work/specs/refactor.md
  skills/do-work/hooks/hooks.json
  skills/do-work/hooks/session-start.sh
  skills/do-work/tools/checks/preflight.sh
  skills/do-work/tools/checks/qualify.sh
  skills/do-work/tools/checks/record-commit-hash.sh
  skills/do-work/tools/do-work-update.sh
  skills/do-work/tools/install-do-work-suite.sh
  skills/do-work/tools/validate-suite-manifest.sh
  skills/do-work/tools/replace-text-section.sh
)

for cutover_export_path in VERSION suite skills; do
  if git -C "$repo_root" check-attr export-ignore -- "$cutover_export_path" \
    | grep -q 'export-ignore: set'; then
    fail "live modular archive must export /$cutover_export_path"
  fi
done

legacy_runtime_paths=(
  SKILL.md
  next-steps.md
  actions
  crew-members
  docs
  hooks
  interviews
  prompts
  specs
  tools/checks
  tools/do-work-update.sh
  tools/queue-kanban
  tools/prime-do-work-update.md
)
for legacy_runtime_path in "${legacy_runtime_paths[@]}"; do
  if [ -f "$repo_root/$legacy_runtime_path" ] \
    || { [ -d "$repo_root/$legacy_runtime_path" ] \
      && find "$repo_root/$legacy_runtime_path" -type f \
        ! -path "$repo_root/tools/queue-kanban/queue-kanban" -print -quit \
        | grep -q .; }; then
    fail "legacy root runtime must be retired at modular cutover: $legacy_runtime_path"
  fi
done

for retained_bootstrap_tool in \
  tools/install-do-work-suite.sh \
  tools/validate-suite-manifest.sh \
  tools/replace-text-section.sh
do
  require_file "$retained_bootstrap_tool"
done

board_files=(
  skills/do-work-board/SKILL.md
  skills/do-work-board/actions/board.md
  skills/do-work-board/actions/help.md
  skills/do-work-board/docs/board-guide.md
  skills/do-work-board/justfile.template
  skills/do-work-board/tools/queue-kanban/go.mod
  skills/do-work-board/tools/queue-kanban/main.go
  skills/do-work-board/tools/queue-kanban/generate_test.go
  skills/do-work-board/tools/queue-kanban/web/board.js
  skills/do-work-board/tools/queue-kanban/web/board.css
)

knowledge_files=(
  skills/do-work-knowledge/SKILL.md
  skills/do-work-knowledge/actions/help.md
  skills/do-work-knowledge/actions/bkb.md
  skills/do-work-knowledge/actions/bkb-reference.md
  skills/do-work-knowledge/actions/dream.md
  skills/do-work-knowledge/actions/memory.md
  skills/do-work-knowledge/actions/memory-reference.md
  skills/do-work-knowledge/actions/memory-value.md
  skills/do-work-knowledge/actions/interview.md
  skills/do-work-knowledge/actions/interview-reference.md
  skills/do-work-knowledge/actions/prompts.md
  skills/do-work-knowledge/actions/setup-memory.md
  skills/do-work-knowledge/hooks/memory-hooks.json
  skills/do-work-knowledge/hooks/memory-session-start.sh
  skills/do-work-knowledge/hooks/memory-stop-capture.sh
  skills/do-work-knowledge/interviews/work-operating-model.md
  skills/do-work-knowledge/prompts/README.md
  skills/do-work-knowledge/crew-members/interviewer.md
)

toolbox_actions=(
  validate-feedback
  code-review
  ui-review
  present-work
  ai-report
  slop-check
  quick-wins
  scan-ideas
  deep-explore
  prime
  inspect
  note
  stray-check
  tidy-repo
  tutorial
  install
)

toolbox_files=(
  skills/do-work-toolbox/SKILL.md
  skills/do-work-toolbox/actions/help.md
  skills/do-work-toolbox/actions/ai-report-reference.md
  skills/do-work-toolbox/actions/deep-explore-reference.md
  skills/do-work-toolbox/docs/code-review-guide.md
  skills/do-work-toolbox/docs/present-work-guide.md
  skills/do-work-toolbox/crew-members/ui-design.md
)

for core_file in "${core_files[@]}"; do
  require_file "$core_file"
done

for retired_pipeline_path in \
  skills/do-work/actions/pipeline.md \
  skills/do-work/actions/pipeline-reference.md \
  skills/do-work/hooks/pipeline-guard.sh
do
  if [ -e "$repo_root/$retired_pipeline_path" ]; then
    fail "stateful pipeline runtime must be retired: $retired_pipeline_path"
  fi
done

if sed -n '/^## Routing/,/^## Dispatch/p' "$repo_root/skills/do-work/SKILL.md" \
  | grep -Eq '^\|[^|]*`(pipeline|full)`'; then
  fail 'core router must not retain the pipeline/full compatibility route'
fi
if grep -Fq 'do-work/pipeline.json' "$repo_root/.gitignore" \
  || grep -Fq 'do-work/pipeline.json' "$repo_root/skills/do-work/hooks/session-start.sh"; then
  fail 'pipeline.json lifecycle must be absent from the live runtime and root ignore file'
fi
if grep -Fq 'pipeline-guard.sh' "$repo_root/skills/do-work/hooks/hooks.json"; then
  fail 'fresh core hooks must not install the retired pipeline Stop guard'
fi

if ! python3 - \
  "$repo_root/README.md" \
  "$repo_root/skills/do-work/actions/help.md" <<'PY'
import pathlib
import sys

approved_prompt = """Use the installed do-work suite to complete this request end to end:

1. Use do-work to capture the request below and record the resulting UR ID.
2. Run do-work verify-requests for that UR. Stop and report if verification fails.
3. Run the UR's REQs through do-work run. Require its built-in tests and review to pass.
4. Use do-work-toolbox present-work for the same UR.
5. Report the implementation, tests, decisions, and deliverable paths.

Request:
<paste request here>"""
missing = [path for path in sys.argv[1:] if approved_prompt not in pathlib.Path(path).read_text()]
if missing:
    raise SystemExit("approved full-cycle prompt is not byte-identical in: " + ", ".join(missing))
PY
then
  fail 'README and core help must carry the approved UR-031 full-cycle prompt byte-for-byte'
fi

if [ -f "$repo_root/skills/do-work/SKILL.md" ]; then
  moved_shim_route_count="$(grep -Ec '^\|.*`\./actions/moved-command-shim\.md`[[:space:]]*\|$' "$repo_root/skills/do-work/SKILL.md" || true)"
  if [ "$moved_shim_route_count" -ne 3 ]; then
    fail "core must route board, knowledge, and toolbox legacy commands through exactly three moved-command rows (found $moved_shim_route_count)"
  fi
fi

if [ -f "$repo_root/skills/do-work/actions/moved-command-shim.md" ]; then
  for replacement_invocation in \
    'do-work-board board' \
    'do-work-knowledge bkb' \
    'do-work-knowledge memory' \
    'do-work-knowledge dream' \
    'do-work-knowledge interview' \
    'do-work-knowledge prompts' \
    'do-work-knowledge setup-memory' \
    'do-work-toolbox <canonical trigger>' \
    'do-work-toolbox install <canonical target>' \
    'just run-kanban'
  do
    grep -Fq "$replacement_invocation" "$repo_root/skills/do-work/actions/moved-command-shim.md" \
      || fail "moved-command shim lacks exact replacement surface: $replacement_invocation"
  done
  grep -Fq 'Do not dispatch another skill' "$repo_root/skills/do-work/actions/moved-command-shim.md" \
    || fail 'moved-command shim must print and stop instead of forwarding'
fi

if [ -f "$repo_root/skills/do-work/tools/install-do-work-suite.sh" ] \
  && ! cmp -s "$repo_root/tools/install-do-work-suite.sh" \
    "$repo_root/skills/do-work/tools/install-do-work-suite.sh"; then
  fail 'staged core installer must be byte-identical to the canonical suite installer'
fi

for board_file in "${board_files[@]}"; do
  require_file "$board_file"
done

for knowledge_file in "${knowledge_files[@]}"; do
  require_file "$knowledge_file"
done

for toolbox_file in "${toolbox_files[@]}"; do
  require_file "$toolbox_file"
done

for toolbox_action in "${toolbox_actions[@]}"; do
  require_file "skills/do-work-toolbox/actions/$toolbox_action.md"
  if [ -f "$repo_root/skills/do-work-toolbox/SKILL.md" ]; then
    toolbox_route_count="$(grep -Fo "\`./actions/$toolbox_action.md\`" "$repo_root/skills/do-work-toolbox/SKILL.md" | wc -l | tr -d ' ' || true)"
    if [ "$toolbox_route_count" -ne 1 ]; then
      fail "toolbox route $toolbox_action must appear exactly once (found $toolbox_route_count)"
    fi
  fi
done

if [ -f "$repo_root/skills/do-work-toolbox/actions/install.md" ] \
  && grep -Eq 'just-kanban|memory-module|do-work-update' "$repo_root/skills/do-work-toolbox/actions/install.md"; then
  fail 'toolbox installer must not own board recipes, memory setup, or core self-update'
fi

if [ -f "$repo_root/skills/do-work/hooks/hooks.json" ] \
  && grep -Fq 'memory-' "$repo_root/skills/do-work/hooks/hooks.json"; then
  fail 'fresh suite core hooks must not enable memory capture'
fi

if [ -f "$repo_root/skills/do-work-knowledge/hooks/memory-hooks.json" ]; then
  if grep -Fq '.claude/skills/do-work/hooks/memory-' "$repo_root/skills/do-work-knowledge/hooks/memory-hooks.json"; then
    fail 'knowledge hook fragment still targets the legacy core hook directory'
  fi
  for knowledge_hook in memory-session-start.sh memory-stop-capture.sh; do
    if ! grep -Fq ".claude/skills/do-work-knowledge/hooks/$knowledge_hook" "$repo_root/skills/do-work-knowledge/hooks/memory-hooks.json"; then
      fail "knowledge hook fragment does not target $knowledge_hook in do-work-knowledge"
    fi
  done
fi

if [ -f "$repo_root/skills/do-work-knowledge/actions/setup-memory.md" ]; then
  for migration_pair in \
    '.claude/skills/do-work/hooks/memory-session-start.sh|.claude/skills/do-work-knowledge/hooks/memory-session-start.sh' \
    '.claude/skills/do-work/hooks/memory-stop-capture.sh|.claude/skills/do-work-knowledge/hooks/memory-stop-capture.sh'
  do
    old_hook="${migration_pair%%|*}"
    new_hook="${migration_pair#*|}"
    grep -Fq "$old_hook" "$repo_root/skills/do-work-knowledge/actions/setup-memory.md" \
      || fail "memory setup lacks legacy migration source: $old_hook"
    grep -Fq "$new_hook" "$repo_root/skills/do-work-knowledge/actions/setup-memory.md" \
      || fail "memory setup lacks modular migration target: $new_hook"
  done
fi

if [ -d "$repo_root/skills/do-work-knowledge" ]; then
  while IFS= read -r knowledge_asset; do
    require_file "skills/do-work-knowledge/$knowledge_asset"
  done < <(git -C "$repo_root" ls-files prompts interviews)
fi

if [ -f "$repo_root/skills/do-work-board/justfile.template" ]; then
  board_template="$repo_root/skills/do-work-board/justfile.template"
  board_recipe_count="$(grep -cF '.claude/skills/do-work-board/tools/queue-kanban' "$board_template" || true)"
  if [ "$board_recipe_count" -ne 4 ]; then
    fail "board Just template must use the do-work-board queue-kanban path in exactly four board recipes (found $board_recipe_count)"
  fi
  if ! grep -Fq '.claude/skills/do-work/tools/do-work-update.sh' "$board_template"; then
    fail 'board Just template must route run-do-work-update to the core updater'
  fi
  if [ "$(grep -cF '# >>> do-work:recipes >>>' "$board_template")" -ne 1 ] \
    || [ "$(grep -cF '# <<< do-work:recipes <<<' "$board_template")" -ne 1 ]; then
    fail 'board Just template must contain exactly one managed do-work:recipes section'
  fi
  for board_recipe_contract in \
    '[ "$wait_count" -lt 320 ]' \
    'remaining_listener_pid=' \
    'listener_executable_name=' \
    'serve --open'
  do
    if ! grep -Fq "$board_recipe_contract" "$board_template"; then
      fail "board Just template lost safety/launch contract: $board_recipe_contract"
    fi
  done
fi

if [ -d "$repo_root/skills/do-work-board/tools/queue-kanban" ]; then
  while IFS= read -r board_source; do
    staged_board_source="skills/do-work-board/tools/queue-kanban/${board_source#tools/queue-kanban/}"
    require_file "$staged_board_source"
  done < <(git -C "$repo_root" ls-files tools/queue-kanban)
fi

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

if [ -f "$repo_root/skills/do-work/actions/work.md" ] \
  && ! grep -Fq -- '--version-file "<skill-root>/actions/version.md"' "$repo_root/skills/do-work/actions/work.md"; then
  fail 'modular core next-version call must name the core version file explicitly instead of using the board tool default'
fi

for staged_router in "$repo_root"/skills/do-work*/SKILL.md; do
  [ -f "$staged_router" ] || continue
  staged_router_root="$(dirname "$staged_router")"
  while IFS= read -r action_path; do
    if [ ! -f "$staged_router_root/${action_path#./}" ]; then
      fail "staged router dispatch target does not resolve: ${staged_router_root##*/}/$action_path"
    fi
  done < <(grep -Eo '`\./actions/[A-Za-z0-9._/-]+\.md`' "$staged_router" | tr -d '`' | sort -u)
done

if [ -d "$repo_root/skills" ]; then
  for staged_skill_root in "$repo_root"/skills/do-work*; do
    [ -d "$staged_skill_root" ] || continue
    if ! python3 - "$staged_skill_root" "$repo_root/suite/modules.tsv" <<'PY'
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
            if reference.startswith(
                (
                    "docs/design/",
                    "docs/handoffs/",
                    "docs/lessons-learned/",
                    "docs/specs/",
                )
            ) or reference == "docs/worklog.md":
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
    sys.stderr.write(f"unresolved staged runtime references in {skill_root.name}:\n")
    sys.stderr.write("\n".join(missing) + "\n")
    raise SystemExit(1)
PY
    then
      fail_count=$((fail_count + 1))
    fi
  done
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
