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

if [ -e "$repo_root/skills/do-work/actions/moved-command-shim.md" ]; then
  fail 'the one-release moved-command shim must be absent after the migration window'
fi
if sed -n '/^## Routing/,/^## Dispatch/p' "$repo_root/skills/do-work/SKILL.md" \
  | grep -Fq 'moved-command-shim.md'; then
  fail 'core routing must not retain moved-command compatibility rows'
fi

sibling_route_contracts=(
  'do-work-board|board'
  'do-work-knowledge|bkb'
  'do-work-knowledge|dream'
  'do-work-knowledge|memory'
  'do-work-knowledge|interview'
  'do-work-knowledge|prompts'
  'do-work-knowledge|setup-memory'
)
for toolbox_action in "${toolbox_actions[@]}"; do
  sibling_route_contracts+=("do-work-toolbox|$toolbox_action")
done

if ! python3 - "$repo_root" "${sibling_route_contracts[@]}" <<'PY'
import pathlib
import re
import sys

repository_root = pathlib.Path(sys.argv[1])
sibling_actions = {contract.split("|", 1)[1] for contract in sys.argv[2:]}
retired_core_forms = {f"do-work {action}" for action in sibling_actions}
retired_core_forms.update(
    {
        "do-work audit codebase",
        "do-work clean up wiki",
        "do-work codebase review",
        "do-work consolidate memory",
        "do-work lint and merge notes",
        "do-work show changes",
        "do-work what changed",
    }
)
retired_pattern = re.compile(
    r"(?<![A-Za-z0-9_-])(?:"
    + "|".join(re.escape(form) for form in sorted(retired_core_forms, key=len, reverse=True))
    + r")(?![A-Za-z0-9_'-])"
)

live_files = [repository_root / "justfile"]
live_files.extend(
    path
    for path in (repository_root / "skills").rglob("*")
    if path.is_file() and path.name not in {"CHANGELOG.md", "queue-kanban"}
)
violations = []
for live_file in live_files:
    for line_number, line in enumerate(live_file.read_text(errors="replace").splitlines(), 1):
        match = retired_pattern.search(line)
        if match:
            violations.append(
                f"{live_file.relative_to(repository_root)}:{line_number}: {match.group(0)}"
            )

prime_file = repository_root / "skills/do-work/tools/prime-do-work-update.md"
transition_fingerprints = (
    "export-ignored through the bridge release",
    "installed bridge validator",
    "unmarked legacy recipe spans",
    "bridge and fresh installs",
    "delete migration branches",
)
prime_text = prime_file.read_text()
for fingerprint in transition_fingerprints:
    if fingerprint in prime_text:
        violations.append(f"{prime_file.relative_to(repository_root)}: {fingerprint}")

if violations:
    raise SystemExit(
        "retired core command or updater-transition restatements remain on live surfaces:\n"
        + "\n".join(violations)
    )
PY
then
  fail 'live shipped surfaces must use sibling-owned commands and permanent updater contracts'
fi

