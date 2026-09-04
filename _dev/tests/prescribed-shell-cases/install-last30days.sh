#!/usr/bin/env bash
# Fixture execution proofs for install-last30days.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

fixture_repo_init "$fixture_root"
export DO_WORK_COMPATIBILITY_REPO_ROOT="$fixture_root"
toolbox_scripts="$fixture_root/toolbox-scripts"
mkdir -p "$toolbox_scripts"
cat > "$toolbox_scripts/install-last30days.sh" <<EOF
#!/usr/bin/env bash
set -euo pipefail
launcher_arguments=(--format text)
if [[ -n "\${DO_WORK_COMPATIBILITY_REPO_ROOT:-}" ]]; then
  launcher_arguments+=(--repo-root "\$DO_WORK_COMPATIBILITY_REPO_ROOT")
fi
exec bash "$repo_root/skills/do-work/tools/do-work-cli.sh" "\${launcher_arguments[@]}" install-last30days "\$@"
EOF
chmod +x "$toolbox_scripts/install-last30days.sh"

# install-last30days: a SKILL.md-only tree fails check, is repaired from a
# complete fixture, and receives the full subtree plus ignore/Python guarantees.
upstream_repo="$fixture_root/last30days-upstream"
fixture_repo_init "$upstream_repo"
mkdir -p "$upstream_repo/skills/last30days/scripts" "$upstream_repo/skills/last30days/support"
printf '# Last30Days\n' > "$upstream_repo/skills/last30days/SKILL.md"
printf 'runtime\n' > "$upstream_repo/skills/last30days/scripts/last30days.py"
printf 'support\n' > "$upstream_repo/skills/last30days/support/data.txt"
fixture_repo_commit_all "$upstream_repo" fixture
last_project="$fixture_root/last-project"
fixture_repo_init "$last_project"
mkdir -p "$last_project/.claude/skills/last30days"
printf '# Sentinel only\n' > "$last_project/.claude/skills/last30days/SKILL.md"
python_bin="$fixture_root/python-bin"
mkdir -p "$python_bin"
printf '%s\n' '#!/usr/bin/env bash' 'exit 0' > "$python_bin/python3.12"
chmod +x "$python_bin/python3.12"
PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" check "$last_project" >/dev/null 2>&1 \
  && fail_case 'install-last30days sentinel-only check accepted a missing runtime script'
last_output="$(TMPDIR="$fixture_root" PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$last_project" "$upstream_repo" 2>&1)" \
  || fail_case "install-last30days complete-source repair returned nonzero: $last_output"
[ "$(cat "$last_project/.claude/skills/last30days/scripts/last30days.py")" = runtime ] \
  && [ "$(cat "$last_project/.claude/skills/last30days/support/data.txt")" = support ] \
  || fail_case 'install-last30days complete-source repair omitted part of the source subtree'
git -C "$last_project" check-ignore -q .claude/skills/last30days/SKILL.md || fail_case 'install-last30days fixture-source case did not resolve the sibling exclude helper'
PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" check "$last_project" >/dev/null 2>&1 \
  || fail_case 'install-last30days repaired tree failed the complete runtime/ignore/Python check'

# install-last30days: reject an incomplete source before publishing any final tree.
incomplete_upstream_repo="$fixture_root/last30days-incomplete-upstream"
fixture_repo_init "$incomplete_upstream_repo"
mkdir -p "$incomplete_upstream_repo/skills/last30days"
printf '# Incomplete\n' > "$incomplete_upstream_repo/skills/last30days/SKILL.md"
fixture_repo_commit_all "$incomplete_upstream_repo" fixture
incomplete_source_project="$fixture_root/last-incomplete-source-project"
mkdir -p "$incomplete_source_project"
TMPDIR="$fixture_root" PATH="$python_bin:$PATH" "$toolbox_scripts/install-last30days.sh" install "$incomplete_source_project" "$incomplete_upstream_repo" >/dev/null 2>&1 \
  && fail_case 'install-last30days incomplete-source case returned success'
[ ! -e "$incomplete_source_project/.claude/skills/last30days" ] \
  || fail_case 'install-last30days incomplete-source case published a final tree'

# Command-interposition used to exercise shell-internal cp/mv branches. The Go
# publisher now owns these faults directly in last30days_test.go, including copy
# rollback, interruption after backup, and a second-writer publication collision.

prescribed_shell_finish
