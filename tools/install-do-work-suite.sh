#!/usr/bin/env bash
# Install and reconcile the complete four-skill do-work suite in one recoverable transaction.
#
# The transaction lives in the do-work-cli `install-suite` command; this file is the
# compatibility launcher that keeps the public argv working and answers
# --print-bootstrap-command from its own literal heredoc. Installing requires Go 1.26.1 or
# newer, which the do-work-cli launcher enforces.
set -euo pipefail

fail() {
  printf 'do-work suite install: %s\n' "$*" >&2
  exit 1
}

print_bootstrap_command() {
  cat <<'BOOTSTRAP'
(
  set -e
  project_root="$(git rev-parse --show-toplevel 2>/dev/null)" || { printf 'do-work bootstrap: run this from inside the target Git repository\n' >&2; exit 1; }
  bootstrap_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-suite-bootstrap.XXXXXX")"
  trap 'rm -rf "$bootstrap_tmp"' EXIT
  archive_file="$bootstrap_tmp/do-work-suite.tar.gz"
  curl -fsSL --retry 3 --retry-delay 2 --retry-max-time 60 -o "$archive_file.download" https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz
  mv "$archive_file.download" "$archive_file"
  mkdir -p "$bootstrap_tmp/source"
  tar xzf "$archive_file" -C "$bootstrap_tmp/source" --strip-components=1
  bash "$bootstrap_tmp/source/tools/install-do-work-suite.sh" --project-root "$project_root" --archive "$archive_file"
)
BOOTSTRAP
}

# The bootstrap snippet is static text that must print before anything is installed, so it is
# answered here rather than by the command, and needs no Go toolchain.
if [ "$#" -eq 1 ] && [ "$1" = '--print-bootstrap-command' ]; then
  print_bootstrap_command
  exit 0
fi

project_root=''
supplied_archive=''
while [ "$#" -gt 0 ]; do
  case "$1" in
    --project-root)
      [ "$#" -ge 2 ] || fail 'usage: install-do-work-suite.sh --project-root <Git-root> [--archive <source.tar.gz>]'
      project_root="$2"
      shift 2
      ;;
    --archive)
      [ "$#" -ge 2 ] || fail 'usage: install-do-work-suite.sh --project-root <Git-root> [--archive <source.tar.gz>]'
      supplied_archive="$2"
      shift 2
      ;;
    *)
      fail 'usage: install-do-work-suite.sh --project-root <Git-root> [--archive <source.tar.gz>]'
      ;;
  esac
done
[ -n "$project_root" ] || fail '--project-root is required'

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
do_work_cli=''
for cli_candidate in \
  "$script_dir/do-work-cli.sh" \
  "$script_dir/../skills/do-work/tools/do-work-cli.sh"; do
  if [ -f "$cli_candidate" ]; then
    do_work_cli="$cli_candidate"
    break
  fi
done
[ -n "$do_work_cli" ] || fail 'do-work-cli.sh is missing beside the installer'

# --project-root becomes the CLI's one repository-root concept, the global --repo-root.
cli_arguments=(--repo-root "$project_root" --format text install-suite)
if [ -n "$supplied_archive" ]; then
  cli_arguments+=(--archive "$supplied_archive")
fi

# Invoked rather than exec'd so the command's exit status passes through this launcher
# unchanged, including the status an interrupted install exits with.
launcher_status=0
bash "$do_work_cli" "${cli_arguments[@]}" || launcher_status=$?
exit "$launcher_status"
