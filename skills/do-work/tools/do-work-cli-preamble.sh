#!/usr/bin/env bash
# shellcheck disable=SC2034 # do_work_cli and launcher_arguments are read by the sourcing launcher.
# Shared do-work-cli launcher preamble. Sourced by every compatibility launcher, never run.
#
# Resolution is anchored on this file's own directory, not the sourcing launcher's, so it
# finds the do-work-cli.sh shipped beside this preamble. do_work_cli stays empty when that
# file is missing, because each launcher owns a different missing-launcher message and exit
# status.

do_work_cli_preamble_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
do_work_cli=''
if [ -f "$do_work_cli_preamble_directory/do-work-cli.sh" ]; then
  do_work_cli="$do_work_cli_preamble_directory/do-work-cli.sh"
fi
unset do_work_cli_preamble_directory

# The argument prefix every launcher passes ahead of its own command and flags.
launcher_arguments=(--format text)
