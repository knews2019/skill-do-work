#!/usr/bin/env bash
# do-work-cli compatibility launcher: protected-inventory
set -euo pipefail

script_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
global_arguments=()
command_arguments=()
while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --repo-root)
      if [[ "$#" -lt 2 || -z "$2" ]]; then
        printf 'protected-inventory: --repo-root needs a value\n' >&2
        exit 2
      fi
      global_arguments+=(--repo-root "$2")
      shift 2
      ;;
    --repo-root=*)
      global_arguments+=("$1")
      shift
      ;;
    *)
      command_arguments+=("$1")
      shift
      ;;
  esac
done

# Only --repo-root is forwarded: the text format is this launcher's contract (its callers
# grep the tab-separated lines), so a caller's --format is not sifted and reaches the
# subcommand as an unknown option, which the CLI refuses with a message rather than
# silently overriding it. The `${array[@]+...}` form is what keeps an empty array from
# being an unbound variable under `set -u` on bash 3.2.
DO_WORK_COMPATIBILITY_SHIM=1 exec bash "$script_directory/../tools/do-work-cli.sh" ${global_arguments[@]+"${global_arguments[@]}"} --format text protected-inventory ${command_arguments[@]+"${command_arguments[@]}"}

