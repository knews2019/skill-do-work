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

screenshot_dispatch_block="$(sed -n '/^## Dispatch/,/^## Safety/p' "$repo_root/skills/do-work/SKILL.md")"
if ! grep -Fq 'mktemp -d do-work/user-requests/.pending-assets/capture.XXXXXX' <<<"$screenshot_dispatch_block" \
  || ! grep -Fq '$screenshot_dispatch_directory/screenshot-{n}.png' <<<"$screenshot_dispatch_block"; then
  fail 'core dispatch must allocate one exclusive staging directory per screenshot-bearing capture'
fi
if ! grep -Eq 'actions/capture\.md.*Step 4' <<<"$screenshot_dispatch_block"; then
  fail 'core dispatch must name capture Step 4 as the owner of staged screenshot cleanup'
fi

if ! python3 - "$repo_root/skills/do-work/actions/capture.md" <<'PY'
import pathlib
import sys

capture_action_text = pathlib.Path(sys.argv[1]).read_text()
try:
    screenshot_step = capture_action_text.split(
        "### Step 4: Handle Screenshots", 1
    )[1].split("### Step 5: Write Files", 1)[0]
except IndexError:
    raise SystemExit("capture action has no bounded screenshot step")
screenshot_shell = screenshot_step.split("```bash", 1)[1].split("```", 1)[0]

required_fragments = (
    'staged_screenshot_path="<exact staged screenshot path supplied by the dispatcher>"',
    'screenshot_staging_directory="$(dirname "$staged_screenshot_path")"',
    'REQ-[num]-screenshot-{n}-[slug].png',
    'screenshot_copy_path="$(mktemp "${screenshot_asset_path}.copying.XXXXXX")"',
    'cp "$staged_screenshot_path" "$screenshot_copy_path"',
    'cmp -s "$staged_screenshot_path" "$screenshot_copy_path"',
    'ln "$screenshot_copy_path" "$screenshot_asset_path"',
    'rm -f "$screenshot_copy_path"',
    'rm "$staged_screenshot_path"',
    'rmdir "$screenshot_staging_directory"',
    "staged source preserved",
)
fragment_positions = []
for required_fragment in required_fragments:
    fragment_position = screenshot_shell.find(required_fragment)
    if fragment_position < 0:
        raise SystemExit(
            f"capture screenshot step is missing {required_fragment!r}"
        )
    fragment_positions.append(fragment_position)

if fragment_positions != sorted(fragment_positions):
    raise SystemExit(
        "capture screenshot cleanup must follow copy, comparison, no-clobber install, and source removal"
    )

directory_cleanup = screenshot_step.split(
    'rmdir "$screenshot_staging_directory"', 1
)[1].split("fi", 1)[0]
if "false" in directory_cleanup:
    raise SystemExit("empty per-dispatch directory cleanup must be best-effort")
PY
then
  fail 'capture Step 4 must isolate, no-clobber, and safely clean staged screenshots'
fi

if ! python3 - "$repo_root/skills/do-work/actions/capture.md" <<'PY'
import pathlib
import subprocess
import sys
import tempfile

capture_action_text = pathlib.Path(sys.argv[1]).read_text()
screenshot_step = capture_action_text.split(
    "### Step 4: Handle Screenshots", 1
)[1].split("### Step 5: Write Files", 1)[0]
screenshot_shell = screenshot_step.split("```bash", 1)[1].split("```", 1)[0].strip()


def rendered_shell(dispatch_name: str, screenshot_number: int) -> str:
    return (
        screenshot_shell.replace(
            "<exact staged screenshot path supplied by the dispatcher>",
            f"do-work/user-requests/.pending-assets/{dispatch_name}/screenshot-{{n}}.png",
        )
        .replace("UR-NNN", "UR-042")
        .replace("[num]", "042")
        .replace("[slug]", "example")
        .replace("{n}", str(screenshot_number))
    )


def run_capture(
    project_root: pathlib.Path,
    dispatch_name: str,
    screenshot_number: int,
    *,
    fail_rmdir: bool = False,
    fail_staged_source_removal: bool = False,
):
    command = rendered_shell(dispatch_name, screenshot_number)
    if fail_rmdir:
        command = "rmdir() { return 1; }\n" + command
    if fail_staged_source_removal:
        staged_source = (
            f"do-work/user-requests/.pending-assets/{dispatch_name}/"
            f"screenshot-{screenshot_number}.png"
        )
        command = (
            'rm() {\n'
            f'  if [ "$#" -eq 1 ] && [ "$1" = "{staged_source}" ]; then return 1; fi\n'
            '  command rm "$@"\n'
            '}\n'
            + command
        )
    return subprocess.run(
        ["bash", "-c", command],
        cwd=project_root,
        text=True,
        capture_output=True,
        check=False,
    )


