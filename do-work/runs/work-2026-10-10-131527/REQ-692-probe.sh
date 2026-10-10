#!/usr/bin/env bash
# GREEN checks for REQ-692 (integrator test gate): validate-feedback gains the opt-in
# --capture [--run] chain, its no-flag Steps 1 to 5 and Output Format fence stay byte-identical,
# every read-only statement names the --capture exception, capture.md admits the hand-off, the
# toolbox help line names the flag, and core SKILL.md routes validate-feedback to the toolbox
# action above the verify row without any live file spelling a retired core trigger. Prose-only
# REQ, so the checks read the files; the two citation/shell-fence contracts run as well.
# With the argument "invariants" it runs only the checks that already hold at base
# (REQ-692-preflight-probe.sh uses that). Fails at base without that argument: the chain does
# not exist yet. Under 15 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
probe_mode="${1:-green}"

python3 - "$root" "$probe_mode" <<'PY'
import hashlib
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
probe_mode = sys.argv[2]
failures = []

action_path = root / "skills/do-work-toolbox/actions/validate-feedback.md"
action_text = action_path.read_text()
action_lines = action_text.split("\n")


def step_section(heading_prefix):
    collected = []
    inside = False
    for line in action_lines:
        if line.startswith(heading_prefix):
            inside = True
            collected.append(line)
            continue
        if inside and (line.startswith("### ") or line.startswith("## ")):
            break
        if inside:
            collected.append(line)
    return "\n".join(collected)


# Invariant 1: Steps 1 to 5 keep their bytes (requirement 3; base hashes taken at bd56c4b0).
base_step_hashes = {
    1: "0b46873932bbf06af431d60ca94698afac7a94d1ad2a1b02a8a06aa712232f46",
    2: "2b337dc5126e340df1f4b271b09517bc60a7158f534cede766e5aa8fd3f6f197",
    3: "18299de67d60785562154183582e89ca936e0c2f04f700c40f94cd4323ba49f8",
    4: "faca58ec02c200b64d7a18257d2f1ba3a10b20826a1498ff3c306d659030f3bb",
    5: "f889c6fc7534d2c2b02dd577d1fc49b5e50ebc4497d52a5bf5c95f9038da7dc4",
}
for step_number, expected_hash in base_step_hashes.items():
    section_text = step_section(f"### Step {step_number}:")
    if hashlib.sha256(section_text.encode()).hexdigest() != expected_hash:
        failures.append(f"validate-feedback Step {step_number} changed bytes (requirement 3)")

# Invariant 2: the no-flag Output Format fence keeps its bytes (requirement 3, decision D-02).
output_heading = action_text.find("\n## Output Format\n")
fence_match = re.search(r"\n```markdown\n.*?\n```\n", action_text[output_heading:], re.S) if output_heading >= 0 else None
base_fence_hash = "27f9c8ad882b7bc0eecbaca7740f11df00757256d2880a7c18db4ffac4812882"
if fence_match is None or hashlib.sha256(fence_match.group(0).encode()).hexdigest() != base_fence_hash:
    failures.append("the Output Format fence under ## Output Format changed bytes (requirement 3, D-02)")

# Invariant 3: no live shipped file spells a retired core validate-feedback trigger
# (the same boundary rule _dev/tests/staged-skills-contract.sh applies, narrowed to this action).
retired_triggers = [
    "validate-feedback", "validate feedback", "triage findings", "triage feedback",
    "feedback review", "review feedback", "assess feedback", "should we push back",
]
retired_pattern = re.compile(
    r"(?<![A-Za-z0-9_-])do-work (" + "|".join(re.escape(t) for t in retired_triggers) + r")(?![A-Za-z0-9_'-])"
)
for live_path in sorted((root / "skills").rglob("*")):
    if not live_path.is_file() or live_path.name in {"CHANGELOG.md", "queue-kanban"}:
        continue
    for line_number, line in enumerate(live_path.read_text(errors="replace").splitlines(), 1):
        for retired_match in retired_pattern.finditer(line):
            failures.append(
                f"{live_path.relative_to(root)}:{line_number}: retired core trigger 'do-work {retired_match.group(1)}'"
            )

