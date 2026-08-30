#!/usr/bin/env bash
# shellcheck disable=SC2016
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fail_count=0
core_root="$repo_root/skills/do-work"
board_root="$repo_root/skills/do-work-board"
knowledge_root="$repo_root/skills/do-work-knowledge"
toolbox_root="$repo_root/skills/do-work-toolbox"

resolve_runtime_file() {
  local relative_path="$1"

  case "$relative_path" in
    SKILL.md|next-steps.md) printf '%s/%s\n' "$core_root" "$relative_path" ;;
    actions/board.md|docs/board-guide.md) printf '%s/%s\n' "$board_root" "$relative_path" ;;
    actions/bkb*|actions/dream.md|actions/interview*|actions/memory*|actions/prompts.md|actions/setup-memory.md|docs/bkb-guide.md|docs/dream-guide.md|docs/interview-guide.md|docs/prompts-guide.md|prompts/*|interviews/*|hooks/memory-*)
      printf '%s/%s\n' "$knowledge_root" "$relative_path" ;;
    actions/ai-report*|actions/architecture-report.md|actions/code-review.md|actions/completed-work-presentation-reference.md|actions/deep-explore*|actions/inspect.md|actions/install.md|actions/note.md|actions/present-video.md|actions/present-work.md|actions/prime.md|actions/quick-wins.md|actions/scan-ideas.md|actions/slop-check.md|actions/stray-check.md|actions/tidy-repo.md|actions/tutorial.md|actions/ui-review.md|actions/validate-feedback.md|docs/ai-report-guide.md|docs/code-review-guide.md|docs/inspect-guide.md|docs/present-video-guide.md|docs/present-work-guide.md|docs/prime-guide.md|docs/quick-wins-guide.md|docs/slop-check-guide.md|docs/stray-check-guide.md|docs/ui-review-guide.md)
      printf '%s/%s\n' "$toolbox_root" "$relative_path" ;;
    actions/*|crew-members/*|docs/*|hooks/*|scripts/*|specs/*|tools/checks/*|tools/estimate-p50.sh|tools/do-work-update.sh|tools/prime-do-work-update.md)
      printf '%s/%s\n' "$core_root" "$relative_path" ;;
    tools/queue-kanban/*) printf '%s/%s\n' "$board_root" "$relative_path" ;;
    *) printf '%s/%s\n' "$repo_root" "$relative_path" ;;
  esac
}

for cutover_export_path in VERSION suite skills; do
  if git -C "$repo_root" check-attr export-ignore -- "$cutover_export_path" \
    | grep -q 'export-ignore: set'; then
    printf 'FAIL: live modular archive still excludes /%s.\n' "$cutover_export_path" >&2
    fail_count=$((fail_count + 1))
  fi
done

for retired_runtime_path in SKILL.md next-steps.md actions crew-members docs hooks interviews prompts specs tools/checks tools/do-work-update.sh tools/queue-kanban tools/prime-do-work-update.md; do
  if [ -f "$repo_root/$retired_runtime_path" ] \
    || { [ -d "$repo_root/$retired_runtime_path" ] \
      && find "$repo_root/$retired_runtime_path" -type f \
        ! -path "$repo_root/tools/queue-kanban/queue-kanban" -print -quit \
        | grep -q .; }; then
    printf 'FAIL: legacy root runtime still exists after modular cutover: %s.\n' "$retired_runtime_path" >&2
    fail_count=$((fail_count + 1))
  fi
done

assert_contains() {
  local file_path="$1"
  local pattern_text="$2"
  local message_text="$3"

  if ! grep -Eq -- "$pattern_text" "$(resolve_runtime_file "$file_path")"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_block_contains() {
  local block_text="$1"
  local pattern_text="$2"
  local message_text="$3"

  if ! grep -Eq "$pattern_text" <<<"$block_text"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_block_not_contains() {
  local block_text="$1"
  local pattern_text="$2"
  local message_text="$3"

  if grep -Eq "$pattern_text" <<<"$block_text"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_file_missing() {
  local file_path="$1"
  local message_text="$2"

  if [ -e "$repo_root/$file_path" ]; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_file_not_contains() {
  local file_path="$1"
  local pattern_text="$2"
  local message_text="$3"

  local resolved_file
  resolved_file="$(resolve_runtime_file "$file_path")"
  if grep -Eiq -- "$pattern_text" "$resolved_file"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    grep -Ein -- "$pattern_text" "$resolved_file" >&2 || true
    fail_count=$((fail_count + 1))
  fi
}

skill_dispatch_block="$(sed -n '/^## Routing/,/^## Dispatch/p' "$core_root/SKILL.md")"
review_archived_input_block="$(sed -n '/^### Step 3: Read the Original Input/,/^### Step 4:/p' "$core_root/actions/review-work.md")"
work_archive_success_block="$(sed -n '/^### Step 8: Archive/,/^\*\*On failure:/p' "$core_root/actions/work.md")"

# Prime creation emits a linked two-file pair, and audit's otherwise read-only routing
# boundary has one narrow exception for shrink-time lesson promotion. Validate the
# operative heading blocks rather than letting nearby templates satisfy the directives.
if ! python3 - "$toolbox_root/actions/prime.md" <<'PY'
import pathlib
import re
import sys

prime_action_text = pathlib.Path(sys.argv[1]).read_text()

def heading_block(source, start, end):
    return source.split(start, 1)[1].split(end, 1)[0]

def without_fences(block):
    return re.sub(r"```.*?```", "", block, flags=re.DOTALL)

def validate(source):
    try:
        create = heading_block(source, "#### Step 5: Write", "#### Step 6: Post-creation checks")
        report_block = heading_block(source, "#### Report", "## Sub-Command: `audit`")
        audit_intro = heading_block(source, "## Sub-Command: `audit`", "### Conventions")
        stakes = heading_block(source, "### Step 6.5: Refresh Stakes", "### Step 6.6: Shrink")
        shrink = heading_block(source, "### Step 6.6: Shrink", "### Output Format")
    except IndexError:
        return ["prime create/audit heading boundary is missing"]

    active_create = without_fences(create)
    errors = []
    if "Create both files in the same operation" not in active_create:
        errors.append("prime create does not require the linked two-file pair")
    for output_path in (
        "{path}/prime-{short-name}.md",
        "{path}/lessons-{short-name}.md",
    ):
        if output_path not in active_create:
            errors.append(f"prime create has no active write directive for {output_path}")
        if output_path not in report_block:
            errors.append(f"prime create report omits {output_path}")
    for token in ("# Lessons: {short-name}", "[`prime-{short-name}.md`](prime-{short-name}.md)"):
        if token not in create:
            errors.append(f"prime create satellite template omits {token}")

    audit_contracts = (
        (audit_intro, ("`Traps`", "Step 6.6", "sole routing-content exception")),
        (stakes, ("`Traps`", "Step 6.6", "existing-lesson promotion exception")),
        (shrink, ("`## Traps`", "existing utility-wide lesson", "sole routing-content write exception", "Promote, don't duplicate")),
    )
    for block, tokens in audit_contracts:
        if not all(token in block for token in tokens):
            errors.append("an operative audit block omits the shrink-only Traps exception")
    if re.search(r"`## Stakes` is the only section audit writes", source):
        errors.append("audit must not retain the unqualified Stakes-only write boundary")
    return errors

baseline_errors = validate(prime_action_text)
if baseline_errors:
    raise SystemExit("; ".join(baseline_errors))

mutants = (
    ("missing satellite-create directive", prime_action_text.replace("Create both files in the same operation", "Create the prime file", 1)),
    ("unqualified Stakes-only boundary", prime_action_text.replace(
        "Outside Step 6.6's shrink operation, `## Stakes` is the only prime section audit authors from current source.",
        "`## Stakes` is the only section audit writes — the routing sections stay read-only.",
        1,
    )),
)
for mutant_name, mutant_text in mutants:
    if mutant_text == prime_action_text or not validate(mutant_text):
        raise SystemExit(f"prime workflow contract accepted {mutant_name}")
PY
then
  printf 'FAIL: prime create/audit workflow boundaries regressed.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Every phrase the version action documents must be reachable from the always-loaded
# router. Derive the phrase set from the action instead of maintaining a second alias
# inventory, then check the one precedence edge that matters: the exact update-check
# phrase must win before generic request verification's `check` alias.
if ! python3 - "$core_root/SKILL.md" "$core_root/actions/version.md" <<'PY'
import pathlib
import re
import sys

router_file = pathlib.Path(sys.argv[1])
version_action_file = pathlib.Path(sys.argv[2])
router_text = router_file.read_text()
version_action_text = version_action_file.read_text()

routing_block = router_text.split("## Routing", 1)[1].split("## Dispatch", 1)[0]
route_rows = []
for line_number, line in enumerate(routing_block.splitlines(), start=1):
    if not line.startswith("|") or "`./actions/" not in line:
        continue
    cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
    if len(cells) != 2:
        continue
    aliases = re.findall(r"`([^`]+)`", cells[0])
    route_rows.append((line_number, aliases, cells[1]))

documented_phrases = []
for mode_name in ("Version request", "Update check", "Recap"):
    mode_match = re.search(
        rf"^- \*\*{re.escape(mode_name)}\*\* — (.+)$",
        version_action_text,
        flags=re.MULTILINE,
    )
    if mode_match is None:
        raise SystemExit(f"version action has no {mode_name} phrase declaration")
    documented_phrases.extend(re.findall(r'"([^"]+)"', mode_match.group(1)))

alias_routes = {}
for row_number, aliases, route in route_rows:
    for alias in aliases:
        alias_routes.setdefault(alias, []).append((row_number, route))

for phrase in dict.fromkeys(documented_phrases):
    matches = alias_routes.get(phrase, [])
    if len(matches) != 1 or "`./actions/version.md`" not in matches[0][1]:
        raise SystemExit(
            f"documented version phrase {phrase!r} must route exactly once to actions/version.md; "
            f"found {matches or 'no route'}"
        )

update_check_row = alias_routes["check for updates"][0][0]
generic_check_matches = alias_routes.get("check", [])
if len(generic_check_matches) != 1:
    raise SystemExit(f"generic check alias must route exactly once; found {generic_check_matches}")
if update_check_row >= generic_check_matches[0][0]:
    raise SystemExit("check for updates must route before the generic check alias")
PY
then
  printf 'FAIL: core version/update phrases drifted from first-match routing.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Public work aliases have one enumerated home in the work guide; the router is
# the executable mirror. The testing-status schema and Go normalizer likewise
# own two representations of one alias map. Compare both seams exactly, then
# mutate either side in memory to prove additions and removals cannot hide.
if ! python3 - \
  "$repo_root/README.md" \
  "$core_root/docs/work-guide.md" \
  "$core_root/SKILL.md" \
  "$core_root/actions/work-reference.md" \
  "$board_root/tools/queue-kanban/testing.go" <<'PY'
import pathlib
import re
import sys


def require_exact(left, right, seam):
    if left == right:
        return
    left_keys = set(left)
    right_keys = set(right)
    details = []
    if left_keys - right_keys:
        details.append(f"left-only={sorted(left_keys - right_keys)!r}")
    if right_keys - left_keys:
        details.append(f"right-only={sorted(right_keys - left_keys)!r}")
    if isinstance(left, dict):
        remapped = sorted(
            key for key in left_keys & right_keys if left[key] != right[key]
        )
        if remapped:
            details.append(
                "remapped="
                + repr([(key, left[key], right[key]) for key in remapped])
            )
    raise AssertionError(f"{seam} drifted: {', '.join(details)}")


def prove_one_sided_mutations(left, right, seam, added_value=None):
    require_exact(left, right, seam)
    existing_key = sorted(left)[0]
    added_key = "contract-only-mutation"

    if isinstance(left, dict):
        additions = (
            ({**left, added_key: added_value}, right),
            (left, {**right, added_key: added_value}),
        )
        removals = (
            ({key: value for key, value in left.items() if key != existing_key}, right),
            (left, {key: value for key, value in right.items() if key != existing_key}),
        )
    else:
        additions = ((left | {added_key}, right), (left, right | {added_key}))
        removals = ((left - {existing_key}, right), (left, right - {existing_key}))

    for mutation_number, mutated_pair in enumerate(additions + removals, start=1):
        try:
            require_exact(*mutated_pair, f"{seam} mutation {mutation_number}")
        except AssertionError:
            continue
        raise AssertionError(
            f"{seam} comparison accepted one-sided mutation {mutation_number}"
        )


readme_text = pathlib.Path(sys.argv[1]).read_text()
guide_text = pathlib.Path(sys.argv[2]).read_text()
router_text = pathlib.Path(sys.argv[3]).read_text()
schema_text = pathlib.Path(sys.argv[4]).read_text()
normalizer_text = pathlib.Path(sys.argv[5]).read_text()
seam_failures = []

if "skills/do-work/docs/work-guide.md#trigger-aliases" not in readme_text:
    raise AssertionError("README must point to the work guide's canonical alias list")
if "Other trigger words:" in readme_text:
    raise AssertionError("README must not carry a second public work-alias inventory")

guide_alias_block = guide_text.split("## Trigger aliases", 1)[1].split("## Tips", 1)[0]
guide_aliases = set(
    re.findall(r"^do-work ([a-z][a-z-]*)$", guide_alias_block, flags=re.MULTILINE)
)

router_block = router_text.split("## Routing", 1)[1].split("## Dispatch", 1)[0]
work_route_rows = []
for line in router_block.splitlines():
    if not line.startswith("|") or "`./actions/work.md`" not in line:
        continue
    cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
    if len(cells) == 2:
        work_route_rows.append(set(re.findall(r"`([^`]+)`", cells[0])))
if len(work_route_rows) != 1:
    raise AssertionError(
        f"work router must have exactly one actions/work.md row; found {len(work_route_rows)}"
    )
router_aliases = work_route_rows[0]
try:
    prove_one_sided_mutations(guide_aliases, router_aliases, "work guide/router aliases")
except AssertionError as error:
    seam_failures.append(str(error))

testing_schema_row = next(
    (
        line
        for line in schema_text.splitlines()
        if line.startswith("|") and line.split("|", 2)[1].strip().startswith("`testing_status`")
    ),
    None,
)
if testing_schema_row is None:
    raise AssertionError("schema table has no testing_status row")
schema_cells = [cell.strip() for cell in testing_schema_row.strip().strip("|").split("|")]
schema_alias_map = {}
for mapping_clause in schema_cells[2].split(";"):
    target_match = re.search(r"→ `([^`]+)`", mapping_clause)
    if target_match is None:
        continue
    for alias in re.findall(r"`([^`]+)`", mapping_clause[: target_match.start()]):
        if alias in schema_alias_map:
            raise AssertionError(f"testing_status schema repeats alias {alias!r}")
        schema_alias_map[alias] = target_match.group(1)

status_constants = dict(
    re.findall(r'^\s*(testingStatus\w+)\s*=\s*"([^"]+)"', normalizer_text, flags=re.MULTILINE)
)
normalizer_function = re.search(
    r"func normalizeTestingStatus\([^)]*\) string \{(.*?)\n\}",
    normalizer_text,
    flags=re.DOTALL,
)
if normalizer_function is None:
    raise AssertionError("testing.go has no normalizeTestingStatus function")
normalizer_alias_map = {}
for aliases_text, target_constant in re.findall(
    r"case ([^:]+):\s*return (testingStatus\w+)", normalizer_function.group(1)
):
    if target_constant not in status_constants:
        raise AssertionError(f"normalizer returns unknown constant {target_constant}")
    for alias in re.findall(r'"([^"]+)"', aliases_text):
        if alias in normalizer_alias_map:
            raise AssertionError(f"normalizeTestingStatus repeats alias {alias!r}")
        normalizer_alias_map[alias] = status_constants[target_constant]

try:
    prove_one_sided_mutations(
        schema_alias_map,
        normalizer_alias_map,
        "testing_status schema/normalizer aliases",
        added_value="in-testing",
    )
except AssertionError as error:
    seam_failures.append(str(error))

if seam_failures:
    raise AssertionError("; ".join(seam_failures))
PY
then
  printf 'FAIL: public work or testing-status alias vocabularies drifted.\n' >&2
  fail_count=$((fail_count + 1))
fi

for retired_pipeline_path in \
  skills/do-work/actions/pipeline.md \
  skills/do-work/actions/pipeline-reference.md \
  skills/do-work/hooks/pipeline-guard.sh
do
  assert_file_missing \
    "$retired_pipeline_path" \
    "stateful pipeline runtime must remain retired: $retired_pipeline_path"
done

assert_block_not_contains \
  "$skill_dispatch_block" \
  '^\|[^|]*`(pipeline|full)`' \
  'Core SKILL.md must not retain a pipeline/full route or compatibility alias.'

assert_file_missing \
  "skills/do-work/actions/moved-command-shim.md" \
  'The one-release moved-command shim must be deleted after client migration.'

assert_block_not_contains \
  "$skill_dispatch_block" \
  'moved-command-shim\.md' \
  'Core routing must not retain compatibility rows for sibling-owned commands.'

assert_file_not_contains \
  "tools/do-work-update.sh" \
  'suite-layout-v2|--capabilities|legacy_shipped_paths|legacy all-in-one skill' \
  'The current updater launcher must not retain bridge capability, monolith, or stale-copy branches.'
assert_file_not_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go" \
  'suite-layout-v2|--capabilities|legacy_shipped_paths|legacy all-in-one skill' \
  'The update transaction must not retain bridge capability, monolith, or stale-copy branches.'

assert_file_not_contains \
  "tools/install-do-work-suite.sh" \
  '--migrate-legacy-do-work|\.claude/skills/do-work/hooks/memory-' \
  'The suite installer launcher must not retain exact recipe or old core memory-hook migrations.'
assert_file_not_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go" \
  '--migrate-legacy-do-work|\.claude/skills/do-work/hooks/memory-' \
  'The install transaction must not retain exact recipe or old core memory-hook migrations.'

assert_file_not_contains \
  "actions/setup-memory.md" \
  '\.claude/skills/do-work/hooks/memory-' \
  'Knowledge memory setup must describe only the current modular hook paths.'

assert_file_not_contains \
  ".gitignore" \
  'do-work/pipeline\.json' \
  'Root .gitignore must not recreate the retired pipeline state-file lifecycle.'

assert_file_not_contains \
  "hooks/session-start.sh" \
  'pipeline\.json|Pipeline active|do-work pipeline' \
  'SessionStart must report queue state without reading or advertising the retired pipeline state machine.'

assert_file_not_contains \
  "hooks/hooks.json" \
  'pipeline-guard\.sh' \
  'Fresh core hook settings must not install the retired pipeline Stop guard.'

if ! python3 - \
  "$repo_root/README.md" \
  "$core_root/actions/help.md" <<'PY'
import pathlib
import sys

approved_prompt = """Use the installed do-work suite to complete this request end to end:

1. Use do-work to capture the request below and record the resulting UR ID.
2. Run do-work verify-requests for that UR. Stop and report if verification fails.
3. Run the UR's REQs through do-work run. Require its built-in tests and review to pass.
4. Use do-work-toolbox ai-report for the same UR.
5. Report the implementation, tests, decisions, and deliverable paths.

Request:
<paste request here>"""
missing = [path for path in sys.argv[1:] if approved_prompt not in pathlib.Path(path).read_text()]
if missing:
    raise SystemExit("approved full-cycle prompt is not byte-identical in: " + ", ".join(missing))
PY
then
  printf 'FAIL: README and core help must carry the approved UR-031 full-cycle prompt byte-for-byte.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Presentation routing has three non-overlapping owners. Keep the public router,
# discovery surfaces, caller guidance, shared mechanics, and the three action
# contracts aligned without snapshotting decorative prose (REQ-192).
if ! python3 - "$repo_root" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
failures = []

def text(relative_path):
    return (root / relative_path).read_text()

def require(relative_path, pattern, message, flags=re.IGNORECASE | re.MULTILINE):
    if not re.search(pattern, text(relative_path), flags):
        failures.append(f"{relative_path}: {message}")

def reject(relative_path, pattern, message, flags=re.IGNORECASE | re.MULTILINE):
    if re.search(pattern, text(relative_path), flags):
        failures.append(f"{relative_path}: {message}")

def require_order(relative_path, earlier_pattern, later_pattern, message):
    source = text(relative_path)
    earlier = re.search(earlier_pattern, source, re.IGNORECASE | re.MULTILINE)
    later = re.search(later_pattern, source, re.IGNORECASE | re.MULTILINE)
    if not earlier or not later:
        missing = []
        if not earlier:
            missing.append(f"earlier predicate /{earlier_pattern}/")
        if not later:
            missing.append(f"later predicate /{later_pattern}/")
        failures.append(f"{relative_path}: {message}; missing " + " and ".join(missing))
    elif earlier.start() >= later.start():
        failures.append(f"{relative_path}: {message}; required load appears after the guarded read")

def executable_markdown_segments(source):
    """Return fenced/inline/command-like text, excluding negative prose outside fences."""
    segments = []
    in_fence = False
    fence_marker = ""
    fence_start = 0
    fence_lines = []
    prohibition_example_context = False
    prohibition = re.compile(r"\b(?:do not|don't|never|must not|does not|without|forbid|reject|avoid)\b", re.IGNORECASE)
    prohibition_example = re.compile(r"^\s*(?:[-*+]|\d+[.)])\s+")
    prohibition_continuation = re.compile(r"^[ \t]+.*`[^`\n]+`")
    command_signal = re.compile(
        r"--no-open|\bremotion\b|\bnpm\s+run\b|\bsleep\b|"
        r"https?://(?:localhost|127\.0\.0\.1):\d+|\bopen\s+https?://|[\"']render[\"']\s*:",
        re.IGNORECASE,
    )
    platform_opener_signal = re.compile(r"(?:^[ \t]*|[;&|][ \t]*)open(?:[ \t]+|$)")

    for line_number, line in enumerate(source.splitlines(), start=1):
        fence = re.match(r"^\s*(```|~~~)", line)
        if fence:
            marker = fence.group(1)
            if not in_fence:
                in_fence = True
                fence_marker = marker
                fence_start = line_number + 1
                fence_lines = []
                prohibition_example_context = False
            elif marker == fence_marker:
                segments.append((f"fenced code lines {fence_start}-{line_number - 1}", "\n".join(fence_lines)))
                in_fence = False
                fence_marker = ""
                fence_lines = []
            continue
        if in_fence:
            fence_lines.append(line)
            continue
        if prohibition.search(line):
            prohibition_example_context = True
            continue
        if prohibition_example_context:
            if not line.strip() or prohibition_example.match(line) or prohibition_continuation.match(line):
                continue
            prohibition_example_context = False
        for inline_number, inline in enumerate(re.findall(r"`([^`\n]+)`", line), start=1):
            segments.append((f"line {line_number} inline code {inline_number}", inline))
        if command_signal.search(line) or platform_opener_signal.search(line):
            segments.append((f"line {line_number} command-like prose", line))

    return segments

def unsafe_executable_video_findings(source):
    unsafe_forms = (
        ("--no-open option", r"(?<![\w-])--no-open(?![\w-])"),
        ("backgrounded Studio/preview command", r"(?:\bremotion\s+studio\b|\bnpm\s+run\s+preview\b)[^\n]*(?<!&)&(?!&)"),
        ("background operator chained to sleep", r"(?<!&)&(?!&)[ \t]*(?:;|\n)?[ \t]*sleep\b"),
        ("numeric readiness sleep", r"(?m)(?:^|[;&|]\s*)sleep\s+\d+(?:\.\d+)?s?\b"),
        (
            "fixed-port Studio command",
            r"\b(?:npx\s+)?remotion\s+studio\b(?:[^\n]*\\\n[ \t]*)*[^\n]*(?<![\w-])--port(?:[ \t]+|=)(?:[\"']\d+[\"']|\d+)(?![\w-])",
        ),
        ("fixed localhost preview URL", r"https?://(?:localhost|127\.0\.0\.1):\d+"),
        ("platform shell opener", r"(?m)(?:^[ \t]*|[;&|][ \t]*)(?-i:open)(?:[ \t]+|$)"),
        ("Remotion render command", r"\b(?:npx\s+)?remotion\s+render\b"),
        ("package render script", r"[\"']render[\"']\s*:"),
    )
    findings = []
    for location, segment in executable_markdown_segments(source):
        for label, pattern in unsafe_forms:
            if re.search(pattern, segment, re.IGNORECASE):
                findings.append((label, location, segment))
    return findings

def reject_executable_video_forms(relative_path):
    for label, location, segment in unsafe_executable_video_findings(text(relative_path)):
        failures.append(f"{relative_path}: executable {label} remains in {location}: {segment!r}")

unsafe_video_mutations = (
    (
        "separated numeric Studio port",
        "```bash\nremotion studio src/Root.tsx --port 3000\n```",
        "fixed-port Studio command",
    ),
    (
        "equals numeric Studio port",
        "```bash\nremotion studio src/Root.tsx --port=3000\n```",
        "fixed-port Studio command",
    ),
    (
        "quoted numeric Studio port",
        "```bash\nremotion studio src/Root.tsx --port \"4311\"\n```",
        "fixed-port Studio command",
    ),
    (
        "equals quoted numeric Studio port",
        "```bash\nremotion studio src/Root.tsx --port='4312'\n```",
        "fixed-port Studio command",
    ),
    (
        "variable platform opener",
        "```bash\nopen \"$REMOTION_PREVIEW_URL\"\n```",
        "platform shell opener",
    ),
    (
        "chained parameter-expansion platform opener",
        "```bash\nnpm run preview && open \"${REMOTION_PREVIEW_URL}\"\n```",
        "platform shell opener",
    ),
    (
        "two-space-indented platform opener",
        "```bash\n  open \"$REMOTION_PREVIEW_URL\"\n```",
        "platform shell opener",
    ),
    (
        "tab-indented platform opener",
        "```bash\n\topen \"${REMOTION_PREVIEW_URL}\"\n```",
        "platform shell opener",
    ),
    (
        "continued separated numeric Studio port",
        "```bash\nremotion studio src/Root.tsx \\\n  --port 4511\n```",
        "fixed-port Studio command",
    ),
    (
        "continued equals numeric Studio port",
        "```bash\nremotion studio src/Root.tsx \\\n  --port=4512\n```",
        "fixed-port Studio command",
    ),
)
for mutation_name, mutation_source, expected_family in unsafe_video_mutations:
    mutation_families = {
        label for label, _location, _segment in unsafe_executable_video_findings(mutation_source)
    }
    if expected_family not in mutation_families:
        failures.append(
            f"unsafe video detector mutation {mutation_name!r} escaped expected family "
            f"{expected_family!r}; found {sorted(mutation_families)!r}"
        )

safe_video_mutations = (
    ("foreground Studio", "```bash\nremotion studio src/Root.tsx\n```"),
    ("package preview", "```bash\nnpm run preview\n```"),
    (
        "negative prohibition prose",
        "Do not add `remotion studio src/Root.tsx --port 3000` or run `open \"$url\"`.",
    ),
    (
        "multiline prohibition examples",
        "Do not use these preview shortcuts:\n\n"
        "- `remotion studio src/Root.tsx --port 3000`\n"
        "- `open \"$REMOTION_PREVIEW_URL\"`",
    ),
    ("ordinary opener prose", "Open the reported preview path in a browser if useful."),
)
for mutation_name, mutation_source in safe_video_mutations:
    mutation_findings = unsafe_executable_video_findings(mutation_source)
    if mutation_findings:
        failures.append(
            f"unsafe video detector rejected safe mutation {mutation_name!r}: "
            f"{mutation_findings!r}"
        )

toolbox_skill = text("skills/do-work-toolbox/SKILL.md")
routing = toolbox_skill.split("## Routing", 1)[1].split("## Dispatch", 1)[0]
route_rows = {}
for line in routing.splitlines():
    match = re.match(r"^\|\s*(.*?)\s*\|\s*`(\./actions/[^`]+)`\s*\|$", line)
    if not match:
        continue
    triggers = tuple(re.findall(r"`([^`]+)`", match.group(1)))
    route_rows[match.group(2)] = triggers

expected_routes = {
    "./actions/ai-report.md": ("ai-report", "showcase", "visual report", "proof of work"),
    "./actions/present-work.md": ("present-work", "portfolio", "work portfolio"),
    "./actions/present-video.md": ("present-video", "remotion", "video walkthrough"),
}
for route, triggers in expected_routes.items():
    if route_rows.get(route) != triggers:
        failures.append(f"skills/do-work-toolbox/SKILL.md: {route} owns {route_rows.get(route)!r}, expected {triggers!r}")

retired_triggers = {
    trigger.casefold()
    for triggers in route_rows.values()
    for trigger in triggers
    if trigger.casefold() in {"present", "client brief"}
}
if retired_triggers:
    failures.append("skills/do-work-toolbox/SKILL.md: broad presentation triggers remain: " + ", ".join(sorted(retired_triggers)))

hint = re.search(r'^argument-hint:\s*"([^"]+)"', toolbox_skill, re.MULTILINE)
hint_commands = {part.strip() for part in hint.group(1).split("|")} if hint else set()
for command in ("ai-report", "present-work", "present-video"):
    if command not in hint_commands:
        failures.append(f"skills/do-work-toolbox/SKILL.md: argument-hint omits {command}")

public_contracts = {
    "skills/do-work-toolbox/actions/help.md": (
        r"ai-report\s+\[REQ\|UR\].*detailed.*stakeholder.*HTML",
        r"present-work\s+(?:all\|portfolio|\[all\|portfolio\]).*portfolio",
        r"present-video\s+\[REQ\|UR\].*(?:source-only|Remotion).*walkthrough",
    ),
    "README.md": (
        r"do-work-toolbox ai-report for the same UR",
        r"do-work-toolbox ai-report",
        r"do-work-toolbox present-work (?:all|portfolio)",
        r"do-work-toolbox present-video",
    ),
    "skills/do-work/actions/help.md": (
        r"do-work-toolbox[\s\S]{0,350}present-video",
        r"do-work-toolbox ai-report for the same UR",
    ),
    "skills/do-work-toolbox/actions/tutorial.md": (
        r"full cycle[\s\S]{0,500}do-work-toolbox ai-report UR-NNN",
        r"Work is done.*[\s\S]{0,250}do-work-toolbox ai-report",
        r"portfolio[\s\S]{0,180}do-work-toolbox present-work (?:all|portfolio)",
        r"video walkthrough[\s\S]{0,180}do-work-toolbox present-video (?:REQ|UR)-NNN",
    ),
}
for relative_path, patterns in public_contracts.items():
    for pattern in patterns:
        require(relative_path, pattern, f"missing presentation discovery predicate /{pattern}/")

for relative_path in (
    "skills/do-work/crew-members/prompt-injection.md",
    "skills/do-work-knowledge/crew-members/prompt-injection.md",
    "skills/do-work-toolbox/crew-members/prompt-injection.md",
):
    require(relative_path, r"condition.*(?:contract|boundary)", "JIT caller rule must lead with the ingestion condition")
    require(relative_path, r"illustrative|not exhaustive", "JIT caller examples must be explicitly illustrative")
    require(relative_path, r"completed-work-presentation-reference\.md", "completed-work presentation readers must point to the shared safety reader")

for relative_path in (
    "skills/do-work/crew-members/anti-slop.md",
    "skills/do-work-knowledge/crew-members/anti-slop.md",
    "skills/do-work-toolbox/crew-members/anti-slop.md",
):
    require(relative_path, r"condition.*contract", "human-facing artifact condition must remain the JIT trigger")
    require(relative_path, r"illustrative", "anti-slop callers must remain illustrative")
    for action in ("ai-report", "present-work", "present-video"):
        require(relative_path, rf"\b{action}\b", f"anti-slop examples omit {action}")
    reject(relative_path, r"client briefs, video scripts, and HTML explainers in present-work", "retired presentation artifact family remains")

completed_work_reference = "skills/do-work-toolbox/actions/completed-work-presentation-reference.md"
target_id_reference = "skills/do-work/actions/work-reference.md"
target_id_match = re.search(
    r"^### Target ID Resolution\s*$[\s\S]*?(?=^## |\Z)",
    text(target_id_reference),
    re.MULTILINE,
)
target_id_contract = target_id_match.group(0) if target_id_match else ""
if not target_id_match:
    failures.append(f"{target_id_reference}: canonical Target ID Resolution section is missing")
for predicate in (
    r"`REQ-` \+ digits and `UR-` \+ digits, \*\*case-insensitive\*\*",
    r"match the digits by \*\*numeric value\*\*",
    r"`req-42`, `REQ-42`, and `REQ-042` all resolve to `REQ-042`",
    r"`Ur-11`/`UR-011` both resolve to `UR-011`",
    r"`user_request:` frontmatter",
):
    if not re.search(predicate, target_id_contract, re.IGNORECASE | re.MULTILINE):
        failures.append(f"{target_id_reference}: canonical Target ID Resolution source seam missing /{predicate}/")
require(completed_work_reference, r"ai-report.*present-video", "shared reference must name both current item-level consumers")
reject(completed_work_reference, r"future completed-work video", "shared reference still describes present-video as future")
for predicate in (r"completed-with-issues", r"Reject `cancelled`, `failed`", r"prompt-injection\.md"):
    require(completed_work_reference, predicate, f"shared reader missing /{predicate}/")
# The shared resolver's own Target ID source-seam assertions live with
# present-work's in the single detector below ("Target ID source-seam").

collision_match = re.search(
    r"^## Collision-Safe Publication\s*$[\s\S]*?(?=^## |\Z)",
    text(completed_work_reference),
    re.MULTILINE,
)
collision_contract = collision_match.group(0) if collision_match else ""
if not collision_match:
    failures.append(f"{completed_work_reference}: canonical Collision-Safe Publication section is missing")
for predicate in (
    r"Before creating any output directory or file[^\n]*complete final path[^\n]*already exists",
    r"first available numeric suffix",
    r"use that one path consistently for the whole artifact",
    r"failed or partial run",
    r"never delete, truncate, merge into, rename, migrate, or overwrite",
    r"Each consumer defines its own preferred path and output shape",
):
    if not re.search(predicate, collision_contract, re.IGNORECASE | re.MULTILINE):
        failures.append(f"{completed_work_reference}: canonical publication source seam missing /{predicate}/")
if re.search(r"\btimestamped\b", collision_contract, re.IGNORECASE):
    failures.append(f"{completed_work_reference}: canonical publication contract must be consumer-neutral, not timestamp-specific")

ai_report = "skills/do-work-toolbox/actions/ai-report.md"
require(ai_report, r"only action that produces detailed stakeholder-facing HTML", "ai-report must retain detailed HTML ownership")
for predicate in (
    r"real screenshots",
    r"SVG callouts",
    r"authentic before/after",
    r"full-page screenshots in both light and dark",
    r"UI captures were not expected for this work",
    r"Never fabricate a screenshot",
):
    require(ai_report, predicate, f"detailed report evidence contract missing /{predicate}/")
require(
    ai_report,
    r"must not create[^\n]*\ba video\b[^\n]*automatic video behavior",
    "ai-report must forbid both video output and automatic video behavior",
)
require_order(
    ai_report,
    r"Read and follow \[`completed-work-presentation-reference\.md`\].*in full \*\*before opening archived user content\*\*",
    r"Build the reference's provenance ledger",
    "ai-report must load the shared completed-work reference before using archived evidence",
)

# architecture-report. REQ-384 replaces the Markdown/carry-forward contract with a
# freeform HTML bundle and authored opening delta, keeping input, evidence, and immutability.
architecture_report = "skills/do-work-toolbox/actions/architecture-report.md"

# 1. Repo-wide input contract. The action takes no completed-work target, which is exactly
# what makes it a sibling of ai-report rather than a mode of it.
require(architecture_report, r"`\$ARGUMENTS` is ignored", "architecture-report must state that it ignores arguments")
require(
    architecture_report,
    r"there is no UR, REQ, or path-scoped form",
    "architecture-report must refuse a UR, REQ, or path-scoped target",
)
require(
    architecture_report,
    r"second, incompatible input contract",
    "architecture-report must say why folding it into ai-report breaks ai-report's input contract",
)
for retired_target_predicate in (
    r"Terminal-Success Target Resolution",
    r"blank is the explicit `most recent` form",
    r"do-work/archive/",
):
    reject(
        architecture_report,
        retired_target_predicate,
        f"architecture-report must not acquire a completed-work target contract /{retired_target_predicate}/",
    )

# 2. Dated-immutable publication. A canonical, in-place `docs/architecture-report.md` is the
# conversion this rejects: it would destroy the baseline every later delta is computed from.
require(
    architecture_report,
    r"ai-reports/<yyyy-mm-dd>_<hhmm>_architecture-report/index\.html",
    "architecture-report must publish the dated HTML entry point",
)
require(
    architecture_report,
    r"Never edit, delete, or regenerate a prior report",
    "architecture-report must forbid touching a prior report",
)
for undated_canonical_predicate in (
    r"ai-reports/architecture-report\.md",
    r"docs/architecture-report\.md",
    r"ai-reports/index\.html",
    r"docs/architecture-report\.html",
):
    reject(
        architecture_report,
        undated_canonical_predicate,
        f"architecture-report must not name an undated canonical report path /{undated_canonical_predicate}/",
    )
require(
    architecture_report,
    r"shared with `actions/ai-report\.md`",
    "architecture-report must publish into the same reports directory as ai-report",
)
require_order(
    architecture_report,
    r"### Step 6: Publish",
    r"Collision-Safe Publication",
    "architecture-report's publish step must delegate to the canonical no-clobber contract",
)
require(
    completed_work_reference,
    r"Collision-Safe Publication\*\* section is consumer-neutral",
    "the shared publication contract must state that a non-completed-work consumer may read it alone",
)

# 3. Freeform authoring. Pin the operative composition step, not stray quality vocabulary.
architecture_source = text(architecture_report)
architecture_compose_match = re.search(
    r"^### Step 4: Compose the Report\n(.*?)(?=^### Step 5:)",
    architecture_source, re.MULTILINE | re.DOTALL,
)
architecture_compose = architecture_compose_match.group(1) if architecture_compose_match else ""
for composition_predicate in (
    r"^Write one self-contained HTML document",
    r"^Redesign the report each run",
    r"layout, sectioning, and visual design.*authoring model",
    r"^Open with an authored.*changed since last report.*section",
    r"reading the previous HTML report",
    r"first report.*no prior HTML baseline",
    r"rendered diagrams.*not fenced code",
    r"clickable section navigation",
    r"Embed all presentation assets",
    r"no CDN or remote runtime dependencies",
    r'<meta name="architecture-report-verified-at" content="<head-hash>">',
    r"metadata.*does not prescribe.*layout",
):
    if not re.search(composition_predicate, architecture_compose, re.IGNORECASE | re.MULTILINE):
        failures.append(f"{architecture_report}: composition contract missing /{composition_predicate}/")
for retired_layout_predicate in (
    r"byte-identical",
    r"Only drifted sections are rewritten",
    r"Sections in this order",
    r"§[0-9Δ]",
    r"Markdown with inline Mermaid only",
    r"No HTML",
    r"two reports (?:can be )?diffed|diff.*between two reports",
):
    reject(
        architecture_report,
        retired_layout_predicate,
        f"architecture-report must retire its fixed Markdown/carry-forward layout /{retired_layout_predicate}/",
    )
require(architecture_report, r"Never write.*`architecture-report\.md`", "architecture-report must prohibit Markdown output")
require(architecture_report, r"Ignore.*Markdown.*prior.*baseline", "architecture-report must ignore Markdown baselines")

# 4. Verification labels and GitHub source links survive freeform composition.
for committed_tree_predicate in (
    r"Verify against a committed tree, never the working tree",
    r"Publish the report as a child commit of the one\s+it watermarks",
    r"never let the report assert its own\s+presence",
    r"VERIFIED.*GitHub.*file.*line",
    r"INFERRED.*basis",
    r"https://github\.com/<owner>/<repo>/blob/<head-hash>/<path>#L<line>",
    r"Never stop to ask",
):
    require(
        architecture_report,
        committed_tree_predicate,
        f"architecture-report must bind the watermark to a committed tree /{committed_tree_predicate}/",
    )

# 5. Audit findings stay with the audit loop, which owns bands and ratchets.
require(
    architecture_report,
    r"it never restates their findings",
    "architecture-report must keep audit findings out of the report body",
)
require(
    architecture_report,
    r"maintainability-audit\.md.*quick-wins\.md|quick-wins\.md.*maintainability-audit\.md",
    "architecture-report must redirect findings to the actions that own them",
)

# 6. Discovery surfaces. A routed action nobody can find is not shipped.
if route_rows.get("./actions/architecture-report.md") != (
    "architecture-report",
    "architecture overview",
    "map the repo",
):
    failures.append(
        "skills/do-work-toolbox/SKILL.md: ./actions/architecture-report.md owns "
        f"{route_rows.get('./actions/architecture-report.md')!r}, expected the three architecture triggers"
    )
if "architecture-report" not in hint_commands:
    failures.append("skills/do-work-toolbox/SKILL.md: argument-hint omits architecture-report")
require(
    "skills/do-work-toolbox/actions/help.md",
    r"architecture-report\s+.*dated.*immutable.*HTML.*architecture",
    "toolbox help must advertise architecture-report as a dated immutable HTML architecture map",
)
require(
    "skills/do-work/actions/help.md",
    r"do-work-toolbox[\s\S]{0,350}architecture-report",
    "core help must list architecture-report among the toolbox commands",
)
require(
    "README.md",
    r"do-work-toolbox architecture-report",
    "README must name the architecture-report command",
)
require("README.md", r"architecture-report.*index\.html", "README must name the HTML entry point")
reject("README.md", r"two reports can be diffed", "README must retire the report-to-report diff promise")


present_work = "skills/do-work-toolbox/actions/present-work.md"
for predicate in (
    r"Only the exact writing arguments `all` and `portfolio`",
    r"Blank input:[\s\S]{0,180}Usage: do-work-toolbox present-work all\|portfolio",
    r"One `UR-NNN` or `REQ-NNN` token[\s\S]{0,260}ai-report <ID>[\s\S]{0,160}present-video <ID>",
    r"Do not delegate them to another action",
    r"No.*canonical-only",
    r"Yes.*canonical-plus-snapshot",
    r"unavailable.*canonical-plus-snapshot",
    r"byte-identical",
    r"never delete a snapshot automatically",
    r"never truncate or replace",
    r"<skill-root>/scripts/publish-portfolio-summary\.sh",
    r"--canonical-only",
    r"--with-snapshot",
):
    require(present_work, predicate, f"portfolio contract missing /{predicate}/")
require(
    present_work,
    r"preserve the supplied (?:ID|token|spelling)[^\n]*both replacement commands",
    "present-work item dispatch must preserve the user's token in both printed commands",
)

# --- Target ID source-seam (REQ-203) ------------------------------------
# Every caller of the canonical Target ID Resolution contract must *inherit*
# it: name the source, actively apply it before its own resolution boundary,
# and never restate or negate the grammar locally. Keyword presence cannot
# express that — "read without applying Target ID Resolution" contains the
# letters `apply` and satisfied the assertion this block replaces — so the
# rules live in one detector, and the mutation matrix below replays the whole
# defect class against the real callers instead of spot-checking one spelling.
target_id_inheritance = re.compile(r"work-reference\.md[^\n]*Target ID Resolution", re.IGNORECASE)
target_id_application = re.compile(r"(?<![\w-])appl(?:y|ies)(?![\w-])[^\n]*Target ID Resolution", re.IGNORECASE)
target_id_negations = (
    r"(?:do not|don't|never|must not|cannot|without|instead of|rather than)\s+"
    r"(?:read(?:ing)?(?:\s+and)?\s+)?appl(?:y|ies|ying)(?![\w-])[^\n]*Target ID Resolution",
    r"(?:do not|don't|never|must not|cannot|without|instead of|rather than)\s+canonicali[sz](?:e|es|ing)(?![\w-])",
    r"(?:do not|don't|never|must not|cannot|without|instead of|rather than)\s+recogni[sz](?:e|es|ing)(?![\w-])",
)
target_id_copied_token_grammar = (
    r"case-insensitive",
    r"numeric[- ]value",
    r"`req-42`|`REQ-42`|`REQ-042`|`Ur-11`|`UR-011`",
)
target_id_copied_membership_grammar = (
    r"`user_request:`|`requests:` array",
)

def target_id_seam_findings(source, boundary_pattern):
    """Defect families in one caller of the canonical Target ID Resolution contract."""
    findings = []
    application = target_id_application.search(source)
    if not target_id_inheritance.search(source):
        findings.append("missing named inheritance")
    if not application:
        findings.append("missing active application directive")
    if any(re.search(negation, source, re.IGNORECASE) for negation in target_id_negations):
        findings.append("semantic negation of the inherited grammar")
    if any(re.search(copied, source, re.IGNORECASE) for copied in target_id_copied_token_grammar):
        findings.append("copied token grammar")
    if any(re.search(copied, source, re.IGNORECASE) for copied in target_id_copied_membership_grammar):
        findings.append("copied UR-membership grammar")
    boundary = re.search(boundary_pattern, source, re.IGNORECASE | re.MULTILINE)
    if not boundary:
        findings.append("missing resolution boundary")
    elif application and application.start() >= boundary.start():
        findings.append("application ordered after the resolution boundary")
    return findings

target_id_active_directive = "read and apply `../../do-work/actions/work-reference.md` → **Target ID Resolution**"

def demote_target_id_directive(demotion_anchor):
    """Move the caller's whole directive paragraph below its resolution boundary."""
    def mutate(source):
        paragraph = next(
            (block for block in source.split("\n\n") if target_id_active_directive in block),
            None,
        )
        if paragraph is None:
            return source  # caught by the "changed nothing" guard, never a traceback
        return source.replace(paragraph + "\n\n", "", 1).replace(
            demotion_anchor, paragraph + "\n\n" + demotion_anchor, 1
        )
    return mutate

target_id_callers = (
    (
        completed_work_reference,
        "shared presentation resolver",
        r"Resolve exactly one target",
        "Never fall back to `do-work/queue/`",
    ),
    (
        present_work,
        "present-work item dispatch",
        r"One `UR-NNN` or `REQ-NNN` token",
        "- **`all` or `portfolio`:** continue.",
    ),
)

target_id_defect_mutations = (
    (
        "semantic negation retaining an apply substring",
        lambda source: source.replace(
            target_id_active_directive,
            target_id_active_directive.replace("read and apply", "read without applying"),
            1,
        ),
        "semantic negation of the inherited grammar",
    ),
    (
        "explicit prohibition of the shared contract",
        lambda source: source.replace(
            target_id_active_directive,
            target_id_active_directive.replace("read and apply", "do not apply"),
            1,
        ),
        "semantic negation of the inherited grammar",
    ),
    (
        "passive citation with no application directive",
        lambda source: source.replace(
            target_id_active_directive,
            "see `../../do-work/actions/work-reference.md` → **Target ID Resolution**",
            1,
        ),
        "missing active application directive",
    ),
    (
        "canonical source no longer named",
        lambda source: source.replace(
            target_id_active_directive, "resolve the supplied token locally", 1
        ),
        "missing named inheritance",
    ),
    (
        "caller-local token grammar",
        lambda source: source + "\n\nSupplied tokens are case-insensitive.\n",
        "copied token grammar",
    ),
    (
        "caller-local UR-membership grammar",
        lambda source: source + "\n\nA REQ belongs to a UR when its `user_request:` frontmatter names it.\n",
        "copied UR-membership grammar",
    ),
)

target_id_safe_mutations = (
    ("unmodified caller", lambda source: source),
    (
        "reworded affirmative directive",
        lambda source: source.replace("read and apply", "read, then apply", 1),
    ),
    (
        "unrelated prohibition prose",
        lambda source: source + "\n\nStop without writing files when the target is unfinished.\n",
    ),
    (
        "caller-specific behavior added",
        lambda source: source + "\n\nReport the resolved archive path in the run summary.\n",
    ),
)

for relative_path, seam_name, boundary_pattern, demotion_anchor in target_id_callers:
    caller_source = text(relative_path)
    # Positive control: the shipped caller is clean under every rule at once.
    shipped_findings = target_id_seam_findings(caller_source, boundary_pattern)
    for finding in shipped_findings:
        failures.append(f"{relative_path}: {seam_name} — {finding}")
    if shipped_findings:
        # Replaying mutations on an already-defective caller only produces
        # derived noise; the findings above are the ones to fix.
        continue
    caller_mutations = target_id_defect_mutations + (
        (
            "directive demoted below the resolution boundary",
            demote_target_id_directive(demotion_anchor),
            "application ordered after the resolution boundary",
        ),
    )
    for mutation_name, mutate, expected_family in caller_mutations:
        mutated_source = mutate(caller_source)
        if mutated_source == caller_source:
            failures.append(
                f"{relative_path}: {seam_name} mutation {mutation_name!r} changed nothing — "
                "the replay no longer matches the shipped text"
            )
            continue
        mutated_families = target_id_seam_findings(mutated_source, boundary_pattern)
        if expected_family not in mutated_families:
            failures.append(
                f"{relative_path}: {seam_name} mutation {mutation_name!r} escaped expected family "
                f"{expected_family!r}; found {mutated_families!r}"
            )
    for mutation_name, mutate in target_id_safe_mutations:
        safe_families = target_id_seam_findings(mutate(caller_source), boundary_pattern)
        if safe_families:
            failures.append(
                f"{relative_path}: {seam_name} rejected safe mutation {mutation_name!r}: "
                f"{safe_families!r}"
            )
require_order(
    present_work,
    r"read `\.\./\.\./do-work/crew-members/prompt-injection\.md`",
    r"Scan archived UR folders and legacy REQs",
    "present-work must load prompt-injection guidance before scanning archive records",
)
require_order(
    present_work,
    r"snapshot[^\n]*success",
    r"atomically refresh[^\n]*canonical[^\n]*same bytes",
    "present-work snapshot branch must publish successfully before canonical refresh",
)
reject(
    present_work,
    r"(?i)(?:write|refresh)[^\n]{0,100}canonical[^\n]{0,180}(?:then|before)[^\n]{0,100}(?:publish|create)[^\n]{0,80}snapshot",
    "present-work must not prescribe canonical-first snapshot publication",
)
for retired_workflow_token in (
    r"\bDetail Mode\b",
    r"\bInteractive Explainer\b",
    r"\bclient brief\b",
    r"\bsibling-link\b",
    r"\bdetail-depth\b",
    r"--with-video",
):
    reject(present_work, retired_workflow_token, f"retired present-work workflow token remains /{retired_workflow_token}/")

for predicate in (
    r"only writing forms are `all` and `portfolio`",
    r"bare invocation.*writes nothing",
    r"item-specific invocation.*writes nothing",
    r"does not silently delegate",
    r"No.*only",
    r"Yes.*byte-identical snapshot",
    r"cannot be asked or answered.*safer preservation branch",
    r"never deletes snapshots automatically",
    r"future REQ.*Lessons Learned.*cite a snapshot",
    r"does not authorize.*back-edit archived REQs or lessons",
    r"snapshot[^\n]*before[^\n]*canonical",
    r"snapshot publication fails[^\n]*prior canonical[^\n]*unchanged",
    r"canonical refresh fails[^\n]*snapshot[^\n]*retained",
):
    require("skills/do-work-toolbox/docs/present-work-guide.md", predicate, f"portfolio guide missing /{predicate}/")

present_video = "skills/do-work-toolbox/actions/present-video.md"
for predicate in (
    r"source-only",
    r"explicit.*present-video.*Remotion.*video walkthrough",
    r"package\.json[\s\S]{0,220}tsconfig\.json[\s\S]{0,220}Root\.tsx[\s\S]{0,220}Video\.tsx[\s\S]{0,300}ProblemScene\.tsx[\s\S]{0,220}ValueScene\.tsx",
    r"registerRoot\(RemotionRoot\)",
    r'"preview": "remotion studio src/Root\.tsx"',
    r"Do not install dependencies",
    r"never renders media",
):
    require(present_video, predicate, f"source-only video contract missing /{predicate}/")
require_order(
    present_video,
    r"Read and follow \[`completed-work-presentation-reference\.md`\].*in full \*\*before opening archived user content\*\*",
    r"Build the reference's provenance ledger",
    "present-video must load the shared completed-work reference before using archived evidence",
)
# --- Publication delegation (REQ-206) -----------------------------------
# Naming the shared section is not delegating to it. A consumer satisfied the
# assertion this block replaces with a passive mention in its Verification
# Checklist, while its execution step carried a paraphrase of the algorithm.
# Delegation now means an active directive *at the step that creates output*,
# with no local restatement of the mechanics anywhere in the file.
publication_application_directive = re.compile(
    r"(?<![\w-])appl(?:y|ies)(?![\w-])[^\n]*Collision-Safe Publication", re.IGNORECASE
)
publication_local_algorithm = (
    r"one final path[^\n]*(?:every|each)[^\n]*file",
    r"(?:every|each)[^\n]*file[^\n]*one final path",
    r"if the preferred[^\n]*exists",
    r"numeric suffix|suffixed sibling",
    r"existing (?:artifacts|deliverables) are immutable",
    r"never (?:delete|truncate|merge|overwrite|rename|migrate)[^\n]*(?:delete|truncate|merge|overwrite|rename|migrate)[^\n]*(?:delete|truncate|merge|overwrite|rename|migrate)",
    r"no pre-existing path|no existing output was changed",
)

def publication_delegation_findings(source, output_creation_pattern):
    """Defect families in one consumer of the shared Collision-Safe Publication section."""
    findings = []
    directive = publication_application_directive.search(source)
    if not directive:
        findings.append("missing active application directive")
    if any(re.search(restatement, source, re.IGNORECASE) for restatement in publication_local_algorithm):
        findings.append("local publication algorithm restated")
    output_creation = re.search(output_creation_pattern, source, re.IGNORECASE | re.MULTILINE)
    if not output_creation:
        findings.append("missing output-creation step")
    elif directive and directive.start() >= output_creation.start():
        findings.append("application ordered after output creation")
    return findings

publication_consumers = (
    (ai_report, "ai-report", r"ai-reports/<report-slug>/", r"Create `screenshots/` only when"),
    (
        present_video,
        "present-video",
        r"do-work/deliverables/<canonical-ID>-video/",
        r"Write only the source tree from Step 4",
    ),
)

publication_delegation_mutations = (
    (
        "active directive removed, passive checklist mention left standing",
        lambda source: publication_application_directive.sub("follow", source, count=1),
        "missing active application directive",
    ),
    (
        "algorithm paraphrased locally",
        lambda source: source
        + "\n\nUse the one final path selected by that contract for every source file.\n",
        "local publication algorithm restated",
    ),
)

for consumer, consumer_name, preferred_path_pattern, output_creation_pattern in publication_consumers:
    require(consumer, preferred_path_pattern, "presentation consumer must retain its preferred output path")
    consumer_source = text(consumer)
    shipped_findings = publication_delegation_findings(consumer_source, output_creation_pattern)
    for finding in shipped_findings:
        failures.append(f"{consumer}: {consumer_name} publication delegation — {finding}")
    if shipped_findings:
        continue
    for mutation_name, mutate, expected_family in publication_delegation_mutations:
        mutated_source = mutate(consumer_source)
        if mutated_source == consumer_source:
            failures.append(
                f"{consumer}: {consumer_name} delegation mutation {mutation_name!r} changed nothing — "
                "the replay no longer matches the shipped text"
            )
            continue
        mutated_families = publication_delegation_findings(mutated_source, output_creation_pattern)
        if expected_family not in mutated_families:
            failures.append(
                f"{consumer}: {consumer_name} delegation mutation {mutation_name!r} escaped expected family "
                f"{expected_family!r}; found {mutated_families!r}"
            )
for video_surface in (
    present_video,
    "skills/do-work-toolbox/docs/present-video-guide.md",
):
    reject_executable_video_forms(video_surface)
require(present_video, r"Never invoke it from `ai-report`, `present-work`, a completion flow", "video action must remain explicit rather than automatic")
require(present_work, r"does not produce per-item briefs.*stakeholder HTML.*video", "portfolio action must reject item-level report and video artifacts")

reject("skills/do-work/actions/review-work.md", r"present-work\.md` parses it for the score", "review still claims portfolio parses a persisted score")
for relative_path in (
    "skills/do-work/actions/capture.md",
    "skills/do-work/actions/work.md",
    "skills/do-work/actions/work-reference.md",
    "skills/do-work/actions/abandon.md",
):
    require(relative_path, r"presentation action|completed-work-presentation-reference", "terminal-success caller guidance must cover the shared presentation family")

retired_patterns = (
    r"Detail Mode",
    r"(?<!non-)Interactive Explainer",
    r"client brief",
    r"--with-video",
    r"--no-open\s*&",
    r"localhost:3000",
    r"remotion render",
)
live_presentation_surfaces = (
    "README.md",
    "skills/do-work-toolbox/SKILL.md",
    "skills/do-work-toolbox/actions/help.md",
    "skills/do-work-toolbox/actions/tutorial.md",
    "skills/do-work-toolbox/actions/completed-work-presentation-reference.md",
    "skills/do-work-toolbox/docs/present-work-guide.md",
    "skills/do-work/actions/help.md",
    "skills/do-work/actions/capture.md",
    "skills/do-work/actions/review-work.md",
    "skills/do-work/actions/work.md",
    "skills/do-work/actions/work-reference.md",
    "skills/do-work/actions/abandon.md",
)
for relative_path in live_presentation_surfaces:
    for pattern in retired_patterns:
        reject(relative_path, pattern, f"retired or unsafe presentation contract remains /{pattern}/")

if failures:
    raise SystemExit("presentation contract failures:\n- " + "\n- ".join(failures))
PY
then
  printf 'FAIL: presentation routing, evidence, portfolio, video, or caller contracts drifted.\n' >&2
  fail_count=$((fail_count + 1))
fi

assert_block_contains \
  "$skill_dispatch_block" \
  '^\| `run`[^|]*\| `\./actions/work\.md`' \
  'Core SKILL.md must route work triggers to actions/work.md so scoped REQ IDs and --wave reach the action input.'

assert_block_contains \
  "$review_archived_input_block" \
  'Archived fallback is context-only' \
  'review-work.md Step 3 must identify the archived-input fallback as context-only.'

assert_block_contains \
  "$review_archived_input_block" \
  'Whenever this archived fallback is used' \
  'review-work.md Step 3 must apply the archived-input authority boundary regardless of review mode.'

assert_block_contains \
  "$review_archived_input_block" \
  'grants no authority to move, reopen, or re-consolidate the closed UR folder.*stays closed and in place' \
  'review-work.md Step 3 must keep the archived UR closed and in place without move, reopen, or re-consolidation authority.'

assert_block_contains \
  "$review_archived_input_block" \
  'follow-up keeps the same `user_request`' \
  'review-work.md Step 3 must keep standalone-review follow-ups on the reviewed REQ user_request.'

assert_block_contains \
  "$review_archived_input_block" \
  'carries `review_generated: true`' \
  'review-work.md Step 3 must retain the review_generated marker on standalone-review follow-ups.'

assert_block_contains \
  "$review_archived_input_block" \
  'goes into `do-work/queue/`' \
  'review-work.md Step 3 must place standalone-review follow-ups in the queue without reopening the UR.'

assert_block_contains \
  "$review_archived_input_block" \
  'Ordinary orchestrated review.*active UR.*still open' \
  'review-work.md Step 3 must identify ordinary orchestrated review while its UR is still active.'

assert_block_contains \
  "$review_archived_input_block" \
  'orchestrated review of a `review_generated: true` follow-up whose UR is already archived.*inherits the same context-only, stays-closed boundary' \
  'review-work.md Step 3 must apply the archived fallback boundary to orchestrated review-generated follow-ups whose UR is already closed.'

assert_block_contains \
  "$review_archived_input_block" \
  'narrow exception.*post-work Review metadata.*archived REQ.*remains unchanged' \
  'review-work.md Step 3 must preserve the narrow archived-REQ review-metadata exception.'

assert_block_contains \
  "$work_archive_success_block" \
  'When and only when the completed REQ carries `review_generated: true` and `do-work/archive/UR-NNN/` already exists for its `user_request`' \
  'actions/work.md Archive success override must require the complete review_generated-and-existing-archived-UR predicate for the same user_request.'

assert_block_contains \
  "$work_archive_success_block" \
  'move the completed REQ into that existing folder in place' \
  'actions/work.md Archive success override must return the completed review follow-up to its existing archived UR folder in place.'

assert_block_contains \
  "$work_archive_success_block" \
  'Never move, reopen, or re-consolidate the archived UR folder' \
  'actions/work.md Archive success override must keep the archived UR folder closed and stationary.'

assert_block_contains \
  "$work_archive_success_block" \
  'Skip the normal active-UR closure branch' \
  'actions/work.md Archive success override must bypass the normal active-UR closure branch.'

assert_block_contains \
  "$work_archive_success_block" \
  'already (set to )?`completed-with-issues`|status is already `completed-with-issues`|preserve[^[:cntrl:]]*`completed-with-issues`' \
  'actions/work.md Archive success path must explicitly preserve completed-with-issues from failed remediation.'

assert_block_not_contains \
  "$work_archive_success_block" \
  '^1\. Update frontmatter: `status: completed`, `completed_at: <timestamp>`$' \
  'actions/work.md Archive success path must not unconditionally overwrite status with completed.'

assert_contains \
  "actions/ai-report.md" \
  'DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND' \
  'actions/ai-report.md must keep sandbox-bypassed agentic image generation behind an explicit opt-in.'

ai_image_backend_block="$(sed -n '/^## Image Generation Backend/,/^## SVG Data-Viz Rules/p' "$toolbox_root/actions/ai-report-reference.md")"
assert_block_contains \
  "$ai_image_backend_block" \
  '<skill-root>/scripts/generate-report-image-batch\.sh' \
  'ai-report-reference.md must delegate the image batch to the shipped batch script.'
assert_block_not_contains \
  "$ai_image_backend_block" \
  'launch_report_image|image_generation_pids|terminate_report_image_batch' \
  'ai-report-reference.md must not restate the batch mechanics it delegates; keep the prompts and the invocation.'
# REQ-198's lock-in survives the move to a script: delegation is not what stops the action
# file from pre-creating generated/, so this guard stays pointed at the action file.
assert_block_not_contains \
  "$ai_image_backend_block" \
  'GEN="ai-reports/<report-slug>/generated"; mkdir -p "\$GEN"' \
  'ai-report-reference.md must not publish generated/ before a current image succeeds.'
assert_block_contains \
  "$ai_image_backend_block" \
  'repository, credentials, network, and external services' \
  'ai-report-reference.md must describe exact agentic opt-in as full-host authority, not containment.'

# REQ-221 moved these mechanics out of the action file into their executable home, so the
# guarantees REQ-198 and REQ-204 locked in are replayed against the script. Deleting one
# must fail here rather than pass because the code changed address.
ai_image_batch_source="$(cat "$toolbox_root/scripts/generate-report-image-batch.sh")"
assert_block_contains \
  "$ai_image_batch_source" \
  'image_generation_pids' \
  'generate-report-image-batch.sh must retain every parallel image helper PID.'
assert_block_contains \
  "$ai_image_batch_source" \
  'wait "\${image_generation_pids\[\$image_index\]}".*image_status=\$\?' \
  'generate-report-image-batch.sh must wait each image PID and retain its individual status.'
assert_block_contains \
  "$ai_image_batch_source" \
  'image_generation_statuses\[\$image_index\]' \
  'generate-report-image-batch.sh must evaluate invocation statuses rather than infer freshness from target presence.'
assert_block_not_contains \
  "$ai_image_batch_source" \
  '^wait$' \
  'generate-report-image-batch.sh must not discard mixed background-job statuses with a bare wait.'
assert_block_contains \
  "$ai_image_batch_source" \
  'mktemp -d "\$report_directory/\.generated\.staging\.XXXXXX"' \
  'generate-report-image-batch.sh must allocate one invocation-private staging directory adjacent to generated/.'
assert_block_contains \
  "$ai_image_batch_source" \
  '\[ ! -e "\$generated_directory" \].*exit 1' \
  'generate-report-image-batch.sh must fail closed instead of clobbering an existing generated/ directory.'
assert_block_contains \
  "$ai_image_batch_source" \
  'image_generation_success_count=0' \
  'generate-report-image-batch.sh must count only status-backed successful images before publication.'
assert_block_contains \
  "$ai_image_batch_source" \
  '\[ "\$image_generation_success_count" -gt 0 \]' \
  'generate-report-image-batch.sh must make publication conditional on at least one status-backed successful image.'
assert_block_contains \
  "$ai_image_batch_source" \
  'mv "\$image_generation_stage" "\$generated_directory"' \
  'generate-report-image-batch.sh must publish the complete verified batch with one adjacent same-filesystem rename.'
assert_block_contains \
  "$ai_image_batch_source" \
  'cleanup_report_image_stage' \
  'generate-report-image-batch.sh must clean the exact invocation-private directory on all-failed and interrupted runs.'
assert_block_contains \
  "$ai_image_batch_source" \
  'nested_image_generation_stage' \
  'generate-report-image-batch.sh must verify the publishing rename instead of trusting a zero-exit mv onto a directory.'

assert_contains \
  "actions/ai-report.md" \
  'mktemp -d' \
  'actions/ai-report.md must run any agentic image fallback from a locked temporary directory, not the repo cwd.'

assert_contains \
  "actions/ai-report.md" \
  'chmod 700' \
  'actions/ai-report.md must lock down the temporary image-generation directory before invoking an agentic backend.'

assert_contains \
  "actions/version.md" \
  'tools/do-work-update\.sh.*--project-root' \
  'actions/version.md update flow must delegate mutation to the shared updater engine.'

assert_contains \
  "actions/version.md" \
  'do not duplicate its archive or overwrite logic here' \
  'actions/version.md must keep archive review and overwrite logic in the shared updater rather than a second agent-only implementation.'

assert_file_not_contains \
  "actions/version.md" \
  'log -1 --format=%H -- actions/version\.md' \
  'actions/version.md must not use the last version.md-touching commit as the committed-customization baseline.'

# No template may emit an instruction addressed to its reader (REQ-080). The Simple REQ template ended
# with a bare "Think carefully before answering." *inside* the fence, so every REQ capture produced
# inherited it — 25 archived REQs carry it, and one builder correctly flagged it as an instruction-like
# artifact and treated it as data (REQ-012 D-02) without anyone looking at the source that kept
# emitting it. A generated work item must contain the request and nothing that reads as a directive to
# whoever processes it; that shape is precisely what crew-members/prompt-injection.md exists to catch,
# and the skill must not manufacture it. Literal negatives on the known phrasings, deliberately not an
# English-grammar detector — the wordings are illustrative, the condition above is the rule.
for capture_artifact_phrase in \
  'Think carefully before answering' \
  'Reason step by step' \
  'Take your time'
do
  assert_file_not_contains \
    "actions/capture-reference.md" \
    "$capture_artifact_phrase" \
    "actions/capture-reference.md templates must not contain an instruction addressed to the reader — it lands inside the fence and is copied into every generated REQ, where it reads as a directive to the builder rather than as the user's request (REQ-080). Matched: $capture_artifact_phrase"
done

assert_contains \
  "actions/capture.md" \
  'maintenance: false' \
  'actions/capture.md base REQ schema must carry maintenance:false so the marker is discoverable, not documented only for complex requests.'

assert_contains \
  "actions/capture.md" \
  'Maintenance assessment' \
  'actions/capture.md Step 1 must assess skill-instruction removal/narrowing and set the maintenance marker (work.md is marker-only and never infers it).'

# One home for the timestamp command (REQ-078). The Timestamp rule in actions/work-reference.md is
# the only place in actions/ that spells a command for obtaining a stamp; every other site cites the
# rule. Eleven inline copies of `date -u +…` is why 0.167.0's new source-preference order — and the
# Windows form in particular — was unreachable from any site an agent actually follows: the citing
# site already handed it a POSIX command, so it never read the rule. Trigger condition, not a file
# list: any file under actions/ other than the rule's own home that spells a `date -u +…` invocation
# is a regression, whatever the format specifier. hooks/ is out of scope on purpose (executable POSIX
# shell, platform-specific by design), as is tools/ (Go source and user-facing strings, a different
# surface with a different tradeoff).
timestamp_command_copies="$(grep -rlE 'date -u \+%' "$core_root/actions" | grep -v 'work-reference\.md' || true)"
if [ -n "$timestamp_command_copies" ]; then
  printf 'FAIL: only actions/work-reference.md may spell a timestamp command; these action files inline a copy, so an agent following them never reaches the rule preference order and a Windows agent gets a command that does not exist on its box — cite the Timestamp rule instead (REQ-078):\n%s\n' \
    "$timestamp_command_copies" >&2
  fail_count=$((fail_count + 1))
fi

# The other half of that same arrangement: because the rule's home is the ONLY place that spells a
# clock command, every other stamp write site has to point at it, or an agent filling a template from
# context never re-reads the rule. That is not hypothetical — three completed_at stamps were once
# written by extrapolating the clock forward instead of reading it, and nothing at those sites said
# otherwise (REQ-244).
#
# Recognition is keyed on shape, never on a spelling or a site list. A site is a placeholder — any
# bracketed or quoted span short enough to name a value rather than be prose — that is assigned to
# something and that either denotes a clock value by name or is the value of an `*_at` key, which is
# the rule's own stated trigger ("every `*_at` field in this schema, and any timestamp a future field
# adds"). The denoting vocabulary below is illustrative and meant to grow; what must not grow is a
# list of files. A first pass at this check keyed on the literal bracket forms it had just fixed and
# consequently could not see `<ISO 8601 timestamp>`, `{timestamp}`, `"<iso>"`, or a bare `*_at` key —
# the enumeration had simply moved up one level (REQ-244 review, F2).
#
# Two requirements at every INSTANT site:
#   1. Spelling. The line carries `<timestamp>` or `<now>` — the two spellings the rule's own sentence
#      says mean it. Recognition stays broad so a new spelling is still caught; the requirement stays
#      narrow so it gets normalized.
#   2. Citation. `Timestamp rule` appears on the site's own line, or — when the site sits inside a
#      fenced block copied verbatim into a generated artifact — on the nearest non-blank line above
#      that fence's opening delimiter. Those are the only two positions a reader of the site sees.
#
# Placeholders naming a *directory* or a path are names, not stamps, and are skipped by shape (a `/`
# inside, or `-`/`/` glued to the outside). DATE-ONLY sites are recognized and counted but are NOT
# required to cite: the rule's Date-only paragraph governs UTC calendar dates while explicitly
# leaving local dates (report slugs, changelog headings) ungoverned, and no one has yet decided which
# of the report-header dates are which. Counting them keeps that open question visible instead of
# silently in scope. The zero-site guard is the anti-vacuity clause, and the success line prints both
# counts so a run that quietly stops recognizing sites is visible rather than green.
if ! python3 - "$repo_root" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
failures = []

placeholder_span = re.compile(r"<([^<>\n]*)>|\{([^{}\n]*)\}|\[([^\[\]\n]*)\]|\"([^\"\n]*)\"")
nested_placeholder = re.compile(r"<[^<>\n]*>|\{[^{}\n]*\}")
instant_word = re.compile(r"\b(?:timestamp|now|iso)\b", re.IGNORECASE)
instant_shape = re.compile(r"YYYY-MM-DDTHH")
date_word = re.compile(r"\b(?:date|today)\b", re.IGNORECASE)
date_shape = re.compile(r"YYYY-MM-DD")
stamp_key_assignment = re.compile(r"[A-Za-z_][A-Za-z0-9_]*_at\"?\s*[:=]\s*$")
assignment = re.compile(r"[:=]")
recognized_spelling = re.compile(r"<(?:timestamp|now)>")
rule_citation = re.compile(r"Timestamp rule")
fence_delimiter = re.compile(r"^\s*(```|~~~)")

instant_site_count = 0
date_only_site_count = 0
for action_file in sorted(root.glob("skills/*/actions/*.md")):
    relative_path = action_file.relative_to(root).as_posix()
    inside_fence = False
    fence_marker = ""
    fence_introduction = ""
    last_nonblank_line = ""
    for line_number, line in enumerate(
        action_file.read_text(encoding="utf-8").splitlines(), start=1
    ):
        fence_match = fence_delimiter.match(line)
        if fence_match:
            marker = fence_match.group(1)
            if not inside_fence:
                inside_fence = True
                fence_marker = marker
                fence_introduction = last_nonblank_line
            elif line.strip().startswith(fence_marker):
                inside_fence = False
                fence_marker = ""
                fence_introduction = ""
            if line.strip():
                last_nonblank_line = line
            continue

        for match in placeholder_span.finditer(line):
            inner_text = next(group for group in match.groups() if group is not None)
            written_form = match.group(0)
            preceding = line[: match.start()]
            following = line[match.end() :]
            if len(inner_text) > 30 and not nested_placeholder.search(inner_text):
                continue
            if "/" in inner_text or preceding[-1:] in ("-", "/") or following[:1] in ("-", "/"):
                continue
            if not assignment.search(preceding):
                continue
            names_an_instant = bool(
                instant_word.search(inner_text)
                or instant_shape.search(inner_text)
                or stamp_key_assignment.search(preceding)
            )
            names_a_date = bool(date_word.search(inner_text) or date_shape.search(inner_text))
            if not names_an_instant:
                date_only_site_count += names_a_date
                continue
            instant_site_count += 1
            if not recognized_spelling.search(line):
                failures.append(
                    f"{relative_path}:{line_number}: instant write site is written "
                    f"{written_form}; the rule recognizes <timestamp> and <now>"
                )
            if not rule_citation.search(line) and not (
                inside_fence and rule_citation.search(fence_introduction)
            ):
                failures.append(
                    f"{relative_path}:{line_number}: stamp write site cites no Timestamp rule "
                    f"on its own line or above its fence: {line.strip()[:110]}"
                )
        if line.strip():
            last_nonblank_line = line

if instant_site_count == 0:
    failures.append(
        "no instant write site was recognized in any shipped action file — the shapes this "
        "check keys on changed and it has gone blind"
    )

if failures:
    raise SystemExit("timestamp citation failures:\n- " + "\n- ".join(failures))
print(
    f"Timestamp rule citation contract: {instant_site_count} instant write sites cited, "
    f"{date_only_site_count} date-only sites recognized."
)
PY
then
  printf 'FAIL: a shipped stamp write site uses an unrecognized spelling or points at no Timestamp rule, so an agent filling it has nothing telling it to read a clock (REQ-244).\n' >&2
  fail_count=$((fail_count + 1))
fi

# REQ-316 — calibration-log arithmetic must project the stamps that actually
# reached the archived REQ, not values an agent still carries from earlier in the
# run. This is deliberately clause-local: generic timestamp prose elsewhere cannot
# satisfy the source-selection boundary at the cross-file write. The in-memory
# mutations replay the REQ-274 stale-stamp class and prove every semantic leg bites.
if ! python3 - "$core_root/actions/work.md" <<'PY'
import pathlib
import re
import sys

work_path = pathlib.Path(sys.argv[1])
work_text = work_path.read_text()


def calibration_source_defects(source):
    section_match = re.search(
        r"^7\.5\. \*\*Append the calibration-log line\.\*\*(.*?)(?=^8\. \*\*Worktree cleanup)",
        source,
        flags=re.DOTALL | re.MULTILINE,
    )
    if section_match is None:
        return {"section"}

    paragraphs = [
        " ".join(paragraph.split())
        for paragraph in re.split(r"\n\s*\n", section_match.group(1))
        if paragraph.strip()
    ]
    directives = [
        paragraph
        for paragraph in paragraphs
        if re.search(r"\bread\b", paragraph, flags=re.IGNORECASE)
        and ("`claimed_at`" in paragraph or "`completed_at`" in paragraph)
        and re.search(
            r"calculation time|just.archived|frontmatter|context",
            paragraph,
            flags=re.IGNORECASE,
        )
    ]
    if len(directives) != 1:
        return {"directive"}

    directive = directives[0]
    defects = set()
    if "`claimed_at`" not in directive:
        defects.add("claimed_at")
    if "`completed_at`" not in directive:
        defects.add("completed_at")
    if re.search(r"\bjust.archived REQ\b", directive, flags=re.IGNORECASE) is None:
        defects.add("just-archived-source")
    if re.search(r"\bfrontmatter\b", directive, flags=re.IGNORECASE) is None:
        defects.add("frontmatter")
    if re.search(
        r"\b(?:never|do not|must not)\b.{0,80}\breuse\b.{0,100}\b(?:held|carried)\b.{0,60}\b(?:context|earlier)\b",
        directive,
        flags=re.IGNORECASE,
    ) is None:
        defects.add("carried-value-ban")
    return defects


live_defects = calibration_source_defects(work_text)
if live_defects:
    raise SystemExit(
        "actions/work.md Step 8 substep 7.5 has an incomplete persisted-stamp "
        f"source contract: {sorted(live_defects)}"
    )

section_match = re.search(
    r"^7\.5\. \*\*Append the calibration-log line\.\*\*(.*?)(?=^8\. \*\*Worktree cleanup)",
    work_text,
    flags=re.DOTALL | re.MULTILINE,
)
paragraphs = [
    paragraph.strip()
    for paragraph in re.split(r"\n\s*\n", section_match.group(1))
    if paragraph.strip()
]
directive = next(
    paragraph
    for paragraph in paragraphs
    if re.search(r"\bread\b", paragraph, flags=re.IGNORECASE)
    and ("`claimed_at`" in paragraph or "`completed_at`" in paragraph)
    and re.search(
        r"calculation time|just.archived|frontmatter|context",
        paragraph,
        flags=re.IGNORECASE,
    )
)


def replace_directive(replacement):
    if directive not in work_text:
        raise AssertionError("live calibration source directive is not replaceable")
    return work_text.replace(directive, replacement, 1)


def replace_once(source, old, new, mutation_name):
    mutated = source.replace(old, new, 1)
    if mutated == source:
        raise AssertionError(f"calibration source mutation {mutation_name!r} changed nothing")
    return mutated


mutations = (
    ("delete clause", replace_directive(""), "directive"),
    (
        "omit claimed_at",
        replace_directive(replace_once(directive, "`claimed_at`", "the claim stamp", "omit claimed_at")),
        "claimed_at",
    ),
    (
        "omit completed_at",
        replace_directive(replace_once(directive, "`completed_at`", "the completion stamp", "omit completed_at")),
        "completed_at",
    ),
    (
        "use generic context source",
        replace_directive(
            re.sub(
                r"from the just.archived REQ(?: file)?(?:'s)? frontmatter",
                "from context",
                directive,
                count=1,
                flags=re.IGNORECASE,
            )
        ),
        "just-archived-source",
    ),
    (
        "remove frontmatter source",
        replace_directive(replace_once(directive, " frontmatter", "", "remove frontmatter source")),
        "frontmatter",
    ),
    (
        "delete carried-value ban",
        replace_directive(re.sub(r";\s*never\b.*?\.\*\*$", ".**", directive, count=1, flags=re.IGNORECASE)),
        "carried-value-ban",
    ),
    (
        "invert carried-value ban",
        replace_directive(replace_once(directive, "never reuse", "may reuse", "invert carried-value ban")),
        "carried-value-ban",
    ),
    (
        "displace clause outside 7.5",
        replace_directive("") + "\n\n" + directive + "\n",
        "directive",
    ),
)

for mutation_name, mutated_source, expected_defect in mutations:
    if mutated_source == work_text:
        raise SystemExit(f"calibration source mutation {mutation_name!r} changed nothing")
    mutation_defects = calibration_source_defects(mutated_source)
    if expected_defect not in mutation_defects:
        raise SystemExit(
            f"calibration source mutation {mutation_name!r} escaped "
            f"{expected_defect!r}; found {sorted(mutation_defects)!r}"
        )
PY
then
  printf 'FAIL: actions/work.md Step 8 substep 7.5 must read claimed_at and completed_at from the just-archived REQ frontmatter at calculation time and forbid carried context values (REQ-316).\n' >&2
  fail_count=$((fail_count + 1))
fi

timestamp_rule_block="$(sed -n '/^\*\*Timestamp rule —/,/^```yaml/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$timestamp_rule_block" \
  'ToUniversalTime' \
  'actions/work-reference.md Timestamp rule must name the .ToUniversalTime() form for Windows — Get-Date -AsUTC is PowerShell 7+ and a stock box ships 5.1 as powershell.exe, where that parameter is unrecognized and the call fails outright (REQ-078).'

assert_block_contains \
  "$timestamp_rule_block" \
  'powershell -NoProfile -Command' \
  'actions/work-reference.md Timestamp rule must give the cmd entry point explicitly — a bare cmdlet is not a command in cmd, which is the shell the Windows clause exists for (REQ-078).'

# write_set's display-only status must never be argued from a one-REQ-at-a-time premise (REQ-075).
# The conclusion is right and permanent — nothing schedules, gates, or dispatches on write_set — but
# REQ-073 falsified that premise: several builders can run at once under a single queue owner. A reader
# who follows the old reasoning concludes the opposite of the contract, since the stated cause is gone.
# The durable reason is that write_set is advisory input to a human's pick and the merge is the
# non-interference proof (actions/work-reference.md → Worktree Dispatch Mode → Fan-Out Dispatch), which
# holds at any builder count.
#
# Granularity is one line, and deliberately so: the premise is only dangerous sitting next to the field
# it purports to explain, and a proximity window false-positives on the canonical Fan-Out Dispatch
# section itself, two lines below the advisory-`write_set` bullet, where "integration runs one REQ at a
# time" is a true statement about integration. Prose files put both on one long line, so the sweep sees
# them; Go and JS comments wrap, so it cannot — those two files are covered by the file-level negatives
# below instead, which is exact because neither has any legitimate reason to state a builder count.
# `|| true` because grep exits 1 on no-match, which is the passing case and must not abort under `set -e`.
#
# TRIGGER CONDITION (REQ-079), not a phrase list: no shipped file may argue write_set's display-only
# status from ANY builder count, in either of the premise's two fingerprints — the thing it *said*
# ("one REQ at a time") and the thing it was *called* ("under the exclusive-session model"). The weak
# form is the more dangerous of the two, because the model it names is still true and only its
# relevance to this conclusion died, so the sentence survives inspection. Both patterns below enumerate
# the wordings seen so far and are ILLUSTRATIVE, NOT EXHAUSTIVE — a new phrasing earns a widened
# pattern, never a second copy of the block. The strong-form pattern is defined once and reused by all
# three of its consumers for the same reason.
#
# tools/queue-kanban/prime-do-kanban.md is deliberately NOT given a file-level negative: its REQ-075
# lesson entry quotes both fingerprints verbatim, which is the one legitimate reason to name them. The
# line sweep's write_set|overlaps filter is what keeps that line out of scope, and it must stay that way.
builder_count_premise_pattern='(one|a single|only one)( [a-z]+){0,2} (REQ|builder|coder|agent)s?( [a-z]+){0,3} (at a time|at once|concurrently)|(one|a single|only one)( [a-z]+){0,2} (REQ|builder|coder|agent)s? is ever (building|running|in flight)'
exclusive_session_premise_pattern='exclusive.session|no other .do-work. session'

stale_write_set_premise_lines="$(
  grep -rIEn "$builder_count_premise_pattern" \
    "$core_root/actions" "$core_root/docs" "$board_root/tools/queue-kanban" 2>/dev/null \
    | grep -E 'write_set|overlaps' || true
)"
if [ -n "$stale_write_set_premise_lines" ]; then
  printf 'FAIL: write_set is display-only at ANY builder count — no shipped file may justify that with a one-REQ-at-a-time premise (REQ-075, pattern widened by REQ-079). Point at actions/work-reference.md → Fan-Out Dispatch instead. Offending lines:\n' >&2
  printf '%s\n' "$stale_write_set_premise_lines" >&2
  fail_count=$((fail_count + 1))
fi

# The weak fingerprint of the same premise (REQ-079). REQ-075 named this form as the more dangerous one
# in its own lesson and then pinned only the strong one, so the guard could not catch the recurrence it
# was written for. Naming the exclusive-session model next to write_set is not automatically wrong —
# the model is still true — but it cannot be the REASON write_set is display-only, because that reason
# has to hold at any builder count and the model says nothing about builders.
stale_write_set_weak_premise_lines="$(
  grep -rIEn "$exclusive_session_premise_pattern" \
    "$core_root/actions" "$core_root/docs" "$board_root/tools/queue-kanban" 2>/dev/null \
    | grep -E 'write_set|overlaps' || true
)"
if [ -n "$stale_write_set_weak_premise_lines" ]; then
  printf 'FAIL: write_set is display-only at ANY builder count — no shipped file may justify that by naming the exclusive-session model either (REQ-079). The model is about queue OWNERS, not builders, so it cannot be the reason; point at actions/work-reference.md → Fan-Out Dispatch instead. Offending lines:\n' >&2
  printf '%s\n' "$stale_write_set_weak_premise_lines" >&2
  fail_count=$((fail_count + 1))
fi

# The wrapped-comment half of the same rule (REQ-075, weak form added by REQ-079). model.go and
# board-cards.js explains write_set in comments and tooltip text that wrap across lines, so the line sweeps
# above cannot see them. Neither file has any business asserting how many REQs run at once, nor any
# business invoking the exclusive-session model at all — so file-level negatives are exact and stable.
assert_file_not_contains \
  "tools/queue-kanban/model.go" \
  "$builder_count_premise_pattern" \
  'tools/queue-kanban/model.go must not explain write_set with a one-REQ-at-a-time premise — it is advisory input to a human pick and the merge is the non-interference proof, at any builder count (REQ-075).'

assert_file_not_contains \
  "tools/queue-kanban/web/board-cards.js" \
  "$builder_count_premise_pattern" \
  'tools/queue-kanban/web/board-cards.js overlaps-badge tooltip must not explain write_set with a one-REQ-at-a-time premise — that reason is false since fan-out dispatch (REQ-075).'

assert_file_not_contains \
  "tools/queue-kanban/model.go" \
  "$exclusive_session_premise_pattern" \
  'tools/queue-kanban/model.go must not invoke the exclusive-session model — it is the weak fingerprint of the retired write_set premise, and the file has no other reason to name it (REQ-079).'

assert_file_not_contains \
  "tools/queue-kanban/web/board-cards.js" \
  "$exclusive_session_premise_pattern" \
  'tools/queue-kanban/web/board-cards.js must not invoke the exclusive-session model — it is the weak fingerprint of the retired write_set premise, and the file has no other reason to name it (REQ-079).'

# write_set parser lock-step (REQ-032, updated by REQ-069). write_set does not gate dispatch; it
# survives as a display-only field the board parser reads for the overlaps badge, so the parser must
# still read it. (The reason it never gates is Fan-Out Dispatch's merge-is-the-proof rule, not any
# builder count — see the REQ-075 assertion above.)
assert_contains \
  "tools/queue-kanban/model.go" \
  'fields\["write_set"\]' \
  'tools/queue-kanban/model.go must keep parsing write_set — the display-only field the board overlaps badge reads.'

# Exclusive-session model (REQ-069). The concurrency machinery — the orchestrator lock, the
# parallel-dispatch gate, the co-dispatch re-validations — was removed for a one-session /
# one-active-REQ / one-coder-context operating rule. None of the lock/claim vocabulary may
# survive in any shipped action file, and the two replacement rules must each be stated exactly
# once: a second copy is drift waiting to diverge, the very failure the removed machinery kept
# re-learning across four REQs.
for removed_concurrency_token in \
  'Concurrent-Orchestrator Lock Guard' \
  'coexisting_sessions' \
  'claimed_reqs' \
  'heartbeat_at' \
  'orchestrator-lock\.json'; do
  concurrency_token_hits="$(grep -rIlE -- "$removed_concurrency_token" "$core_root/actions" 2>/dev/null || true)"
  if [ -n "$concurrency_token_hits" ]; then
    printf 'FAIL: removed concurrency-machinery token "%s" still present in a shipped action file (exclusive-session model, REQ-069):\n%s\n' \
      "$removed_concurrency_token" "$concurrency_token_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

# Reservation removal. The reserve action allocated REQs to a DIFFERENT worktree/cloud session, which
# the then-current exclusive-session model declared outside the product contract. REQ-096 re-grained
# that model to claim-anywhere, so the IDEA is now in contract — but the machinery is not, and this
# ratchet does not loosen: what stays dead is the reserve VERB, the `reserved` STATUS, and the
# frontmatter fields, because a claim is marked with the advisory `assigned_to` field instead (no
# staleness clock, no router-budget cost). None of the retired vocabulary may survive in shipped prose
# or the board tool. Tokens are underscore/path-shaped on purpose so ordinary English "reserved for"
# never false-positives — and `assigned_to` deliberately matches none of them.
for removed_reservation_token in \
  'status: reserved' \
  'reserved_for' \
  'reserved_at' \
  'do-work reserve' \
  'actions/reserve\.md'; do
  reservation_token_hits="$(grep -rIlE -- "$removed_reservation_token" \
    "$core_root/actions" "$core_root/docs" "$board_root/tools/queue-kanban" \
    "$core_root/SKILL.md" "$core_root/next-steps.md" 2>/dev/null || true)"
  # Per-file exemption (same pattern as the maintainer-doc allowlist): the update
  # flow in actions/version.md must NAME the removed reserve files on its Step 5
  # deletion line — tar extraction never deletes what upstream removed, so without
  # that rm every pre-0.161.0 install keeps the orphaned files forever. Only lines
  # containing `rm -f` are exempt; any other mention in version.md still fails.
  exempt_update_flow_file="$core_root/actions/version.md"
  filtered_reservation_hits=""
  for reservation_hit_file in $reservation_token_hits; do
    if [ "$reservation_hit_file" = "$exempt_update_flow_file" ] && \
       [ -z "$(grep -E -- "$removed_reservation_token" "$reservation_hit_file" | grep -vE -- 'rm -f' || true)" ]; then
      continue
    fi
    filtered_reservation_hits="$filtered_reservation_hits$reservation_hit_file
"
  done
  reservation_token_hits="$(printf '%s' "$filtered_reservation_hits")"
  if [ -n "$reservation_token_hits" ]; then
    printf 'FAIL: removed reservation-workflow token "%s" still present in a shipped file — the reserve verb/status stay dead; the advisory assigned_to field is how a claim is marked (REQ-096/REQ-097):\n%s\n' \
      "$removed_reservation_token" "$reservation_token_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

assert_file_missing \
  "actions/reserve.md" \
  'the reserve action must stay removed — a claim is marked with the advisory assigned_to field, not a reserve verb or a reserved status (REQ-096 re-grained ownership to claim-anywhere without reviving either).'

assert_file_missing \
  "actions/prime-req-reservation.md" \
  'the reservation prime doc must stay removed along with the reserve action.'

# REQ-298 — a size query's failure must never arrive as a size.
#
# THE CONDITION, stated where shipped shell is governed
# (`_dev/primes/prime-shell-commands.md` § Unchecked Exit Status): a command
# substitution whose exit status is discarded while only its content is judged
# lets a tool that never ran read as a tool that found nothing.
#
# WHY THIS CHECK IS NARROW, and it is a finding rather than a compromise. The
# broad shape — any substitution with a `|| true` / `|| echo 0` fallback — was
# measured across every shipped script: 15 sites, and every one of them is a
# CORRECT use. `ps -o pgid=` returning empty means the process is gone; `git
# rev-parse --show-toplevel` returning empty means this is not a repository;
# `mktemp` returning empty is guarded and reported on the next line. In each the
# emptiness IS the information, so a check on that shape would flag fifteen
# correct lines to catch zero defects and would be muted within a week.
#
# What separates the defect from the correct use is semantic — does a guard then
# make a safety decision on the value, where the collapsed value silently
# satisfies the safe branch — and that is not greppable. So this pins the one
# query where the distinction is unambiguous: `git cat-file -s` answers "how big
# is this blob", the answer is only ever used in a size guard, and a failure
# collapsed to empty or 0 is exactly the incident (a 13,900-byte archived REQ
# truncated to 57 bytes passed the truncation floor and was written and staged).
# A display-only caller may fall back to '?', which no reader can mistake for a
# size. Anything numeric or empty may not.
size_query_collapse_hits="$(grep -rnE "git cat-file -s[^|]*\|\|[[:space:]]*(true|echo[[:space:]]+0|echo[[:space:]]+\"\"|echo[[:space:]]+'')" \
  --include="*.sh" "$repo_root/skills" 2>/dev/null | grep -v '_test' || true)"
if [ -n "$size_query_collapse_hits" ]; then
  printf 'FAIL: a `git cat-file -s` discards its exit status into a value indistinguishable from a real size (REQ-298) — check the blob exists first with `git cat-file -e` and refuse when an existing blob will not size, or fall back to a token no reader can mistake for a size:\n%s\n' \
    "$size_query_collapse_hits" >&2
  fail_count=$((fail_count + 1))
fi


# Impact/effort separation (REQ-289). `effort_estimate` had two writers with two
# different meanings: capture judged SIZE, while review MUST-stamped it from an
# IMPACT gate. The split gave each axis its own field, and every value on both
# axes a token that is unique repo-wide by plain-text search.

# Check A — no action file derives effort_estimate from an impact judgment. This
# is the defect the REQ names: review-work.md's follow-up step and
# work-reference.md's Discovered Tasks Classification both stamped the SIZE field
# from an IMPACT verdict, which is what let work.md's mechanical-effort
# short-circuit forecast a three-hour fix at five minutes.
#
# The pin is the derivation RELATION applied to this field, not one verb and not
# mere co-occurrence. After the split only `impact:` is ever derived from a
# verdict — `effort_estimate` is judged as size by whoever writes it.
#
# REQ-293 widened this on both axes, because the original pinned a spelling:
#   * the VERB SET was the single literal `stamp*`, so "derive", "set from",
#     "comes from" and "map to" all walked past. Measured at the time: 1 of 6
#     realistic re-drifts caught. Worse, the mutation test that "confirmed" it
#     used the word "stamping" — the one verb it greps — so the test was
#     self-confirming.
#   * the PROXIMITY WINDOW was `[^.]{0,80}`, which treats any period as a
#     sentence end. A file path between the verb and the field breaks the window,
#     and a cited path is this repo's dominant style — so the very sentences most
#     likely to carry the defect were the ones it could not see.
# The window is now "same line, either order", which is what the enclosing loop
# already iterates by, and the verb set is a class.
#
# `gate` is the retired routing field and is dead beside `effort_estimate` in any
# spelling.
if ! python3 - "$core_root/actions" "$toolbox_root/actions" "$core_root/docs" "$core_root/crew-members" \
  "$toolbox_root/actions" "$knowledge_root/actions" "$board_root/actions" \
  "$core_root/SKILL.md" "$toolbox_root/SKILL.md" "$knowledge_root/SKILL.md" "$board_root/SKILL.md" <<'PY'
import pathlib
import re
import sys

# A derivation is any of these applied to effort_estimate, in either order on the
# line. The set is the CLASS of "this field's value comes from that judgment",
# deliberately not one verb — a re-drift is written by whoever writes it next, in
# their own words.
derivation_verbs = (
    r"stamp\w*|deriv\w*|set\s+(?:it\s+)?from|comes?\s+from|taken\s+from|read\s+from"
    r"|map(?:s|ped|ping)?\s+(?:it\s+)?(?:to|from)|infer\w*|follow\w*\s+from|based\s+on"
    r"|translat\w*|copie[sd]?\s+from|mirror\w*"
)
# The window is `.{0,60}` — bounded, but NOT `[^.]`. Excluding the period was the
# original bug: a cited file path between the verb and the field ended the window
# early, and a cited path is this repo's dominant style. Sixty characters is
# tight enough that a long schema line describing the separation correctly does
# not trip it, and wide enough to span "derived from the review gate's recorded
# disposition into effort_estimate".
derivation_pattern = re.compile(
    r"(?:" + derivation_verbs + r").{0,60}effort_estimate"
    r"|effort_estimate.{0,60}(?:" + derivation_verbs + r")"
    r"|\bgate\b",
    flags=re.IGNORECASE,
)
# An impact verdict named on the same line is what makes a derivation a
# CROSS-AXIS one. Without it, "effort_estimate is set from your own judgment of
# size" would fail, which is the correct behaviour being described.
impact_axis_mention = re.compile(r"impact[-_ ]|verdict|gate", flags=re.IGNORECASE)
negated_derivation = re.compile(
    r"\b(?:never|not|no longer|rather than|instead of|cannot|must not|is a different axis)\b",
    flags=re.IGNORECASE,
)
offending_lines = []
for scan_root in sys.argv[1:]:
    root_path = pathlib.Path(scan_root)
    scanned_files = [root_path] if root_path.is_file() else sorted(root_path.glob("*.md"))
    for action_file in scanned_files:
        if not action_file.is_file():
            continue
        for line_number, line in enumerate(action_file.read_text().splitlines(), start=1):
            if "effort_estimate" not in line:
                continue
            if not derivation_pattern.search(line):
                continue
            if not impact_axis_mention.search(line):
                continue
            # A NEGATED derivation is the rule being stated, not broken: the
            # action files say "effort_estimate is never derived from that token"
            # in as many words, and that sentence is the contract. Skipping it is
            # not a hole — a re-drift asserts the derivation, it does not deny it.
            if negated_derivation.search(line):
                continue
            offending_lines.append(f"  {action_file}:{line_number}")
if offending_lines:
    raise SystemExit("\n".join(offending_lines))
PY
then
  printf 'FAIL: an action file stamps or gates effort_estimate (REQ-289) — judge effort as effort, and stamp impact: from the two questions in actions/review-work.md Step 10. Offending lines listed above.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Check B leg 1 — the six tokens are defined, not aspirational. A vocabulary the
# schema names but no action file uses is a rename that never landed.
for impact_effort_token in \
  'impact-critical' \
  'impact-user-visible' \
  'impact-rule-change' \
  'impact-negligible' \
  'effort-mechanical' \
  'effort-substantive'; do
  if ! grep -rqI -- "$impact_effort_token" "$core_root/actions"; then
    printf 'FAIL: impact/effort token "%s" appears nowhere under skills/do-work/actions (REQ-289) — the split vocabulary must be defined where it is written, not only described.\n' \
      "$impact_effort_token" >&2
    fail_count=$((fail_count + 1))
  fi
done

# Check B leg 2 — the bare axis words never survive un-prefixed. Without this a
# stray `gate: user-visible` re-enters and one grep starts returning two axes
# again, which is the whole reason the tokens were renamed.
unprefixed_axis_word_hits="$(grep -rInoE -- '(^|[^-])(user-visible|rule-change)' \
  "$core_root/actions" "$toolbox_root/actions" "$board_root/tools/queue-kanban" 2>/dev/null || true)"
if [ -n "$unprefixed_axis_word_hits" ]; then
  printf 'FAIL: the axis words user-visible/rule-change appear without the impact- prefix (REQ-289) — every value on both axes must be a repo-unique token, or two axes collapse back under one grep:\n%s\n' \
    "$unprefixed_axis_word_hits" >&2
  fail_count=$((fail_count + 1))
fi

# Check C — the read-only legacy aliases resolve on BOTH sides. Over forty
# archived REQs carry the literal `effort_estimate: trivial|normal`; dropping the
# aliases would invalidate every one of them, which the REQ forbids.
assert_contains \
  "actions/work-reference.md" \
  '`trivial` → `effort-mechanical`' \
  'actions/work-reference.md Schema Read Contract must carry the read-only alias `trivial` → `effort-mechanical` (REQ-289) — without it every archived REQ written before the rename stops resolving.'
assert_contains \
  "actions/work-reference.md" \
  '`normal` → `effort-substantive`' \
  'actions/work-reference.md Schema Read Contract must carry the read-only alias `normal` → `effort-substantive` (REQ-289) — without it every archived REQ written before the rename stops resolving.'

if ! python3 - "$board_root/tools/queue-kanban/model.go" <<'PY'
import pathlib
import re
import sys

model_source = pathlib.Path(sys.argv[1]).read_text()
constant_values = dict(
    re.findall(r'^\s*(\w+)\s*=\s*"((?:impact|effort)-[a-z-]+)"\s*$', model_source, flags=re.MULTILINE)
)


def schema_entry_body(field_name):
    entry_match = re.search(
        r'"' + field_name + r'":\s*\{(.*?)\n\t\},', model_source, flags=re.DOTALL
    )
    if entry_match is None:
        raise SystemExit(f"no schemaReadContractFields entry for {field_name!r}")
    return entry_match.group(1)


effort_entry = schema_entry_body("effort_estimate")
alias_pairs = {
    alias: constant_values.get(target, target)
    for alias, target in re.findall(r'"([a-z_]+)":\s*(\w+|"[a-z-]+")', effort_entry)
}
alias_pairs = {alias: target.strip('"') for alias, target in alias_pairs.items()}
missing = [
    f"{alias} -> {expected}"
    for alias, expected in (("trivial", "effort-mechanical"), ("normal", "effort-substantive"))
    if alias_pairs.get(alias) != expected
]
if missing:
    raise SystemExit("effort_estimate aliases missing or wrong: " + ", ".join(missing))
PY
then
  printf 'FAIL: tools/queue-kanban/model.go effort_estimate entry must map the read-only legacy aliases trivial -> effort-mechanical and normal -> effort-substantive (REQ-289) — the board is the reader that makes an archived REQ still resolve.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Check D — the board parser and the schema row agree in BOTH directions.
# actions/work-reference.md tells its own impact: line to keep the parser in
# lock-step "both changing in the same commit"; until now that sentence had no
# enforcement, so either side could drift alone and stay green.
if ! python3 - "$core_root/actions/work-reference.md" "$board_root/tools/queue-kanban/model.go" <<'PY'
import pathlib
import re
import sys

schema_text = pathlib.Path(sys.argv[1]).read_text()
model_source = pathlib.Path(sys.argv[2]).read_text()

schema_line = next(
    (line for line in schema_text.splitlines() if line.startswith("impact: ")), None
)
if schema_line is None:
    raise SystemExit("actions/work-reference.md has no `impact:` schema line")
schema_tokens = set(re.findall(r"impact-[a-z-]+", schema_line))

constant_values = dict(
    re.findall(r'^\s*(\w+)\s*=\s*"(impact-[a-z-]+)"\s*$', model_source, flags=re.MULTILINE)
)
impact_entry = re.search(r'"impact":\s*\{(.*?)\n\t\},', model_source, flags=re.DOTALL)
if impact_entry is None:
    raise SystemExit('model.go has no schemaReadContractFields entry for "impact"')
canonical_match = re.search(
    r"canonicalValues:\s*\[\]string\{([^}]*)\}", impact_entry.group(1)
)
if canonical_match is None:
    raise SystemExit('model.go "impact" entry declares no canonicalValues')
parser_tokens = {
    constant_values.get(item.strip(), item.strip().strip('"'))
    for item in canonical_match.group(1).split(",")
    if item.strip()
}

# REQ-293 F4 — the DEFAULT, not only the token set. Check D compared the two
# vocabularies and never read defaultValue, so nothing held the schema line's
# "Absent or unrecognized reads as `impact-user-visible`" to anything.
#
# This is the highest-value pin in the REQ because REQ-290 depends on it: if the
# default ever became `impact-negligible`, `do-work run --skip-impact-negligible`
# would invert into "skip everything, including every REQ predating the field" —
# and every check in this file would stay green. Absence must never be mistakable
# for the user's stop signal.
default_match = re.search(r'defaultValue:\s*(\w+|"[^"]+")', impact_entry.group(1))
if default_match is None:
    raise SystemExit('model.go "impact" entry declares no defaultValue')
parser_default = constant_values.get(
    default_match.group(1), default_match.group(1).strip('"')
)
if parser_default != "impact-user-visible":
    raise SystemExit(
        f'model.go\'s impact defaultValue is "{parser_default}", not "impact-user-visible" — '
        "an absent or unrecognized impact must never read as the user's stop signal, "
        "or --skip-impact-negligible inverts into skip-everything (REQ-290)"
    )
# ...and the schema line has to say the same thing, in the file a reader reads.
if not re.search(r"[Aa]bsent or unrecognized reads as `impact-user-visible`", schema_line):
    raise SystemExit(
        "the impact schema line no longer states that absent or unrecognized reads as "
        "`impact-user-visible` — the parser default and the documented default must agree, "
        "and this is the half a reader acts on"
    )

if schema_tokens != parser_tokens:
    raise SystemExit(
        "impact vocabulary disagrees — only in the schema row: "
        f"{sorted(schema_tokens - parser_tokens)}; only in model.go: "
        f"{sorted(parser_tokens - schema_tokens)}"
    )
PY
then
  printf 'FAIL: the impact vocabulary in actions/work-reference.md and tools/queue-kanban/model.go disagree (REQ-289) — the schema line requires the parser in lock-step, both changing in the same commit.\n' >&2
  fail_count=$((fail_count + 1))
fi

# The retired impact vocabularies. `gate:` had its own three-word vocabulary and
# the Discovered Tasks ladder had a third; both collapse onto the four impact-
# tokens, so no shipped file may keep either alive. Narrowing this loop to make a
# straggler pass is exactly how a retired vocabulary survives.
# REQ-293 F3: the ladder tokens are matched by TOKEN, not by markup. The
# original patterns required bold — `**[critical]**` — so `- [low] a style nit`
# and the backticked `` `[critical]` `` both walked past clean, and the
# backticked form is one the tree actually carried (review-work.md:201 before
# REQ-289). The bracketed word IS the retired vocabulary; whatever emphasis
# surrounds it is incidental. `[ ]`-style checkboxes are unaffected: the tokens
# are words, not spaces.
for retired_impact_token in \
  'gate: user-visible' \
  'gate: rule-change' \
  'gate: trivial' \
  '\[critical\]' \
  '\[normal\]' \
  '\[low\]'; do
  retired_impact_hits="$(grep -rIlE -- "$retired_impact_token" \
    "$core_root/actions" "$core_root/docs" "$toolbox_root/actions" "$board_root/tools/queue-kanban" 2>/dev/null || true)"
  if [ -n "$retired_impact_hits" ]; then
    printf 'FAIL: retired impact vocabulary "%s" still present in a shipped file (REQ-289 — the gate: field and the [critical]/[normal]/[low] ladder both collapse onto the four impact- tokens):\n%s\n' \
      "$retired_impact_token" "$retired_impact_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

# REQ-293 F6 — `--skip-impact-negligible` is declared in several places that must
# agree, and nothing held them together. This is not hypothetical: REQ-290's own
# review found THREE already-stale restatements of the ready-set conditions
# inside these same two files, one of them thirteen lines from the condition it
# contradicted.
#
# The pin is presence at each declaration site, in the file that owns it. A site
# that loses the flag is a reader that will not know about it; a usage string
# that loses it tells the user the flag does not exist.
if ! skip_negligible_missing="$(python3 - "$core_root/actions/work.md" "$core_root/actions/work-reference.md" <<'SKIPSITES'
import pathlib
import sys

work_text = pathlib.Path(sys.argv[1]).read_text()
reference_text = pathlib.Path(sys.argv[2]).read_text()

# Each entry is (where it lives, which file, a phrase unique to that site).
required_sites = [
    ("work.md `## Input` bullet", work_text, "**`--skip-impact-negligible`** (boolean flag)"),
    ("work.md argument-strip list", work_text,
     "After stripping `--wave N`, `--fan-out [N]`, and `--skip-impact-negligible`"),
    ("work.md usage string, default branch", work_text,
     "do-work run [REQ-NNN|UR-NNN ...] [--skip-impact-negligible]"),
    ("work.md usage string, --wave branch", work_text,
     "do-work run --wave N [--skip-impact-negligible]"),
    ("work.md Step 1 skip paragraph", work_text,
     "**`--skip-impact-negligible` skips negligible REQs and reports them, never silently.**"),
    ("work.md Orchestrator Checklist Step 0", work_text,
     "--skip-impact-negligible); reject unrecognized residue"),
    ("work-reference.md auto-wave condition 5", reference_text,
     "**Not dropped by `--skip-impact-negligible`**"),
]
missing = [where for where, text, phrase in required_sites if phrase not in text]
if missing:
    raise SystemExit("\n".join("  " + where for where in missing))
SKIPSITES
)"; then
  printf 'FAIL: --skip-impact-negligible lost a declaration site (REQ-290, pinned by REQ-293 F6) — every site below must agree, or a reader acts on a flag one of them does not mention:\n%s\n' \
    "$skip_negligible_missing" >&2
  fail_count=$((fail_count + 1))
fi

# The same shape for the impact title tag's emitter set: every flow that mints a
# REQ title must write the `[<impact token>] ` tag, or the board's title-matching
# search box cannot find a non-default verdict — which is the whole reason the
# tag exists beside the field.
for title_tag_emitter in \
  "$core_root/actions/capture-reference.md" \
  "$core_root/actions/review-work.md" \
  "$core_root/actions/work-reference.md"; do
  if ! grep -qF -- '[<impact token>] ' "$title_tag_emitter"; then
    printf 'FAIL: %s no longer names the `[<impact token>] ` title tag (REQ-290, pinned by REQ-293 F6) — a follow-up carrying the field but not the tag is invisible to the board title search.\n' \
      "$title_tag_emitter" >&2
    fail_count=$((fail_count + 1))
  fi
done


assert_contains \
  "actions/work-reference.md" \
  '^## Execution Model — Claim Anywhere, One Releaser' \
  'actions/work-reference.md must define the Execution Model section that replaced the orchestrator-lock/parallel-dispatch machinery. Renamed from "Exclusive Session" by REQ-096: claiming is no longer exclusive to one checkout, only the release tail is, and a heading naming the retired boundary is exactly the stale-name drift this suite exists to catch.'

# The invariant is about OWNERSHIP, not build count (REQ-073), and REQ-096 moved which
# ownership it names. REQ-069's wording ("one active REQ, one coder context") conflated
# ownership with build count, so lifting the builder cap required rewording it to
# "one queue owner per checkout"; REQ-096 then re-grained the model to claim-anywhere,
# which makes the *claiming* half false and leaves the release tail as the only thing
# that must not run twice. So the invariant is now "one releaser per queue" — still
# exactly once, because a second copy is drift waiting to diverge.
# `|| true` is load-bearing under `set -euo pipefail`: grep exits 1 on no match, and
# with pipefail that aborts the whole suite silently — a missing invariant would read
# as a crash with no FAIL line rather than as the failure it is.
releaser_invariant_count="$( { grep -roh 'one releaser per queue' "$core_root/actions" || true; } | wc -l | tr -d ' ')"
if [ "$releaser_invariant_count" != "1" ]; then
  printf 'FAIL: the ownership invariant ("one releaser per queue") must be stated exactly once across actions/ (found %s) — every other mention is a pointer, not a restatement.\n' \
    "$releaser_invariant_count" >&2
  fail_count=$((fail_count + 1))
fi

# The superseded wording must be gone, not merely outnumbered — same ratchet the
# "one active REQ, one coder context" check below applies to its predecessor. It bounds
# claiming to a single checkout, which is exactly what claim-anywhere makes false.
superseded_owner_invariant_hits="$(grep -rIlE -- 'one queue owner per checkout' "$core_root/actions" "$core_root/docs" "$core_root/SKILL.md" 2>/dev/null || true)"
if [ -n "$superseded_owner_invariant_hits" ]; then
  printf 'FAIL: the superseded invariant wording "one queue owner per checkout" still appears (REQ-096 replaced it with the one-releaser-per-queue formulation — any checkout may claim):\n%s\n' \
    "$superseded_owner_invariant_hits" >&2
  fail_count=$((fail_count + 1))
fi

# The old wording must be gone, not merely outnumbered: it says one active REQ and
# one coder context, which is exactly what fan-out makes false.
retired_invariant_hits="$(grep -rIlE -- 'one active REQ, one coder context' "$core_root/actions" "$core_root/docs" "$core_root/SKILL.md" 2>/dev/null || true)"
if [ -n "$retired_invariant_hits" ]; then
  printf 'FAIL: the retired invariant wording "one active REQ, one coder context" still appears (REQ-073 replaced it with the one-queue-owner formulation):\n%s\n' \
    "$retired_invariant_hits" >&2
  fail_count=$((fail_count + 1))
fi

# Fan-out dispatch (REQ-073). The builder cap was two sentences in the Worktree
# Dispatch Mode opening; both must stay gone, and the section that replaced them
# must keep the three things that make N builders safe without coordination: the
# human picks, the merge is the proof, and a named set of steps never parallelises.
for retired_builder_cap_phrase in \
  'The single active builder' \
  'only one builder is ever in flight'; do
  builder_cap_hits="$(grep -rIlE -- "$retired_builder_cap_phrase" "$core_root/actions" 2>/dev/null || true)"
  if [ -n "$builder_cap_hits" ]; then
    printf 'FAIL: the retired builder-cap phrase "%s" still appears — REQ-073 raised worktree dispatch from one builder to N under a single queue owner:\n%s\n' \
      "$retired_builder_cap_phrase" "$builder_cap_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

assert_contains \
  "actions/work-reference.md" \
  '\*\*Fan-Out Dispatch' \
  'actions/work-reference.md must define Fan-Out Dispatch inside Worktree Dispatch Mode — several builders under one releaser, with no new coordination state.'

fan_out_block="$(sed -n '/\*\*Fan-Out Dispatch/,/^## Composed Exit Summary/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$fan_out_block" \
  'Serial-only' \
  'Fan-Out Dispatch must name what never parallelises — queue transitions, REQ id allocation, and the version/changelog files.'

assert_block_contains \
  "$fan_out_block" \
  'CHANGELOG\.md' \
  'the Serial-only list must name CHANGELOG.md explicitly — one entry per REQ written by the owner, because unique version numbers do not make a shared prepend safe.'

assert_block_contains \
  "$fan_out_block" \
  'advisory input|never a gate' \
  'Fan-Out Dispatch must keep write_set advisory input to the humans pick and never a gate — nothing schedules on it under any builder count.'

assert_block_contains \
  "$fan_out_block" \
  'line proximity, not meaning' \
  'Fan-Out Dispatch must state the merge gates honest limit: git detects conflicts by line proximity, so two REQs appending to a shared registry merge cleanly and can still be jointly wrong.'

assert_block_contains \
  "$fan_out_block" \
  'survivable, not prevented' \
  'Fan-Out Dispatch must carry crew-members/background-agents.md own ceiling note — the run-directory pattern makes fan-out failures survivable, never prevented.'

assert_block_contains \
  "$fan_out_block" \
  'absolute main-tree path' \
  'Fan-Out Dispatch must state the brief-delivery trap: a repo-relative path resolves inside the worktree against its own stale copy of do-work/.'

# Overlapping parallel-writer isolation and hand-back (0.186.37). The shared guide's
# former "reach for it" wording made worktrees optional and said nothing about bringing
# completed branches home, so overlapping writers could interleave edits or strand
# accepted work on side branches. Pin the trigger, its non-trigger cases, and the full
# safe hand-back in every package-local copy that fan-out actions load.
for background_agents_path in \
  "$core_root/crew-members/background-agents.md" \
  "$knowledge_root/crew-members/background-agents.md" \
  "$toolbox_root/crew-members/background-agents.md"
do
  background_worktree_block="$(sed -n '/^\*\*Worktree isolation is a separate axis\./,/^## Manifest Format/p' "$background_agents_path")"

  assert_block_contains \
    "$background_worktree_block" \
    'explicitly declared' \
    "$background_agents_path must key worktree isolation to explicit ownership declarations, not guessed overlap."

  assert_block_contains \
    "$background_worktree_block" \
    'file lists or globs overlap' \
    "$background_agents_path must make declared file-list or glob overlap the worktree trigger."

  assert_block_contains \
    "$background_worktree_block" \
    'worktree\*\* on a separate branch before it writes' \
    "$background_agents_path must isolate every overlapping writer before the first write."

  assert_block_contains \
    "$background_worktree_block" \
    'Read-only agents' \
    "$background_agents_path must leave read-only fan-out outside the worktree trigger."

  assert_block_contains \
    "$background_worktree_block" \
    'declarations are disjoint need no extra worktree' \
    "$background_agents_path must leave declared-disjoint writers outside the worktree trigger."

  assert_block_contains \
    "$background_worktree_block" \
    'missing declarations' \
    "$background_agents_path must not infer overlap from a missing declaration."

  assert_block_contains \
    "$background_worktree_block" \
    'remains display-only' \
    "$background_agents_path must not promote write_set into a scheduling or isolation gate."

  assert_block_contains \
    "$background_worktree_block" \
    'integrates one branch at a time' \
    "$background_agents_path must serialize completed branch integration."

  assert_block_contains \
    "$background_worktree_block" \
    'text conflicts, integration seams, and semantic composition' \
    "$background_agents_path must reconcile both textual and semantic overlap during integration."

  assert_block_contains \
    "$background_worktree_block" \
    'merged state before starting the next merge' \
    "$background_agents_path must verify each integrated state before merging the next branch."

  assert_block_contains \
    "$background_worktree_block" \
    'preserve and report the branch and worktree' \
    "$background_agents_path must preserve unsafe work for recovery instead of stranding or deleting it."

  assert_block_contains \
    "$background_worktree_block" \
    'never force-merge' \
    "$background_agents_path must forbid force-merging or force-deleting an unsafe hand-back."
done

work_implementation_block="$(sed -n '/^### Step 6: Implementation/,/^### Step 6\.25:/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_implementation_block" \
  '^\*\*Overlapping parallel writers:\*\*' \
  'actions/work.md Step 6 must surface the overlapping-writer worktree rule where implementation agents are dispatched.'

assert_block_contains \
  "$work_implementation_block" \
  'hand every completed branch back for serial reconciliation and merged-state verification' \
  'actions/work.md Step 6 must require completed parallel-writer branches to return through verified serial integration.'

assert_block_contains \
  "$work_implementation_block" \
  'one worktree per builder regardless of overlap' \
  'actions/work.md Step 6 must preserve auto-waves stronger worktree-per-builder rule.'

# REQ-309 — focused tests answer whether the changed area regressed; they do not prove
# that a repository's declared whole-repo pass/fail contract is green. Pin the added
# Step 6.5 lane and downstream Error Handling row by meaning, then mutate each behavior
# that made REQ-283's hand-back possible so a broad vocabulary match cannot satisfy the
# check. Extract the two sections independently so one cannot lend vocabulary to the other.
if ! python3 - "$core_root/actions/work.md" <<'PY'
import pathlib
import re
import sys

work_text = pathlib.Path(sys.argv[1]).read_text()


def extract_canonical_gate_lane(source):
    testing_match = re.search(
        r"^### Step 6\.5: Testing\n(?P<body>.*?)^### Step 7: Review$",
        source,
        flags=re.DOTALL | re.MULTILINE,
    )
    if testing_match is None:
        return None
    lane_match = re.search(
        r"^\d+\. \*\*Run (?:the )?declared canonical repository gate\*\*.*?"
        r"(?=^\d+\. \*\*|\Z)",
        testing_match.group("body"),
        flags=re.DOTALL | re.MULTILINE | re.IGNORECASE,
    )
    return None if lane_match is None else lane_match.group(0)


def extract_repeated_failure_row(source):
    error_handling_match = re.search(
        r"^## Error Handling\n(?P<body>.*?)^## Progress Reporting$",
        source,
        flags=re.DOTALL | re.MULTILINE,
    )
    if error_handling_match is None:
        return None
    for line in error_handling_match.group("body").splitlines():
        if line.startswith("|") and "fail repeatedly" in line.lower():
            return line
    return None


def canonical_gate_defects(source):
    lane = extract_canonical_gate_lane(source)
    if lane is None:
        return {"missing canonical-gate lane"}

    normalized = " ".join(lane.lower().split())
    predicates = {
        "explicit project-guidance trigger":
            r"if project guidance explicitly declares a canonical repository-wide pass/fail gate",
        "additive to focused tests": r"in addition to (?:the )?focused tests",
        "project-root execution": r"from (?:the )?project root",
        "final or merged-state execution":
            r"final implementation state.*(?:post-merge|merged) state",
        "direct status verdict": r"direct exit status",
        "zero required": r"(?:must|required? to|require) (?:be )?zero",
        "no focused-test baseline exemption":
            r"focused-test baseline (?:exclusions|exemptions).*do not apply",
        "current-diff remediation":
            r"failure (?:is )?caused by the current diff.*existing remediation loop",
        "unrelated or pre-existing stop path":
            r"(?:unrelated or pre-existing|pre-existing or unrelated) failure",
        "claim and checkpoint preserved":
            r"preserve the claimed req and (?:its )?checkpoint",
        "stop before successful completion":
            r"stop before successful archive, commit, or hand-back",
    }
    return {
        defect_name
        for defect_name, predicate in predicates.items()
        if re.search(predicate, normalized) is None
    }


def repeated_failure_defects(source):
    row = extract_repeated_failure_row(source)
    if row is None:
        return {"missing repeated-failure row"}

    columns = [column.strip() for column in row.strip().strip("|").split("|")]
    if len(columns) != 2:
        return {"malformed repeated-failure row"}
    trigger = " ".join(columns[0].lower().split())
    action = " ".join(columns[1].lower().split())
    predicates = {
        "focused-test trigger": (trigger, r"focused tests?"),
        "current-diff gate trigger":
            (trigger, r"current-diff canonical repository gate"),
        "three-attempt remediation": (action, r"after 3 fix attempts"),
        "Code classification": (action, r"classify as code failure"),
        "follow-up failure details":
            (action, r"create (?:a )?follow-up req with .*failure details"),
        "failed archive path": (action, r"archive as failed"),
        "canonical hold exception":
            (action, r"excludes an unrelated or pre-existing canonical repository gate failure"),
        "claim and checkpoint preserved":
            (action, r"preserve the claimed req and its checkpoint"),
        "never archive hold": (action, r"never archive (?:it|the req)"),
        "Step 6.5 hold owner": (action, r"step 6\.5 owns (?:that|the) hold"),
    }
    defects = {
        defect_name
        for defect_name, (text, predicate) in predicates.items()
        if re.search(predicate, text) is None
    }
    if trigger == "tests fail repeatedly":
        defects.add("broad repeated-test trigger")
    return defects


live_defects = canonical_gate_defects(work_text)
if live_defects:
    raise SystemExit(
        "actions/work.md Step 6.5 canonical repository gate contract is incomplete: "
        + ", ".join(sorted(live_defects))
    )
live_lane = extract_canonical_gate_lane(work_text)

mutations = (
    ("implicit declaration", "explicitly declares", "implicitly suggests", "explicit project-guidance trigger"),
    ("focused tests replaced", "in addition to focused tests", "instead of focused tests", "additive to focused tests"),
    ("wrong working directory", "from the project root", "from the changed area's directory", "project-root execution"),
    ("pre-final state", "final implementation state", "pre-flight state", "final or merged-state execution"),
    ("summary verdict", "direct exit status", "reported summary", "direct status verdict"),
    ("nonzero accepted", "must be zero", "may be nonzero", "zero required"),
    ("baseline exemption restored", "do not apply", "still apply", "no focused-test baseline exemption"),
    ("current diff waived", "use the existing remediation loop", "waive it", "current-diff remediation"),
    ("unrelated failure omitted", "unrelated or pre-existing failure", "environmental note", "unrelated or pre-existing stop path"),
    ("claim discarded", "preserve the claimed REQ and its checkpoint", "discard the claimed REQ and its checkpoint", "claim and checkpoint preserved"),
    ("hand-back allowed", "stop before successful archive, commit, or hand-back", "continue through successful archive, commit, and hand-back", "stop before successful completion"),
)

for mutation_name, old, new, expected_defect in mutations:
    mutated_lane = live_lane.replace(old, new, 1)
    if mutated_lane == live_lane:
        raise SystemExit(
            f"canonical repository gate mutation {mutation_name!r} changed nothing"
        )
    mutated_text = work_text.replace(live_lane, mutated_lane, 1)
    mutation_defects = canonical_gate_defects(mutated_text)
    if expected_defect not in mutation_defects:
        raise SystemExit(
            f"canonical repository gate mutation {mutation_name!r} escaped "
            f"{expected_defect!r}; found {sorted(mutation_defects)!r}"
        )

live_error_defects = repeated_failure_defects(work_text)
if live_error_defects:
    raise SystemExit(
        "actions/work.md Error Handling repeated-failure contract is incomplete: "
        + ", ".join(sorted(live_error_defects))
    )
live_error_row = extract_repeated_failure_row(work_text)

error_mutations = (
    ("broad trigger restored", "Focused tests or a current-diff canonical repository gate fail repeatedly", "Tests fail repeatedly", "broad repeated-test trigger"),
    ("focused trigger removed", "Focused tests or a current-diff canonical repository gate", "Tests or a current-diff canonical repository gate", "focused-test trigger"),
    ("current-diff trigger removed", "current-diff canonical repository gate", "canonical repository gate", "current-diff gate trigger"),
    ("three attempts removed", "After 3 fix attempts", "After remediation", "three-attempt remediation"),
    ("Code classification changed", "classify as Code failure", "classify as Environment failure", "Code classification"),
    ("follow-up failure details removed", "create a follow-up REQ with the focused-test or current-diff gate failure details", "note the failure", "follow-up failure details"),
    ("failed archive path removed", "archive as failed", "stop processing", "failed archive path"),
    ("canonical hold exception removed", "excludes an unrelated or pre-existing canonical repository gate failure", "also consumes an unrelated or pre-existing canonical repository gate failure", "canonical hold exception"),
    ("claim and checkpoint discarded", "preserve the claimed REQ and its checkpoint", "discard the claimed REQ and its checkpoint", "claim and checkpoint preserved"),
    ("never-archive hold inverted", "never archive it", "always archive it", "never archive hold"),
    ("never-archive hold removed", " and never archive it", "", "never archive hold"),
    ("Step 6.5 hold owner removed", "Step 6.5 owns that hold", "the table owns that hold", "Step 6.5 hold owner"),
)

for mutation_name, old, new, expected_defect in error_mutations:
    mutated_row = live_error_row.replace(old, new, 1)
    if mutated_row == live_error_row:
        raise SystemExit(
            f"Error Handling mutation {mutation_name!r} changed nothing"
        )
    mutated_text = work_text.replace(live_error_row, mutated_row, 1)
    mutation_defects = repeated_failure_defects(mutated_text)
    if expected_defect not in mutation_defects:
        raise SystemExit(
            f"Error Handling mutation {mutation_name!r} escaped "
            f"{expected_defect!r}; found {sorted(mutation_defects)!r}"
        )
PY
then
  printf 'FAIL: actions/work.md Step 6.5 and Error Handling must agree on the hold for an unrelated or pre-existing canonical repository gate failure (REQ-309/REQ-317).\n' >&2
  fail_count=$((fail_count + 1))
fi

three_attempt_count="$(grep -roh 'consecutive fix attempts' "$core_root/actions" | wc -l | tr -d ' ')"
if [ "$three_attempt_count" != "1" ]; then
  printf 'FAIL: the three-attempt stop condition ("consecutive fix attempts ... in its current context only") must be stated exactly once across actions/ (found %s).\n' \
    "$three_attempt_count" >&2
  fail_count=$((fail_count + 1))
fi

# Claim-respecting crash recovery (REQ-071). Recovery resets frontmatter and strips thirteen
# generated sections, and the pipeline does not commit until Step 9 — so an unconditional recovery
# destroys a finished Plan/Exploration/Scope that exists nowhere else. The premise that licensed
# running it on every working/ file is gone and must stay gone; what replaced it is a classification,
# a human-authorized takeover, and an ask-on-ambiguity default. Each assertion below pins one of those.
#
# The premise sweep (widened by REQ-077). Trigger condition, not a wording list: any sentence in
# either pipeline file telling the reader that every file in do-work/working/ belongs to the current
# session, or that the pipeline keeps no claim record at all, restates that premise — and since
# REQ-077 the pipeline does keep a claim record (CHECKPOINT.md's In Progress list, written at claim
# time by Step 2). The fingerprints below are the ones seen in the tree so far; they are illustrative,
# not exhaustive, so a newly discovered wording earns a generalized pattern rather than a fourth
# literal. Whole-file scope over both files is also deliberate: the two predecessor assertions were
# each narrower than the premise — one read only actions/work.md's Step 1 block while the live
# restatement sat in Step 2, the other was scoped to actions/work-reference.md by argument — so
# broadening the pattern alone would have left both of them green.
for premise_file in actions/work.md actions/work-reference.md; do
  for premise_fingerprint in \
    'no other live session whose in-flight claim a recovery could disturb' \
    "every .working/. file is this session's" \
    'no lock or claim record'
  do
    assert_file_not_contains \
      "$premise_file" \
      "$premise_fingerprint" \
      "$premise_file must not restate the retired premise that every do-work/working/ file is this session's to recover, nor claim the pipeline keeps no claim record at all: Crash Recovery classifies each working/ file against CHECKPOINT.md's In Progress list, which Step 2 writes at claim time (REQ-071, REQ-077). Matched fingerprint: $premise_fingerprint"
  done
done

crash_recovery_block="$(sed -n '/^## Crash Recovery (Step 1)/,/^## Worktree Dispatch Mode/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$crash_recovery_block" \
  'foreign claim' \
  'actions/work-reference.md Crash Recovery must classify a claimed working/ REQ the checkpoint does not name as a foreign claim, not as this sessions own leftover (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'absent checkpoint is ambiguous' \
  'actions/work-reference.md Crash Recovery must state that a missing CHECKPOINT.md is ambiguous, never permission to recover — a hard crash usually leaves no checkpoint at all (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'never authorizes' \
  'actions/work-reference.md Crash Recovery must state that the staleness threshold gates only the offer and never authorizes a takeover by itself, or a later edit simplifies it into an automatic one (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'unparseable, future-dated, or absent' \
  'actions/work-reference.md Crash Recovery must guard a bad claimed_at toward asking (immediately eligible), or a REQ carrying a corrupt stamp is protected from takeover forever (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'no human to answer' \
  'actions/work-reference.md Crash Recovery must keep the unattended path non-blocking — a foreign claim is left alone, reported, and the run continues (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'substep 1 removes it' \
  'actions/work-reference.md Crash Recovery must say claimed_at is read while classifying, before substep 1 discards it — the same ordering trap as the Scope/write_set decision (REQ-071).'

# A label-less in-progress entry is never auto-recovered (REQ-104). The retired bullet made a
# locally modified do-work/CHECKPOINT.md stand in for authorship. REQ-095's two-clone acceptance run
# falsified that premise: F-06 shows the checkpoint conflicting on *every* concurrent claim, even two
# that overlap in nothing, so a merge-resolving checkout holds a dirty checkpoint for reasons unrelated
# to who wrote which entry — and F-07 reproduced the consequence, a label-less *foreign* entry
# classified as an own crash and its live claim stripped. That strip is the 2026-07-01 incident,
# reached through the label-less door. Both pins are needed: the positive one keeps the report-only
# classification stated, the negative one keeps the "locally modified ⇒ mine" inference from being
# reintroduced as an optimization for pre-0.170.0 checkpoints.
assert_block_contains \
  "$crash_recovery_block" \
  'claim of unknown origin, always report-only' \
  'actions/work-reference.md Crash Recovery must classify a label-less In Progress entry as a claim of unknown origin that is always report-only — never recovered, whatever state the local checkpoint is in (REQ-104).'

assert_block_not_contains \
  "$crash_recovery_block" \
  'locally modified or otherwise uncommitted' \
  'actions/work-reference.md Crash Recovery must not treat a locally modified CHECKPOINT.md as evidence this checkout authored a label-less entry: under claim-anywhere every concurrent claim conflicts on that file, so a merge-resolving checkout is dirty for reasons unrelated to authorship and the inference strips a live foreign claim (REQ-104).'

# Recovery stamps the flip instant (REQ-074). The automatic reset had been silently out of compliance
# with status_changed_at's own trigger condition since the field was introduced, and nothing caught it
# for that entire span — the manual reset in actions/forensics.md stamped, this one did not. Since the
# same substep removes claimed_at, an unstamped recovery leaves no trace of the reset at all and the
# board dates a just-recovered REQ from created_at.
assert_block_contains \
  "$crash_recovery_block" \
  'stamp `status_changed_at' \
  'actions/work-reference.md Crash Recovery substep 1 must stamp status_changed_at on the reset — the field is written on any status flip with no dedicated *_at of its own, and this substep also removes claimed_at, so it is the only surviving trace of when recovery happened (REQ-074).'

# The hand-back file's one legal write location (REQ-082). Fan-Out Dispatch makes a per-builder
# REQ-NNN-handback.md mandatory and background-agents.md has the sub-agent write it itself, while
# Sole integrator says a builder never writes the main tree and do-work/ is main-tree-only — so the
# mandatory file had no legal home, and an agent hitting that had three moves, two of which corrupt
# the run (write the main tree anyway, or write the worktree's copy where it lands in the branch and
# the orchestrator reads nothing). The failure this guards is a later maintenance pass reading "the
# builder never writes the main tree" as absolute and deleting the carve-out as redundant, silently
# restoring the contradiction.
worktree_dispatch_block="$(sed -n '/^## Worktree Dispatch Mode (Step 1)/,/^## Composed Exit Summary/p' "$core_root/actions/work-reference.md")"

# actions/work.md is the executable, condensed hand-back path. It must retain the two
# pre-merge guards from the canonical reference: isolate owner bookkeeping from any
# unrelated staged work, and reject builder commits under do-work/ while the branch
# diff can still see them. A trailing pointer to the reference is too late because a
# reader following the numbered commands has already merged by then.
work_action_handback_block="$(sed -n '/^\*\*Hand-back merge/,/^### Step 6\.25:/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_action_handback_block" \
  'git diff --cached --name-only' \
  'actions/work.md hand-back step 0 must inspect the index before committing owner bookkeeping, or a plain commit can hide unrelated staged work below the evidence range.'

assert_block_contains \
  "$work_action_handback_block" \
  'git commit -- do-work/' \
  'actions/work.md hand-back step 0 must commit owner bookkeeping path-limited, never sweep the whole index into the bookkeeping commit.'

assert_block_contains \
  "$work_action_handback_block" \
  'clean index' \
  'actions/work.md hand-back step 0 must end with a clean index so unrelated staged work cannot leak into the merge commit.'

# Hand-back scratch is expected main-tree state, but it must never enter the
# bookkeeping commit. Preserve all three categories in both the canonical and
# condensed instructions: stage owner bookkeeping, tolerate the named scratch
# file without staging it, and stop on everything else.
for handback_staging_block in "$worktree_dispatch_block" "$work_action_handback_block"
do
  assert_block_contains \
    "$handback_staging_block" \
    'stage.*`manifest\.md`.*`REQ-NNN-brief\.md`' \
    'Hand-back step 0 must stage the owner-written manifest and per-builder briefs explicitly instead of admitting the whole run directory.'

  assert_block_contains \
    "$handback_staging_block" \
    'allow but never stage.*`REQ-NNN-handback\.md`' \
    'Hand-back step 0 must treat each expected REQ-NNN-handback.md as allowed scratch that remains unstaged.'

  assert_block_contains \
    "$handback_staging_block" \
    'stop and surface.*every other `do-work/` path' \
    'Hand-back step 0 must stop on do-work paths outside the explicit bookkeeping and scratch categories.'
done

assert_block_not_contains \
  "$worktree_dispatch_block" \
  'the run directory if fan-out created one' \
  'actions/work-reference.md must not restore the whole fan-out run directory to the stage set, because it contains never-staged hand-back scratch.'

assert_block_not_contains \
  "$work_action_handback_block" \
  "this run's directory" \
  'actions/work.md must not restore the whole run directory to the stage set, because it contains never-staged hand-back scratch.'

assert_block_contains \
  "$work_action_handback_block" \
  'git diff --name-only <pre>\.\.\.<operative_name> -- do-work/' \
  'actions/work.md must reject builder queue-state commits before merging, while the three-dot branch diff can still see them.'

assert_block_contains \
  "$worktree_dispatch_block" \
  'exactly one exception' \
  'actions/work-reference.md Sole integrator must carry the bounded hand-back exception — Fan-Out Dispatch mandates REQ-NNN-handback.md, and without a named legal write location the builder must either violate sole-integrator or drop the file the durability pattern exists for (REQ-082).'

assert_block_contains \
  "$worktree_dispatch_block" \
  'never staged, committed, or merged' \
  'actions/work-reference.md Sole integrator must say the hand-back file is never staged/committed/merged — "you may write it" naturally reads as including committing it, which would turn run scratch into branch content and trip the builder-wrote-do-work probe (REQ-082).'

# Scope, not a list: the exception is one path derived from the builder's own REQ id, and do-work/'s
# main-tree-only rule must be stated as a condition. Both were closed enumerations away from silently
# growing — "the run directory" instead of one file, and a three-item list that predated do-work/runs/.
assert_block_contains \
  "$worktree_dispatch_block" \
  'Every path under `do-work/` exists in the main tree only' \
  'actions/work-reference.md State stays home must state the condition rather than enumerate the queue/working//CHECKPOINT.md trio — do-work/runs/ postdated that list and a reader could argue it was out of scope (REQ-082).'

assert_block_not_contains \
  "$worktree_dispatch_block" \
  'the queue, `working/`, `CHECKPOINT.md` — exists in the main tree only' \
  'actions/work-reference.md must not restore State stays home three-item enumeration — a hand-maintained list of what lives under do-work/ goes stale the moment a directory is added (REQ-082).'

# REQ-299 — every `##` section Step 6 names is classified, and the routed ones match the
# hand-back contents exactly. REQ-270 fixed `## Discovered Tasks` and left `## Decisions`
# unqualified: Step 6 told a worktree builder to write a file `State stays home` forbids it
# to touch, and the two readers outside Step 8 (review-work's traceability check, the
# Decision Brief's HANDLED block) could not inherit a rule scoped to Step 8's substeps. The
# failure was silent both ways — review reported clean, the brief rendered empty.
#
# The check carries no list of sections, deliberately. It classifies whatever Step 6
# mentions: routed to the hand-back, or explicitly `not yours to write`. A section added to
# Step 6 later that says neither fails here rather than shipping the same defect again.
if ! python3 - "$core_root/actions/work.md" "$core_root/actions/work-reference.md" <<'PY'
import pathlib
import re
import sys

work_text = pathlib.Path(sys.argv[1]).read_text()
reference_text = pathlib.Path(sys.argv[2]).read_text()

builder_instruction_block = re.search(
    r"^All routes include these instructions to the agent.*?^\*\*Hand-back merge",
    work_text,
    flags=re.DOTALL | re.MULTILINE,
)
if builder_instruction_block is None:
    raise SystemExit(
        "actions/work.md Step 6 no longer has an 'All routes include these instructions to "
        "the agent' block ending at the hand-back merge — the extraction anchor moved"
    )

# One bullet per top-level `- ` item; continuation lines belong to the bullet above them.
bullets = []
for line in builder_instruction_block.group(0).splitlines():
    if line.startswith("- "):
        bullets.append(line)
    elif bullets:
        bullets[-1] += " " + line.strip()

section_token = re.compile(r"`(## [A-Z][A-Za-z -]*)`")
routing_clause = re.compile(r"that section goes in your hand-back", re.IGNORECASE)
disclaimer_clause = re.compile(r"not yours to write", re.IGNORECASE)

# Classification is per section, not per mention: several bullets may name the same
# section, and one clear statement of who writes it is enough for a reader. What must
# never happen is a section Step 6 names that no bullet classifies at all.
routed_sections = set()
disclaimed_sections = set()
mentioned_sections = set()
for bullet in bullets:
    mentioned = set(section_token.findall(bullet))
    mentioned_sections |= mentioned
    if routing_clause.search(bullet):
        routed_sections |= mentioned
    if disclaimer_clause.search(bullet):
        disclaimed_sections |= mentioned

unclassified = mentioned_sections - routed_sections - disclaimed_sections
if unclassified:
    raise SystemExit(
        f"actions/work.md Step 6 names {', '.join(sorted(unclassified))} without saying "
        "whether the builder authors it: every `##` section the block names must either be "
        "routed to the hand-back for worktree dispatch mode or be marked 'not yours to "
        "write', or a builder is told to write a file it may not touch"
    )

if not routed_sections:
    raise SystemExit(
        "actions/work.md Step 6 routes no `##` section to the builder's hand-back — the "
        "check would pass vacuously, so the routing clause itself must have changed"
    )

handback_row = next(
    (line for line in reference_text.splitlines() if line.startswith("| per-builder output |")),
    None,
)
if handback_row is None:
    raise SystemExit(
        "actions/work-reference.md Fan-Out Dispatch has no per-builder output row — the "
        "hand-back contents contract moved"
    )
named_sections = set(section_token.findall(handback_row))

if routed_sections != named_sections:
    raise SystemExit(
        "the sections Step 6 tells the builder to author and the sections the hand-back "
        "contract names disagree — routed by Step 6 but not carried by the hand-back: "
        f"{sorted(routed_sections - named_sections)}; named by the hand-back but not routed "
        f"by Step 6: {sorted(named_sections - routed_sections)}"
    )
PY
then
  printf 'FAIL: actions/work.md Step 6 and actions/work-reference.md disagree about which `##` sections a worktree builder hands back (REQ-299) — a section Step 6 tells the builder to author that the hand-back never carries is lost silently.\n' >&2
  fail_count=$((fail_count + 1))
fi

# REQ-299 — the rule's home. REQ-270 stated it in actions/work.md Step 8's preamble, opening
# "Some substeps below", so review-work's traceability check and the Decision Brief — both
# outside Step 8 — could not inherit it. The rule now lives in the reference every reader
# already loads, keyed on the condition, with its reader list explicitly illustrative.
builder_section_rule_block="$(sed -n '/^## Reading a Builder-Authored Section (any step)/,/^## Composed Exit Summary/p' "$core_root/actions/work-reference.md")"

if [ -z "$builder_section_rule_block" ]; then
  printf 'FAIL: actions/work-reference.md has no `## Reading a Builder-Authored Section (any step)` section (REQ-299) — the rule must live outside actions/work.md Step 8, where readers at other steps can inherit it.\n' >&2
  fail_count=$((fail_count + 1))
fi

assert_block_contains \
  "$builder_section_rule_block" \
  'The condition carries the rule, not any list of readers' \
  'actions/work-reference.md Reading a Builder-Authored Section must key on the condition and mark its reader list illustrative (REQ-299) — a closed list of readers is how the rule missed review-work and the Decision Brief the first time.'

assert_block_contains \
  "$builder_section_rule_block" \
  'relative to the project root' \
  'actions/work-reference.md Reading a Builder-Authored Section must say which root the hand-back path resolves against (REQ-299) — a reader resolving do-work/runs/ against the vendored .claude/skills/ directory of a consumer install finds nothing.'

assert_block_contains \
  "$builder_section_rule_block" \
  'Absence is only silence when you know you looked' \
  'actions/work-reference.md Reading a Builder-Authored Section must carry the absence-vs-silence rule (REQ-299) — without it a reader cannot distinguish an unread hand-back from an empty one.'

assert_contains \
  'actions/work.md' \
  'Reading a Builder-Authored Section \(any step\)' \
  'actions/work.md Step 8 must point at actions/work-reference.md → Reading a Builder-Authored Section rather than restating the rule (REQ-299) — a Step-8-local copy is what kept readers outside Step 8 from inheriting it.'

# Both readers outside Step 8 must inherit the rule and must say which absence they found.
assert_contains \
  'actions/review-work.md' \
  'Reading a Builder-Authored Section \(any step\)' \
  "actions/review-work.md Step 4's traceability check must read the REQ's ## Decisions per the shared rule (REQ-299) — under fan-out the section is in the hand-back, and reading the REQ file alone scores a missing section as clean."

# The unreadable case is the half whose loss collapses the two facts: an empty hand-back
# already reads as "nothing recorded" everywhere, so only the unread one needs naming.
for builder_section_reader_path in actions/review-work.md actions/work-reference.md; do
  if ! grep -q 'could not be read\|hand-back unread' "$(resolve_runtime_file "$builder_section_reader_path")"; then
    printf 'FAIL: %s must distinguish "no section anywhere" from "the builder recorded nothing" (REQ-299) — an unread hand-back and an empty one are different facts and must never render the same.\n' "$builder_section_reader_path" >&2
    fail_count=$((fail_count + 1))
  fi
done

# REQ-308 / REQ-314 — every REQ writer judges effort_estimate by the same
# three-way contract it judges impact.
# capture.md required a judged `impact:` in a rule written to close exactly this hole
# ("an absent impact: must not be the common case") while the neighbouring field was
# only ever MAY-set. That stopped being cosmetic when `do-work run-simple-reqs` began
# selecting work on effort_estimate: at capture time 14 of 22 pending REQs carried the
# field, so 8 read as effort-substantive by default and were invisible to that verb —
# not because anyone judged them substantive, but because nobody judged them.
#
# The check pins the PROPERTY, not the wording (REQ-293's lesson): the two checklist
# lines must be the same sentence apart from the field they name. A rule that drifts on
# one field and not the other fails here, whatever either sentence happens to say.
if ! python3 - \
  "$core_root/actions/capture.md" \
  "$core_root/actions/review-work.md" \
  "$core_root/actions/work-reference.md" \
  "$board_root/tools/queue-kanban/model.go" <<'PY'
import pathlib
import re
import sys

capture_text = pathlib.Path(sys.argv[1]).read_text()
review_text = pathlib.Path(sys.argv[2]).read_text()
work_reference_text = pathlib.Path(sys.argv[3]).read_text()
board_model_text = pathlib.Path(sys.argv[4]).read_text()

checklist_block = re.search(
    r"^## Verification Checklist$(.*?)(?=^## |\Z)",
    capture_text,
    flags=re.DOTALL | re.MULTILINE,
)
if checklist_block is None:
    raise SystemExit("actions/capture.md has no '## Verification Checklist' section — the anchor moved")

judged_field_names = ("impact:", "effort_estimate:")
lines_by_field = {field_name: [] for field_name in judged_field_names}
for checklist_line in checklist_block.group(1).splitlines():
    if not checklist_line.startswith("- [ ] "):
        continue
    for field_name in judged_field_names:
        if f"`{field_name}`" in checklist_line:
            lines_by_field[field_name].append(checklist_line)

# The judged-verdict line is the one that states the three-way contract. A field can
# carry other checklist lines (impact carries the title-mirror one), so the judged line
# is identified by naming the field with no other field beside it.
def judged_verdict_line(field_name):
    candidates = [
        checklist_line
        for checklist_line in lines_by_field[field_name]
        if checklist_line.count("`") == 2
    ]
    if len(candidates) != 1:
        raise SystemExit(
            f"actions/capture.md's checklist must carry exactly one line stating how "
            f"`{field_name}` is decided; found {len(candidates)}: {candidates}"
        )
    return candidates[0]

# Same sentence, different field. Strip the field token and compare what is left.
skeletons = {
    field_name: judged_verdict_line(field_name).replace(f"`{field_name}`", "<field>")
    for field_name in judged_field_names
}
if skeletons["impact:"] != skeletons["effort_estimate:"]:
    raise SystemExit(
        "actions/capture.md judges impact: and effort_estimate: by different standards — "
        "the two checklist lines must be one sentence apart from the field they name.\n"
        f"  impact:          {skeletons['impact:']}\n"
        f"  effort_estimate: {skeletons['effort_estimate:']}"
    )

# And the standard has to be the three-way one, or a matching pair of weak rules passes.
judged_line = judged_verdict_line("impact:")
for required_alternative, description in (
    ("judged", "judge it yourself"),
    ("put to the user", "ask the user"),
    ("absent", "leave it absent"),
):
    if required_alternative not in judged_line:
        raise SystemExit(
            f"the judged-verdict checklist line no longer offers '{description}' — the "
            f"three-way contract is what stops a copied default: {judged_line}"
        )

def markdown_section(document_text, heading_pattern, next_heading_pattern=r"^#{1,3} "):
    section_match = re.search(
        rf"^{heading_pattern}$(.*?)(?={next_heading_pattern}|\Z)",
        document_text,
        flags=re.DOTALL | re.MULTILINE,
    )
    if section_match is None:
        raise SystemExit(f"missing contract section matching {heading_pattern!r}")
    return section_match.group(1)

review_followup_section = markdown_section(
    review_text,
    r"### Step 10: Create Follow-up REQs",
    next_heading_pattern=r"^### ",
)
discovered_tasks_section = markdown_section(
    work_reference_text,
    r"## Discovered Tasks Classification \(Step 8\)",
    next_heading_pattern=r"^## ",
)

def require_three_way_effort_judgment(section_name, section_text):
    required_properties = (
        (r"effort-mechanical", "emits the mechanical judgment"),
        (r"effort-substantive", "emits the non-small judgment"),
        (r"judg", "requires a size judgment"),
        (r"(?:ask|put).{0,80}(?:user|question)|(?:user|question).{0,80}(?:ask|put)", "puts an unclear judgment to the user"),
        (r"(?:omit|absent).{0,120}(?:neither|cannot|could not|not possible)|(?:neither|cannot|could not|not possible).{0,120}(?:omit|absent)", "permits omission only when neither judging nor asking is possible"),
        (r"(?:never|not).{0,80}(?:copied|default)|(?:copied|default).{0,80}(?:never|not)", "forbids a copied default"),
    )
    # Template payload is checked separately at the emission seam. It must not
    # satisfy the instruction contract when the surrounding writer rule is weak.
    instruction_text = re.sub(r"^```.*?^```$", "", section_text, flags=re.DOTALL | re.MULTILINE)
    effort_directives = [
        paragraph
        for paragraph in re.split(r"\n\s*\n", instruction_text)
        if "`effort_estimate`" in paragraph and "effort-mechanical" in paragraph
    ]
    if len(effort_directives) != 1:
        raise SystemExit(
            f"{section_name} must carry exactly one effort_estimate writer directive; "
            f"found {len(effort_directives)}"
        )
    # Anchor every semantic leg at the field's own directive. Even the impact/title
    # prefix in the same paragraph has unrelated default/never language that must not
    # make a weakened effort rule pass.
    effort_directive = effort_directives[0]
    effort_directive = effort_directive[effort_directive.index("`effort_estimate`"):]
    flattened = " ".join(effort_directive.split())
    for property_pattern, description in required_properties:
        if re.search(property_pattern, flattened, flags=re.IGNORECASE) is None:
            raise SystemExit(f"{section_name} no longer {description}")

require_three_way_effort_judgment("actions/review-work.md Step 10", review_followup_section)
require_three_way_effort_judgment(
    "actions/work-reference.md Discovered Tasks Classification",
    discovered_tasks_section,
)

# The emitted follow-up must carry the result. A prose judgment beside a template that
# drops the field recreates the same absent-by-accident population at the write seam.
review_template_match = re.search(
    r"```markdown\n---\n(.*?)\n---",
    review_followup_section,
    flags=re.DOTALL,
)
if review_template_match is None:
    raise SystemExit("actions/review-work.md Step 10 has no follow-up frontmatter template")
if re.search(r"^effort_estimate:\s*\[", review_template_match.group(1), flags=re.MULTILINE) is None:
    raise SystemExit("actions/review-work.md Step 10's follow-up template drops the judged effort_estimate")

# Schema surfaces name all current writers. Keeping a capture-only gloss here teaches
# the weaker contract even when the two writer steps themselves are correct.
schema_effort_line = next(
    (line for line in work_reference_text.splitlines() if line.startswith("effort_estimate:")),
    None,
)
if schema_effort_line is None:
    raise SystemExit("actions/work-reference.md has no effort_estimate schema line")
for writer_name in ("capture", "review", "discovered"):
    if writer_name not in schema_effort_line.lower():
        raise SystemExit(f"actions/work-reference.md's effort_estimate schema line omits the {writer_name} writer")

effort_read_contract_line = next(
    (line for line in work_reference_text.splitlines() if line.startswith("| `effort_estimate`")),
    None,
)
if effort_read_contract_line is None:
    raise SystemExit("actions/work-reference.md has no effort_estimate Schema Read Contract row")
for writer_name in ("capture", "review", "discovered"):
    if writer_name not in effort_read_contract_line.lower():
        raise SystemExit(f"actions/work-reference.md's effort_estimate read-contract row omits the {writer_name} writer")

model_effort_comment = re.search(
    r"// effort_estimate —.*?\n\s*EffortEstimate\s+string",
    board_model_text,
    flags=re.DOTALL,
)
if model_effort_comment is None:
    raise SystemExit("queue-kanban/model.go has no effort_estimate schema-mirror comment")
for writer_name in ("capture", "review", "discovered"):
    if writer_name not in model_effort_comment.group(0).lower():
        raise SystemExit(f"queue-kanban/model.go's effort_estimate schema mirror omits the {writer_name} writer")
PY
then
  printf 'FAIL: a REQ writer does not carry the three-way effort_estimate judgment contract (REQ-308 / REQ-314) — a field nobody judged reads as effort-substantive by default and is invisible to `do-work run-simple-reqs`.\n' >&2
  fail_count=$((fail_count + 1))
fi

# No shipped file may still say capture MAY set the field — the weaker rule survives
# wherever it was restated, and a reader who lands on the restatement never sees the rule.
capture_may_set_effort="$(cd "$repo_root" && grep -rIn 'apture MAY \(set\|emit\)' skills/ 2>/dev/null || true)"
if [ -n "$capture_may_set_effort" ]; then
  printf 'FAIL: a shipped file still says capture MAY set a field it must now judge (REQ-308):\n%s\n' "$capture_may_set_effort" >&2
  fail_count=$((fail_count + 1))
fi

# The claim-time write (REQ-077) is the other half of REQ-071's gate. REQ-071 made recovery consume
# the checkpoint's In Progress record, but Step 10 (session end) was its only write site — so a hard
# crash left no record at all, every crashed REQ classified as a foreign claim, and the own-crash
# branch became unreachable by the exact event it exists to handle. Each assertion below pins one
# numbered requirement of the fix: write it at claim time, keep it a list, keep it from becoming a
# lock, and remove it when the REQ leaves working/.
work_step_two_block="$(sed -n '/^### Step 2: Claim the Request/,/^### Step 3: Triage/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_step_two_block" \
  'In Progress \(interrupted\)' \
  'actions/work.md Step 2 must record the claim in CHECKPOINT.md In Progress (interrupted) at claim time — with Step 10 (session end) as the only write site, a hard crash leaves recovery no classification input and the REQ strands in working/ forever (REQ-077).'

in_progress_record_block="$(sed -n '/^## In-Progress Record (Step 2)/,/^## Triage Section Template/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$in_progress_record_block" \
  'record is a list' \
  'actions/work-reference.md In-Progress Record must specify the record as a list, one entry per claimed REQ — fan-out claims several REQs under one owner, so a singular record classifies every claim but the newest as a foreign claim after a crash (REQ-077).'

assert_block_contains \
  "$in_progress_record_block" \
  'never grow into one' \
  'actions/work-reference.md In-Progress Record must state that the record is a classification input and must never grow into a lock, heartbeat, or liveness check — that is the machinery REQ-069 deleted and REQ-073 declined to revive (REQ-077).'

# The writer label (REQ-094). The checkpoint is a tracked file wherever a consumer commits do-work/,
# so another checkout's live claim arrives by an ordinary git pull looking exactly like a local one —
# recovery then reads it as an own crash and strips a REQ someone is actively building. Each
# assertion below pins one half of the fix: the entry carries the label, and a labeled foreign entry
# is reported rather than aged into the takeover ladder.
assert_block_contains \
  "$in_progress_record_block" \
  'writer: <hostname>:<absolute-checkout-path>' \
  'actions/work-reference.md In-Progress Record must give the entry format a writer: <hostname>:<absolute-checkout-path> label — the path alone collides across machines and the hostname alone collides across checkouts on one machine, so neither half identifies a checkout by itself (REQ-094).'

assert_block_contains \
  "$crash_recovery_block" \
  'claim held by' \
  'actions/work-reference.md Crash Recovery must report a foreign-label entry as a claim held by that writer and leave it untouched — a label is positive evidence of another checkouts claim, so it is classified by the label and never aged into the three-hour takeover ladder (REQ-094).'

# The tripwire had to be reworded to admit a static writer label without admitting liveness
# machinery, and the naive reword drops the ban to make room for the carve-out. Pinning both phrases
# to the SAME paragraph is the point: a carve-out that drifts into its own paragraph reads as a
# general permission rather than as the one exception to the ban standing beside it.
in_progress_tripwire_paragraph="$(printf '%s\n' "$in_progress_record_block" | grep -F 'never grow into one' || true)"
if [ -z "$in_progress_tripwire_paragraph" ] || ! printf '%s\n' "$in_progress_tripwire_paragraph" | grep -qF 'never refreshed'; then
  printf 'FAIL: actions/work-reference.md In-Progress Record must keep the tripwire ("never grow into one") and the static-label carve-out ("never refreshed") in one paragraph — the label is the single exception to the lock/heartbeat/liveness ban, and it only reads as an exception standing next to it (REQ-077, REQ-094).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The two label-destruction paths in actions/work.md Step 10 (REQ-102). Both were scoped to entries
# "carrying another checkout's writer: label", which silently excludes the label-less legacy entry
# that Crash Recovery classifies report-only in a clean, committed checkpoint — the wholesale rewrite
# could drop it and the session-start delete could remove the whole file, after which the next run
# sees a working/ REQ "not named there" and ages it into the three-hour takeover ladder. Nothing
# pinned either clause, so a later "simplify the checkpoint rewrite" pass could reopen the hole with
# the suite green. Both must scope preservation to every entry this checkout did not write.
work_session_checkpoint_block="$(sed -n '/^#### Session Checkpoint/,/^## Clarify Questions/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_session_checkpoint_block" \
  'entry this checkout did not write through verbatim' \
  'actions/work.md Step 10 Session Checkpoint must carry every In Progress entry this checkout did not write through verbatim — scoping the preserve rule to labeled foreign entries lets the wholesale rewrite drop a label-less entry that Crash Recovery had classified report-only (REQ-102).'

assert_block_contains \
  "$work_session_checkpoint_block" \
  'no entry this checkout did not write remains' \
  'actions/work.md session-start step 3 must gate deleting CHECKPOINT.md on no entry this checkout did not write remaining — gating on labeled foreign entries alone authorizes deleting a checkpoint whose only surviving claim is label-less, which the next run then ages into the takeover ladder (REQ-102).'

assert_block_contains \
  "$work_archive_success_block" \
  'In Progress \(interrupted\)' \
  'actions/work.md Step 8 must remove the REQ In Progress entry as part of the archive move — a REQ still listed there after it leaves working/ is the contradiction the next run is told to report, and a report that fires on every normal completion trains readers to ignore it (REQ-077).'

work_step_one_block="$(sed -n '/^### Step 1: Find Next Request/,/^### Step 2\.0/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_step_one_block" \
  "Crash Recovery's input" \
  'actions/work.md Step 1 must name do-work/CHECKPOINT.md as Crash Recoverys input, so the read is understood as a precondition rather than resume convenience (REQ-071).'

# Ordering, not just presence: the checkpoint read has to appear ahead of the Crash Recovery
# paragraph in Step 1. Recovery consumes the checkpoint, so a later edit that moves the read below
# it would leave recovery classifying against a file it has not opened.
checkpoint_read_line="$(printf '%s\n' "$work_step_one_block" | grep -n 'CHECKPOINT\.md' | head -1 | cut -d: -f1)"
crash_recovery_line="$(printf '%s\n' "$work_step_one_block" | grep -n '\*\*Crash Recovery:\*\*' | head -1 | cut -d: -f1)"
if [ -z "$checkpoint_read_line" ] || [ -z "$crash_recovery_line" ] || [ "$checkpoint_read_line" -ge "$crash_recovery_line" ]; then
  printf 'FAIL: actions/work.md Step 1 must read do-work/CHECKPOINT.md BEFORE the Crash Recovery paragraph (checkpoint line: %s, recovery line: %s) — the checkpoint is recovery input, not resume decoration (REQ-071).\n' \
    "${checkpoint_read_line:-none}" "${crash_recovery_line:-none}" >&2
  fail_count=$((fail_count + 1))
fi

# Allocator / verifier subcommands (REQ-072). The tool grew a second write surface
# (one version line, outside do-work/) and three release-ritual subcommands. Two
# things must not drift: the never-write-CHANGELOG boundary, and the three prose
# call sites — a subcommand nothing calls is dead weight, and a call site with no
# stated fallback turns the Go toolchain into a hard dependency of the pipeline.
assert_file_not_contains \
  "tools/queue-kanban/release.go" \
  'os\.(WriteFile|Create|OpenFile)\(.*CHANGELOG' \
  'tools/queue-kanban/release.go must never write CHANGELOG.md — unique version numbers do not make a shared prepend safe, so the changelog stays an owner-only human write.'

assert_file_not_contains \
  "tools/queue-kanban/verify.go" \
  'os\.(WriteFile|Create|OpenFile|Remove|Rename)\(' \
  'tools/queue-kanban/verify.go must stay read-only — verify reports and routes, and repairs belong to actions/cleanup.md, which asks first.'

# The unknown-subcommand message must list every subcommand main() dispatches, or the
# error text lies about what exists. Derived from the dispatch switch rather than pinned
# to a hand-written chain: a frozen list only ever verifies the subcommands that already
# existed when it was written, which is exactly how it would miss the next one added
# (Closed Enumerations Go Stale).
unknown_subcommand_message="$(grep -F 'unknown subcommand %q' "$board_root/tools/queue-kanban/main.go" || true)"
if [ -z "$unknown_subcommand_message" ]; then
  printf 'FAIL: tools/queue-kanban/main.go has no unknown-subcommand message — a mistyped subcommand must name the valid ones, not exit silently.\n' >&2
  fail_count=$((fail_count + 1))
else
  dispatched_subcommands="$(awk '
    /^\tswitch subcommand \{$/ { in_dispatch_switch=1; next }
    in_dispatch_switch && /^\t}$/ { exit }
    in_dispatch_switch && /^\tcase / {
      case_labels=$0
      sub(/^[[:space:]]*case[[:space:]]+/, "", case_labels)
      sub(/:.*/, "", case_labels)
      while (match(case_labels, /"[a-z-]+"/)) {
        print substr(case_labels, RSTART + 1, RLENGTH - 2)
        case_labels=substr(case_labels, RSTART + RLENGTH)
      }
    }
  ' "$board_root/tools/queue-kanban/main.go" | sort -u)"

  advertised_subcommands="$(printf '%s\n' "$unknown_subcommand_message" \
    | sed -n 's/.*(want \(.*\))\\n".*/\1/p' \
    | tr '|' '\n' \
    | sed 's/^[[:space:]]*//; s/[[:space:]]*$//' \
    | sort -u)"
  if [ -z "$dispatched_subcommands" ]; then
    printf 'FAIL: could not read the subcommand dispatch switch in tools/queue-kanban/main.go — this check derives the expected list from it, so a shape change here silently disables it.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  for dispatched_subcommand in $dispatched_subcommands; do
    if ! grep -qxF -- "$dispatched_subcommand" <<<"$advertised_subcommands"; then
      printf 'FAIL: tools/queue-kanban/main.go dispatches %s but its unknown-subcommand message does not list it — the error text lies about what exists.\n' \
        "$dispatched_subcommand" >&2
      fail_count=$((fail_count + 1))
    fi
  done
fi

# Project release instructions must never point next-version at the managed core
# package. That action file is the suite runtime's version source, not the consumer
# application's selected release source; bumping it poisons later update comparisons.
assert_file_not_contains \
  "actions/work.md" \
  'queue-kanban next-version|queue-kanban/queue-kanban next-version|--version-file.*<skill-root>/actions/version\.md' \
  'actions/work.md must leave project version bumps to the Changelog Entry Procedure and never mutate the installed suite version.'

assert_contains \
  "actions/work-reference.md" \
  'Exclude every installed skill, dependency, vendored package, and generated tree.*\.claude/skills/.*installed do-work suite.*VERSION.*actions/version\.md.*must never be selected or bumped' \
  'The canonical project version-source procedure must exclude installed suite metadata, not only remove one unsafe accelerator call site.'

assert_contains \
  "tools/queue-kanban/main.go" \
  'func rejectLeftoverArguments' \
  'tools/queue-kanban/main.go must keep the shared leftover-argument rejection — an unconsumed token is an error for ANY subcommand, not silence, which is how next-version writing the wrong tree stayed invisible (REQ-081).'

for release_subcommand_call_site in \
  'actions/capture.md:queue-kanban next-req' \
  'actions/forensics.md:queue-kanban verify'; do
  call_site_file="${release_subcommand_call_site%%:*}"
  call_site_pattern="${release_subcommand_call_site#*:}"
  assert_contains \
    "$call_site_file" \
    "$call_site_pattern" \
    "$call_site_file must call \`$call_site_pattern\` — REQ-072 wired the allocator/verifier into the three existing actions instead of adding a new action or a SKILL.md routing row."
  # Every call site must name what to do when the toolchain is missing. The board
  # action reports and stops because there the compiler IS the capability; these
  # three must fall back to the procedure they accelerate.
  assert_contains \
    "$call_site_file" \
    'If .go. is absent|when .go. is absent|absent or the build fails' \
    "$call_site_file must state the fallback for a missing Go toolchain next to its queue-kanban call — the compiler is an accelerator here, never a dependency of the pipeline."
done

assert_contains \
  "CLAUDE.md" \
  'three write surfaces' \
  'CLAUDE.md must state the tool has exactly three write surfaces once next-req reserves ids — testing fields, next-version, and reservation markers are the complete set, and nothing but this sentence records the count.'

assert_contains \
  "actions/capture.md" \
  '\.req-reservations/REQ-NNNNNN' \
  'actions/capture.md must stage next-req reservation markers with the UR/REQ capture — an uncommitted marker would reserve only one checkout and pollute git status.'

assert_contains \
  "docs/forensics-guide.md" \
  'Release and queue invariants' \
  'docs/forensics-guide.md must list the release/queue invariants check — a forensics check absent from the user-facing guide is invisible to the person who would run it.'

assert_contains \
  "docs/ai-report-guide.md" \
  'completed-with-issues' \
  'docs/ai-report-guide.md must reflect the terminal-success set (completed | completed-with-issues), not only completed.'

assert_contains \
  "docs/ai-report-guide.md" \
  'DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND' \
  'docs/ai-report-guide.md must document agentic image backends as opt-in via the env flag, not the opportunistic default.'

assert_contains \
  "docs/cleanup-guide.md" \
  'completed-with-issues' \
  'docs/cleanup-guide.md sweep wording must include completed-with-issues in the terminal-status set.'

assert_file_missing \
  "prompts/ultracode-fable-workflow.md" \
  'retired ultracode/fable prompt file must be removed from the active prompt library.'

active_runtime_docs=(
  "SKILL.md"
  "README.md"
  "next-steps.md"
  "actions/work.md"
  "actions/work-reference.md"
  "prompts/README.md"
)

for runtime_doc in "${active_runtime_docs[@]}"; do
  if [ -f "$repo_root/$runtime_doc" ]; then
    assert_file_not_contains \
      "$runtime_doc" \
      'ultracode|fable' \
      "active runtime doc $runtime_doc must not mention the retired ultracode/fable workflow."
  fi
done

# Router word budget (REQ-020 ratchet). The 2026-07 bloat cleanup cut SKILL.md
# from 5,507 to 2,396 words by deleting duplicate enumerations of the action set
# (the router loads on EVERY invocation — words here tax all 30+ verbs). Budget =
# post-diet count + ~10% headroom. If you hit this limit: the fix is a merge or a
# lazy-load (see actions/help.md for the pattern), not a bigger budget — raise it
# only with an accompanying decisions/ note saying why routing itself had to grow.
router_word_budget=2650
router_word_count="$(wc -w < "$core_root/SKILL.md")"
if [ "$router_word_count" -gt "$router_word_budget" ]; then
  printf 'FAIL: SKILL.md is %s words — over the %s-word router budget. Merge or lazy-load; do not grow the always-loaded router.\n' \
    "$router_word_count" "$router_word_budget" >&2
  fail_count=$((fail_count + 1))
fi

# Hardened checks (REQ-018): the work.md prose pointers and the shipped scripts
# must not drift apart — a pointer at a missing script silently un-hardens the step.
# Each entry is "<script>|<action file that must reference it>". The referencing file is
# per-script because not every hardened check belongs to the work pipeline —
# blanked-req-scan.sh is called from forensics (and, for its restore mode, cleanup). Pinning
# the right caller keeps the assertion strong; hardcoding actions/work.md for all of them
# would have forced either a bogus reference or dropping the check entirely.
hardened_check_scripts=(
  "tools/checks/archive-collision.sh|actions/work.md"
  "tools/checks/preflight.sh|actions/work.md"
  "tools/checks/scope-drift.sh|actions/work.md"
  "tools/checks/qualify.sh|actions/work.md"
  "tools/checks/record-commit-hash.sh|actions/work.md"
  "tools/checks/blanked-req-scan.sh|actions/forensics.md"
  "tools/checks/blanked-req-scan.sh|actions/cleanup.md"
  "scripts/protected-inventory.sh|actions/commit.md"
  "scripts/protected-inventory.sh|actions/inspect.md"
  "tools/estimate-p50.sh|actions/work.md"
  "tools/estimate-p50.sh|actions/verify-requests.md"
)

for check_script_entry in "${hardened_check_scripts[@]}"; do
  check_script="${check_script_entry%%|*}"
  referencing_action_file="${check_script_entry##*|}"
  if [ ! -x "$(resolve_runtime_file "$check_script")" ]; then
    printf 'FAIL: %s must exist and be executable (%s points at it).\n' "$check_script" "$referencing_action_file" >&2
    fail_count=$((fail_count + 1))
  fi
  assert_contains \
    "$referencing_action_file" \
    "$(basename "$check_script")" \
    "$referencing_action_file must reference $check_script — the hardened step's pointer was removed without un-hardening."
done

for delegated_inventory_check in uncommitted-inventory.sh associate-files.sh; do
  assert_contains \
    "scripts/protected-inventory.sh" \
    "$delegated_inventory_check" \
    "protected-inventory.sh must continue delegating to $delegated_inventory_check."
done

# Pre-flight dirty-tree relevance (feedback 2026-08-04). The serial qualifier and
# reviewer inspect the repository-wide working/staged diff, so changes that predate
# the active REQ can contaminate its evidence even though Step 9 stages explicit
# paths. Keep the broad detection, but keep its warning subordinate to the canonical
# Current-REQ relevance rule: preserve/exclude unexpected state and continue unless
# it prevents this REQ from completing. The old "may stage unrelated files" rationale
# contradicted the explicit staging contract and invited removal of a useful check.
current_req_relevance_block="$(sed -n '/^\*\*Current-REQ relevance\./,/^\*\*Three-attempt stop\./p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$current_req_relevance_block" \
  'only.*prevents the active REQ.*implemented, tested, archived, or committed' \
  'actions/work-reference.md must keep unexpected repository state gated by Current-REQ relevance — otherwise pre-flight warnings can become blockers or cleanup work.'

assert_block_contains \
  "$current_req_relevance_block" \
  'preserve it, exclude it from this REQ.s staging, and continue' \
  'actions/work-reference.md must keep the preserve/exclude/continue response for repository state that does not prevent the active REQ.'

# The rule covers SESSION state too, and says the non-action out loud (feedback
# 2026-08-04). A session enumerated its sibling claude PIDs, noticed a commit that
# landed 20 seconds into its run, and blocked on a four-option "how should I proceed?"
# prompt instead of starting its REQ. Nothing shipped told it to — but the contract
# only said the *pipeline* does not detect a concurrent run, which an agent can read
# as leaving the question open for it to raise. Exclusivity is the user's guarantee;
# asking them to re-confirm it rebuilds the coordination the exclusive-session model
# deleted, one prompt at a time.
assert_block_contains \
  "$current_req_relevance_block" \
  '[Nn]ever probe for a concurrent session' \
  'actions/work-reference.md must keep Current-REQ relevance covering session state with the explicit never-probe/never-ask clause — without it an agent re-derives the concurrency check the exclusive-session model removed and stalls the loop on a prompt.'

# The secret-shaped exclusion is a behavior probe, not a grep, because the bug it
# guards was a glob that LOOKED right. `.env|.env.*` reads as covering ".env*" but
# matches neither `.envrc` (direnv — routinely full of exported secrets) nor
# `.environment`: both are suffixes with no dot, so each fell through to `A` and the
# callers would read and stage them. Codex caught it on PR #134. Assert the tags,
# not the pattern — the next wrong pattern will look right too.
inventory_probe_dir="$(mktemp -d)"
cleanup_inventory_probe() {
  rm -rf -- "$inventory_probe_dir"
}
trap cleanup_inventory_probe EXIT
inventory_probe_setup_error="$inventory_probe_dir/setup-error.txt"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$inventory_probe_dir" || exit 1
  git init -q .
  git config user.email probe@test
  git config user.name probe
  mkdir -p nested uppercase
  for probe_name in .env .env.local .envrc .environment production.env \
                    credentials.json server.pem .ENV.PRODUCTION AuthCredentials.json \
                    private.PEM UPPER-SECRET.txt ordinary.js; do
    echo probe > "nested/$probe_name"
  done
  echo probe > uppercase/.ENV
  git add nested/.env.local
  git commit -qm 'seed tracked secret deletion'
  rm nested/.env.local
  echo probe > .env
  git add .env
  git commit -qm 'seed tracked secret rename'
  git mv .env visible-config.txt
) 2>"$inventory_probe_setup_error"; then
  printf 'FAIL: could not set up the uncommitted-inventory behavior probe:\n' >&2
  sed 's/^/  /' "$inventory_probe_setup_error" >&2
  fail_count=$((fail_count + 1))
else
  inventory_probe_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_probe_dir" 2>/dev/null || true)"
  for must_be_excluded in .env .envrc .environment production.env credentials.json server.pem \
                          .ENV.PRODUCTION AuthCredentials.json private.PEM UPPER-SECRET.txt; do
    if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tnested/%s' "$must_be_excluded")"; then
      printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag nested/%s as X (secret-shaped) — it is reachable by the advertised exclusion globs.\n' "$must_be_excluded" >&2
      fail_count=$((fail_count + 1))
    fi
  done
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tuppercase/.ENV')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag uppercase/.ENV as X (case-insensitive secret-shaped basename).\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'XD\tnested/.env.local')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag a deleted secret-shaped path as XD so its deletion can be associated and committed without reading its former contents.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'XD\t.env')" || ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tvisible-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must fail closed when a secret-shaped rename origin moves to an ordinary-looking destination: XD for the source and X for the destination.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tnested/ordinary.js')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must quarantine every A as X when an excluded path makes addition provenance ambiguous.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  # A caller can disable rename detection in repository config. The shipped
  # inventory must override that setting explicitly; otherwise the same staged
  # rename above degrades to an XD + A pair before the action has any provenance
  # to retain.
  git -C "$inventory_probe_dir" config status.renames false
  inventory_renames_disabled_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$inventory_renames_disabled_output" | grep -qxF "$(printf 'XD\t.env')" || \
      ! printf '%s\n' "$inventory_renames_disabled_output" | grep -qxF "$(printf 'X\tvisible-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must force rename detection even when status.renames=false: XD for .env and X for visible-config.txt.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  # commit Step 1 resets a pre-staged X destination. That erases rename
  # provenance from porcelain status and leaves an already-staged source
  # deletion plus an untracked destination. The destination must remain X.
  if ! git -C "$inventory_probe_dir" reset -q -- visible-config.txt; then
    printf 'FAIL: could not reset the secret-rename destination for the re-inventory probe.\n' >&2
    fail_count=$((fail_count + 1))
  else
    inventory_after_reset_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
      "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_probe_dir" 2>/dev/null || true)"
    if ! printf '%s\n' "$inventory_after_reset_output" | grep -qxF "$(printf 'XD\t.env')" || \
        ! printf '%s\n' "$inventory_after_reset_output" | grep -qxF "$(printf 'X\tvisible-config.txt')"; then
      printf 'FAIL: reset-and-reinventory must fail closed: XD for .env and X, never A, for visible-config.txt.\n' >&2
      fail_count=$((fail_count + 1))
    fi

    staged_deletion_metadata="$inventory_probe_dir/staged-deletion-metadata.bin"
    git -C "$inventory_probe_dir" diff --cached --name-status --no-renames -z -- .env > "$staged_deletion_metadata"
    staged_deletion_status=''
    staged_deletion_path=''
    staged_deletion_extra=''
    {
      IFS= read -r -d '' staged_deletion_status
      IFS= read -r -d '' staged_deletion_path
      IFS= read -r -d '' staged_deletion_extra || true
    } < "$staged_deletion_metadata"
    if [ "$staged_deletion_status" != 'D' ] || [ "$staged_deletion_path" != '.env' ] || [ -n "$staged_deletion_extra" ]; then
      printf 'FAIL: reset probe must leave one exact cached deletion for .env; got status=%s path=%s extra=%s.\n' \
        "$staged_deletion_status" "$staged_deletion_path" "$staged_deletion_extra" >&2
      fail_count=$((fail_count + 1))
    fi
    if git -C "$inventory_probe_dir" add -u -- .env 2>/dev/null; then
      printf 'FAIL: probe no longer reproduces Git rejecting git add -u for an already-staged rename-source deletion.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi
fi
cleanup_inventory_probe
trap - EXIT

# An ordinary addition remains readable when no secret-shaped deletion makes
# its provenance ambiguous. Keep this in a separate repository: combining it
# with the XD fixture above would assert the unsafe behavior REQ-128 removes.
ordinary_addition_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$ordinary_addition_probe_dir" || exit 1
  git init -q .
  echo ordinary > ordinary.js
); then
  printf 'FAIL: could not set up the ordinary-addition inventory probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  ordinary_addition_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$ordinary_addition_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$ordinary_addition_output" | grep -qxF "$(printf 'A\tordinary.js')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must leave an ordinary addition as A when no XD exists.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$ordinary_addition_probe_dir"

# REQ-148: awk's common `NR == FNR` first-file discriminator fails when the
# quarantine file is empty: every record from the inventory then also satisfies
# NR == FNR and is swallowed as if it belonged to the exclusion table. Pin the
# portable filename discriminator in every active bridge/modular copy, then
# exercise both an empty quarantine and a populated once-X-always-X set.
association_candidate_action_files=(
  scripts/protected-inventory.sh
)
for association_candidate_action_file in "${association_candidate_action_files[@]}"; do
  assert_contains \
    "$association_candidate_action_file" \
    'FILENAME == ARGV\[1\] \{ excluded\[\$0\] = 1; next \}' \
    "$association_candidate_action_file must distinguish the quarantine by filename so an empty first file cannot swallow every inventory candidate."
  assert_file_not_contains \
    "$association_candidate_action_file" \
    'NR == FNR \{ excluded\[\$0\] = 1; next \}' \
    "$association_candidate_action_file must not use NR == FNR for the possibly-empty quarantine merge."
done

association_candidate_probe_dir="$(mktemp -d)"
association_empty_quarantine="$association_candidate_probe_dir/empty-quarantine.txt"
association_populated_quarantine="$association_candidate_probe_dir/populated-quarantine.txt"
association_inventory="$association_candidate_probe_dir/inventory.txt"
association_expected_empty="$association_candidate_probe_dir/expected-empty.txt"
association_expected_populated="$association_candidate_probe_dir/expected-populated.txt"
association_actual_output="$association_candidate_probe_dir/actual.txt"
: > "$association_empty_quarantine"
printf 'M\tsrc/modified.js\nA\tsrc/added.js\nD\tsrc/deleted.js\nXD\t.env.deleted\nX\tcurrent-secret.txt\n' > "$association_inventory"
printf '%s\n' \
  'src/modified.js' \
  'src/added.js' \
  'src/deleted.js' \
  '.env.deleted' > "$association_expected_empty"

filter_association_candidates() {
  awk -F '\t' '
    FILENAME == ARGV[1] { excluded[$0] = 1; next }
    {
      tag = $1
      sub(/^[^\t]*\t/, "")
      if (tag != "X" && !($0 in excluded)) print
    }
  ' "$1" "$2"
}

filter_association_candidates "$association_empty_quarantine" "$association_inventory" > "$association_actual_output"
if ! cmp -s "$association_expected_empty" "$association_actual_output"; then
  printf 'FAIL: an empty secret quarantine must preserve every safe M/A/D/XD association candidate; expected/actual diff:\n' >&2
  diff -u "$association_expected_empty" "$association_actual_output" >&2 || true
  fail_count=$((fail_count + 1))
fi

printf '%s\n' 'src/added.js' 'previous-secret.txt' > "$association_populated_quarantine"
awk -F '\t' '$1 == "X" { sub(/^[^\t]*\t/, ""); print }' "$association_inventory" >> "$association_populated_quarantine"
printf '%s\n' \
  'src/modified.js' \
  'src/deleted.js' \
  '.env.deleted' > "$association_expected_populated"
filter_association_candidates "$association_populated_quarantine" "$association_inventory" > "$association_actual_output"
if ! cmp -s "$association_expected_populated" "$association_actual_output"; then
  printf 'FAIL: a populated secret quarantine must exclude only retained/current X paths and preserve every other safe candidate; expected/actual diff:\n' >&2
  diff -u "$association_expected_populated" "$association_actual_output" >&2 || true
  fail_count=$((fail_count + 1))
fi
rm -rf -- "$association_candidate_probe_dir"

# Git has no tracked source blob to compare when both a secret-shaped source and
# its copied destination are untracked. The inventory must therefore quarantine
# the ordinary-looking destination too, rather than trusting copy detection
# that cannot exist for this shape.
untracked_secret_copy_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$untracked_secret_copy_probe_dir" || exit 1
  git init -q .
  printf 'fixture-secret\n' > .env
  cp .env application-config.txt
); then
  printf 'FAIL: could not set up the untracked secret-copy inventory probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  untracked_secret_copy_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$untracked_secret_copy_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$untracked_secret_copy_output" | grep -qxF "$(printf 'X\t.env')" || \
      ! printf '%s\n' "$untracked_secret_copy_output" | grep -qxF "$(printf 'X\tapplication-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must quarantine an ordinary-looking untracked copy beside an untracked secret-shaped source.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$untracked_secret_copy_probe_dir"

# Copy detection must be requested explicitly rather than inherited from
# repository configuration. Without it, a secret-derived copy is reported as an
# ordinary A and both action callers are allowed to read it. Keep an ordinary
# rename in the same repository to prove copy-aware detection does not change
# the established M classification for non-secret moves.
copy_inventory_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$copy_inventory_probe_dir" || exit 1
  git init -q .
  git config user.email probe@test
  git config user.name probe
  copy_line=1
  while [ "$copy_line" -le 100 ]; do
    printf 'secret-source-%03d\n' "$copy_line"
    copy_line=$((copy_line + 1))
  done > .env.copy-source
  echo ordinary-source > ordinary-source.txt
  git add .env.copy-source ordinary-source.txt
  git commit -qm 'seed copy and rename sources'
  cp .env.copy-source copied-config.txt
  echo changed >> .env.copy-source
  git add .env.copy-source copied-config.txt
  git mv ordinary-source.txt ordinary-destination.txt
  git config status.renames false
); then
  printf 'FAIL: could not set up the copy-aware inventory probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  copy_inventory_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$copy_inventory_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$copy_inventory_output" | grep -qxF "$(printf 'X\tcopied-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag a secret-derived copy destination as X, never A, even when status.renames=false.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$copy_inventory_output" | grep -qxF "$(printf 'M\tordinary-destination.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must retain M for an ordinary rename while copy-aware detection is active.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$copy_inventory_probe_dir"

# The other Step 5 branch remains live: a tracked secret deletion that is not
# cached yet still needs git add -u, followed by deletion-only metadata.
unstaged_deletion_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$unstaged_deletion_probe_dir" || exit 1
  git init -q .
  git config user.email probe@test
  git config user.name probe
  echo probe > .env
  git add .env
  git commit -qm 'seed unstaged secret deletion'
  rm .env
); then
  printf 'FAIL: could not set up the unstaged secret-deletion probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  if [ -n "$(git -C "$unstaged_deletion_probe_dir" diff --cached --name-status --no-renames -- .env)" ]; then
    printf 'FAIL: unstaged secret-deletion probe unexpectedly began with cached metadata.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! git -C "$unstaged_deletion_probe_dir" add -u -- .env || \
       ! git -C "$unstaged_deletion_probe_dir" diff --cached --name-status --no-renames -- .env \
         | grep -qxF "$(printf 'D\t.env')"; then
    printf 'FAIL: an unstaged tracked secret deletion must still stage as one exact cached D entry.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$unstaged_deletion_probe_dir"

# The action prose is executable behavior. Both consumers must retain a
# run-level quarantine before their second inventory reaches content or REQ
# association, and the manual fallback must force copy-aware detection too.
for inventory_consumer in actions/commit.md actions/inspect.md; do
  assert_contains \
    "$inventory_consumer" \
    'once X.*always X|once `X`.*always `X`' \
    "$inventory_consumer must preserve a once-X-always-X quarantine across every re-inventory in the action run."
  assert_contains \
    "$inventory_consumer" \
    'git -c status.renames=copies status --porcelain=v1 --untracked-files=all -z' \
    "$inventory_consumer manual fallback must force copy-aware detection instead of inheriting status.renames=false."
done

assert_contains \
  'actions/commit.md' \
  'git rev-parse --git-path do-work-commit-secret-quarantine' \
  'actions/commit.md must re-derive a deterministic Git-private quarantine across separate command blocks.'

assert_contains \
  'actions/inspect.md' \
  'git rev-parse --git-path do-work-inspect-secret-quarantine' \
  'actions/inspect.md must re-derive a deterministic Git-private quarantine across separate command blocks.'

assert_contains \
  'actions/commit.md' \
  'already staged.*exact.*deletion|exact.*deletion.*already staged' \
  'actions/commit.md Step 5 must recognize an exact already-staged XD deletion before deciding whether git add -u is needed.'

assert_contains \
  'actions/commit.md' \
  'skip.*git add -u|without.*git add -u' \
  'actions/commit.md Step 5 must skip git add -u when cached name/status already proves the XD deletion.'

assert_contains \
  'actions/commit.md' \
  'otherwise.*stages.*verif|otherwise it stages.*verif' \
  'actions/commit.md Step 5 must retain the unstaged-deletion branch and verify deletion-only metadata after git add -u.'

# Argument parsing must fail fast. The watchdog keeps a regression from hanging
# the entire suite forever: the original `shift 2` with one argument left made
# the loop reread --repo-root indefinitely on Bash 3.2 and newer alike.
associate_missing_value_output="$(mktemp)"
"$core_root/tools/checks/associate-files.sh" --repo-root </dev/null \
  >"$associate_missing_value_output" 2>&1 &
associate_missing_value_pid=$!
(
  sleep 2
  if kill -0 "$associate_missing_value_pid" 2>/dev/null; then
    kill "$associate_missing_value_pid" 2>/dev/null || true
  fi
) &
associate_missing_value_watchdog_pid=$!
if wait "$associate_missing_value_pid"; then
  associate_missing_value_exit=0
else
  associate_missing_value_exit=$?
fi
kill "$associate_missing_value_watchdog_pid" 2>/dev/null || true
wait "$associate_missing_value_watchdog_pid" 2>/dev/null || true
if [ "$associate_missing_value_exit" -ne 2 ]; then
  printf 'FAIL: tools/checks/associate-files.sh --repo-root with no value must fail promptly with exit 2, got %s.\n' \
    "$associate_missing_value_exit" >&2
  fail_count=$((fail_count + 1))
fi
rm -f "$associate_missing_value_output"

# The status alias belongs to the Schema Read Contract, but this shell reader
# owns its own terminal-success predicate. Exercise the new alias against a
# real REQ fixture so the prose and helper cannot drift back apart.
assert_contains \
  "actions/work-reference.md" \
  '`complete`/`done`/`finished`/`closed`.*`completed`' \
  'actions/work-reference.md must document complete with the other aliases that normalize to completed before terminal-success readers consume the status.'

associate_complete_probe_dir="$(mktemp -d)"
mkdir -p "$associate_complete_probe_dir/do-work/archive/UR-301"
cat > "$associate_complete_probe_dir/do-work/archive/UR-301/REQ-501-legacy-complete.md" <<'EOF'
---
id: REQ-501
status: complete
completed_at: 2026-08-07T12:00:00Z
---

## Implementation Summary

**Files changed:**
- `legacy-file.txt`, `second-file.txt` (modified)
- Notes mention `phantom-file.txt`, but this prose bullet claims no file.
EOF
if ! associate_complete_output="$(printf 'legacy-file.txt\nsecond-file.txt\nphantom-file.txt\n' | "$core_root/tools/checks/associate-files.sh" --repo-root "$associate_complete_probe_dir")" \
    || ! printf '%s\n' "$associate_complete_output" | grep -qxF "$(printf 'REQ-501\tlegacy-file.txt')" \
    || ! printf '%s\n' "$associate_complete_output" | grep -qxF "$(printf 'REQ-501\tsecond-file.txt')" \
    || ! printf '%s\n' "$associate_complete_output" | grep -qxF -- "$(printf -- '-\tphantom-file.txt')"; then
  printf 'FAIL: tools/checks/associate-files.sh must associate every path on a multi-path bullet, preserve root-level filenames, ignore prose-only backticks, and normalize the documented terminal-success alias.\n' >&2
  fail_count=$((fail_count + 1))
fi
rm -rf -- "$associate_complete_probe_dir"

associate_unmatched_probe_dir="$(mktemp -d)"
mkdir -p "$associate_unmatched_probe_dir/do-work/archive/UR-301"
cat > "$associate_unmatched_probe_dir/do-work/archive/UR-301/REQ-502-unmatched-summary.md" <<'EOF'
---
id: REQ-502
status: completed
completed_at: 2026-08-07T12:00:00Z
---

## Implementation Summary

**Files changed:**
- `legacy-file.txt`, `unclosed-file.txt
EOF
if associate_unmatched_output="$(printf 'legacy-file.txt\n' | "$core_root/tools/checks/associate-files.sh" --repo-root "$associate_unmatched_probe_dir" 2>&1)"; then
  associate_unmatched_exit=0
else
  associate_unmatched_exit=$?
fi
if [ "$associate_unmatched_exit" -ne 2 ] || ! grep -qF 'PARSE-FAILED:' <<<"$associate_unmatched_output"; then
  printf 'FAIL: tools/checks/associate-files.sh must fail loudly with exit 2 when a path-led Implementation Summary bullet has an unmatched backtick; got exit %s: %s\n' \
    "$associate_unmatched_exit" "$associate_unmatched_output" >&2
  fail_count=$((fail_count + 1))
fi
rm -rf -- "$associate_unmatched_probe_dir"

# A git-status failure is not a clean tree. Process substitution used to hide
# the producer's failure, making a bare repository return the clean-tree exit 1.
inventory_failure_probe_dir="$(mktemp -d)"
if ! GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
  git init -q --bare "$inventory_failure_probe_dir"; then
  printf 'FAIL: could not set up the uncommitted-inventory git-status failure probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  if inventory_failure_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_failure_probe_dir" 2>&1)"; then
    inventory_failure_exit=0
  else
    inventory_failure_exit=$?
  fi
  if [ "$inventory_failure_exit" -ne 2 ]; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must return exit 2 when git status fails, got %s.\n' \
      "$inventory_failure_exit" >&2
    fail_count=$((fail_count + 1))
  fi
  if ! grep -qF 'STATUS-FAILED:' <<<"$inventory_failure_output"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must diagnose a git-status failure instead of reporting a clean tree.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$inventory_failure_probe_dir"

assert_contains \
  "tools/checks/preflight.sh" \
  'git -c status\.renames=copies status --porcelain=v1 --untracked-files=all -z' \
  'tools/checks/preflight.sh must keep NUL-safe repository-wide dirty-file detection — serial qualification and review read the repository-wide working/staged diff.'

preflight_dirty_warning_line="$(grep -E '^[[:space:]]*echo "WARN: .*uncommitted changes' "$core_root/tools/checks/preflight.sh" || true)"

assert_block_contains \
  "$preflight_dirty_warning_line" \
  'unless they prevent the active REQ' \
  'tools/checks/preflight.sh dirty-file warning must stay subordinate to Current-REQ relevance instead of turning unexpected state into a blocker.'

assert_block_contains \
  "$preflight_dirty_warning_line" \
  'qualification/review evidence' \
  'tools/checks/preflight.sh must explain dirty files as possible qualification/review evidence contamination, not as files the explicit commit step will automatically stage.'

assert_block_not_contains \
  "$preflight_dirty_warning_line" \
  'may stage unrelated files|swept into (the )?commit' \
  'tools/checks/preflight.sh must not claim dirty files are automatically staged — Step 9 stages explicit paths.'

preflight_template_block="$(sed -n '/^## Pre-Flight Template/,/^## Implementation Summary Template/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$preflight_template_block" \
  'preserve/exclude from this REQ.*qualification/review evidence' \
  'actions/work-reference.md Pre-Flight template must frame pre-existing dirty files as qualification/review evidence contamination.'

assert_file_not_contains \
  "actions/work.md" \
  'unrelated dirty files get swept into the commit' \
  'actions/work.md must not teach that pre-existing dirty files are swept into the explicit per-REQ commit.'

# Behavioral probes for tools/checks/scope-drift.sh. The Scope header may carry
# trailing annotations before the colon (REQ-178 wrote "**Files I will touch
# (all new, …):**"); both match sites — path extraction and the unparseable-header
# guard — must recognize it, so an annotated header either parses into a real
# comparison or FAILs loudly. It must never degrade to the Route A "no Scope
# list" SKIP, which silently disables the check.
scope_drift_probe_dir="$(mktemp -d)"
cat > "$scope_drift_probe_dir/annotated-header.md" <<'EOF'
---
id: REQ-900
status: working
---

## Scope

**Files I will touch (all new):**

- `tools/example-check.sh` (create)

**Files I will NOT touch:**

- `README.md` (out of scope)

## Implementation Summary

**Files changed:**

- `tools/example-check.sh` (created)
EOF
if scope_drift_annotated_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/annotated-header.md" 2>&1)"; then
  scope_drift_annotated_exit=0
else
  scope_drift_annotated_exit=$?
fi
if [ "$scope_drift_annotated_exit" -ne 0 ] || ! grep -qF 'OK:' <<<"$scope_drift_annotated_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must run a real comparison on an annotated touch-list header (exit 0 on matching sets, never SKIP); got exit %s: %s\n' \
    "$scope_drift_annotated_exit" "$scope_drift_annotated_output" >&2
  fail_count=$((fail_count + 1))
fi

cat > "$scope_drift_probe_dir/annotated-unparseable.md" <<'EOF'
---
id: REQ-901
status: working
---

## Scope

**Files I will touch (all new):** tools/example-check.sh without backticks

## Implementation Summary

**Files changed:**

- `tools/example-check.sh` (created)
EOF
if scope_drift_unparseable_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/annotated-unparseable.md" 2>&1)"; then
  scope_drift_unparseable_exit=0
else
  scope_drift_unparseable_exit=$?
fi
if [ "$scope_drift_unparseable_exit" -ne 1 ] || ! grep -qF 'FAIL:' <<<"$scope_drift_unparseable_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must FAIL (exit 1) when an annotated touch-list header yields zero parseable paths, never SKIP; got exit %s: %s\n' \
    "$scope_drift_unparseable_exit" "$scope_drift_unparseable_output" >&2
  fail_count=$((fail_count + 1))
fi

cat > "$scope_drift_probe_dir/no-scope-section.md" <<'EOF'
---
id: REQ-902
status: working
---

## Implementation Summary

**Files changed:**

- `tools/example-check.sh` (created)
EOF
if scope_drift_absent_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/no-scope-section.md" 2>&1)"; then
  scope_drift_absent_exit=0
else
  scope_drift_absent_exit=$?
fi
if [ "$scope_drift_absent_exit" -ne 2 ] || ! grep -qF 'SKIP:' <<<"$scope_drift_absent_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must keep the SKIP exit 2 contract when no Scope touch-list exists (Route A REQs rely on it); got exit %s: %s\n' \
    "$scope_drift_absent_exit" "$scope_drift_absent_output" >&2
  fail_count=$((fail_count + 1))
fi

# REQ-344 touched nine files against two declarations. Its seven undeclared paths
# were grouped behind two bullet-leading paths, so the old first-token parser
# reported only two of seven. Pin the complete set, not just a nonzero exit.
cat > "$scope_drift_probe_dir/req-344-multi-path.md" <<'EOF'
---
id: REQ-903
status: working
---

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/frontmatter.go` (modify)
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go` (modify)

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/code-review.md`, `skills/do-work/actions/capture-reference.md`, `skills/do-work/actions/capture.md`, `skills/do-work/actions/clarify.md`, `skills/do-work-board/tools/queue-kanban/frontmatter.go` (modified)
- `skills/do-work/actions/work-reference.md`, `skills/do-work/actions/work.md`, `skills/do-work/docs/capture-guide.md`, `skills/do-work-board/tools/queue-kanban/frontmatter_test.go` (modified)
EOF
if scope_drift_req344_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/req-344-multi-path.md" 2>&1)"; then
  scope_drift_req344_exit=0
else
  scope_drift_req344_exit=$?
fi
scope_drift_req344_expected_paths=(
  'skills/do-work-toolbox/actions/code-review.md'
  'skills/do-work/actions/capture-reference.md'
  'skills/do-work/actions/capture.md'
  'skills/do-work/actions/clarify.md'
  'skills/do-work/actions/work-reference.md'
  'skills/do-work/actions/work.md'
  'skills/do-work/docs/capture-guide.md'
)
scope_drift_req344_missing=0
for scope_drift_req344_path in "${scope_drift_req344_expected_paths[@]}"; do
  if ! grep -qxF "  $scope_drift_req344_path" <<<"$scope_drift_req344_output"; then
    scope_drift_req344_missing=1
  fi
done
if [ "$scope_drift_req344_exit" -ne 1 ] \
    || [ "$scope_drift_req344_missing" -ne 0 ] \
    || ! grep -qF 'DRIFT: touched but never declared in ## Scope:' <<<"$scope_drift_req344_output" \
    || grep -qF 'DRIFT: declared in ## Scope but never touched:' <<<"$scope_drift_req344_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must report REQ-344\047s exact seven undeclared paths from multi-path bullets and no false unused declarations; got exit %s: %s\n' \
    "$scope_drift_req344_exit" "$scope_drift_req344_output" >&2
  fail_count=$((fail_count + 1))
fi

# Both sides must consume every closed backtick pair. Matching first tokens make
# this fixture a false OK if either parser falls back to its first-token behavior;
# root-level filenames prove that a slash heuristic cannot silently disarm it.
cat > "$scope_drift_probe_dir/symmetric-multi-path.md" <<'EOF'
---
id: REQ-904
status: working
---

## Scope

**Files I will touch:**
- `src/shared.sh`, `scope-only.txt`, `justfile` (modify)

## Implementation Summary

**Files changed:**
- `src/shared.sh`, `summary-only.txt`, `README.md` (modified)
EOF
if scope_drift_symmetric_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/symmetric-multi-path.md" 2>&1)"; then
  scope_drift_symmetric_exit=0
else
  scope_drift_symmetric_exit=$?
fi
scope_drift_symmetric_expected=(
  '  README.md'
  '  summary-only.txt'
  '  justfile'
  '  scope-only.txt'
)
scope_drift_symmetric_missing=0
for scope_drift_symmetric_path in "${scope_drift_symmetric_expected[@]}"; do
  if ! grep -qxF "$scope_drift_symmetric_path" <<<"$scope_drift_symmetric_output"; then
    scope_drift_symmetric_missing=1
  fi
done
if [ "$scope_drift_symmetric_exit" -ne 1 ] \
    || [ "$scope_drift_symmetric_missing" -ne 0 ] \
    || ! grep -qF 'DRIFT: touched but never declared in ## Scope:' <<<"$scope_drift_symmetric_output" \
    || ! grep -qF 'DRIFT: declared in ## Scope but never touched:' <<<"$scope_drift_symmetric_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must compare every later path on both multi-path bullets, including filename-only paths; got exit %s: %s\n' \
    "$scope_drift_symmetric_exit" "$scope_drift_symmetric_output" >&2
  fail_count=$((fail_count + 1))
fi

cat > "$scope_drift_probe_dir/matching-multi-path.md" <<'EOF'
---
id: REQ-905
status: working
---

## Scope

**Files I will touch:**
- `src/shared.sh`, `.gitignore`, `justfile` (modify)

## Implementation Summary

**Files changed:**
- `src/shared.sh`, `.gitignore`, `justfile` (modified)
EOF
if scope_drift_matching_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/matching-multi-path.md" 2>&1)"; then
  scope_drift_matching_exit=0
else
  scope_drift_matching_exit=$?
fi
if [ "$scope_drift_matching_exit" -ne 0 ] || ! grep -qF 'OK:' <<<"$scope_drift_matching_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must accept identical multi-path lists containing root-level filenames; got exit %s: %s\n' \
    "$scope_drift_matching_exit" "$scope_drift_matching_output" >&2
  fail_count=$((fail_count + 1))
fi

cat > "$scope_drift_probe_dir/prose-backticks.md" <<'EOF'
---
id: REQ-906
status: working
---

## Scope

**Files I will touch:**
- `src/shared.sh` (modify)

## Implementation Summary

**Files changed:**
- `src/shared.sh` (modified)
- Notes mention `README.md` and `sort -u`, but this prose bullet claims no file.
EOF
if scope_drift_prose_output="$("$core_root/tools/checks/scope-drift.sh" \
    "$scope_drift_probe_dir/prose-backticks.md" 2>&1)"; then
  scope_drift_prose_exit=0
else
  scope_drift_prose_exit=$?
fi
if [ "$scope_drift_prose_exit" -ne 0 ] || ! grep -qF 'OK:' <<<"$scope_drift_prose_output"; then
  printf 'FAIL: tools/checks/scope-drift.sh must ignore backticked spans on prose-only bullets; got exit %s: %s\n' \
    "$scope_drift_prose_exit" "$scope_drift_prose_output" >&2
  fail_count=$((fail_count + 1))
fi

for scope_drift_unmatched_side in scope summary; do
  cat > "$scope_drift_probe_dir/unmatched-$scope_drift_unmatched_side.md" <<EOF
---
id: REQ-907
status: working
---

## Scope

**Files I will touch:**
$(if [ "$scope_drift_unmatched_side" = scope ]; then printf '%s\n' '- `src/shared.sh`, `unclosed.txt'; else printf '%s\n' '- `src/shared.sh` (modify)'; fi)

## Implementation Summary

**Files changed:**
$(if [ "$scope_drift_unmatched_side" = summary ]; then printf '%s\n' '- `src/shared.sh`, `unclosed.txt'; else printf '%s\n' '- `src/shared.sh` (modified)'; fi)
EOF
  if scope_drift_unmatched_output="$("$core_root/tools/checks/scope-drift.sh" \
      "$scope_drift_probe_dir/unmatched-$scope_drift_unmatched_side.md" 2>&1)"; then
    scope_drift_unmatched_exit=0
  else
    scope_drift_unmatched_exit=$?
  fi
  if [ "$scope_drift_unmatched_exit" -ne 1 ] || ! grep -qF 'FAIL:' <<<"$scope_drift_unmatched_output" \
      || ! grep -qF 'unmatched backtick' <<<"$scope_drift_unmatched_output"; then
    printf 'FAIL: tools/checks/scope-drift.sh must fail loudly when the %s path list has an unmatched backtick; got exit %s: %s\n' \
      "$scope_drift_unmatched_side" "$scope_drift_unmatched_exit" "$scope_drift_unmatched_output" >&2
    fail_count=$((fail_count + 1))
  fi
done
rm -rf -- "$scope_drift_probe_dir"

# Review regressions: prescribed shell and roadmap classification are runtime
# contracts even though they live in Markdown/just recipes rather than compiled code.
extract_kanban_shutdown_line() {
  awk '
    /^run-kanban \$port=/ {
      getline
      getline
      sub(/^[[:space:]]*/, "")
      print
      exit
    }
  ' "$repo_root/$1"
}

root_kanban_shutdown_line="$(extract_kanban_shutdown_line justfile)"
installer_kanban_shutdown_line="$(extract_kanban_shutdown_line skills/do-work-board/justfile.template)"
if [ "$root_kanban_shutdown_line" != "$installer_kanban_shutdown_line" ]; then
  printf 'FAIL: Justfile and the board-owned template must carry one identical run-kanban shutdown line.\n' >&2
  fail_count=$((fail_count + 1))
fi

for kanban_recipe_file in "skills/do-work-board/justfile.template" "justfile"; do
  assert_file_not_contains \
    "$kanban_recipe_file" \
    'case "\$listener_command" in \*queue-kanban\*' \
    "$kanban_recipe_file must not identify a stale board from arbitrary argv text."
  assert_contains \
    "$kanban_recipe_file" \
    'lsof -a -p "\$listener_pid" -d txt -Fn' \
    "$kanban_recipe_file must identify a stale board from its executable, preserving cross-repo binary names without matching unrelated arguments."
  assert_contains \
    "$kanban_recipe_file" \
    'wait_count.*-lt 320.*lsof -a -p "\$listener_pid".*tcp:"\$port".*-sTCP:LISTEN' \
    "$kanban_recipe_file must bound a listener-specific wait at 320 iterations instead of polling process existence."
  assert_contains \
    "$kanban_recipe_file" \
    'remaining_listener_pid=.*lsof -ti tcp:"\$port".*-sTCP:LISTEN' \
    "$kanban_recipe_file must query the port again after the bounded shutdown wait."
  assert_contains \
    "$kanban_recipe_file" \
    'ps -p "\$remaining_listener_pid" -o args=' \
    "$kanban_recipe_file must resolve the remaining listener's full command for the refusal."
  assert_file_not_contains \
    "$kanban_recipe_file" \
    'while kill -0 "\$listener_pid"' \
    "$kanban_recipe_file must wait on listener ownership, not process existence."
  assert_contains \
    "$kanban_recipe_file" \
    '^run-do-work-update:' \
    "$kanban_recipe_file must ship the project-local do-work update shortcut."
done

# Execute the canonical shutdown line with command seams. A queue-kanban PID that remains a
# listener throughout the bounded wait must make the recipe line fail before build+serve, and
# the diagnosis must name both the PID and full command. sleep is a no-op in the fixture, so all
# 320 iterations run without slowing the suite.
if stuck_listener_output="$(
  port=8090
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/tmp/queue-kanban\n'
    else
      printf '4242\n'
    fi
  }
  ps() { printf '/tmp/queue-kanban serve --port 8090\n'; }
  kill() { return 0; }
  sleep() { return 0; }
  eval "$root_kanban_shutdown_line" 2>&1
)"; then
  stuck_listener_status=0
else
  stuck_listener_status=$?
fi
if [ "$stuck_listener_status" -ne 1 ]; then
  printf 'FAIL: run-kanban shutdown must refuse startup when a listener remains after the bounded wait; got exit %s.\n' "$stuck_listener_status" >&2
  fail_count=$((fail_count + 1))
fi
if ! printf '%s\n' "$stuck_listener_output" | grep -qF 'pid 4242'; then
  printf 'FAIL: run-kanban stuck-listener refusal must name pid 4242.\n' >&2
  fail_count=$((fail_count + 1))
fi
if ! printf '%s\n' "$stuck_listener_output" | grep -qF '/tmp/queue-kanban serve --port 8090'; then
  printf 'FAIL: run-kanban stuck-listener refusal must name the listener command.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Preserve the older safety boundary: a foreign executable is refused immediately and never
# passed to kill. This behavior probe complements the executable-identity contract above.
foreign_kill_marker="$(mktemp)"
if foreign_listener_output="$(
  port=8090
  # The shutdown recipe reads this through eval, which ShellCheck cannot trace.
  : "$port"
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/usr/bin/python3\n'
    else
      printf '3131\n'
    fi
  }
  ps() { printf '/usr/bin/python3 -m http.server 8090\n'; }
  kill() { printf 'called\n' > "$foreign_kill_marker"; }
  sleep() { return 0; }
  eval "$root_kanban_shutdown_line" 2>&1
)"; then
  foreign_listener_status=0
else
  foreign_listener_status=$?
fi
# This fixture asserts the refusal status and kill boundary; its captured text is intentionally ignored.
: "$foreign_listener_output"
if [ "$foreign_listener_status" -ne 1 ] || [ -s "$foreign_kill_marker" ]; then
  printf 'FAIL: run-kanban must refuse a foreign listener without calling kill.\n' >&2
  fail_count=$((fail_count + 1))
fi
rm -f "$foreign_kill_marker"

if [ ! -x "$core_root/tools/do-work-update.sh" ]; then
  printf 'FAIL: modular core tools/do-work-update.sh must be executable for the Just shortcut.\n' >&2
  fail_count=$((fail_count + 1))
fi
assert_contains \
  "tools/do-work-update.sh" \
  "--project-root" \
  'tools/do-work-update.sh must derive and validate the consuming project root before updating.'
assert_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go" \
  '"tar", "xzf", upstreamTarball, "-C", freshUpstream' \
  'The update transaction must extract only into staging; behavioral probes verify runtime do-work data is outside every managed destination.'
assert_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go" \
  'Install this complete four-skill suite' \
  'The installed suite transaction must require confirmation after showing its reviewed diff.'
assert_file_not_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go" \
  '"cp", "-R", .*skillRoot' \
  'The update transaction must not reintroduce the pre-update rollback copy — git is the undo, and a duplicated tree on every run buys nothing git does not already hold. A mid-update failure reports the partial install instead; see _dev/tests/update-script-behavior.sh.'
assert_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go" \
  'runRecoveryIfNeeded' \
  'The installed suite transaction must automatically recover its validated managed paths after a destructive-region failure.'

assert_file_not_contains \
  "scripts/run-blocked-check.sh" \
  'else probe_wrapper=""' \
  'run-blocked-check.sh must not drop the blocked-check time limit when timeout/gtimeout is unavailable.'

assert_contains \
  "scripts/run-blocked-check.sh" \
  'exit 124' \
  'run-blocked-check.sh must preserve a bounded portable fallback and report a timed-out blocked check as exit 124.'

assert_contains \
  "scripts/run-blocked-check.sh" \
  'probe_group_is_safe' \
  'run-blocked-check.sh must verify an isolated process group before executing fallback probe code.'
assert_contains \
  "scripts/run-blocked-check.sh" \
  'kill -"\$probe_signal" -- "-\$probe_process_group_id"' \
  'run-blocked-check.sh must signal the verified probe group rather than only its wrapper PID.'

assert_contains \
  "actions/install.md" \
  'non-empty `SKILL.md` and `<project-root>/\.claude/skills/last30days/scripts/last30days\.py`' \
  'actions/install.md must define last30days health as a complete runnable payload.'
assert_contains \
  "actions/install.md" \
  'never merges files from different upstream versions' \
  'actions/install.md must document last30days publication as a validated directory transaction.'

if [ ! -s "$repo_root/_dev/lessons/validated-runtime-boundaries.md" ] \
  || ! grep -Fq '../lessons/validated-runtime-boundaries.md' "$repo_root/_dev/primes/prime-shell-commands.md"; then
  printf 'FAIL: prime-shell-commands.md must link the durable validated-runtime-boundaries lesson.\n' >&2
  fail_count=$((fail_count + 1))
fi

if ! bash -n "$core_root/scripts/run-blocked-check.sh"; then
  printf 'FAIL: scripts/run-blocked-check.sh must remain syntactically valid.\n' >&2
  fail_count=$((fail_count + 1))
fi
assert_contains \
  "actions/work.md" \
  'run-blocked-check.sh' \
  'actions/work.md must invoke the shipped blocked-check timeout runner.'

# Target ID Resolution contract (REQ-067). The run tokenizer recognized exactly one token
# shape (REQ- + digits); every other token was residue by construction, so a UR argument was
# handled by whatever the reading agent improvised. The shared contract in work-reference.md is
# the single definition of both token shapes and UR->REQ expansion that run/abandon/
# roadmap cite instead of restating; work.md's Input and usage string must offer the UR- shape.
assert_contains \
  "actions/work-reference.md" \
  '### Target ID Resolution' \
  'actions/work-reference.md must define the shared Target ID Resolution contract (REQ-/UR- token shapes + UR->REQ expansion by user_request: scan) that the id-taking actions cite instead of restating.'

work_input_block="$(sed -n '/^## Input/,/^## Steps/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_input_block" \
  'UR-NNN' \
  'actions/work.md Input must accept UR-NNN targeting tokens (usage string offers UR-NNN and the tokenizer recognizes the UR- shape) per the Target ID Resolution contract — a UR argument must no longer fall through to the unrecognized-argument guard.'

# UR ids accepted by abandon (REQ-068). The action keyed entirely on REQ-NNN tokens; a UR
# argument had no defined handling (abandon's globs substitute a REQ number). Its Input section
# must name the UR- shape and cite the shared Target ID Resolution contract rather than restating
# the resolution rule. Scope the UR- check to the Input section — abandon.md already names
# archive/UR-NNN/ folders elsewhere, so a file-wide grep would pass vacuously without the action
# actually accepting a UR argument.
abandon_input_block="$(sed -n '/^## Input/,/^## Steps/p' "$core_root/actions/abandon.md")"

assert_block_contains \
  "$abandon_input_block" \
  'Target ID Resolution' \
  'actions/abandon.md Input must cite the Target ID Resolution contract so it accepts UR-NNN targeting tokens (a UR cancels its cancellable members) rather than only REQ-NNN.'

# A legacy failed REQ can already sit inside a closed archive/UR-NNN/ folder. The board still
# treats that status as active, so explicitly targeting the REQ must reach abandon's existing
# confirmed in-place failed->cancelled write without moving or reopening the closed UR.
assert_contains \
  "actions/abandon.md" \
  'status: failed.*archive/UR-NNN/.*cancellable in place' \
  'actions/abandon.md must let an explicitly targeted failed REQ inside archive/UR-NNN/ use the confirmed in-place cancellation path; refusing it leaves a permanent active board card.'

assert_file_not_contains \
  "actions/abandon.md" \
  'status: failed.*archive/UR-NNN/.*refuse' \
  'actions/abandon.md must not restore the legacy nested-failure refusal after the in-place resolution path was extended to that shape.'

assert_contains \
  "actions/forensics.md" \
  'archive/UR-NNN/.*cancel.*in place' \
  'actions/forensics.md must route a failed REQ already inside archive/UR-NNN/ to explicit in-place abandon instead of claiming no action is needed.'

assert_contains \
  "actions/roadmap.md" \
  '^-[[:space:]]+\*\*Ready\*\*[[:space:]]+— normalized `status` is `pending`' \
  'actions/roadmap.md must require pending status before classifying a queued REQ as Ready.'

# REQ scoping in roadmap (REQ-070). roadmap accepted UR-NNN as a scope token but not REQ-NNN, so
# `do-work roadmap REQ-067` silently fell through to a whole-queue survey — the inverse of the
# UR-run asymmetry the batch fixed. Its Input must now name a REQ-NNN scope token and cite the
# shared Target ID Resolution contract for the token shape (rather than restating it).
roadmap_input_block="$(sed -n '/^## Input/,/^## Steps/p' "$core_root/actions/roadmap.md")"

assert_block_contains \
  "$roadmap_input_block" \
  'REQ-NNN' \
  'actions/roadmap.md Input must accept a REQ-NNN scope token (single-REQ survey), not only UR-NNN — a recognized REQ id must no longer fall through to the full-survey default.'

assert_block_contains \
  "$roadmap_input_block" \
  'Target ID Resolution' \
  'actions/roadmap.md Input must cite the shared Target ID Resolution contract for the REQ-/UR- token shapes rather than restating them.'

# ADR-017 memory engine contracts: the Stop capture must never block a session
# end, destructive consolidation must never be hook-wired, and the engine's
# core guardrails (cap, prompt-injection load, compose-don't-clobber install)
# must stay stated where agents read them.
assert_contains \
  "hooks/memory-stop-capture.sh" \
  'stop_hook_active' \
  'hooks/memory-stop-capture.sh must keep the stop_hook_active loop guard.'

assert_file_not_contains \
  "hooks/memory-stop-capture.sh" \
  '"decision":[[:space:]]*"block"' \
  'hooks/memory-stop-capture.sh must never emit a blocking decision — capture cannot hold a session open.'

for hooks_json_file in "hooks/hooks.json" "hooks/memory-hooks.json"; do
  assert_file_not_contains \
    "$hooks_json_file" \
    'dream' \
    "$hooks_json_file must never wire the destructive dream action to a hook."
done

assert_contains \
  "actions/memory.md" \
  '2,500' \
  'actions/memory.md must state the 2,500-character working-memory hard cap.'

assert_contains \
  "actions/memory.md" \
  'prompt-injection\.md' \
  'actions/memory.md recall must load the prompt-injection guardrail before reading hook-captured log content.'

assert_contains \
  "actions/setup-memory.md" \
  'install memory-module' \
  'knowledge setup-memory must retain the documented memory-module migration alias (ADR-017).'

assert_contains \
  "actions/setup-memory.md" \
  'settings\.json\.pre-memory-module' \
  'knowledge setup-memory hook merge must back up settings.json before composing entries.'

assert_contains \
  "hooks/memory-stop-capture.sh" \
  'REDACTED' \
  'hooks/memory-stop-capture.sh must redact credential-shaped text before persisting captures — defense in depth behind the machine-local store.'

# Raw captures and the per-machine ledger must never become committable: the installer
# adds them to .git/info/exclude (machine-local), never the project's .gitignore.
assert_contains \
  "actions/setup-memory.md" \
  '\*\*/memory/logs/' \
  'knowledge setup-memory must add memory/logs/ to .git/info/exclude — verbatim captures must not be committable.'

assert_contains \
  "actions/setup-memory.md" \
  '\*\*/memory/usage-ledger\.jsonl' \
  'knowledge setup-memory must add memory/usage-ledger.jsonl to .git/info/exclude.'

assert_contains \
  "hooks/memory-session-start.sh" \
  'session capture' \
  'hooks/memory-session-start.sh must strip raw session-capture sections from the startup injection — unvetted transcript text must not enter context before a prompt-injection guard can load.'

# CLAUDE.md/AGENTS.md are the maintainer doc, export-ignored since 0.136.0 so they never
# land in consumer installs (nested CLAUDE.md is auto-loaded into consumer agents' context).
assert_contains \
  ".gitattributes" \
  '^/CLAUDE\.md[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /CLAUDE.md — the maintainer doc must not ship to consumer installs.'

assert_contains \
  ".gitattributes" \
  '^/AGENTS\.md[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /AGENTS.md — the redirect stub must not ship to consumer installs.'

# do-work/ and kb/ are TRACKED in this repo (they are the same Trail of Intent the skill tells
# consumers to commit, and the tracked-path-only data-loss guards in tools/checks/record-commit-hash.sh
# plus cleanup Pass 6 blanked-REQ recovery only work on tracked REQs). That makes these two
# export-ignore lines the only barrier between the maintainer's queue/KB and every consumer
# install — the pre-0.157.0 blanket .git/info/exclude entry no longer backstops them.
assert_contains \
  ".gitattributes" \
  '^/do-work[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /do-work — this repo tracks its own queue, so without this line the maintainer archive ships to every consumer install.'

assert_contains \
  ".gitattributes" \
  '^/kb[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /kb — this repo tracks its own knowledge base, so without this line it ships to every consumer install.'

assert_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go" \
  'moduleInstallPlan' \
  'The installed suite transaction must construct an explicit managed module plan; behavioral probes verify the project knowledge base is outside it.'

# Shipped files must not cite the skill's own CLAUDE.md/AGENTS.md — those files are absent
# downstream, so a citation dangles. The full rule lives in _dev/primes/prime-action-files.md → Cross-Referencing.
#
# The check is INVERTED on purpose: it flags ANY mention of CLAUDE.md/AGENTS.md in a shipped
# path, and exempts a short per-file allowlist. It used to enumerate citation idioms
# (`see CLAUDE.md`, `per CLAUDE.md`, `CLAUDE.md →`) and that shape failed exactly as
# "Closed Enumerations Go Stale" predicts: scored against the seven real occurrences in the
# tree it caught ZERO — including the `actions/memory-reference.md` line REQ-088 was filed to
# fix, and six `CLAUDE.md § Before Every Commit` comments in tools/queue-kanban that shipped
# for months. Extending the idiom list would have caught 4 of 6 and false-positived on
# actions/prime.md, where the colon is sentence punctuation. A mention is cheap to detect and
# impossible to phrase around; the judgement lives in the allowlist, where it is visible.
#
# The allowlist is PER-FILE, never per-directory: allowlisting `actions/` wholesale would have
# exempted actions/memory-reference.md, the file that started this. Entries are files whose
# subject genuinely IS a consumer project's own CLAUDE.md/AGENTS.md (prime registries, the KB
# schema file, and tidy-repo's layout rules) — those
# references are correct and must not be "fixed". A new shipped file that mentions the
# maintainer doc fails this check until someone decides which of the two it is; that decision
# is the point.
shipped_citation_paths=(SKILL.md next-steps.md README.md actions crew-members prompts interviews specs docs hooks tools)
maintainer_doc_mention_allowlist=(
  tools/install-do-work-suite.sh
  actions/prime.md
  docs/prime-guide.md
  actions/version.md
  actions/tidy-repo.md
  actions/bkb.md
  actions/bkb-reference.md
  docs/bkb-guide.md
  actions/capture.md
  actions/validate-feedback.md
  actions/prompts.md
  README.md
  prompts/README.md
  prompts/prompt-kit-step2-personal-context-doc.md
)
maintainer_doc_mentions="$(cd "$repo_root" && grep -rIn 'CLAUDE\.md\|AGENTS\.md' "${shipped_citation_paths[@]}" 2>/dev/null || true)"
unallowed_maintainer_doc_mentions=""
while IFS= read -r maintainer_doc_hit; do
  [ -n "$maintainer_doc_hit" ] || continue
  maintainer_doc_hit_path="${maintainer_doc_hit%%:*}"
  maintainer_doc_hit_allowed=0
  for allowlisted_citation_file in "${maintainer_doc_mention_allowlist[@]}"; do
    if [ "$maintainer_doc_hit_path" = "$allowlisted_citation_file" ]; then
      maintainer_doc_hit_allowed=1
      break
    fi
  done
  if [ "$maintainer_doc_hit_allowed" -eq 0 ]; then
    unallowed_maintainer_doc_mentions="${unallowed_maintainer_doc_mentions}${maintainer_doc_hit}"$'\n'
  fi
done <<< "$maintainer_doc_mentions"
if [ -n "$unallowed_maintainer_doc_mentions" ]; then
  printf 'FAIL: shipped files must not mention the skill'\''s own CLAUDE.md/AGENTS.md (export-ignored — absent in consumer installs). Restate the rule inline, point at a shipped home, or — if the mention really is about a *consumer project'\''s* CLAUDE.md — add the file to maintainer_doc_mention_allowlist in this script:\n%s' "$unallowed_maintainer_doc_mentions" >&2
  fail_count=$((fail_count + 1))
fi

# Prospective half of the defensive-surface discipline (REQ-169). validate-feedback must
# challenge only remedies that ADD long-lived defense, make the cost call affect the verdict,
# and expose that call per finding; ordinary fixes/deletions/simplifications stay unchanged.
assert_contains \
  "actions/validate-feedback.md" \
  'what incident earned this, and is the fix still cheaper than the surface it added\?' \
  'validate-feedback must ask the user-supplied surface-cost rubric verbatim for remedies that add defense.'

assert_contains \
  "actions/validate-feedback.md" \
  'guard, fallback, retry, validation layer, rule, or warning apparatus' \
  'validate-feedback must define the surface-adding remedy boundary explicitly.'

assert_contains \
  "actions/validate-feedback.md" \
  'Direct bug fixes, deletions, and simplifications.*N/A' \
  'validate-feedback must leave non-surface-adding remedies outside the new skepticism pass.'

assert_contains \
  "actions/validate-feedback.md" \
  'must not receive.*Accept.*Push back.*Discuss' \
  'validate-feedback must prevent an unearned/net-costly added defense from receiving a plain Accept verdict.'

assert_contains \
  "actions/validate-feedback.md" \
  '\*\*Surface-cost:\*\* N/A / Earned / Flagged' \
  'validate-feedback per-finding output must expose the rubric result to the reader.'

# Review finding-follow-up closure seam (REQ-170 remediation). The review consumer now
# rejects a finding-origin REQ without matching closure evidence, so its own follow-up
# producer must emit the canonical proof shape instead of creating work that cannot pass.
review_finding_closure_gate_block="$(sed -n '/^6\. \*\*Enforce finding closure\*\*/,/^\*\*What NOT to do:\*\*/p' "$core_root/actions/review-work.md")"
review_generated_followup_template_block="$(sed -n '/^For each finding that routes to its own REQ/,/^\*\*When the root cause is ambiguous requirements/p' "$core_root/actions/review-work.md")"

assert_block_contains \
  "$review_finding_closure_gate_block" \
  'captured GREEN.*fails before/passes after.*exact named finding surface was deleted' \
  'review-work finding-closure consumer must keep the matching captured-GREEN-or-deletion gate.'

assert_block_contains \
  "$review_generated_followup_template_block" \
  '^review_generated: true$' \
  'review-work review-follow-up template probe must stay scoped to the review_generated producer.'

for proof_field in 'RED prompt/case' 'Why RED now' 'GREEN when' 'Validation'; do
  assert_block_contains \
    "$review_generated_followup_template_block" \
    "^\\*\\*${proof_field}:\\*\\*" \
    "review-work review-generated follow-up template must emit the canonical Red-Green Proof field: ${proof_field}."
done

assert_block_contains \
  "$review_generated_followup_template_block" \
  'RED prompt/case:.*Named regression test/check that fails before the fix.*exact finding surface to delete' \
  'review-work review-generated follow-up template must name a fail-before regression check or exact deletion surface.'

assert_block_contains \
  "$review_generated_followup_template_block" \
  'GREEN when:.*same named test/check passes after the fix.*exact named finding surface is absent' \
  'review-work review-generated follow-up template must pair the producer RED with matching pass-after or deletion GREEN.'

assert_block_contains \
  "$review_generated_followup_template_block" \
  'Validation:.*actions/work-reference\.md.*Finding-Closure Ratchet' \
  'review-work review-generated follow-up template must cite the canonical Finding-Closure Ratchet from its proof block.'

# Every shipped action template that emits the exact marker inherits the same consumer
# gate. Enumerate fenced producer blocks from the shipped actions so a future package or
# producer cannot appear outside this ratchet merely because its filename was not listed.
review_generated_field_count="$(
  find "$repo_root/skills" -type f -path '*/actions/*.md' -name '*.md' \
    -exec grep -h '^review_generated: true$' {} + \
    | wc -l \
    | tr -d '[:space:]'
)"
review_generated_template_count=0

while IFS=$'\t' read -r producer_file producer_line has_proof_shape has_named_red has_matching_green has_destination_finding_closure_citation has_pau_block; do
  producer_relative_path="${producer_file#"$repo_root/"}"
  review_generated_template_count=$((review_generated_template_count + 1))

  # REQ-264: every template that mints a buildable REQ must emit the P-A-U block. Both of
  # these produce `status: pending` work, and qualify.sh's Check 4 is DISARMED — not passed
  # — on a REQ that carries no such section, because each of its FAIL branches keys on the
  # box's state and a missing box satisfies all of them by absence. qualify now WARNs, and
  # this is the other half: stop minting the REQs that trip it. Six queued REQs were
  # created without the block before this assertion existed.
  if [ "$has_pau_block" -ne 1 ]; then
    printf 'FAIL: %s:%s review_generated template must emit the AI Execution State (P-A-U Loop) block with all three boxes; without it qualify.sh Check 4 is disarmed on every REQ it mints (REQ-264).\n' \
      "$producer_relative_path" "$producer_line" >&2
    fail_count=$((fail_count + 1))
  fi

  if [ "$has_proof_shape" -ne 1 ]; then
    printf 'FAIL: %s:%s review_generated template must emit the canonical four-field Red-Green Proof shape.\n' \
      "$producer_relative_path" "$producer_line" >&2
    fail_count=$((fail_count + 1))
  fi
  if [ "$has_named_red" -ne 1 ]; then
    printf 'FAIL: %s:%s review_generated template must name a fail-before regression check or exact deletion surface.\n' \
      "$producer_relative_path" "$producer_line" >&2
    fail_count=$((fail_count + 1))
  fi
  if [ "$has_matching_green" -ne 1 ]; then
    printf 'FAIL: %s:%s review_generated template must pair its RED with matching pass-after or deletion GREEN.\n' \
      "$producer_relative_path" "$producer_line" >&2
    fail_count=$((fail_count + 1))
  fi
  if [ "$has_destination_finding_closure_citation" -ne 1 ]; then
    printf 'FAIL: %s:%s review_generated template must cite the Finding-Closure Ratchet with a destination-safe path and contain no source-relative `../do-work/actions/` citation.\n' \
      "$producer_relative_path" "$producer_line" >&2
    fail_count=$((fail_count + 1))
  fi
done < <(
  while IFS= read -r shipped_action_file; do
    awk '
      /^```/ {
        if (!inside_fence) {
          inside_fence = 1
          fence_start = NR
          fence_block = $0 ORS
          next
        }

        fence_block = fence_block $0 ORS
        if (fence_block ~ /(^|\n)review_generated: true(\n|$)/) {
          has_proof_shape = \
            fence_block ~ /## Red-Green Proof/ \
            && fence_block ~ /\*\*RED prompt\/case:\*\*/ \
            && fence_block ~ /\*\*Why RED now:\*\*/ \
            && fence_block ~ /\*\*GREEN when:\*\*/ \
            && fence_block ~ /\*\*Validation:\*\*/
          has_named_red = fence_block ~ /RED prompt\/case:.*Named regression test\/check that fails before the fix.*exact finding surface to delete/
          has_matching_green = fence_block ~ /GREEN when:.*same named test\/check passes after the fix.*exact named finding surface is absent/
          has_destination_finding_closure_citation = \
            fence_block ~ /`(actions\/work-reference\.md|\.claude\/skills\/do-work\/actions\/work-reference\.md)`.*Finding-Closure Ratchet/ \
            && fence_block !~ /(\.\.\/)+do-work\/actions\//
          has_pau_block = \
            fence_block ~ /## AI Execution State \(P-A-U Loop\)/ \
            && fence_block ~ /\*\*\[PLAN\]:\*\*/ \
            && fence_block ~ /\*\*\[APPLY\]:\*\*/ \
            && fence_block ~ /\*\*\[UNIFY\]:\*\*/
          printf "%s\t%d\t%d\t%d\t%d\t%d\t%d\n", FILENAME, fence_start, has_proof_shape, has_named_red, has_matching_green, has_destination_finding_closure_citation, has_pau_block
        }
        inside_fence = 0
        fence_block = ""
        next
      }

      inside_fence { fence_block = fence_block $0 ORS }
    ' "$shipped_action_file"
  done < <(find "$repo_root/skills" -type f -path '*/actions/*.md' -name '*.md' -print | sort)
)

if [ "$review_generated_field_count" -eq 0 ]; then
  printf 'FAIL: shipped actions must retain at least one exact review_generated: true producer for this closure seam probe.\n' >&2
  fail_count=$((fail_count + 1))
elif [ "$review_generated_template_count" -ne "$review_generated_field_count" ]; then
  printf 'FAIL: found %s exact review_generated fields but only %s fenced producer templates; every exact field must live in a checked template.\n' \
    "$review_generated_field_count" "$review_generated_template_count" >&2
  fail_count=$((fail_count + 1))
fi

# Common Rationalizations regrowth ratchet (REQ-027). The four "earned" template
# sections (Rules / Common Rationalizations / Red Flags / Verification Checklist)
# drifted from "included when they'd help" to "included because the template listed
# them" — 20 of 42 action files carried all four (24 carry a Common Rationalizations
# table at all), most filled with generic engineering
# advice a capable model already follows. This check catches regrowth in Common
# Rationalizations specifically: a table whose rows carry no do-work-specific noun is
# exactly that generic filler (see _dev/primes/prime-action-files.md for the full
# omission test). Scoped to files added after REQ-027 — the baseline below grandfathers
# the existing tree (as of REQ-027) so this check lands green without a mass rewrite in
# the same commit; REQ-025/028/029/030/031 clean up the backlog under the new rule.
common_rationalizations_baseline_action_files=(
  abandon.md ai-report.md bkb-reference.md bkb.md board.md capture.md clarify.md
  cleanup.md code-review.md commit.md deep-explore-reference.md deep-explore.md
  dream.md forensics.md help.md inspect.md install.md interview-reference.md
  interview.md kb-lessons-handoff.md note.md present-work.md
  prime.md prompts.md quick-wins.md
  review-work.md roadmap.md sample-archived-req.md scan-ideas.md slop-check.md
  stray-check.md tidy-repo.md tutorial.md ui-review.md validate-feedback.md
  verify-requests.md version.md work-reference.md work.md
)

is_grandfathered_rationalizations_file() {
  local candidate_file_name="$1"
  local baseline_entry
  for baseline_entry in "${common_rationalizations_baseline_action_files[@]}"; do
    if [ "$baseline_entry" = "$candidate_file_name" ]; then
      return 0
    fi
  done
  return 1
}

# Illustrative, not exhaustive (Closed Enumerations Go Stale) — grep for do-work
# machinery vocabulary, not general software-engineering nouns. Split in two so "REQ"
# stays case-sensitive (case-insensitive it collides with "required"/"requires", which
# show up in ordinary generic advice and would silently defeat the check). The memory
# subsystem (working memory / daily log / usage ledger / bootstrap / Stop hook) is
# first-class do-work machinery too, so its nouns are recognized alongside the rest.
rationalization_noun_pattern_case_sensitive='\bREQ-?[0-9]*\b|\bUR-[0-9]+\b'
rationalization_noun_pattern_case_insensitive='queue|frontmatter|pipeline|\barchive|do-work|\bdomain\b|\bblocked\b|\bkb/|\bprime\b|\bclarify|working/|crew-member|\bschema\b|status:|working memory|daily log|\bledger\b|\bbootstrap\b|stop hook'

for action_file_path in "$repo_root"/skills/do-work*/actions/*.md; do
  action_file_name="$(basename "$action_file_path")"
  if is_grandfathered_rationalizations_file "$action_file_name"; then
    continue
  fi
  if ! grep -q '^## Common Rationalizations' "$action_file_path"; then
    continue
  fi
  rationalization_block="$(awk '/^## Common Rationalizations/{flag=1;next}/^## /{flag=0}flag' "$action_file_path" || true)"
  rationalization_rows="$(grep '^|' <<<"$rationalization_block" | grep -viE "if you're thinking" | grep -vE '^\|[-: |]+\|$' || true)"
  if [ -z "$rationalization_rows" ]; then
    continue
  fi
  if ! grep -qE "$rationalization_noun_pattern_case_sensitive" <<<"$rationalization_rows" \
     && ! grep -qiE "$rationalization_noun_pattern_case_insensitive" <<<"$rationalization_rows"; then
    printf 'FAIL: %s Common Rationalizations table has no do-work-specific noun (REQ, UR, queue, frontmatter, pipeline, archive, domain, blocked, kb/, prime, clarify, working/, crew-member, schema, status, working memory, daily log, ledger, bootstrap, stop hook — illustrative list) in any row — every row reads as generic engineering advice a capable model already follows. Add rows naming a specific do-work failure mode, or omit the section entirely (see _dev/primes/prime-action-files.md for the omission test).\n' \
      "$action_file_name" >&2
    fail_count=$((fail_count + 1))
  fi
done

# Capture redaction must run on the FULL extracted text before any byte-budget
# cut: every credential pattern needs a complete token shape, so redacting after
# truncation lets a severed token (`ghp_1234567`) persist as an unmatched
# fragment (validate-feedback 2026-07-28, reproduced).
assert_contains \
  "actions/memory-reference.md" \
  'redaction runs BEFORE truncation' \
  'actions/memory-reference.md Stop-Capture spec must order redaction before truncation.'

if command -v jq &>/dev/null; then
  redaction_order_workdir="$(mktemp -d)"
  mkdir -p "$redaction_order_workdir/memory/logs"
  # Pad so the GitHub token straddles the side-truncation cut: with a short
  # assistant side the user side is cut at (1500-15) - 11 - 12 bytes, which
  # lands mid-token for a token starting at byte 1451.
  redaction_order_pad="$(printf 'A%.0s' $(seq 1 1450))"
  printf '{"type":"user","message":{"content":"%sghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ012345 tail"}}\n{"type":"assistant","message":{"content":[{"type":"text","text":"shortreply"}]}}\n' \
    "$redaction_order_pad" > "$redaction_order_workdir/transcript.jsonl"
  printf '{"transcript_path":"%s/transcript.jsonl"}' "$redaction_order_workdir" \
    | CLAUDE_PROJECT_DIR="$redaction_order_workdir" bash "$knowledge_root/hooks/memory-stop-capture.sh" >/dev/null 2>&1 || true
  if ! ls "$redaction_order_workdir/memory/logs/"*.md >/dev/null 2>&1; then
    printf 'FAIL: hooks/memory-stop-capture.sh redaction-order probe wrote no capture — the probe transcript should be captured, not dropped.\n' >&2
    fail_count=$((fail_count + 1))
  elif grep -rqE 'ghp_[A-Za-z0-9_]+' "$redaction_order_workdir/memory/logs/" 2>/dev/null; then
    printf 'FAIL: hooks/memory-stop-capture.sh persisted a truncation-severed credential fragment — redaction must run on the full extracted text before any byte-budget cut.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$redaction_order_workdir"
fi

# Worktree dispatch mode (REQ-033). The mode is only safe because three pieces
# hold together: a name that correlates a leftover with its REQ, a merged-ness
# assertion that is never forced, and a verification of the *merged* state before
# archive. Each can be quietly dropped without breaking anything visible, so pin
# them individually — plus cleanup's ownership of the unmerged leftovers, without
# which a crashed parallel run leaves branches nothing in the skill ever removes.
assert_contains \
  "actions/work-reference.md" \
  'worktree-agent-REQ-' \
  'actions/work-reference.md must keep the worktree/branch naming convention that correlates a leftover with its REQ id.'

assert_contains \
  "actions/work-reference.md" \
  'git branch -d' \
  'actions/work-reference.md must keep git branch -d (never -D) as the free merged-ness assertion on worktree cleanup.'

assert_contains \
  "actions/work.md" \
  '[Pp]ost-merge verification' \
  'actions/work.md Step 8 must gate archival on re-running the REQ acceptance checks against the merged tree.'

assert_contains \
  "actions/cleanup.md" \
  'worktree-agent-' \
  'actions/cleanup.md must keep the consent-gated orphaned-worktree pass that owns unmerged builder branches.'

# Board overlap annotation stays display-only (REQ-052, follow-up to REQ-034). The
# invariant is structural: `annotateWriteSetOverlap` runs AFTER `bucketColumns` in
# buildBoard, so no column-placement code can read the annotation. That ordering is one
# movable line — a future edit that hoists the call above bucketing (to "reuse" the
# overlap set while placing cards) turns a display badge into scheduling logic, and the
# co-dispatch decision belongs to actions/work.md Step 1's gate, not the viewer. Assert
# the real call-site order inside buildBoard rather than a comment: a comment survives
# the move. Paired with the instruction-side claim in actions/board.md, whose rewording
# is the other way the invariant quietly stops being a promise.
build_board_block="$(sed -n '/^func buildBoard(/,/^}/p' "$board_root/tools/queue-kanban/model.go")"
bucket_columns_call_line="$(grep -nF 'bucketColumns(board.AllRequests' <<<"$build_board_block" | head -1 | cut -d: -f1 || true)"
overlap_annotation_call_line="$(grep -nF 'annotateWriteSetOverlap(board.AllRequests)' <<<"$build_board_block" | head -1 | cut -d: -f1 || true)"

if [ -z "$bucket_columns_call_line" ] || [ -z "$overlap_annotation_call_line" ]; then
  printf 'FAIL: tools/queue-kanban/model.go buildBoard must call both bucketColumns(board.AllRequests ...) and annotateWriteSetOverlap(board.AllRequests) — one call site was renamed or removed, so the display-only ordering of the write-set overlap annotation can no longer be verified. Fix: restore the call sites (annotation last) or update this anchor in the same commit.\n' >&2
  fail_count=$((fail_count + 1))
elif [ "$overlap_annotation_call_line" -lt "$bucket_columns_call_line" ]; then
  printf 'FAIL: tools/queue-kanban/model.go calls annotateWriteSetOverlap BEFORE bucketColumns in buildBoard — the write-set overlap annotation is display only and must stay structurally unable to affect column placement. Fix: move the annotateWriteSetOverlap call back below bucketColumns; if the board really must schedule on write sets, that decision belongs to actions/work.md Step 1 dispatch gate.\n' >&2
  fail_count=$((fail_count + 1))
fi

board_rules_block="$(sed -n '/^## Rules/,/^## Common Rationalizations/p' "$board_root/actions/board.md")"

assert_block_contains \
  "$board_rules_block" \
  '`annotateWriteSetOverlap` in `tools/queue-kanban/model\.go` runs \*after\* column bucketing' \
  'actions/board.md Rules must keep the claim that annotateWriteSetOverlap runs after column bucketing — the ordering ratchet above pins the code, this pins the promise a parser-editing agent reads. Fix: restate it in the parser lock-step rule.'

assert_block_contains \
  "$board_rules_block" \
  'never column logic' \
  'actions/board.md Rules must keep the overlap annotation display-only claim (drives the overlaps badge and drawer row, never column logic, never blocking) — without it nothing tells a parser-editing agent that which REQs may co-dispatch stays with actions/work.md Step 1 gate.'

# Mid-run answer durability. A question asked interactively can satisfy every wording rule
# and still produce an answer that dies with the session: a consumer's long-running
# orchestrator asked two user-owned decisions, got detailed answers, wrote nothing to disk,
# and the next builder re-decided one of them in a fresh context off the stored
# `Recommended:` rationale the user had just rejected (validate-feedback 2026-08-01). Three
# halves hold the fix together — the principle, the pipeline branch that obeys it, and the
# one named format both cite — and each reads as removable prose on its own.
assert_contains \
  "crew-members/clear-questions.md" \
  'outlive the transcript' \
  'crew-members/clear-questions.md must keep the principle that an interactively obtained answer is written into the durable record before it is acted on — wording rules alone let a compliant question produce an answer that dies with the session.'

work_open_questions_block="$(sed -n '/^### Step 3\.5: Open Questions/,/^### Step 3\.7/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_open_questions_block" \
  'escalated to the user mid-run and answered' \
  'actions/work.md Step 3.5 must keep the third branch for a decision escalated mid-run and answered (orchestrator writes it in before dispatch) — with only builder-decides and defer-to-clarify, a mid-run answer has nowhere to land.'

assert_block_contains \
  "$work_open_questions_block" \
  '\*\*not\*\* `- \[~\]`' \
  'actions/work.md Step 3.5 must keep a mid-run user answer recorded as - [x] and never - [~] — filing it as a builder decision sends a settled question back through clarify and leaves the rejected Recommended: rationale standing as its reason.'

assert_contains \
  "actions/clarify.md" \
  'Canonical answered-question format' \
  'actions/clarify.md Step 4 must keep the named entry point declaring its - [x] form canonical for any caller that obtains a user answer — cited by name (not step number) from clear-questions.md Principle 8 and work.md Step 3.5.'

# Decision revalidation (0.187.0). A reversed decision used to update its own ADR or
# follow-up REQ while pending siblings stayed runnable against the dead assumption. The
# v1 repair is intentionally report-only: one evidence-backed queue scan, automatically
# composed by clarify only for real overrides, with a cost gate before large scans.
verify_input_block="$(sed -n '/^## Input/,/^## Capture QA Workflow/p' "$core_root/actions/verify-requests.md")"
verify_revalidation_block="$(sed -n '/^## Decision Revalidation Workflow/,/^## What NOT To Do/p' "$core_root/actions/verify-requests.md")"
clarify_revalidation_block="$(sed -n '/^### Step 5\.25: Revalidate queued work after reversals/,/^### Step 5\.5:/p' "$core_root/actions/clarify.md")"

# REQ-288 K2/K3/K4 — the three shipped contradictions in clarify's Step 4. Each check
# names the defect it pins, and each pins a rule that a plausible edit would undo:
# K3 and K4 both destroy work when they regress (an archived REQ holding an open
# question; an approved follow-up archived completed without ever being built).
clarify_step4_block="$(sed -n '/^### Step 4: Collect answers/,/^### Step 5\.25:/p' "$core_root/actions/clarify.md")"
clarify_checklist_block="$(sed -n '/^## Verification Checklist/,$p' "$core_root/actions/clarify.md")"

# --- K2: the durable record is the answer line PLUS a dated reasoning note ---
assert_block_contains \
  "$clarify_step4_block" \
  'durable record is the .- \[x\] \[question\] → \[answer\]. form below \*\*together with a dated note' \
  "K2: clarify's canonical answered-question block must define the durable record as the answer line PLUS the dated reasoning note (clear-questions.md Principle 8), not the answer line alone."

assert_block_contains \
  "$clarify_step4_block" \
  'put out of scope' \
  'K2: the dated note must be required to carry anything the answer put out of scope, or the next reader re-derives the decision from Recommended:.'

# K2's date constraint: cite the rule, never copy a clock command into an action file.
assert_block_contains \
  "$clarify_step4_block" \
  'Timestamp rule.s date-only paragraph .`actions/work-reference\.md`. — cite it, never spell a clock command' \
  "K2: clarify's dated note must cite the Timestamp rule's date-only paragraph rather than spelling a command; ungoverned prose dates get fabricated (UR-055)."

assert_block_not_contains \
  "$(cat "$core_root/actions/clarify.md")" \
  'date -u \+%F|date -u \+%Y|Get-Date' \
  'K2: clarify.md must never spell a clock command — work-reference.md is the only place in actions/ that does.'

# The date-only paragraph itself must key on the condition, not on a consumer list.
assert_block_contains \
  "$(cat "$core_root/actions/work-reference.md")" \
  'any UTC calendar date written into a durable record' \
  "K2: the Timestamp rule's date-only paragraph must key on the condition; an enumerated consumer list goes stale and left REQ prose notes ungoverned (CLAUDE.md 'State conditions, not lists')."

# --- K3: per-question verbs set no REQ-level state; Step 5 aggregates once ---
assert_block_contains \
  "$clarify_step4_block" \
  'never sets the REQ.s status and never archives' \
  'K3: a per-question verb must not set whole-file state — discarding one question while skipping another made two branches unfollowable at once.'

assert_block_contains \
  "$clarify_checklist_block" \
  'stayed .pending-answers. in .do-work/queue/' \
  'K3: a REQ holding even one remaining unanswered question must never be archived; the checklist is what makes that auditable.'

assert_block_contains \
  "$(sed -n '/^### Step 5: Resolve each REQ/,/^### Step 5\.25:/p' "$core_root/actions/clarify.md")" \
  'Any remaining .- \[ \]. wins' \
  "K3: Step 5 must compute status once from every question's outcome, with any open question holding the REQ in pending-answers."

# --- K4: the completed fast path routes on the marker, never on question prose ---
assert_block_contains \
  "$(sed -n '/^### Step 5: Resolve each REQ/,/^### Step 5\.25:/p' "$core_root/actions/clarify.md")" \
  'builder_decided: true. follow-up whose questions were all confirmed' \
  'K4: the completed fast path must key on the builder_decided marker; keyed on question prose, a reworded consent question archives an approved follow-up without ever building it.'

assert_block_contains \
  "$(sed -n '/^### Step 5: Resolve each REQ/,/^### Step 5\.25:/p' "$core_root/actions/clarify.md")" \
  'Never infer this branch from question prose' \
  'K4: the marker must be stated as the entire test, so a future edit cannot reintroduce prose matching as a fallback.'

assert_block_not_contains \
  "$(cat "$core_root/actions/clarify.md")" \
  'whose question is "Should I process this as a new task\?"' \
  'K4: clarify must not route on the literal discriminator phrase any more — that is the defense review-work.md predicted would fail.'

# REQ-342 — the Canonical answered-question format writes text obtained from outside the
# file straight into a REQ body, and until this rule landed nothing said the text was
# neutralized first. Three consequences, all reproduced on fixtures: a line-leading
# `- [ ]` inside an answer becomes a real open question that Step 5's "any remaining
# unchecked item wins" then pins the REQ on forever; one unbalanced fence renders every
# section below it inside <pre><code> on the board while the prose greps still see those
# sections; and text landing above the file's first line drops the whole frontmatter, so
# status, title and user_request read empty and the REQ leaves its UR silently.
#
# The rule is prose, so a keyword match would pass a gutted version of it. This extracts
# the named entry point's own blockquote, grades the two branch lines in isolation, and
# pairs positive meaning with scoped contradiction checks. Replacement trials keep the
# earlier deletion coverage; insertion trials preserve that vocabulary while adding the
# captured narrowings, including the closed-list failure named by
# prime-shell-commands.md § Closed Enumerations Go Stale.
if ! python3 - "$core_root/actions/clarify.md" <<'PY'
import pathlib
import re
import sys

clarify_text = pathlib.Path(sys.argv[1]).read_text()


def extract_named_entry_point_block(source):
    """The contiguous blockquote that opens Step 4 — the named entry point itself."""
    step_match = re.search(
        r"^### Step 4: Collect answers\n(?P<body>.*?)^### Step 5: ",
        source,
        flags=re.DOTALL | re.MULTILINE,
    )
    if step_match is None:
        return None
    quoted_lines = []
    started = False
    for line in step_match.group("body").splitlines():
        if line.startswith(">"):
            started = True
            quoted_lines.append(line)
        elif started:
            break
    return "\n".join(quoted_lines) if quoted_lines else None


def neutralization_defects(source):
    block = extract_named_entry_point_block(source)
    if block is None:
        return {"missing named entry point blockquote"}

    normalized = " ".join(block.lower().split())
    branch_lines = [
        " ".join(line.lower().split())
        for line in block.splitlines()
        if re.match(r"^>\s*-\s*", line)
    ]
    single_line_branch = next(
        (line for line in branch_lines if "inline after the `→`" in line), ""
    )
    contained_passage_branch = next(
        (line for line in branch_lines if "blockquote" in line), ""
    )
    delimiter_judgment_clause = next(
        (
            " ".join(line.lower().split())
            for line in block.splitlines()
            if "judge" in line.lower() and "condition" in line.lower()
        ),
        "",
    )
    predicates = {
        # The rule exists at all, and at the named entry point rather than in
        # loose prose a caller citing the name would never reach.
        "neutralization stated at the named entry point":
            r"neutraliz\w+ before it is written",
        # The contract protects every do-work Markdown record body that accepts
        # outside text, not only the answer line that first exposed the class.
        "do-work record-body reach":
            r"any do-work markdown record body",
        # Markdown delimiter containment cannot repair a file that has already
        # become binary/non-text. The preflight is deliberately a condition,
        # not another list of current writers.
        "text-byte preflight":
            r"c0.*?del.*?lf.*?tab.*?refuse.*?report",
        "byte identity after containment":
            r"byte-identical apart from containment bytes",
        # The whole point of the REQ: a condition, never a character list.
        "delimiter condition, not a character list":
            r"could this line be read as one of this file.s own delimiters",
        "shapes marked illustrative":
            r"illustrative, not a checklist",
        # The four shapes the fixture proved. Present as examples of the
        # condition — their absence means the condition lost its worked cases.
        "checkbox shape named": r"line-leading `- \[ \]` or `- \[x\]`",
        "bare thematic-break shape named": r"bare `---` line",
        "heading shape named": r"`## ` line",
        "unbalanced fence shape named": r"unbalanced code fence",
        # Neutralizing must not cost the answer trail the record exists for.
        "answer preserved intact": r"nothing is edited or dropped",
        # Branch two, half one: the quote prefix is what a line-based scan sees,
        # and a fence alone does not stop one — the sibling incident in this
        # file's own history was a regex keying on the first line-leading `- [ ]`.
        "line-based scan defended":
            r"prefix takes every line start away from a line-based scan",
        # Branch two, half two: containment the pasted text provably cannot close.
        "container cannot be closed from inside":
            r"code fence longer than the longest backtick run anywhere in the text",
        # The file's first line is a delimiter too, and no fence can guard it.
        "opening frontmatter fence never written above":
            r"never the file.s first line",
    }
    defects = {
        defect_name
        for defect_name, predicate in predicates.items()
        if re.search(predicate, normalized) is None
    }
    scoped_predicates = {
        # Isolate the two list branches before grading them, so vocabulary in
        # the surrounding blockquote cannot lend a weakened branch its force.
        "single-line branch binds both triggers": (
            single_line_branch,
            r"(?:the|an) answer (?:is one|has a single) line\b.*?cannot be a delimiter where it lands.*?inline after the `→`",
        ),
        # Accept equivalent universal grammar while keeping all three triggers
        # and the action in one branch: own passage, line break, delimiter line.
        "contained branch covers every passage and risky answer": (
            contained_passage_branch,
            r"(?:anything\b.*?|every\b.*?)body passage.*?(?:or|plus).*?(?:any|each) answer.*?line break.*?(?:or|plus).*?delimiter-shaped line.*?(?:put.*?inside|place.*?in) a blockquote",
        ),
        "delimiter judgment covers every line": (
            delimiter_judgment_clause,
            r"judge (?:every|all) line.*?(?:one|a single) condition",
        ),
    }
    defects.update(
        defect_name
        for defect_name, (clause, predicate) in scoped_predicates.items()
        if re.search(predicate, clause) is None
    )

    # Presence predicates above prove the broad rule still exists. These scoped
    # contradictions prove it has not been narrowed by an inserted sentence
    # that leaves every broad phrase untouched. Keep each pattern tied to the
    # subject it restricts: this contract legitimately says ONLY containment
    # bytes are added, so a file-wide ban on `only` would reject the live rule.
    narrowing_patterns = {
        "delimiter judgment narrowed to the first line":
            r"judge (?:only )?(?:its|the) first line",
        "delimiter condition narrowed by a closed list":
            r"only (?:the )?(?:four|named|listed).*?shapes.*?count as delimiters",
        "record-body reach narrowed to answers":
            r"(?:sub-?contract|containment|rule).*?appl(?:y|ies) only to (?:req )?answer",
        "record-body reach narrowed to pasted code":
            r"(?:containment )?rule appl(?:y|ies) only to pasted code",
        "record-body reach narrowed to clarify":
            r"(?:containment )?rule appl(?:y|ies) only in clarify.*?orchestrator answers? (?:are )?exempt",
        "illustrative shapes made exhaustive":
            r"(?:treat|use|regard).*?shapes as (?:the )?(?:complete|exhaustive) (?:check)?list",
        "blockquote prefix made optional":
            r"`> ` prefix is optional.*?fence is (?:enough|sufficient)",
        "inline branch narrowed by length":
            r"(?:answer|one-line answers?).*?(?:short|length).*?(?:long|longer) one-line answer",
        "contained branch narrowed by a line threshold":
            r"(?:passage|contained) branch appl(?:y|ies) only to .*?(?:ten|10) lines? (?:or more|and above)",
        "contained branch made advisory":
            r"(?:use|apply).*?(?:passage|contained) branch only where practical",
        "neutralization narrowed to risky-looking text":
            r"(?:apply|use) neutralization only when .*?(?:look|seem)s? risky",
        "dynamic fence narrowed to a fixed fence":
            r"(?:triple-backtick|three-backtick|fixed) fence is (?:enough|sufficient)",
        "first-line prohibition narrowed to answers":
            r"(?:first-line prohibition|prohibition).*?appl(?:y|ies) only to (?:answered questions|answers)",
    }
    defects.update(
        defect_name
        for defect_name, pattern in narrowing_patterns.items()
        if re.search(pattern, normalized) is not None
    )
    return defects


live_defects = neutralization_defects(clarify_text)
if live_defects:
    raise SystemExit(
        "actions/clarify.md Step 4 Canonical answered-question format is missing its "
        "neutralization contract: " + ", ".join(sorted(live_defects))
    )
live_block = extract_named_entry_point_block(clarify_text)

mutations = (
    ("rule removed entirely",
     "neutralized before it is written",
     "recorded as the user typed it",
     "neutralization stated at the named entry point"),
    ("condition replaced by a closed list",
     "could this line be read as one of this file's own delimiters",
     "does this line start with `- [ ]`, `- [x]`, `---`, `## ` or a backtick fence",
     "delimiter condition, not a character list"),
    ("reach narrowed back to REQ answers",
     "any do-work Markdown record body",
     "a REQ answer line",
     "do-work record-body reach"),
    ("control-byte preflight removed",
     "C0 control or DEL except LF and TAB",
     "unusual character",
     "text-byte preflight"),
    ("byte identity weakened to normalization",
     "byte-identical apart from containment bytes",
     "readable after normalization",
     "byte identity after containment"),
    ("examples hardened into a checklist",
     "illustrative, not a checklist",
     "the complete set to check",
     "shapes marked illustrative"),
    ("checkbox shape dropped",
     "a line-leading `- [ ]` or `- [x]`",
     "a stray marker",
     "checkbox shape named"),
    ("thematic-break shape dropped",
     "a bare `---` line",
     "a separator",
     "bare thematic-break shape named"),
    ("heading shape dropped",
     "a `## ` line",
     "a title line",
     "heading shape named"),
    ("fence shape dropped",
     "an unbalanced code fence",
     "a stray backtick",
     "unbalanced fence shape named"),
    ("characters silently stripped",
     "nothing is edited or dropped",
     "the delimiter characters are removed",
     "answer preserved intact"),
    ("quote prefix dropped, leaving line-based scans exposed",
     "The `> ` prefix takes every line start away from a line-based scan",
     "The fence is enough on its own",
     "line-based scan defended"),
    ("single-line branch loosened to a judgement call",
     "The answer is one line",
     "the answer is short",
     "single-line branch binds both triggers"),
    ("container weakened to a fixed fence",
     "a code fence longer than the longest backtick run anywhere in the text",
     "a triple-backtick code fence",
     "container cannot be closed from inside"),
    ("frontmatter placement rule dropped",
     "never the file's first line",
     "wherever the answer reads best",
     "opening frontmatter fence never written above"),
)

for mutation_name, old, new, expected_defect in mutations:
    mutated_block = live_block.replace(old, new, 1)
    if mutated_block == live_block:
        raise SystemExit(
            f"answer-neutralization mutation {mutation_name!r} changed nothing — "
            "the anchor it rewrites is no longer in the named entry point block"
        )
    mutated_text = clarify_text.replace(live_block, mutated_block, 1)
    mutation_defects = neutralization_defects(mutated_text)
    if expected_defect not in mutation_defects:
        raise SystemExit(
            f"answer-neutralization mutation {mutation_name!r} escaped "
            f"{expected_defect!r}; found {sorted(mutation_defects)!r}"
        )

# Semantic checks must accept ordinary prose maintenance, not freeze the live
# sentences byte-for-byte. This trial keeps both branch conditions and actions
# intact while changing their grammar and verb choice.
legitimate_rewording_block = live_block
for old, new in (
    ("The answer is one line", "An answer has a single line"),
    ("put the outside text inside a blockquote", "place the outside text in a blockquote"),
):
    reworded_block = legitimate_rewording_block.replace(old, new, 1)
    if reworded_block == legitimate_rewording_block:
        raise SystemExit(
            f"answer-neutralization legitimate-rewording anchor {old!r} disappeared"
        )
    legitimate_rewording_block = reworded_block
legitimate_rewording_defects = neutralization_defects(
    clarify_text.replace(live_block, legitimate_rewording_block, 1)
)
if legitimate_rewording_defects:
    raise SystemExit(
        "answer-neutralization legitimate rewording was rejected: "
        + repr(sorted(legitimate_rewording_defects))
    )

universal_trigger_rewording = live_block.replace(
    "Anything written as its own body passage, or any answer with a line break or delimiter-shaped line",
    "Every outside-text body passage, plus each answer containing a line break or delimiter-shaped line",
    1,
)
if universal_trigger_rewording == live_block:
    raise SystemExit("answer-neutralization universal-trigger rewording changed nothing")
universal_trigger_rewording_defects = neutralization_defects(
    clarify_text.replace(live_block, universal_trigger_rewording, 1)
)
if universal_trigger_rewording_defects:
    raise SystemExit(
        "answer-neutralization equivalent universal-trigger rewording was rejected: "
        + repr(sorted(universal_trigger_rewording_defects))
    )


def insert_after(block, anchor, addition, mutation_name):
    if anchor not in block:
        raise SystemExit(
            f"answer-neutralization insertion {mutation_name!r} changed nothing — "
            "the anchor it follows is no longer in the named entry point block"
        )
    return block.replace(anchor, anchor + addition, 1)


# Replacement mutations can prove that required words are load-bearing and still
# miss the more dangerous edit: preserving every broad sentence while adding a
# narrower instruction beside it. Each trial below keeps the current positive
# vocabulary intact and inserts one plausible contradiction.
insertion_mutations = (
    (
        "delimiter condition narrowed to the four examples",
        "The shapes already proved are illustrative, not a checklist",
        ". For this rule, only the four shapes named below count as delimiters",
        "delimiter condition narrowed by a closed list",
    ),
    (
        "record-body reach narrowed to answers",
        "any do-work Markdown record body",
        "; this sub-contract applies only to REQ answer lines",
        "record-body reach narrowed to answers",
    ),
    (
        "illustrative examples made exhaustive",
        "The shapes already proved are illustrative, not a checklist",
        "; treat those four shapes as the complete checklist",
        "illustrative shapes made exhaustive",
    ),
    (
        "blockquote prefix made optional",
        "The `> ` prefix takes every line start away from a line-based scan",
        ". For fenced text the `> ` prefix is optional because the fence is sufficient",
        "blockquote prefix made optional",
    ),
    (
        "inline branch narrowed to short answers",
        "The answer is one line",
        " and short; long one-line answers use the passage branch",
        "inline branch narrowed by length",
    ),
    (
        "dynamic fence narrowed to triple backticks",
        "cannot be closed from inside",
        ". A triple-backtick fence is sufficient regardless of the pasted text",
        "dynamic fence narrowed to a fixed fence",
    ),
    (
        "first-line prohibition narrowed to answers",
        "It is never the file's first line",
        "; this prohibition applies only to answered questions",
        "first-line prohibition narrowed to answers",
    ),
)

for mutation_name, anchor, addition, expected_defect in insertion_mutations:
    mutated_block = insert_after(live_block, anchor, addition, mutation_name)
    mutated_text = clarify_text.replace(live_block, mutated_block, 1)
    mutation_defects = neutralization_defects(mutated_text)
    if expected_defect not in mutation_defects:
        raise SystemExit(
            f"answer-neutralization insertion {mutation_name!r} escaped "
            f"{expected_defect!r}; found {sorted(mutation_defects)!r}"
        )

# Replay the independently reproduced RED set that motivated REQ-361. These
# qualifiers use deliberately different grammar from the useful insertion
# trials above; the checker must pin the property, not memorize one mutation.
captured_narrowing_mutations = (
    (
        "delimiter judgment narrowed to the first line",
        "Then judge every line against one condition:",
        "Then judge its first line against one condition:",
        "delimiter judgment narrowed to the first line",
    ),
    (
        "record-body reach narrowed to pasted code",
        "Those sites are illustrative, not a writer checklist",
        "Those sites are illustrative, not a writer checklist; this containment rule applies only to pasted code snippets",
        "record-body reach narrowed to pasted code",
    ),
    (
        "reach narrowed to clarify with orchestrator exemption",
        "Those sites are illustrative, not a writer checklist",
        "Those sites are illustrative, not a writer checklist; this rule applies only in clarify and mid-run orchestrator answers are exempt",
        "record-body reach narrowed to clarify",
    ),
    (
        "passage branch narrowed to ten lines",
        "The destination decides which containment branch applies:",
        "The destination decides which containment branch applies; the passage branch applies only to passages of ten lines or more:",
        "contained branch narrowed by a line threshold",
    ),
    (
        "passage branch made conditional on practicality",
        "The destination decides which containment branch applies:",
        "The destination decides which containment branch applies; use the passage branch only where practical:",
        "contained branch made advisory",
    ),
    (
        "whole rule narrowed to risky-looking answers",
        "The destination decides which containment branch applies:",
        "The destination decides which containment branch applies; apply neutralization only when the answer looks risky:",
        "neutralization narrowed to risky-looking text",
    ),
)

for mutation_name, old, new, expected_defect in captured_narrowing_mutations:
    mutated_block = live_block.replace(old, new, 1)
    if mutated_block == live_block:
        raise SystemExit(
            f"answer-neutralization captured RED {mutation_name!r} changed nothing"
        )
    mutation_defects = neutralization_defects(
        clarify_text.replace(live_block, mutated_block, 1)
    )
    if expected_defect not in mutation_defects:
        raise SystemExit(
            f"answer-neutralization captured RED {mutation_name!r} escaped "
            f"{expected_defect!r}; found {sorted(mutation_defects)!r}"
        )
PY
then
  printf 'FAIL: the Canonical answered-question format (actions/clarify.md Step 4) must state how text obtained from outside the file is neutralized before it is written into a REQ body, keyed on the delimiter condition rather than a character list (REQ-342).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Stated once, inherited by name. Every caller cites the format instead of restating it,
# so the condition sentence must exist in exactly one shipped file — a second copy is the
# drift that leaves one caller neutralizing and another not.
neutralization_condition_count="$({ grep -rhoF "could this line be read as one of this file's own delimiters" "$repo_root/skills" || true; } | wc -l | tr -d ' ')"
if [ "$neutralization_condition_count" != "1" ]; then
  printf 'FAIL: the answer-neutralization condition must be stated exactly once across skills/ — callers cite the Canonical answered-question format by name rather than restating it (found %s).\n' \
    "$neutralization_condition_count" >&2
  fail_count=$((fail_count + 1))
fi

# REQ-360 — every action that writes outside text into a do-work Markdown body
# inherits the one canonical containment contract by name. Isolate each writer's
# own step: nearby citations in another step must not lend it the rule. The same
# lock-in closes every hand-authored frontmatter example that still taught a form
# the Frontmatter Quoting contract forbids.
if ! python3 - \
  "$core_root/actions/verify-requests.md" \
  "$core_root/actions/capture.md" \
  "$core_root/actions/stakeholder-answers.md" \
  "$core_root/actions/abandon.md" \
  "$core_root/actions/work-reference.md" \
  "$core_root/actions/capture-reference.md" \
  "$core_root/actions/review-work.md" \
  "$core_root/actions/sample-archived-req.md" \
  "$core_root/docs/capture-guide.md" \
  "$core_root/docs/work-guide.md" <<'PY'
import pathlib
import re
import sys

(
    verify_path,
    capture_path,
    stakeholder_path,
    abandon_path,
    work_reference_path,
    capture_reference_path,
    review_path,
    sample_path,
    capture_guide_path,
    work_guide_path,
) = map(pathlib.Path, sys.argv[1:])


def section(source, start, end):
    match = re.search(start + r"(?P<body>.*?)" + end, source,
                      flags=re.DOTALL | re.MULTILINE)
    if match is None:
        raise SystemExit(f"cannot isolate writer block {start!r} -> {end!r}")
    return match.group("body")


def outside_inheritance_defect(block):
    normalized = " ".join(block.lower().split())
    if "outside-text containment" not in normalized:
        return "missing named Outside-text containment inheritance"
    # Bind the action verb to this repository's two deliberate directive
    # phrases. A generic "per" or "using" elsewhere in a whole Step block must
    # not lend a passive citation authority it does not have.
    active_pattern = (
        r"(?:apply `actions/clarify\.md` step 4.s |"
        r"contained per that format.s )\*\*outside-text containment\*\*"
    )
    if re.search(active_pattern, normalized) is None:
        return "passive Outside-text containment citation"
    return None


def require_active_inheritance(label, block):
    defect = outside_inheritance_defect(block)
    if defect:
        raise SystemExit(f"{label}: {defect}")


verify = verify_path.read_text()
capture = capture_path.read_text()
stakeholder = stakeholder_path.read_text()
abandon = abandon_path.read_text()

outside_writer_blocks = (
    ("verify-requests Step 7 resolved-answer writer",
     section(verify, r"^### Step 7: Offer Fixes\n", r"^## Scoring Guidelines")),
    ("capture queued-addendum writer",
     section(capture, r"^### Step 2: ", r"^### Step 3: ")),
    ("capture UR input writer",
     section(capture, r"^### Step 5: Write Files\n", r"^### Step 6: ")),
    ("stakeholder-answers new-UR writer",
     section(stakeholder, r"^### Step 4: ", r"^### Step 5: ")),
    ("abandon cancellation-reason writer",
     section(abandon, r"^### Step 3: ", r"^### Step 4: ")),
)

for writer_label, writer_block in outside_writer_blocks:
    require_active_inheritance(writer_label, writer_block)

    missing_mutation = writer_block.replace(
        "Outside-text containment", "local containment", 1
    )
    if missing_mutation == writer_block:
        raise SystemExit(f"{writer_label}: missing-name mutation changed nothing")
    if outside_inheritance_defect(missing_mutation) is None:
        raise SystemExit(f"{writer_label}: missing-name mutation escaped")

    if "Apply `actions/clarify.md` Step 4's" in writer_block:
        passive_mutation = writer_block.replace(
            "Apply `actions/clarify.md` Step 4's", "See `actions/clarify.md` Step 4's", 1
        )
    else:
        passive_mutation = writer_block.replace(
            "contained per that format's", "see that format's", 1
        )
    if passive_mutation == writer_block:
        raise SystemExit(f"{writer_label}: passive-citation mutation changed nothing")
    if outside_inheritance_defect(passive_mutation) != "passive Outside-text containment citation":
        raise SystemExit(f"{writer_label}: passive-citation mutation escaped")


def frontmatter_inheritance_defect(block):
    normalized = " ".join(block.lower().split())
    if "frontmatter quoting" not in normalized:
        return "missing named Frontmatter Quoting inheritance"
    if re.search(r"written per \*\*frontmatter quoting\*\*", normalized) is None:
        return "passive Frontmatter Quoting citation"
    return None


stakeholder_step5 = section(stakeholder, r"^### Step 5: ", r"^### Step 6: ")
stakeholder_step5_defect = frontmatter_inheritance_defect(stakeholder_step5)
if stakeholder_step5_defect:
    raise SystemExit(
        "stakeholder-answers Step 5 blocked_by writer: " + stakeholder_step5_defect
    )

mutated_step5 = stakeholder_step5.replace("Frontmatter Quoting", "local quoting", 1)
if mutated_step5 == stakeholder_step5:
    raise SystemExit(
        "stakeholder-answers Step 5 Frontmatter Quoting mutation changed nothing"
    )
if frontmatter_inheritance_defect(mutated_step5) is None:
    raise SystemExit(
        "stakeholder-answers Step 5 lost Frontmatter Quoting but the writer check stayed green"
    )

passive_step5 = stakeholder_step5.replace(
    "written per **Frontmatter Quoting**", "see **Frontmatter Quoting**", 1
)
if passive_step5 == stakeholder_step5:
    raise SystemExit(
        "stakeholder-answers Step 5 passive-citation mutation changed nothing"
    )
if frontmatter_inheritance_defect(passive_step5) != "passive Frontmatter Quoting citation":
    raise SystemExit(
        "stakeholder-answers Step 5 passive Frontmatter Quoting citation escaped the writer check"
    )

def frontmatter_contract_defects(source):
    block = section(
        source,
        r"^\*\*Named contract — Frontmatter Quoting\.\*\*",
        r"^\*\*Do not wrap user text in a double-quoted scalar\.\*\*",
    )
    normalized = " ".join(block.lower().split())
    predicates = {
        "outside-authorship condition governs every field":
            r"condition is the rule.*?whenever a frontmatter value carries text nobody in this pipeline composed",
        "current field names stay illustrative":
            r"fields.*?illustrative, never a list to check against",
        "outside-text preflight inherited": r"outside-text containment.s accepted-text preflight",
        "failed preflight refused and reported": r"refuses and reports.*?instead of normalizing",
        "every physical line indented": r"every physical content line indented.*?blank lines included",
        "zero-terminal-LF strip form": r"`key: \|-` for zero terminal lf bytes",
        "one-terminal-LF clip form": r"`key: \|` for exactly one",
        "multiple-terminal-LF keep form": r"`key: \|\+` for multiple",
    }
    defects = {name for name, pattern in predicates.items() if re.search(pattern, normalized) is None}
    narrowing_patterns = {
        "outside-authorship condition narrowed to named fields":
            r"(?:condition|contract|rule).*?appl(?:y|ies) only to (?:the )?(?:fields?|names?).*?(?:named|listed|above|today)",
    }
    defects.update(
        name for name, pattern in narrowing_patterns.items()
        if re.search(pattern, normalized) is not None
    )
    return defects


def markdown_fenced_blocks(source):
    """Return Markdown fence bodies; inline code is deliberately invisible."""
    lines = source.splitlines()
    blocks = []
    line_index = 0
    while line_index < len(lines):
        opening = re.match(r"^\s*(`{3,}|~{3,})([^`]*)$", lines[line_index])
        if opening is None:
            line_index += 1
            continue
        fence = opening.group(1)
        fence_character = re.escape(fence[0])
        closing = re.compile(r"^\s*" + fence_character + "{" + str(len(fence)) + r",}\s*$")
        body_start = line_index + 1
        line_index = body_start
        while line_index < len(lines) and closing.match(lines[line_index]) is None:
            line_index += 1
        if line_index == len(lines):
            raise SystemExit("unclosed Markdown fence while deriving Frontmatter Quoting examples")
        blocks.append((opening.group(2).strip().lower(), "\n".join(lines[body_start:line_index])))
        line_index += 1
    return blocks


def canonical_frontmatter_schema(source):
    candidates = [
        body for info, body in markdown_fenced_blocks(source)
        if info in {"yaml", "yml"} and re.search(r"^id:\s*REQ-", body, flags=re.MULTILINE)
    ]
    if len(candidates) != 1:
        raise SystemExit(
            f"Frontmatter Quoting inventory found {len(candidates)} canonical YAML schema fences, want one"
        )
    return candidates[0]


def frontmatter_schema_inventory(source):
    schema = canonical_frontmatter_schema(source)
    governed_fields = set()
    encoder_owned_fields = set()
    annotation_defects = set()
    for schema_line in schema.splitlines():
        field_match = re.match(r"^([a-z_][a-z0-9_]*):(?P<rest>.*)$", schema_line)
        if field_match is None or "raw user text" not in field_match.group("rest").lower():
            continue
        field_name = field_match.group(1)
        governed_fields.add(field_name)
        annotation = field_match.group("rest").lower()
        encoder_owned = "escaping encoder" in annotation
        if encoder_owned:
            encoder_owned_fields.add(field_name)
        if "frontmatter quoting" not in annotation and not encoder_owned:
            annotation_defects.add(
                f"governed schema field {field_name} lacks Frontmatter Quoting or an encoder discriminator"
            )
    # Count, not a copied field list: the schema is the inventory. The floor is
    # only a deletion guard for today's seven annotated fields; additions grow
    # the set without editing this check.
    if len(governed_fields) < 7:
        annotation_defects.add(
            f"governed schema inventory has {len(governed_fields)} fields, want at least seven"
        )
    return schema, governed_fields, encoder_owned_fields, annotation_defects


def shipped_markdown_sources(repository_root):
    """Read Markdown from every source package declared by the suite manifest."""
    manifest_path = repository_root / "suite/modules.tsv"
    module_roots = []
    for manifest_line in manifest_path.read_text().splitlines()[1:]:
        if not manifest_line.strip():
            continue
        source_path, _ = manifest_line.split("\t", 1)
        module_root = repository_root / source_path
        if not module_root.is_dir():
            raise SystemExit(f"shipped module root is missing: {module_root}")
        module_roots.append(module_root)
    if not module_roots:
        raise SystemExit("suite manifest yielded no shipped module roots")
    return {
        str(markdown_path): markdown_path.read_text()
        for module_root in module_roots
        for markdown_path in module_root.rglob("*.md")
    }


def is_do_work_record_schema(fenced_body):
    """Distinguish do-work records from unrelated YAML in shipped packages."""
    req_record = (
        re.search(r"^\s*id:\s*REQ-", fenced_body, flags=re.MULTILINE) is not None
        and re.search(r"^\s*status:\s*", fenced_body, flags=re.MULTILINE) is not None
    )
    ur_record = (
        re.search(r"^\s*id:\s*UR-", fenced_body, flags=re.MULTILINE) is not None
        and re.search(r"^\s*requests:\s*", fenced_body, flags=re.MULTILINE) is not None
    )
    return req_record or ur_record


def frontmatter_fenced_example_defects(
    markdown_sources,
    governed_fields,
    encoder_owned_fields,
    inventory_source_name,
    inventory_schema,
):
    defects = set()
    independent_example_count = 0
    fenced_texts = []
    for source_name, source in markdown_sources.items():
        for _, fenced_body in markdown_fenced_blocks(source):
            fenced_texts.append(fenced_body)
            if not is_do_work_record_schema(fenced_body):
                continue
            inventory_block = (
                source_name == inventory_source_name and fenced_body == inventory_schema
            )
            for block_line in fenced_body.splitlines():
                field_match = re.match(r"^\s*([a-z_][a-z0-9_]*):\s*(.*)$", block_line)
                if field_match is None or field_match.group(1) not in governed_fields:
                    continue
                field_name, scalar = field_match.group(1), field_match.group(2).lstrip()
                if not inventory_block:
                    independent_example_count += 1
                safe_hand_authored = scalar.startswith("'") or scalar.startswith("|")
                safe_encoder_owned = field_name in encoder_owned_fields and scalar.startswith('"')
                if not safe_hand_authored and not safe_encoder_owned:
                    defects.add(
                        f"{source_name}: fenced {field_name} example is not single-quoted, literal, or encoder-owned"
                    )
    if independent_example_count == 0:
        defects.add("no governed field appears in an independent do-work record example")

    inline_counterexample = 'title: "Fix: A " # B"'
    work_reference_text = markdown_sources.get(str(work_reference_path), "")
    if inline_counterexample not in work_reference_text:
        defects.add("Frontmatter Quoting inline counterexample disappeared")
    if any(inline_counterexample in fenced_text for fenced_text in fenced_texts):
        defects.add("Frontmatter Quoting inline counterexample leaked into the fenced-example scan")
    return defects


work_reference = work_reference_path.read_text()
live_frontmatter_defects = frontmatter_contract_defects(work_reference)
if live_frontmatter_defects:
    raise SystemExit(
        "Frontmatter Quoting block-scalar contract defects: "
        + ", ".join(sorted(live_frontmatter_defects))
    )

(
    live_schema,
    governed_frontmatter_fields,
    encoder_owned_frontmatter_fields,
    schema_annotation_defects,
) = frontmatter_schema_inventory(work_reference)
if schema_annotation_defects:
    raise SystemExit(
        "Frontmatter Quoting schema annotation defects: "
        + ", ".join(sorted(schema_annotation_defects))
    )

repository_root = work_reference_path.parents[3]
all_shipped_markdown_sources = shipped_markdown_sources(repository_root)
live_fenced_example_defects = frontmatter_fenced_example_defects(
    all_shipped_markdown_sources,
    governed_frontmatter_fields,
    encoder_owned_frontmatter_fields,
    str(work_reference_path),
    live_schema,
)
if live_fenced_example_defects:
    raise SystemExit(
        "Frontmatter Quoting fenced-example defects: "
        + ", ".join(sorted(live_fenced_example_defects))
    )

schema_only_sources = {
    str(work_reference_path): (
        'The inline counterexample remains title: "Fix: A " # B"\n\n'
        "```yaml\n" + live_schema + "\n```\n"
    )
}
schema_only_defects = frontmatter_fenced_example_defects(
    schema_only_sources,
    governed_frontmatter_fields,
    encoder_owned_frontmatter_fields,
    str(work_reference_path),
    live_schema,
)
if "no governed field appears in an independent do-work record example" not in schema_only_defects:
    raise SystemExit(
        "Frontmatter Quoting schema-only non-vacuity trial escaped; found "
        + repr(sorted(schema_only_defects))
    )

frontmatter_mutations = (
    ("universal field condition removed", "The condition is the rule: whenever a frontmatter value carries text nobody in this pipeline composed", "Use this rule for frontmatter values written by the current actions", "outside-authorship condition governs every field"),
    ("preflight citation removed", "Outside-text containment's accepted-text preflight", "ordinary text", "outside-text preflight inherited"),
    ("normalization allowed", "refuses and reports text that fails it instead of normalizing bytes", "normalizes text that fails it", "failed preflight refused and reported"),
    ("blank-line indentation omitted", "(blank lines included)", "(non-blank lines only)", "every physical line indented"),
    ("strip used for every multiline value", "`key: |-` for zero terminal LF bytes, `key: |` for exactly one, and `key: |+` for multiple", "`key: |-` for every multiline value", "one-terminal-LF clip form"),
)
for mutation_name, old, new, expected_defect in frontmatter_mutations:
    mutated = work_reference.replace(old, new, 1)
    if mutated == work_reference:
        raise SystemExit(f"Frontmatter Quoting mutation {mutation_name!r} changed nothing")
    defects = frontmatter_contract_defects(mutated)
    if expected_defect not in defects:
        raise SystemExit(
            f"Frontmatter Quoting mutation {mutation_name!r} escaped {expected_defect!r}; "
            f"found {sorted(defects)!r}"
        )

frontmatter_narrowed = work_reference.replace(
    "they are illustrative, never a list to check against",
    "they are illustrative, never a list to check against. Despite that condition, "
    "this contract applies only to the fields named above",
    1,
)
if frontmatter_narrowed == work_reference:
    raise SystemExit("Frontmatter Quoting insertion narrowing changed nothing")
narrowed_defects = frontmatter_contract_defects(frontmatter_narrowed)
if "outside-authorship condition narrowed to named fields" not in narrowed_defects:
    raise SystemExit(
        "Frontmatter Quoting current-fields-only insertion escaped; found "
        + repr(sorted(narrowed_defects))
    )

schema_without_title_inheritance = live_schema.replace(
    "**Frontmatter Quoting** contract above", "local quoting rule", 1
)
if schema_without_title_inheritance == live_schema:
    raise SystemExit("Frontmatter Quoting schema-annotation deletion changed nothing")
work_reference_without_title_inheritance = work_reference.replace(
    live_schema, schema_without_title_inheritance, 1
)
_, _, _, deleted_annotation_defects = frontmatter_schema_inventory(
    work_reference_without_title_inheritance
)
if not any("governed schema field title lacks" in defect for defect in deleted_annotation_defects):
    raise SystemExit(
        "Frontmatter Quoting schema-annotation deletion escaped; found "
        + repr(sorted(deleted_annotation_defects))
    )

unsafe_title_schema = live_schema.replace(
    "title: 'Short descriptive title'", "title: Short descriptive title", 1
)
if unsafe_title_schema == live_schema:
    raise SystemExit("Frontmatter Quoting unsafe-scalar mutation changed nothing")
unsafe_title_reference = work_reference.replace(live_schema, unsafe_title_schema, 1)
unsafe_title_sources = dict(all_shipped_markdown_sources)
unsafe_title_sources[str(work_reference_path)] = unsafe_title_reference
_, unsafe_fields, unsafe_encoder_fields, unsafe_annotation_defects = frontmatter_schema_inventory(
    unsafe_title_reference
)
if unsafe_annotation_defects:
    raise SystemExit(
        "Frontmatter Quoting unsafe-scalar mutation damaged the inventory: "
        + repr(sorted(unsafe_annotation_defects))
    )
unsafe_title_defects = frontmatter_fenced_example_defects(
    unsafe_title_sources,
    unsafe_fields,
    unsafe_encoder_fields,
    str(work_reference_path),
    unsafe_title_schema,
)
if not any("fenced title example" in defect for defect in unsafe_title_defects):
    raise SystemExit(
        "Frontmatter Quoting unsafe-scalar mutation escaped; found "
        + repr(sorted(unsafe_title_defects))
    )

toolbox_review_path = repository_root / "skills/do-work-toolbox/actions/code-review.md"
toolbox_review_key = str(toolbox_review_path)
toolbox_review = all_shipped_markdown_sources.get(toolbox_review_key)
if toolbox_review is None:
    raise SystemExit("Frontmatter Quoting shipped toolbox review source is missing")
unsafe_toolbox_review = toolbox_review.replace(
    "title: '[<impact token>] Code review: [brief description]'",
    "title: [<impact token>] Code review: [brief description]",
    1,
)
if unsafe_toolbox_review == toolbox_review:
    raise SystemExit("Frontmatter Quoting toolbox unsafe-title mutation changed nothing")
unsafe_toolbox_sources = dict(all_shipped_markdown_sources)
unsafe_toolbox_sources[toolbox_review_key] = unsafe_toolbox_review
unsafe_toolbox_defects = frontmatter_fenced_example_defects(
    unsafe_toolbox_sources,
    governed_frontmatter_fields,
    encoder_owned_frontmatter_fields,
    str(work_reference_path),
    live_schema,
)
if not any(
    toolbox_review_key in defect and "fenced title example" in defect
    for defect in unsafe_toolbox_defects
):
    raise SystemExit(
        "Frontmatter Quoting shipped toolbox unsafe-title mutation escaped; found "
        + repr(sorted(unsafe_toolbox_defects))
    )

future_field_line = (
    "future_user_text: 'safe future value'   # raw user text — "
    "**Frontmatter Quoting** contract above\n"
)
future_schema = live_schema.replace(
    "title: 'Short descriptive title'", future_field_line + "title: 'Short descriptive title'", 1
)
if future_schema == live_schema:
    raise SystemExit("Frontmatter Quoting future-field mutation changed nothing")
future_reference = work_reference.replace(live_schema, future_schema, 1)
future_sources = dict(all_shipped_markdown_sources)
future_sources[str(work_reference_path)] = future_reference
future_toolbox_review = toolbox_review.replace(
    "title: '[<impact token>] Code review: [brief description]'",
    "title: '[<impact token>] Code review: [brief description]'\n"
    "future_user_text: unsafe future value",
    1,
)
if future_toolbox_review == toolbox_review:
    raise SystemExit("Frontmatter Quoting future-field independent example changed nothing")
future_sources[toolbox_review_key] = future_toolbox_review
_, future_fields, future_encoder_fields, future_annotation_defects = frontmatter_schema_inventory(
    future_reference
)
if "future_user_text" not in future_fields or future_annotation_defects:
    raise SystemExit(
        "Frontmatter Quoting future field did not join the derived inventory cleanly: "
        + repr(sorted(future_fields)) + " / " + repr(sorted(future_annotation_defects))
    )
future_field_defects = frontmatter_fenced_example_defects(
    future_sources,
    future_fields,
    future_encoder_fields,
    str(work_reference_path),
    future_schema,
)
if not any("fenced future_user_text example" in defect for defect in future_field_defects):
    raise SystemExit(
        "Frontmatter Quoting future unsafe field escaped; found "
        + repr(sorted(future_field_defects))
    )

required_single_quoted_examples = (
    (capture_reference_path, r"^title: 'Add keyboard shortcuts'(?:\s+#.*)?$", "UR input title"),
    (review_path, r"^title: '\[<impact token>\] Review fix: \[brief description\]'", "review follow-up title"),
    (sample_path, r"^title: 'Add user avatar component'(?:\s+#.*)?$", "sample archived title"),
    (capture_guide_path, r"^title: 'Brief descriptive title'$", "capture-guide title"),
    (work_guide_path, r"^assigned_to: 'cloud-alpha'$", "work-guide earmark"),
)
for path, pattern, label in required_single_quoted_examples:
    example_text = path.read_text()
    example_match = re.search(pattern, example_text, flags=re.MULTILINE)
    if example_match is None:
        raise SystemExit(f"{path}: {label} is not in the required single-quoted form")
    citation_window = example_text[example_match.start():example_match.end() + 500]
    if "Frontmatter Quoting" not in citation_window:
        raise SystemExit(f"{path}: {label} does not cite Frontmatter Quoting at the example")

ur_template = capture_reference_path.read_text()
full_input = section(ur_template, r"^## Full Verbatim Input\n", r"^---\n\*Captured:")
if re.search(r"^> `+", full_input, flags=re.MULTILINE) is None:
    raise SystemExit("capture-reference UR Full Verbatim Input example is still live Markdown rather than quoted fenced text")

PY
then
  printf 'FAIL: outside-text body writers and hand-authored frontmatter examples must inherit the canonical containment/quoting contracts at their own write sites (REQ-360).\n' >&2
  fail_count=$((fail_count + 1))
fi

assert_block_contains \
  "$verify_input_block" \
  'Repeating `--against` batches several reversals into one queue scan' \
  'verify-requests decision mode must accept repeated --against sources and share one queue scan instead of multiplying semantic cost.'

assert_block_contains \
  "$verify_input_block" \
  'Reject a capture-QA target mixed with any `--against`' \
  'verify-requests must keep capture QA targets mutually exclusive with decision-revalidation sources.'

assert_block_contains \
  "$verify_revalidation_block" \
  'status: superseded.*exactly one.*rel: superseded-by' \
  'decision-file sources must resolve one explicit superseded-by successor rather than guess a replacement from similar prose.'

assert_block_contains \
  "$verify_revalidation_block" \
  'reject a successor Decision that only confirms or restates the old choice without a semantic reversal' \
  'a superseded relation whose Decision text did not actually reverse the choice is not a revalidation source.'

assert_block_contains \
  "$verify_revalidation_block" \
  'reject absolute paths, `\.\.` escapes, symlink escapes, missing files, and directories' \
  'decision-file sources must stay inside the repository even through symlink resolution.'

assert_block_contains \
  "$verify_revalidation_block" \
  'builder_decided: true.*old `Recommended:` choice and new answer' \
  'REQ sources must be answered builder-decision follow-ups that preserve both sides of the reversal.'

assert_block_contains \
  "$verify_revalidation_block" \
  'exact `do-work/queue/REQ-\*\.md` files' \
  'decision revalidation must scan the canonical queue rather than a UR capture-time requests array or archive history.'

assert_block_contains \
  "$verify_revalidation_block" \
  'exclude every terminal status: `completed`, `completed-with-issues`, `failed`, and `cancelled`' \
  'decision revalidation must exclude all terminal queue records, including failed REQs that are not terminally resolved for UR closure.'

assert_block_contains \
  "$verify_revalidation_block" \
  'Scan every other canonical state, including `blocked`' \
  'decision revalidation must include blocked queued work because it remains unfinished.'

assert_block_contains \
  "$verify_revalidation_block" \
  'Exclude every source follow-up REQ by id' \
  'decision revalidation must exclude answered follow-up sources from matching their own rejected recommendation.'

assert_block_contains \
  "$verify_revalidation_block" \
  'Do not scan `do-work/working/`, `do-work/archive/`, or legacy REQs' \
  'decision revalidation v1 must exclude claimed and archived bodies from semantic cost.'

assert_block_contains \
  "$verify_revalidation_block" \
  'list every claimed REQ id as \*\*excluded from v1\*\*' \
  'decision revalidation must disclose claimed work even though it does not semantically scan those living logs.'

assert_block_contains \
  "$verify_revalidation_block" \
  'complete file.*not only `## Scope`' \
  'decision revalidation must inspect complete queued REQs because Scope normally does not exist before claim.'

assert_block_contains \
  "$verify_revalidation_block" \
  '^\- \*\*Likely affected:\*\*' \
  'decision revalidation must retain a high-evidence class for explicit citations and direct restatements.'

assert_block_contains \
  "$verify_revalidation_block" \
  '^\- \*\*Possibly affected:\*\*' \
  'decision revalidation must retain an evidence-backed semantic class for copied assumptions without citations.'

assert_block_contains \
  "$verify_revalidation_block" \
  'mentions the old decision only as history, a rejected alternative, superseded context' \
  'decision revalidation must not flag historical or rejected mentions as live dependencies.'

assert_block_contains \
  "$verify_revalidation_block" \
  'a short exact excerpt from the REQ' \
  'every decision-revalidation candidate must quote evidence from the queued REQ.'

assert_block_contains \
  "$verify_revalidation_block" \
  'old → replacement conflict in plain language' \
  'every decision-revalidation candidate must explain the old-to-new conflict.'

assert_block_contains \
  "$verify_revalidation_block" \
  'copyable, provenance-preserving next step' \
  'every decision-revalidation candidate must include a reconciliation command.'

assert_block_contains \
  "$verify_revalidation_block" \
  'changes no REQ body, frontmatter, status, or location' \
  'decision revalidation must remain report-only and never turn a candidate into pending-answers.'

assert_block_contains \
  "$verify_revalidation_block" \
  'explicit `--against` invocation always proceeds.*10,000-word confirmation threshold belongs only to clarify' \
  'explicit decision scans must show cost and proceed; only the automatic clarify caller gets the 10,000-word gate.'

assert_block_contains \
  "$clarify_revalidation_block" \
  'If `overturned_decision_sources` is empty, skip this step' \
  'clarify must invoke decision revalidation only when a builder decision was actually overturned.'

assert_contains \
  "actions/clarify.md" \
  'source enters it only when its REQ carries `builder_decided: true` and the user.s answer is semantically different' \
  'clarify must recognize a reversal from the stored builder recommendation and the genuinely different user answer.'

assert_contains \
  "actions/clarify.md" \
  'Confirmation, discard, discovery approval.*never enter the set' \
  'clarify must not scan after confirmations, discards, or discovered-task approvals.'

assert_block_contains \
  "$clarify_revalidation_block" \
  'once with every source id in the set' \
  'clarify must batch all reversals into one queue scan.'

assert_block_contains \
  "$clarify_revalidation_block" \
  '10,000 queued words or fewer.*automatically' \
  'clarify must auto-scan ordinary queues through the approved 10,000-word threshold.'

assert_block_contains \
  "$clarify_revalidation_block" \
  'Above \*\*10,000 queued words\*\*.*file count, word count, approximate 1\.3–1\.6-tokens-per-word input range.*Ask one choice' \
  'clarify must show the large-queue estimate and ask before an automatic over-threshold semantic scan.'

assert_block_contains \
  "$clarify_revalidation_block" \
  'combined command with repeated flags' \
  'declining a large automatic scan must preserve every reversal in one copyable explicit command.'

assert_contains \
  "actions/verify-requests.md" \
  'default mode is \*\*capture QA\*\*' \
  'the existing verify-requests capture-QA path must remain the default after --against is added.'

assert_contains \
  "crew-members/prompt-injection.md" \
  'verify-requests \(re-reads UR input\.md verbatim or compares decision sources with complete queued REQs\)' \
  'the prompt-injection caller inventory must cover both verify-requests ingestion modes.'

# The canonical repository-maintainer gate owns its production command inventory.
# Exercise only its recursion-safe focused contract here; invoking normal mode would
# recursively run this aggregate.
maintainer_verify_probe="$repo_root/_dev/tests/maintainer-verify.sh"
if [ ! -x "$maintainer_verify_probe" ]; then
  printf 'FAIL: _dev/tests/maintainer-verify.sh is missing or not executable — repository-native verification has no canonical gate.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$maintainer_verify_probe" --self-test; then
  printf 'FAIL: canonical maintainer verification self-test failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The suite manifest is executable input to both update paths and the fresh installer. Keep
# its path-safety and exact-module contract in one behavioral probe rather than duplicating
# grep assertions for each caller.
suite_manifest_probe="$repo_root/_dev/tests/suite-manifest-contract.sh"
if [ ! -f "$suite_manifest_probe" ]; then
  printf 'FAIL: _dev/tests/suite-manifest-contract.sh is missing — the suite layout has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$suite_manifest_probe"; then
  printf 'FAIL: suite manifest contract probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Published package links have two valid execution contexts: this source tree and the
# manifest-mapped install under .claude/skills/. The local probe parses rendered Markdown
# outside code, validates first-party raw/blob paths without network access, and keeps the
# installed core changelog byte-identical to the release source.
shipped_reference_probe="$repo_root/_dev/tests/shipped-package-reference-contract.sh"
if [ ! -f "$shipped_reference_probe" ]; then
  printf 'FAIL: _dev/tests/shipped-package-reference-contract.sh is missing — shipped package references have no source/install coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$shipped_reference_probe"; then
  printf 'FAIL: shipped package reference contract failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Markdown-prescribed shell is shipped executable guidance, while hook scripts are shipped
# executable code. Exercise both through one attributable syntax/lint probe, including its
# negative fixture so a silently-neutered extractor cannot report a decorative green result.
action_shell_probe="$repo_root/_dev/tests/action-shell-blocks.sh"
if [ ! -x "$action_shell_probe" ]; then
  printf 'FAIL: _dev/tests/action-shell-blocks.sh is missing or not executable — shipped shell guidance has no lint coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$action_shell_probe" --self-test; then
  printf 'FAIL: action shell-block negative self-test failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$action_shell_probe"; then
  printf 'FAIL: shipped shell-block lint failed (see the attributed FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The SessionStart banner is deliberately fail-soft: malformed or missing runtime inputs must
# produce its unknown/zero defaults, not let shell options terminate the hook before output.
session_start_probe="$repo_root/_dev/tests/session-start-hook-behavior.sh"
if [ ! -x "$session_start_probe" ]; then
  printf 'FAIL: _dev/tests/session-start-hook-behavior.sh is missing or not executable — the startup banner fallback has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$session_start_probe"; then
  printf 'FAIL: SessionStart hook behavior probes failed (see the fixture FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Shared prescribed-shell rationale has one shipped home; callers retain only local intent,
# commands, and explicit pointers so a primitive fix cannot drift across prose copies again.
prescribed_shell_probe="$repo_root/_dev/tests/prescribed-shell-canonicalization.sh"
if [ ! -x "$prescribed_shell_probe" ]; then
  printf 'FAIL: _dev/tests/prescribed-shell-canonicalization.sh is missing or not executable — shell primitive restatements have no ratchet.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$prescribed_shell_probe"; then
  printf 'FAIL: prescribed shell primitive canonicalization failed (see the attributed FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The prescribed-shell suite reports how many cases it holds, and that figure is only as
# honest as the rule that produces it. Before REQ-339 the rule matched `^# <one-token>: `,
# so `# generate-report-image caller contract:` and `# generate-report-image, interrupted
# directly:` were real cases no figure ever mentioned — nothing failed, the number just
# read low. This fixture carries three headers in three spellings plus four near-misses
# that must stay uncounted, so a narrowed rule and an over-broad one both fail here.
prescribed_shell_counter="$repo_root/_dev/tests/prescribed-shell-case-count.sh"
if [ ! -f "$prescribed_shell_counter" ]; then
  printf 'FAIL: _dev/tests/prescribed-shell-case-count.sh is missing — the reported case count has no definition to hold to.\n' >&2
  fail_count=$((fail_count + 1))
else
  case_count_fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/prescribed-shell-case-count.XXXXXX")"
  cat > "$case_count_fixture_root/probe-counter.sh" <<'PROBE_COUNTER_FIXTURE'
# Fixture execution proofs for probe-counter.
# probe-counter: a bare header counts.
# probe-counter, interrupted: a comma-qualified header counts.
# probe-counter caller contract: a space-qualified header counts.
# The wrapped description below is prose, not a header: it only carries a colon.
# reaching a colon mid-sentence: a continuation line must not count.
# qualifying a request: a word that only starts with the script name must not count.
# probe-counter.sh named inside a sentence (REQ-000: still not a header).
PROBE_COUNTER_FIXTURE
  # 'none' rather than 0 so a counter that could not run stays distinguishable from one
  # that legitimately found nothing.
  if ! observed_header_count="$(
    # shellcheck source=_dev/tests/prescribed-shell-case-count.sh
    source "$prescribed_shell_counter"
    count_named_case_headers "$case_count_fixture_root/probe-counter.sh"
  )"; then
    observed_header_count='none'
  fi
  if [ "$observed_header_count" != '3' ]; then
    printf 'FAIL: count_named_case_headers reported %s of the fixture 3 case headers — the suite'"'"'s reported case count no longer matches what the case files hold.\n' \
      "$observed_header_count" >&2
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$case_count_fixture_root"
fi

# REQ-168 removed specific generic guidance and an arbitrary commit-size heuristic. Keep
# those exact deletions without turning its historical audit into a living registry or
# banning future incident-backed defensive sections.
defensive_surface_probe="$repo_root/_dev/tests/defensive-surface-audit.sh"
if [ ! -x "$defensive_surface_probe" ]; then
  printf 'FAIL: _dev/tests/defensive-surface-audit.sh is missing or not executable — REQ-168 exact deletion regressions have no ratchet.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$defensive_surface_probe"; then
  printf 'FAIL: defensive-surface exact deletion regression failed (see the attributed FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Behavioral probes for tools/checks/record-commit-hash.sh. Kept in their own file because
# they build a throwaway git repo and run the real script rather than grepping prose — but
# invoked from here, since nothing auto-discovers _dev/tests/*.sh and an uninvoked probe file
# is dead weight that reads as coverage.
record_commit_hash_probe="$repo_root/_dev/tests/record-commit-hash-guards.sh"
if [ ! -f "$record_commit_hash_probe" ]; then
  printf 'FAIL: _dev/tests/record-commit-hash-guards.sh is missing — the write-back guards have no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash <(sed \
  -e "s|^repo_root=.*|repo_root=\"$repo_root\"|" \
  -e 's|\$repo_root/tools/checks/|\$repo_root/skills/do-work/tools/checks/|g' \
  "$record_commit_hash_probe"); then
  printf 'FAIL: record-commit-hash guard probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Behavioral probes for tools/do-work-update.sh — same reasoning as above, and the same
# no-auto-discovery caveat. These build a synthetic install plus a stubbed upstream fetch
# because the no-rollback-copy contract is a runtime property (no `.bak` left behind, a
# mid-update failure that reports instead of restoring) that no grep can assert.
update_script_probe="$repo_root/_dev/tests/update-script-behavior.sh"
if [ ! -f "$update_script_probe" ]; then
  printf 'FAIL: _dev/tests/update-script-behavior.sh is missing — the updater has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$update_script_probe"; then
  printf 'FAIL: update-script behavior probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Behavioral probes for skills/do-work/tools/do-work-cli.sh, the launcher that builds the
# Go command on demand. Same no-auto-discovery caveat as the probes above: nothing walks
# _dev/tests/*.sh, so an uninvoked probe file is dead weight that reads as coverage. It
# covers build-on-demand when the binary is absent, argv passthrough including a spaced
# argument, refusal to run stale output after a failed rebuild, the Go version floor, and
# the absence of a leftover build temp — all runtime properties no grep can assert, proved
# against a fake toolchain so the probe never depends on the installed Go.
do_work_cli_launcher_probe="$repo_root/_dev/tests/do-work-cli-launcher-behavior.sh"
if [ ! -x "$do_work_cli_launcher_probe" ]; then
  printf 'FAIL: _dev/tests/do-work-cli-launcher-behavior.sh is missing or not executable — the do-work-cli launcher has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$do_work_cli_launcher_probe"; then
  printf 'FAIL: do-work-cli launcher behavior probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The live modular suite package boundaries need a contract independent from the active
# root bootstrap tools. This checks the staged router, required core
# runtime, hook targets, and the ban on leaking repository maintainer instructions.
staged_skills_probe="$repo_root/_dev/tests/staged-skills-contract.sh"
if [ ! -f "$staged_skills_probe" ]; then
  printf 'FAIL: _dev/tests/staged-skills-contract.sh is missing — staged skill packages have no boundary coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$staged_skills_probe"; then
  printf 'FAIL: staged skills contract probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Fresh installation is a four-module/configuration transaction with its own hermetic
# bootstrap. Keep these probes separate because they build several isolated Git clients,
# including malformed, fallback, rollback, and interruption states.
suite_installer_probe="$repo_root/_dev/tests/install-suite-behavior.sh"
if [ ! -f "$suite_installer_probe" ]; then
  printf 'FAIL: _dev/tests/install-suite-behavior.sh is missing — the full-suite installer has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$suite_installer_probe"; then
  printf 'FAIL: suite installer behavior probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The P50 estimator's contracts — deterministic output, nearest-5 rounding, the
# 10-minute floor, dependency-graph critical-path math, and the print-only
# backwards-compatibility guarantee — are runtime properties no grep can assert.
p50_estimator_probe="$repo_root/_dev/tests/p50-estimator-determinism.sh"
if [ ! -f "$p50_estimator_probe" ]; then
  printf 'FAIL: _dev/tests/p50-estimator-determinism.sh is missing — the P50 estimator has no lock-in coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$p50_estimator_probe"; then
  printf 'FAIL: P50 estimator probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The cheaper-model selector decides which queued REQs a smaller model may build.
# Its contracts are runtime predicates no grep can assert: the Schema Read
# Contract's `trivial` alias (most pre-rename REQs spell it that way), the
# dependency-readiness filter that keeps its caller's explicit `do-work run
# REQ-NNN` handoff from becoming a `depends_on` bypass, the vetoes for work with
# no objective gate, and the deliberate NON-veto for `tdd: true`.
select_simple_reqs_probe="$repo_root/_dev/tests/select-simple-reqs-behavior.sh"
if [ ! -f "$select_simple_reqs_probe" ]; then
  printf 'FAIL: _dev/tests/select-simple-reqs-behavior.sh is missing — the cheaper-model selector has no lock-in coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$select_simple_reqs_probe"; then
  printf 'FAIL: cheaper-model selector probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Managed Just sections are a byte-preserving ownership boundary, not a prose convention.
# Exercise the real utility across replacement, append, creation, malformed
# markers, filename variants, spaces, modes, idempotence, and Just parsing.
replace_section_tool="$repo_root/tools/replace-text-section.sh"
if [ ! -x "$replace_section_tool" ]; then
  printf 'FAIL: tools/replace-text-section.sh is missing or not executable — managed recipe ownership has no implementation.\n' >&2
  fail_count=$((fail_count + 1))
else
  section_workdir="$(mktemp -d)"
  section_file="$section_workdir/managed-section.just"
  template_file="$section_workdir/complete-template.just"
  printf '# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\n' > "$section_file"
  printf 'set shell := ["bash", "-cu"]\n\n# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\n' > "$template_file"
  chmod 750 "$template_file"

  byte_target="$section_workdir/project with spaces/Justfile"
  mkdir -p "$(dirname "$byte_target")"
  printf 'prefix\000byte\n# >>> do-work:recipes >>>\nold:\n    echo old\n# <<< do-work:recipes <<<\nsuffix\n' > "$byte_target"
  chmod 640 "$byte_target"
  expected_target="$section_workdir/expected-byte-target"
  printf 'prefix\000byte\n# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\nsuffix\n' > "$expected_target"
  if ! "$replace_section_tool" --target "$byte_target" --section-file "$section_file"; then
    printf 'FAIL: replace-text-section could not replace one valid managed section.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$byte_target" "$expected_target"; then
    printf 'FAIL: replace-text-section changed bytes outside the managed section or wrote the wrong replacement.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  target_mode="$(stat -c '%a' "$byte_target" 2>/dev/null || stat -f '%Lp' "$byte_target" 2>/dev/null || true)"
  if [ "$target_mode" != 640 ]; then
    printf 'FAIL: replace-text-section changed the existing target mode (got %s, want 640).\n' "$target_mode" >&2
    fail_count=$((fail_count + 1))
  fi
  cp "$byte_target" "$section_workdir/idempotent-snapshot"
  if ! "$replace_section_tool" --target "$byte_target" --section-file "$section_file" \
     || ! cmp -s "$byte_target" "$section_workdir/idempotent-snapshot"; then
    printf 'FAIL: replace-text-section is not byte-idempotent on repeated execution.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  absent_target="$section_workdir/absent project/.justfile"
  mkdir -p "$(dirname "$absent_target")"
  if ! "$replace_section_tool" --target "$absent_target" --section-file "$section_file" --template-file "$template_file" \
     || ! cmp -s "$absent_target" "$template_file"; then
    printf 'FAIL: replace-text-section did not create an absent target from the complete supplied template.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  absent_mode="$(stat -c '%a' "$absent_target" 2>/dev/null || stat -f '%Lp' "$absent_target" 2>/dev/null || true)"
  if [ "$absent_mode" != 750 ]; then
    printf 'FAIL: replace-text-section did not preserve the complete template mode on create (got %s, want 750).\n' "$absent_mode" >&2
    fail_count=$((fail_count + 1))
  fi

  for justfile_variant in justfile Justfile .justfile; do
    variant_target="$section_workdir/variants/$justfile_variant"
    mkdir -p "$(dirname "$variant_target")"
    printf 'custom-%s:\n    echo untouched\n' "$justfile_variant" > "$variant_target"
    if ! "$replace_section_tool" --target "$variant_target" --section-file "$section_file"; then
      printf 'FAIL: replace-text-section could not append to marker-free %s.\n' "$justfile_variant" >&2
      fail_count=$((fail_count + 1))
    elif [ "$(grep -c '^# >>> do-work:recipes >>>$' "$variant_target")" -ne 1 ] \
      || ! grep -q "^custom-$justfile_variant:" "$variant_target"; then
      printf 'FAIL: replace-text-section duplicated ownership or changed custom content in %s.\n' "$justfile_variant" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  reserved_section_file="$section_workdir/reserved-section.just"
  awk '
    $0 == "# >>> do-work:recipes >>>" { inside=1 }
    inside { print }
    $0 == "# <<< do-work:recipes <<<" { found_end=1; exit }
    END { if (!inside || !found_end) exit 1 }
  ' "$repo_root/skills/do-work-board/justfile.template" > "$reserved_section_file"

  collision_index=0
  for reserved_recipe_name in run-kanban run-kanban-cli kanban-static kanban-summary run-do-work-update; do
    collision_index=$((collision_index + 1))
    collision_target="$section_workdir/collision-$collision_index.just"
    case "$reserved_recipe_name" in
      run-kanban)
        printf 'run-kanban:\n    echo collision\n' > "$collision_target"
        ;;
      run-kanban-cli)
        printf '@run-kanban-cli $view="open:all": dependency-recipe\n    echo collision\n' > "$collision_target"
        ;;
      kanban-static)
        printf 'kanban-static destination="build/output": dependency-recipe\n    echo collision\n' > "$collision_target"
        ;;
      kanban-summary)
        printf 'alias kanban-summary := custom-summary\n' > "$collision_target"
        ;;
      run-do-work-update)
        printf 'run-do-work-update:\r\n    echo collision\r\n' > "$collision_target"
        ;;
    esac
    cp "$collision_target" "$collision_target.before"
    collision_output="$section_workdir/collision-$collision_index.out"
    if "$replace_section_tool" --target "$collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$collision_output" 2>&1; then
      printf 'FAIL: replace-text-section accepted external reserved recipe or alias %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fq "reserved Just recipe or alias outside managed section: $reserved_recipe_name" "$collision_output"; then
      printf 'FAIL: replace-text-section collision error did not name %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$collision_target" "$collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collision %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  multiline_header_index=0
  for multiline_header_kind in \
    ordinary-single ordinary-double ordinary-backtick \
    triple-single triple-double triple-backtick; do
    multiline_header_index=$((multiline_header_index + 1))
    multiline_header_payload='payload'
    case "$multiline_header_kind" in
      ordinary-single) multiline_header_delimiter="'" ;;
      ordinary-double) multiline_header_delimiter='"' ;;
      ordinary-backtick) multiline_header_delimiter='`' ;;
      triple-single)
        multiline_header_delimiter="'''"
        multiline_header_payload="payload's"
        ;;
      triple-double)
        multiline_header_delimiter='"""'
        multiline_header_payload='payload"s'
        ;;
      triple-backtick)
        multiline_header_delimiter='```'
        multiline_header_payload='payload`s'
        ;;
    esac
    multiline_header_target="$section_workdir/multiline-header-$multiline_header_index.just"
    printf 'run-kanban value=%s\n%s\n%s:\n    echo collision\n' \
      "$multiline_header_delimiter" "$multiline_header_payload" "$multiline_header_delimiter" \
      > "$multiline_header_target"
    if command -v just >/dev/null 2>&1 \
      && ! just --justfile "$multiline_header_target" --list >/dev/null 2>&1; then
      printf 'FAIL: %s multiline-default recipe-header fixture is not valid Just syntax.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
      continue
    fi
    cp "$multiline_header_target" "$multiline_header_target.before"
    multiline_header_output="$section_workdir/multiline-header-$multiline_header_index.out"
    if "$replace_section_tool" --target "$multiline_header_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$multiline_header_output" 2>&1; then
      printf 'FAIL: replace-text-section accepted reserved recipe with a %s multiline default.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fq \
      'reserved Just recipe or alias outside managed section: run-kanban' \
      "$multiline_header_output"; then
      printf 'FAIL: replace-text-section did not name the reserved recipe with a %s multiline default.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$multiline_header_target" "$multiline_header_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting a %s multiline-default collision.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  raw_multiline_target="$section_workdir/raw-multiline-string.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      "raw_value := '''" \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      "\\'''" \
      'custom-after:' \
      '    echo after'
  } > "$raw_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$raw_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: raw multiline-string collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$raw_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking raw multiline-string content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$raw_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting a raw multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  cooked_multiline_target="$section_workdir/cooked-multiline-string-crlf.just"
  {
    printf '%s\r\n' \
      'custom-before:' \
      '    echo before' \
      'cooked_value := """' \
      'run-kanban-cli:' \
      'odd escaped delimiter remains: \"""' \
      'alias kanban-summary := ignored' \
      'even escaped delimiter closes: \\"""' \
      'joined_value := """closed""" + """open' \
      'run-do-work-update:' \
      '"""' \
      'custom-after:' \
      '    echo after'
  } > "$cooked_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$cooked_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: cooked multiline-string collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$cooked_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking cooked multiline-string content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$cooked_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting a cooked multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_and_command_multiline_target="$section_workdir/ordinary-and-command-multiline-literals.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      "raw_value := '" \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      "raw backslash is literal: \\'" \
      'cooked_value := "' \
      'run-kanban-cli:' \
      'escaped double quote remains: \"' \
      'alias kanban-summary := ignored' \
      'even backslashes close: \\"' \
      'command_value := ```' \
      '  printf "%s\n" safe' \
      'run-do-work-update:' \
      'alias kanban-static := ignored' \
      '```' \
      'custom-after:' \
      '    echo after'
  } > "$ordinary_and_command_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary-quote and triple-backtick collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$ordinary_and_command_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking ordinary-quote or triple-backtick content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting ordinary-quote and triple-backtick literals.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_backtick_multiline_target="$section_workdir/ordinary-backtick-multiline-command.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      'command_value := `' \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      'raw backslash closes: \` + `' \
      'printf safe' \
      '`' \
      'custom-after:' \
      '    echo after'
  } > "$ordinary_backtick_multiline_target"
  cp "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_target.before"
  ordinary_backtick_multiline_expected="$section_workdir/ordinary-backtick-multiline-command.expected"
  cp "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_expected"
  printf '\n' >> "$ordinary_backtick_multiline_expected"
  cat "$reserved_section_file" >> "$ordinary_backtick_multiline_expected"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary multiline-backtick collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$ordinary_backtick_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking ordinary multiline-backtick content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_expected"; then
    printf 'FAIL: replace-text-section wrote unexpected bytes after accepting an ordinary multiline-backtick command.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting an ordinary multiline-backtick command.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_backtick_collision_target="$section_workdir/ordinary-backtick-nearby-collisions.just"
  {
    printf '%s\n' \
      '# unmatched ordinary backtick in comment: `' \
      'one_line_command := `printf safe`' \
      'custom-body:' \
      '    echo unmatched-backtick `' \
      'custom-summary:' \
      '    echo custom' \
      'run-kanban:' \
      '    echo real collision before command' \
      'command_value := `' \
      'kanban-static:' \
      'raw backslash closes: \` + `' \
      'alias kanban-static := ignored' \
      '`' \
      'alias kanban-summary := custom-summary' \
      '@run-kanban-cli:' \
      '    echo real collision after command' \
      'run-do-work-update:' \
      '    echo real collision after inactive forms'
  } > "$ordinary_backtick_collision_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_collision_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary multiline-backtick nearby-collision control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$ordinary_backtick_collision_target" "$ordinary_backtick_collision_target.before"
    ordinary_backtick_collision_output="$section_workdir/ordinary-backtick-nearby-collisions.out"
    expected_ordinary_backtick_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update, run-kanban, run-kanban-cli'
    if "$replace_section_tool" --target "$ordinary_backtick_collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$ordinary_backtick_collision_output" 2>&1; then
      printf 'FAIL: replace-text-section ignored real reserved definitions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fxq -- "$expected_ordinary_backtick_collision" "$ordinary_backtick_collision_output"; then
      printf 'FAIL: replace-text-section did not report only exact sorted collisions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$ordinary_backtick_collision_target" "$ordinary_backtick_collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collisions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  ordinary_and_command_collision_target="$section_workdir/ordinary-and-command-nearby-collisions.just"
  {
    printf '%s\n' \
      'custom-summary:' \
      '    echo custom' \
      'run-kanban:' \
      '    echo real collision before raw string' \
      "raw_value := '" \
      'alias kanban-static := ignored' \
      "raw backslash is literal: \\'" \
      'alias kanban-summary := custom-summary' \
      'cooked_value := "' \
      'kanban-static:' \
      'escaped double quote remains: \"' \
      'alias kanban-static := ignored' \
      'even backslashes close: \\"' \
      '@run-kanban-cli:' \
      '    echo real collision before command literal' \
      'command_value := ```' \
      '  printf "%s\n" safe' \
      'kanban-static:' \
      '```' \
      'run-do-work-update:' \
      '    echo real collision after command literal'
  } > "$ordinary_and_command_collision_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_collision_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary-quote and triple-backtick nearby-collision control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$ordinary_and_command_collision_target" "$ordinary_and_command_collision_target.before"
    ordinary_and_command_collision_output="$section_workdir/ordinary-and-command-nearby-collisions.out"
    expected_ordinary_and_command_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update, run-kanban, run-kanban-cli'
    if "$replace_section_tool" --target "$ordinary_and_command_collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$ordinary_and_command_collision_output" 2>&1; then
      printf 'FAIL: replace-text-section ignored real reserved definitions around ordinary-quote or triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fxq -- "$expected_ordinary_and_command_collision" "$ordinary_and_command_collision_output"; then
      printf 'FAIL: replace-text-section did not report only exact sorted collisions around ordinary-quote and triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$ordinary_and_command_collision_target" "$ordinary_and_command_collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collisions around ordinary-quote or triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  inactive_literal_forms_target="$section_workdir/inactive-literal-forms-nearby-collision.just"
  {
    printf '%s\n' \
      '# delimiter-looking comment: ```' \
      'one_line_command := `printf safe`' \
      'custom-body:' \
      "    printf '%s\\n' '\`\`\`'" \
      'kanban-static:' \
      '    echo real collision after inactive forms'
  } > "$inactive_literal_forms_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$inactive_literal_forms_target" --list >/dev/null 2>&1; then
    printf 'FAIL: inactive multiline-literal opener control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$inactive_literal_forms_target" "$inactive_literal_forms_target.before"
    inactive_literal_forms_output="$section_workdir/inactive-literal-forms-nearby-collision.out"
    expected_inactive_literal_forms_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-static'
    if "$replace_section_tool" --target "$inactive_literal_forms_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$inactive_literal_forms_output" 2>&1; then
      printf 'FAIL: replace-text-section let an inactive comment, recipe body, or one-line backtick hide a real collision.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fxq -- "$expected_inactive_literal_forms_collision" "$inactive_literal_forms_output"; then
      printf 'FAIL: replace-text-section did not report the exact collision after inactive multiline-literal opener forms.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$inactive_literal_forms_target" "$inactive_literal_forms_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting a collision after inactive literal forms.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  multiline_collision_target="$section_workdir/multiline-nearby-collisions.just"
  {
    printf '%s\n' \
      '# delimiter-looking comment: """' \
      "single_token := '\"\"\"'" \
      "double_token := \"'''\"" \
      'backtick_token := `printf "\"\"\""`' \
      'custom-shell:' \
      "    echo '\"\"\"'" \
      'custom-summary:' \
      '    echo custom' \
      'alias kanban-summary := custom-summary' \
      'payload := """' \
      'run-kanban:' \
      'alias run-kanban-cli := ignored' \
      'even escaped delimiter closes: \\"""' \
      'run-do-work-update:' \
      '    echo real collision'
  } > "$multiline_collision_target"
  cp "$multiline_collision_target" "$multiline_collision_target.before"
  multiline_collision_output="$section_workdir/multiline-nearby-collisions.out"
  expected_multiline_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update'
  if "$replace_section_tool" --target "$multiline_collision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions >"$multiline_collision_output" 2>&1; then
    printf 'FAIL: replace-text-section ignored real reserved definitions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! grep -Fxq -- "$expected_multiline_collision" "$multiline_collision_output"; then
    printf 'FAIL: replace-text-section did not report only exact sorted real collisions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$multiline_collision_target" "$multiline_collision_target.before"; then
    printf 'FAIL: replace-text-section changed the target after rejecting collisions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  noncollision_target="$section_workdir/noncollisions.just"
  {
    printf '# run-kanban:\n'
    printf 'reserved_value := "run-kanban-cli:"\n'
    printf "[doc('kanban-static: is reserved')]\n"
    printf 'custom-recipe: kanban-summary\n'
    printf '    echo run-do-work-update:\n'
    printf 'run-kanban-extra:\n    echo prefix\n'
    printf 'alias custom-summary := kanban-summary\n\n'
    cat "$reserved_section_file"
  } > "$noncollision_target"
  if ! "$replace_section_tool" --target "$noncollision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated comments, variables, attributes, dependencies, bodies, prefixes, aliases to reserved recipes, or managed definitions as collisions.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  external_collision_target="$section_workdir/external-after-managed.just"
  cat "$reserved_section_file" > "$external_collision_target"
  printf '\nrun-kanban:\n    echo external collision\n' >> "$external_collision_target"
  cp "$external_collision_target" "$external_collision_target.before"
  if "$replace_section_tool" --target "$external_collision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section ignored a reserved recipe outside an existing managed span.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$external_collision_target" "$external_collision_target.before"; then
    printf 'FAIL: replace-text-section mutated a managed target before rejecting its external collision.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  retired_flag_target="$section_workdir/retired-flag.just"
  printf 'custom-only:\n    echo untouched\n' > "$retired_flag_target"
  cp "$retired_flag_target" "$retired_flag_target.before"
  if "$replace_section_tool" --target "$retired_flag_target" --section-file "$section_file" --migrate-legacy-do-work >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section still accepts the retired legacy-migration flag.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$retired_flag_target" "$retired_flag_target.before"; then
    printf 'FAIL: replace-text-section changed the target after rejecting the retired flag.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  malformed_index=0
  for malformed_content in \
    '# >>> do-work:recipes >>>|one|# >>> do-work:recipes >>>|# <<< do-work:recipes <<<' \
    '# <<< do-work:recipes <<<|one|# >>> do-work:recipes >>>' \
    '# >>> do-work:recipes >>>|one' \
    'one|# <<< do-work:recipes <<<' ; do
    malformed_index=$((malformed_index + 1))
    malformed_target="$section_workdir/malformed-$malformed_index.just"
    printf '%s\n' "$malformed_content" | tr '|' '\n' > "$malformed_target"
    cp "$malformed_target" "$malformed_target.before"
    if "$replace_section_tool" --target "$malformed_target" --section-file "$section_file" >/dev/null 2>&1; then
      printf 'FAIL: replace-text-section accepted malformed marker case %s.\n' "$malformed_index" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$malformed_target" "$malformed_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting malformed marker case %s.\n' "$malformed_index" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$section_workdir/variants/justfile" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced a Justfile that does not parse.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$section_workdir"
fi

# Root assertion inputs must use the tracked lowercase justfile spelling. This
# source-level replay stays case-sensitive even on developer filesystems that
# resolve Justfile and justfile to the same inode, and it pins the four late
# assertions so deleting one cannot make the casing check vacuously green.
if ! python3 - "$repo_root/_dev/tests/contract-regressions.sh" <<'PY'
import pathlib
import re
import sys

contract_source = pathlib.Path(sys.argv[1]).read_text()
file_assertion_pattern = re.compile(
    r'^assert_(?:contains|file_not_contains) \\\n  "([^"]+)" \\\n  \'([^\']*)\' \\',
    flags=re.MULTILINE,
)
file_assertions = file_assertion_pattern.findall(contract_source)
expected_root_patterns = (
    r'^# >>> do-work:recipes >>>$',
    r'^# <<< do-work:recipes <<<$',
    r'skills/do-work-board/tools/queue-kanban',
    r'skill_root="\$project_root/skills/do-work".*\$skill_root/tools/do-work-update\.sh',
)

problems = []
for expected_pattern in expected_root_patterns:
    matching_paths = [
        file_path
        for file_path, pattern_text in file_assertions
        if pattern_text == expected_pattern
    ]
    if matching_paths != ["justfile"]:
        problems.append(f"{expected_pattern!r} uses {matching_paths!r}, want ['justfile']")

wrong_case_inputs = sorted({
    file_path
    for file_path, _ in file_assertions
    if file_path.casefold() == "justfile" and file_path != "justfile"
})
if wrong_case_inputs:
    problems.append(f"live root-file inputs use non-tracked casing: {wrong_case_inputs!r}")

if problems:
    raise SystemExit("; ".join(problems))
PY
then
  printf 'FAIL: late root-justfile assertions must exist and use the tracked lowercase path.\n' >&2
  fail_count=$((fail_count + 1))
fi

assert_contains \
  "justfile" \
  '^# >>> do-work:recipes >>>$' \
  'root justfile must open the exact managed do-work recipe section.'
assert_contains \
  "justfile" \
  '^# <<< do-work:recipes <<<$' \
  'root justfile must close the exact managed do-work recipe section.'
assert_contains \
  "skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go" \
  'managedsection\.ReplaceSection' \
  'suite installer must reconcile recipes through the managed-section utility.'
assert_contains \
  "justfile" \
  'skills/do-work-board/tools/queue-kanban' \
  'root Justfile must build the canonical board sibling source.'
assert_contains \
  "justfile" \
  'skill_root="\$project_root/skills/do-work".*\$skill_root/tools/do-work-update\.sh' \
  'root Justfile fallback must invoke the canonical modular core updater.'
assert_contains \
  "skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go" \
  'os\.CreateTemp\(parentDirectory' \
  'the managed-section replacer must create its temporary file in the target directory for atomic replacement.'
assert_contains \
  "skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go" \
  'os\.Rename\(temporaryPath, path\)' \
  'the managed-section replacer must atomically rename the validated temporary over the target.'
assert_contains \
  "tools/replace-text-section.sh" \
  'do-work-cli' \
  'replace-text-section.sh must remain a launcher over the do-work-cli replace-section command.'

if [ "$fail_count" -gt 0 ]; then
  exit 1
fi

printf 'Contract regression checks passed.\n'