with tempfile.TemporaryDirectory(prefix="screenshot-capture-contract.") as temporary_root:
    project_root = pathlib.Path(temporary_root)
    dispatch_directory = project_root / "do-work/user-requests/.pending-assets/capture.multi"
    dispatch_directory.mkdir(parents=True)
    (dispatch_directory / "screenshot-1.png").write_bytes(b"first")
    (dispatch_directory / "screenshot-2.png").write_bytes(b"second")

    for screenshot_number in (1, 2):
        capture_result = run_capture(project_root, "capture.multi", screenshot_number)
        if capture_result.returncode != 0:
            raise SystemExit(
                f"multi-screenshot capture {screenshot_number} failed: {capture_result.stderr}"
            )

    asset_directory = project_root / "do-work/user-requests/UR-042/assets"
    if (asset_directory / "REQ-042-screenshot-1-example.png").read_bytes() != b"first":
        raise SystemExit("first screenshot did not survive the second screenshot capture")
    if (asset_directory / "REQ-042-screenshot-2-example.png").read_bytes() != b"second":
        raise SystemExit("second screenshot was not installed at its distinct asset path")
    if dispatch_directory.exists():
        raise SystemExit("empty per-dispatch staging directory was not removed")

    race_a_directory = project_root / "do-work/user-requests/.pending-assets/capture.race-a"
    race_b_directory = project_root / "do-work/user-requests/.pending-assets/capture.race-b"
    race_a_directory.mkdir(parents=True)
    race_b_directory.mkdir(parents=True)
    race_a_source = race_a_directory / "screenshot-6.png"
    race_b_source = race_b_directory / "screenshot-6.png"
    race_a_bytes = b"verified-dispatch-a"
    race_b_bytes = b"competing-dispatch-b"
    race_a_source.write_bytes(race_a_bytes)
    race_b_source.write_bytes(race_b_bytes)
    race_destination = asset_directory / "REQ-042-screenshot-6-example.png"
    race_a_verified = project_root / "race-a-verified"
    race_b_copied = project_root / "race-b-copied"

    race_a_prefix = f'''cmp() {{
  command cmp "$@"
  comparison_status=$?
  if [ "$comparison_status" -eq 0 ]; then
    : > "{race_a_verified}"
    while [ ! -e "{race_b_copied}" ]; do sleep 0.01; done
  fi
  return "$comparison_status"
}}
'''
    race_b_prefix = f'''cp() {{
  while [ ! -e "{race_a_verified}" ]; do sleep 0.01; done
  command cp "$@"
  copy_status=$?
  if [ "$copy_status" -eq 0 ]; then : > "{race_b_copied}"; fi
  return "$copy_status"
}}
cmp() {{
  command cmp "$@"
  comparison_status=$?
  while [ ! -e "{race_destination}" ]; do sleep 0.01; done
  return "$comparison_status"
}}
'''
    race_a_process = subprocess.Popen(
        ["bash", "-c", race_a_prefix + rendered_shell("capture.race-a", 6)],
        cwd=project_root,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    race_b_process = subprocess.Popen(
        ["bash", "-c", race_b_prefix + rendered_shell("capture.race-b", 6)],
        cwd=project_root,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    try:
        race_a_stdout, race_a_stderr = race_a_process.communicate(timeout=10)
        race_b_stdout, race_b_stderr = race_b_process.communicate(timeout=10)
    except subprocess.TimeoutExpired:
        race_a_process.terminate()
        race_b_process.terminate()
        raise SystemExit("coordinated screenshot collision timed out")

    if race_a_process.returncode != 0 or race_b_process.returncode == 0:
        raise SystemExit(
            "coordinated collision did not install exactly one dispatch: "
            f"A={race_a_process.returncode} ({race_a_stderr}), "
            f"B={race_b_process.returncode} ({race_b_stderr})"
        )
    if race_destination.read_bytes() != race_a_bytes:
        raise SystemExit("winning dispatch published bytes from another dispatch")
    if race_a_source.exists():
        raise SystemExit("winning dispatch retained its installed staged source")
    if race_b_source.read_bytes() != race_b_bytes:
        raise SystemExit("losing dispatch did not retain its staged source")
    if list(asset_directory.glob(f"{race_destination.name}.copying*")):
        raise SystemExit("coordinated collision left a screenshot temporary copy behind")

    collision_directory = project_root / "do-work/user-requests/.pending-assets/capture.collision"
    collision_directory.mkdir(parents=True)
    collision_source = collision_directory / "screenshot-3.png"
    collision_source.write_bytes(b"new")
    collision_destination = asset_directory / "REQ-042-screenshot-3-example.png"
    collision_destination.write_bytes(b"existing")
    collision_result = run_capture(project_root, "capture.collision", 3)
    if collision_result.returncode == 0:
        raise SystemExit("capture overwrote an existing permanent screenshot destination")
    if collision_source.read_bytes() != b"new" or collision_destination.read_bytes() != b"existing":
        raise SystemExit("destination collision did not preserve both source and existing asset")

    cleanup_directory = project_root / "do-work/user-requests/.pending-assets/capture.cleanup"
    cleanup_directory.mkdir(parents=True)
    cleanup_source = cleanup_directory / "screenshot-4.png"
    cleanup_source.write_bytes(b"cleanup")
    cleanup_result = run_capture(project_root, "capture.cleanup", 4, fail_rmdir=True)
    if cleanup_result.returncode != 0:
        raise SystemExit("best-effort staging-directory cleanup invalidated a verified capture")
    if cleanup_source.exists():
        raise SystemExit("verified staged source was not removed before best-effort directory cleanup")
    if (asset_directory / "REQ-042-screenshot-4-example.png").read_bytes() != b"cleanup":
        raise SystemExit("verified destination was not retained after rmdir failure")
    if "could not be removed" not in cleanup_result.stderr:
        raise SystemExit("best-effort rmdir failure was not reported")

    source_cleanup_directory = (
        project_root / "do-work/user-requests/.pending-assets/capture.source-cleanup"
    )
    source_cleanup_directory.mkdir(parents=True)
    source_cleanup_source = source_cleanup_directory / "screenshot-5.png"
    source_cleanup_source.write_bytes(b"source-cleanup")
    source_cleanup_destination = asset_directory / "REQ-042-screenshot-5-example.png"
    source_cleanup_result = run_capture(
        project_root,
        "capture.source-cleanup",
        5,
        fail_staged_source_removal=True,
    )
    if source_cleanup_result.returncode != 0:
        raise SystemExit("staged-source cleanup failure invalidated a verified capture")
    if source_cleanup_source.read_bytes() != b"source-cleanup":
        raise SystemExit("failed staged-source cleanup did not preserve the source")
    if source_cleanup_destination.read_bytes() != b"source-cleanup":
        raise SystemExit("failed staged-source cleanup did not preserve the verified destination")
    if pathlib.Path(f"{source_cleanup_destination}.copying").exists():
        raise SystemExit("temporary screenshot copy remained after staged-source cleanup failure")
    if "staged source could not be removed" not in source_cleanup_result.stderr:
        raise SystemExit("best-effort staged-source cleanup failure was not reported")

    later_collision_result = run_capture(project_root, "capture.source-cleanup", 5)
    if later_collision_result.returncode == 0:
        raise SystemExit("later destination collision did not retain no-clobber behavior")
    if source_cleanup_source.read_bytes() != b"source-cleanup":
        raise SystemExit("later collision did not preserve the staged source")
    if source_cleanup_destination.read_bytes() != b"source-cleanup":
        raise SystemExit("later collision changed the verified destination")
PY
then
  fail 'capture Step 4 executable screenshot lifecycle regressions failed'
fi

if ! python3 - "$repo_root/skills/do-work/tools/checks/preflight.sh" <<'PY'
import pathlib
import subprocess
import sys
import tempfile

preflight_script = pathlib.Path(sys.argv[1])
pathological_names = (
    'space name.txt',
    'quote"name.txt',
    'star*.txt',
    'starZZ.txt',
    'bracket[ab].txt',
    'bracketa.txt',
)
with tempfile.TemporaryDirectory(prefix="preflight-path-contract.") as temporary_root:
    project_root = pathlib.Path(temporary_root)
    subprocess.run(["git", "init", "-q", str(project_root)], check=True)
    for pathological_name in pathological_names:
        (project_root / pathological_name).write_text(pathological_name)
    (project_root / "do-work").mkdir()
    (project_root / "do-work/ignored file.txt").write_text("ignored")

    preflight_result = subprocess.run(
        ["bash", str(preflight_script)],
        cwd=project_root,
        text=True,
        capture_output=True,
        check=False,
    )
    if preflight_result.returncode != 0:
        raise SystemExit(f"preflight exited nonzero: {preflight_result.stderr}")
    output_lines = preflight_result.stdout.splitlines()
    for pathological_name in pathological_names:
        rendered_path = f"  {pathological_name}"
        if output_lines.count(rendered_path) != 1:
            raise SystemExit(
                f"preflight did not preserve {pathological_name!r} exactly once:\n"
                + preflight_result.stdout
            )
    if any("ignored file.txt" in output_line for output_line in output_lines):
        raise SystemExit("preflight reported a dirty path inside do-work/")
PY
then
  fail 'preflight must preserve spaces, quotes, and glob characters in dirty filenames'
fi

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

if ! python3 - \
  "$repo_root" \
  "$repo_root/_dev/tests/fixtures/retired-core-moved-command-triggers.tsv" \
  "${sibling_route_contracts[@]}" <<'PY'
import collections
import csv
import pathlib
import re
import sys
import tempfile

repository_root = pathlib.Path(sys.argv[1])
trigger_fixture = pathlib.Path(sys.argv[2])
declared_owner_actions = {
    tuple(contract.split("|", 1)) for contract in sys.argv[3:]
}

fixture_lines = [
    line
    for line in trigger_fixture.read_text().splitlines()
    if line and not line.startswith("#")
]
fixture_reader = csv.DictReader(fixture_lines, delimiter="\t")
expected_header = ["owner", "canonical_action", "match_kind", "legacy_trigger"]
if fixture_reader.fieldnames != expected_header:
    raise SystemExit(
        f"invalid retired-trigger fixture header: {fixture_reader.fieldnames!r}"
    )
trigger_rows = list(fixture_reader)
for row_number, row in enumerate(trigger_rows, 1):
    if None in row or any(not value or value != value.strip() for value in row.values()):
        raise SystemExit(f"incomplete retired-trigger fixture row {row_number}: {row!r}")

legacy_triggers = [row["legacy_trigger"] for row in trigger_rows]
duplicate_triggers = sorted(
    trigger
    for trigger, count in collections.Counter(legacy_triggers).items()
    if count != 1
)
if duplicate_triggers:
    raise SystemExit(
        "duplicate retired-trigger fixture entries: " + ", ".join(duplicate_triggers)
    )

invalid_owner_actions = sorted(
    {
        (row["owner"], row["canonical_action"])
        for row in trigger_rows
        if (row["owner"], row["canonical_action"]) not in declared_owner_actions
    }
)
if invalid_owner_actions:
    raise SystemExit(
        "invalid retired-trigger owner/action pairs: "
        + ", ".join("/".join(pair) for pair in invalid_owner_actions)
    )

allowed_match_kinds = {
    "direct",
    "install-space",
    "install-prefix",
    "setup-space",
    "install-head",
}
invalid_match_kinds = sorted(
    {row["match_kind"] for row in trigger_rows} - allowed_match_kinds
)
if invalid_match_kinds:
    raise SystemExit(
        "invalid retired-trigger match kinds: " + ", ".join(invalid_match_kinds)
    )

expected_direct_counts = {
    ("do-work-board", "board"): 6,
    ("do-work-knowledge", "bkb"): 4,
    ("do-work-knowledge", "memory"): 5,
    ("do-work-knowledge", "dream"): 5,
    ("do-work-knowledge", "interview"): 3,
    ("do-work-knowledge", "prompts"): 2,
    ("do-work-toolbox", "validate-feedback"): 8,
    ("do-work-toolbox", "code-review"): 5,
    ("do-work-toolbox", "ui-review"): 6,
    ("do-work-toolbox", "present-work"): 7,
    ("do-work-toolbox", "ai-report"): 7,
    ("do-work-toolbox", "slop-check"): 3,
    ("do-work-toolbox", "quick-wins"): 7,
    ("do-work-toolbox", "scan-ideas"): 8,
    ("do-work-toolbox", "deep-explore"): 5,
    ("do-work-toolbox", "prime"): 6,
    ("do-work-toolbox", "inspect"): 5,
    ("do-work-toolbox", "note"): 3,
    ("do-work-toolbox", "stray-check"): 8,
    ("do-work-toolbox", "tidy-repo"): 10,
    ("do-work-toolbox", "tutorial"): 4,
}
actual_direct_counts = collections.Counter(
    (row["owner"], row["canonical_action"])
    for row in trigger_rows
    if row["match_kind"] == "direct"
)
if actual_direct_counts != collections.Counter(expected_direct_counts):
    raise SystemExit(
        "retired-trigger direct inventory counts drifted: "
        f"expected {expected_direct_counts!r}, found {dict(actual_direct_counts)!r}"
    )

install_kind_prefixes = {
    "install-space": "install ",
    "install-prefix": "install-",
    "setup-space": "setup ",
}
install_target_maps = {}
for match_kind, trigger_prefix in install_kind_prefixes.items():
    kind_rows = [row for row in trigger_rows if row["match_kind"] == match_kind]
    malformed = [
        row["legacy_trigger"]
        for row in kind_rows
        if not row["legacy_trigger"].startswith(trigger_prefix)
        or row["legacy_trigger"] == trigger_prefix
    ]
    if malformed:
        raise SystemExit(
            f"malformed {match_kind} retired triggers: " + ", ".join(malformed)
        )
    target_map = {
        row["legacy_trigger"][len(trigger_prefix) :]: (
            row["owner"],
            row["canonical_action"],
        )
        for row in kind_rows
    }
    if len(target_map) != len(kind_rows):
        raise SystemExit(f"duplicate install targets in retired-trigger kind {match_kind}")
    install_target_maps[match_kind] = target_map

baseline_install_targets = install_target_maps["install-space"]
for match_kind, target_map in install_target_maps.items():
    if target_map != baseline_install_targets:
        raise SystemExit(
            f"retired install families disagree for {match_kind}: {target_map!r}"
        )
if len(baseline_install_targets) != 22:
    raise SystemExit(
        f"retired install inventory must contain 22 targets, found {len(baseline_install_targets)}"
    )
expected_install_ownership = collections.Counter(
    {
        ("do-work-board", "board"): 4,
        ("do-work-knowledge", "setup-memory"): 3,
        ("do-work-toolbox", "install"): 15,
    }
)
if collections.Counter(baseline_install_targets.values()) != expected_install_ownership:
    raise SystemExit(
        "retired install target ownership drifted: "
        f"{dict(collections.Counter(baseline_install_targets.values()))!r}"
    )

install_head_rows = [
    row for row in trigger_rows if row["match_kind"] == "install-head"
]
actual_install_heads = {
    (row["owner"], row["canonical_action"], row["legacy_trigger"])
    for row in install_head_rows
}
expected_install_heads = {
    ("do-work-toolbox", "install", "install"),
    ("do-work-toolbox", "install", "setup"),
    ("do-work-toolbox", "install", "install-"),
}
if actual_install_heads != expected_install_heads or len(install_head_rows) != 3:
    raise SystemExit(
        f"retired bare install heads drifted: {sorted(actual_install_heads)!r}"
    )

sorted_triggers = sorted(legacy_triggers, key=lambda trigger: (-len(trigger), trigger))
trigger_match_kinds = {
    row["legacy_trigger"]: row["match_kind"] for row in trigger_rows
}
retired_command_head = re.compile(r"(?<![A-Za-z0-9_-])do-work ")
forbidden_right_boundary = re.compile(r"[A-Za-z0-9_'-]")
install_target = re.compile(r"[A-Za-z0-9_]+(?:-[A-Za-z0-9_]+)*")
exempt_command_fragments = (
    "PROJECT_NAME — do-work queue board",
    '<title> text "do-work queue board"',
    'strings.Contains(bodyText, "do-work queue board")',
)

def collect_live_files(root):
    live_files = [root / "justfile"]
    live_files.extend(
        path
        for path in sorted((root / "skills").rglob("*"))
        if path.is_file() and path.name not in {"CHANGELOG.md", "queue-kanban"}
    )
    return [path for path in live_files if path.is_file()]


def find_exempt_occurrence_spans(text):
    spans = set()
    exempt_command = "do-work queue board"
    for fragment in exempt_command_fragments:
        command_offset = fragment.index(exempt_command)
        search_start = 0
        while True:
            fragment_start = text.find(fragment, search_start)
            if fragment_start < 0:
                break
            occurrence_start = fragment_start + command_offset
            spans.add(
                (occurrence_start, occurrence_start + len(exempt_command))
            )
            search_start = fragment_start + 1
    return spans


def find_retired_matches(text):
    matches = []
    exempt_occurrence_spans = find_exempt_occurrence_spans(text)
    for command_match in retired_command_head.finditer(text):
        command_start = command_match.end()
        command_remainder = text[command_start:]
        trigger = None
        boundary_index = None
        for candidate in sorted_triggers:
            if not command_remainder.startswith(candidate):
                continue
            candidate_boundary = command_start + len(candidate)
            if candidate == "install-" and candidate_boundary < len(text):
                target_match = install_target.match(text, candidate_boundary)
                if target_match:
                    candidate_boundary = target_match.end()
            if candidate_boundary < len(text) and forbidden_right_boundary.match(
                text[candidate_boundary]
            ):
                if (
                    text[candidate_boundary] == "'"
                    or trigger_match_kinds[candidate]
                    not in {"install-space", "install-prefix", "setup-space"}
                ):
                    break
                continue
            trigger = candidate
            boundary_index = candidate_boundary
            break
        if trigger is None or boundary_index is None:
            continue
        occurrence_span = (command_match.start(), boundary_index)
        if trigger == "queue board" and occurrence_span in exempt_occurrence_spans:
            continue
        matches.append((trigger, *occurrence_span))
    return matches


def matched_triggers(text):
    return [trigger for trigger, _, _ in find_retired_matches(text)]


adversarial_controls = {
    '<title> text "do-work queue board"': [],
    'strings.Contains(bodyText, "do-work queue board")': [],
    "PROJECT_NAME — do-work queue board; do-work queue board": ["queue board"],
    '<title> text "do-work queue board"; do-work queue board': ["queue board"],
    'strings.Contains(bodyText, "do-work queue board"); do-work queue board': [
        "queue board"
    ],
    "do-work install ui-design2": ["install"],
    "do-work setup ui-design2": ["setup"],
    "do-work install-custom-target": ["install-"],
}
adversarial_failures = []
for adversarial_control, expected_triggers in adversarial_controls.items():
    matches = matched_triggers(adversarial_control)
    if matches != expected_triggers:
        adversarial_failures.append(
            f"{adversarial_control!r}: expected {expected_triggers!r}, found {matches!r}"
        )
if adversarial_failures:
    raise SystemExit(
        "retired-trigger matcher missed adversarial controls:\n"
        + "\n".join(adversarial_failures)
    )


negative_controls = (
    "Do-Work Board Skill",
    "PROJECT_NAME — do-work queue board",
    '<title> text "do-work queue board"',
    'strings.Contains(bodyText, "do-work queue board")',
    "do-work board's testing view",
    "The board shows knowledge and toolbox status as ordinary nouns.",
    "This work pipeline runs after the CI pipeline and data pipeline.",
    "do-work run REQ-157",
    "do-work review code",
    "do-work help",
    "do-work install-custom-target's",
    "do-work-board board",
    "do-work-knowledge memory recall",
    "do-work-toolbox code-review",
)
for negative_control in negative_controls:
    if matched_triggers(negative_control):
        raise SystemExit(
            f"retired-trigger matcher rejected negative control: {negative_control!r}"
        )

positive_controls = {
    "Deprecated — do-work kanban": ["kanban"],
}
for positive_control, expected_triggers in positive_controls.items():
    matches = matched_triggers(positive_control)
    if matches != expected_triggers:
        raise SystemExit(
            f"retired-trigger matcher missed positive control {positive_control!r}: "
            f"expected {expected_triggers!r}, found {matches!r}"
        )

for row in trigger_rows:
    trigger = row["legacy_trigger"]
    negative_forms = [
        f"undo-work {trigger}",
        f"do-work {trigger}'s",
    ]
    if row["match_kind"] == "direct" or trigger == "setup":
        negative_forms.extend(
            (f"do-work {trigger}-suffix", f"do-work {trigger}suffix")
        )
    elif trigger == "install":
        negative_forms.append(f"do-work {trigger}suffix")
    elif trigger == "install-":
        negative_forms.append(f"do-work {trigger}-suffix")
    for negative_form in negative_forms:
        if matched_triggers(negative_form):
            raise SystemExit(
                f"retired-trigger matcher crossed a command boundary for {row!r}: "
                f"{negative_form!r}"
            )

with tempfile.TemporaryDirectory(prefix="retired-trigger-contract-") as temp_directory:
    mutation_root = pathlib.Path(temp_directory)
    root_mutation_file = mutation_root / "justfile"
    module_mutation_file = mutation_root / "skills/do-work-example/SKILL.md"
    module_mutation_file.parent.mkdir(parents=True)
    root_mutation_file.write_text(
        "\n".join(
            f"row-{row_number}: do-work {row['legacy_trigger']} --example"
            for row_number, row in enumerate(trigger_rows, 1)
        )
        + "\n"
    )
    module_mutation_file.write_text(
        "\n".join(
            f"row-{row_number}: `do-work {row['legacy_trigger']}:`"
            for row_number, row in enumerate(trigger_rows, 1)
        )
        + "\n"
    )

    historical_mutations = {
        mutation_root / "CHANGELOG.md": "do-work kanban\n",
        mutation_root / "skills/do-work/CHANGELOG.md": "do-work recall\n",
        mutation_root / "do-work/archive/REQ-historical.md": "do-work code review\n",
        mutation_root / "_dev/tests/fixtures/negative.md": "do-work describe changes\n",
    }
    for historical_file, historical_text in historical_mutations.items():
        historical_file.parent.mkdir(parents=True, exist_ok=True)
        historical_file.write_text(historical_text)

    collected_mutation_files = collect_live_files(mutation_root)
    expected_mutation_files = [root_mutation_file, module_mutation_file]
    if collected_mutation_files != expected_mutation_files:
        raise SystemExit(
            "live-surface collector included history/fixtures or missed a live surface: "
            f"{collected_mutation_files!r}"
        )
    for mutation_file in collected_mutation_files:
        mutation_lines = mutation_file.read_text().splitlines()
        if len(mutation_lines) != len(trigger_rows):
            raise SystemExit(f"incomplete mutation file: {mutation_file}")
        for row, mutation_line in zip(trigger_rows, mutation_lines):
            matches = matched_triggers(mutation_line)
            if len(matches) != 1:
                raise SystemExit(
                    f"expected exactly one retired match for {row!r} in "
                    f"{mutation_file.name}, found {len(matches)}"
                )
            if matches[0] != row["legacy_trigger"]:
                raise SystemExit(
                    f"retired-trigger row identity mismatch for {row!r}: "
                    f"matched {matches[0]!r}"
                )

violations = []
for live_file in collect_live_files(repository_root):
    for line_number, line in enumerate(live_file.read_text(errors="replace").splitlines(), 1):
        for legacy_trigger, _, _ in find_retired_matches(line):
            violations.append(
                f"{live_file.relative_to(repository_root)}:{line_number}: "
                f"do-work {legacy_trigger}"
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

board_help_contract="$(sed -n 's/^argument-hint: "\(board \[[^]]*\]\).*$/\1/p' \
  "$repo_root/skills/do-work-board/SKILL.md")"
if [ -z "$board_help_contract" ]; then
  fail 'do-work-board argument-hint must expose its board modes for core help'
elif ! grep -Fq "$board_help_contract" "$repo_root/skills/do-work/actions/help.md"; then
  fail "core help must mirror the board command and modes: $board_help_contract"
fi

for mirrored_tool_name in \
  install-do-work-suite.sh \
  replace-text-section.sh \
  validate-suite-manifest.sh; do
  canonical_tool_path="$repo_root/tools/$mirrored_tool_name"
  staged_tool_path="$repo_root/skills/do-work/tools/$mirrored_tool_name"
  if [ -f "$staged_tool_path" ] \
    && ! cmp -s "$canonical_tool_path" "$staged_tool_path"; then
    fail "staged core $mirrored_tool_name must be byte-identical to canonical $mirrored_tool_name"
  fi
done

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
  && grep -Eq -- 'queue-kanban(/queue-kanban)? next-version|--version-file.*<skill-root>/actions/version\.md' "$repo_root/skills/do-work/actions/work.md"; then
  fail 'core work must not use the managed suite version file as a consumer-project release source'
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