if probe_mode != "invariants":
    # Flags and input (requirements 1, 2).
    for required_text in ("--capture", "--run"):
        if required_text not in action_text:
            failures.append(f"validate-feedback does not mention {required_text}")
    # Wrong-repo question (requirement 4), exact wording from the REQ.
    if "These findings cite paths that are not in this repo. Continue capture here, or stop?" not in action_text:
        failures.append("validate-feedback lacks the wrong-repo question text (requirement 4)")
    # Discuss question options (requirement 5).
    for option_text in ("Accept: capture as a REQ", "Park:", "Drop: no work"):
        if option_text not in action_text:
            failures.append(f"validate-feedback lacks the Discuss option '{option_text}' (requirement 5)")
    # Capture, verify and run chain (requirements 6, 9, 10).
    for chain_text in ("actions/capture.md", "actions/verify-requests.md", "do-work verify-requests UR-NNN", "do-work run UR-NNN"):
        if chain_text not in action_text:
            failures.append(f"validate-feedback lacks '{chain_text}' (requirements 6, 9, 10)")
    # Every read-only statement names the --capture exception (requirement 3, plus the
    # description blockquote and the Red Flags line, which restate the same rule).
    read_only_statement = re.compile(
        r"read-only|no files were modified|create no reqs|capture the accepts to save|created or edited files",
        re.I,
    )
    for line_number, line in enumerate(action_lines, 1):
        if read_only_statement.search(line) and "--capture" not in line:
            failures.append(
                f"validate-feedback.md:{line_number}: read-only statement without the --capture exception"
            )
    # capture.md admits the hand-off (requirement 7, decision D-03).
    capture_text = (root / "skills/do-work/actions/capture.md").read_text()
    if "validate-feedback --capture" not in capture_text:
        failures.append("capture.md does not name validate-feedback --capture (requirement 7, D-03)")
    # Toolbox help names the flag (decision D-04).
    help_lines = (root / "skills/do-work-toolbox/actions/help.md").read_text().splitlines()
    if not any(line.lstrip().startswith("validate-feedback") and "--capture" in line for line in help_lines):
        failures.append("toolbox help.md validate-feedback line does not mention --capture (D-04)")
    # Core route (requirement 11, decision D-05).
    core_lines = (root / "skills/do-work/SKILL.md").read_text().splitlines()
    route_rows = [
        index for index, line in enumerate(core_lines)
        if line.startswith("| `validate-feedback`")
        and "`triage feedback`" in line
        and "`../do-work-toolbox/actions/validate-feedback.md`" in line
    ]
    verify_rows = [index for index, line in enumerate(core_lines) if line.startswith("| `verify`")]
    capture_rows = [index for index, line in enumerate(core_lines) if line.startswith("| `capture-request:`")]
    if len(route_rows) != 1:
        failures.append(f"core SKILL.md needs exactly one validate-feedback route row, found {len(route_rows)}")
    elif "review feedback" in core_lines[route_rows[0]]:
        failures.append("core validate-feedback row must leave 'review feedback' out (requirement 11)")
    elif not verify_rows or not capture_rows or not route_rows[0] < verify_rows[0] < capture_rows[0]:
        failures.append("core validate-feedback row must sit above the verify row and the capture fallback (D-05)")

if failures:
    print("\n".join(failures))
    sys.exit(1)
print(f"REQ-692 file checks passed ({probe_mode})")
PY

citation_output="$(bash "$root/_dev/tests/shipped-package-reference-contract.sh" 2>&1)" || {
  printf '%s\n' "$citation_output"
  echo "shipped-package-reference-contract.sh failed"
  exit 1
}
shell_block_output="$(bash "$root/_dev/tests/action-shell-blocks.sh" 2>&1)" || {
  printf '%s\n' "$shell_block_output"
  echo "action-shell-blocks.sh failed"
  exit 1
}
echo "REQ-692 probe passed ($probe_mode)"