core_routing_section="$(sed -n '/^## Routing/,/^## Dispatch/p' "$repo_root/skills/do-work/SKILL.md")"
for sibling_route_contract in "${sibling_route_contracts[@]}"; do
  IFS='|' read -r sibling_owner public_action <<< "$sibling_route_contract"
  expected_action_path="$repo_root/skills/$sibling_owner/actions/$public_action.md"
  ownership_matches="$(find "$repo_root"/skills/do-work*/actions -maxdepth 1 -type f -name "$public_action.md" -print)"
  ownership_count="$(printf '%s\n' "$ownership_matches" | awk 'NF { count++ } END { print count + 0 }')"
  if [ "$ownership_count" -ne 1 ] || [ "$ownership_matches" != "$expected_action_path" ]; then
    fail "public sibling action $public_action must be owned only by $sibling_owner (found $ownership_count: ${ownership_matches:-none})"
  fi

  sibling_routing_section="$(sed -n '/^## Routing/,/^## /p' "$repo_root/skills/$sibling_owner/SKILL.md")"
  sibling_route_count="$(printf '%s\n' "$sibling_routing_section" | grep -cF "\`./actions/$public_action.md\`" || true)"
  if [ "$sibling_route_count" -ne 1 ]; then
    fail "$sibling_owner must route $public_action exactly once (found $sibling_route_count)"
  fi

  core_route_count="$(printf '%s\n' "$core_routing_section" | grep -cF "\`./actions/$public_action.md\`" || true)"
  if [ "$core_route_count" -ne 0 ]; then
    fail "core must not route sibling-owned action $public_action"
  fi
done

for sibling_owner in do-work-board do-work-knowledge do-work-toolbox; do
  expected_route_count=0
  for sibling_route_contract in "${sibling_route_contracts[@]}"; do
    case "$sibling_route_contract" in
      "$sibling_owner"'|'*) expected_route_count=$((expected_route_count + 1)) ;;
    esac
  done
  sibling_routing_section="$(sed -n '/^## Routing/,/^## /p' "$repo_root/skills/$sibling_owner/SKILL.md")"
  routed_action_paths="$(printf '%s\n' "$sibling_routing_section" \
    | grep -Eo '`\./actions/[A-Za-z0-9._/-]+\.md`' \
    | grep -vF '`./actions/help.md`' || true)"
  actual_route_count="$(printf '%s\n' "$routed_action_paths" | awk 'NF { count++ } END { print count + 0 }')"
  if [ "$actual_route_count" -ne "$expected_route_count" ]; then
    fail "$sibling_owner must expose exactly its $expected_route_count declared public action routes (found $actual_route_count)"
  fi
done

router_behavior_contracts=(
  'do-work-board|Pass `serve`, `static`, `summary`, `cli`, `--port N`, and `--out DIR` through to the board action.|An unknown command prints board help and stops.'
  'do-work-knowledge|Pass the complete remainder through.|An unknown command prints help and stops.'
  'do-work-toolbox|Pass all remaining arguments through.|Unknown single words print help;'
)
for router_behavior_contract in "${router_behavior_contracts[@]}"; do
  IFS='|' read -r sibling_owner pass_through_contract unknown_help_contract <<< "$router_behavior_contract"
  sibling_router="$repo_root/skills/$sibling_owner/SKILL.md"
  if [ "$(grep -cF "$pass_through_contract" "$sibling_router" || true)" -ne 1 ]; then
    fail "$sibling_owner must state its argument pass-through contract exactly once"
  fi
  if [ "$(grep -cF "$unknown_help_contract" "$sibling_router" || true)" -ne 1 ]; then
    fail "$sibling_owner must state its unknown-command help contract exactly once"
  fi
done

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
  if grep -Fq '.claude/skills/do-work/hooks/memory-' "$repo_root/skills/do-work-knowledge/actions/setup-memory.md"; then
    fail 'memory setup must not retain old core-path hook migration instructions'
  fi
  for knowledge_hook in memory-session-start.sh memory-stop-capture.sh; do
    grep -Fq ".claude/skills/do-work-knowledge/hooks/$knowledge_hook" \
      "$repo_root/skills/do-work-knowledge/actions/setup-memory.md" \
      || fail "memory setup must target the current modular hook: $knowledge_hook"
  done
fi

if [ -d "$repo_root/skills/do-work-knowledge" ]; then
  while IFS= read -r knowledge_asset; do
    require_file "skills/do-work-knowledge/$knowledge_asset"
  done < <(git -C "$repo_root" ls-files prompts interviews)
fi

if [ -f "$repo_root/skills/do-work-board/justfile.template" ]; then
  board_template="${STAGED_SKILLS_BOARD_TEMPLATE:-$repo_root/skills/do-work-board/justfile.template}"
  board_recipe_count="$(grep -cF '.claude/skills/do-work-board/tools/queue-kanban' "$board_template" || true)"
  if [ "$board_recipe_count" -ne 4 ]; then
    fail "board Just template must use the do-work-board queue-kanban path in exactly four board recipes (found $board_recipe_count)"
  fi
  if ! python3 - "$board_template" <<'PY'
import pathlib
import re
import sys

template_file = pathlib.Path(sys.argv[1])
template_text = template_file.read_text()
template_lines = template_text.splitlines()
recipe_header = re.compile(r"^run-do-work-update(?:\s+.*)?:\s*$")
next_recipe_header = re.compile(r"^[A-Za-z0-9][A-Za-z0-9_-]*(?:\s+.*)?:\s*$")

header_indexes = [
    index for index, line in enumerate(template_lines) if recipe_header.match(line)
]
if len(header_indexes) != 1 or template_text.count("run-do-work-update") != 1:
    raise SystemExit("run-do-work-update must appear as exactly one recipe")

recipe_body_lines = []
for line in template_lines[header_indexes[0] + 1 :]:
    if line == "# <<< do-work:recipes <<<" or next_recipe_header.match(line):
        break
    recipe_body_lines.append(line)
recipe_body = "\n".join(recipe_body_lines)

core_updater_path = ".claude/skills/do-work/tools/do-work-update.sh"
core_updater_invocation = re.compile(
    r'\bbash\s+"\$project_root/' + re.escape(core_updater_path) + r'"'
)
if recipe_body.count(core_updater_path) != 1 or len(core_updater_invocation.findall(recipe_body)) != 1:
    raise SystemExit("run-do-work-update must invoke the core updater exactly once")

for forbidden_board_call in ("run-kanban", "do-work-board", "queue-kanban"):
    if forbidden_board_call in recipe_body:
        raise SystemExit(
            "run-do-work-update must not invoke a board recipe or binary: "
            + forbidden_board_call
        )
PY
  then
    fail 'board Just template has an invalid isolated run-do-work-update recipe'
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
