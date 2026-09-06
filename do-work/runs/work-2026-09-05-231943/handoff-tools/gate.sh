#!/usr/bin/env bash
# Cloud-container wrapper for the maintainer gate.
#
# Three pieces of harness-injected environment make the gate red for reasons that have
# nothing to do with the repository's own code:
#   NODE_OPTIONS                 - Claude Code node heap flag; heavy-runtime-fingerprint.py
#                                  refuses it as an "opaque runtime extension".
#   GIT_CONFIG_COUNT/KEY_*/VALUE_* - proxy-era github URL rewriting + credential.interactive=false;
#                                  the same fingerprint refuses them as "opaque Git configuration
#                                  override".
#   commit.gpgsign=true in the global config, with an EMPTY signing key - every `git commit`
#                                  inside a test fixture repository fails, so lane-mutation tests
#                                  see a dirty tree where they expect a new revision.
# The gate's own probes force GIT_CONFIG_GLOBAL=/dev/null, so pointing it at a sanitized
# config here does not change what the fingerprint seals.
#
# The checkout to gate is DO_WORK_GATE_ROOT when set, so a builder in a worktree gates its own
# branch rather than the main checkout. A REQ-557 builder ran this wrapper from its worktree,
# got a green verdict about the main tree, and had to notice that itself.
set -uo pipefail
cd "${DO_WORK_GATE_ROOT:-/home/user/skill-do-work}" || exit 1
exec env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/gitconfig-gate" \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  bash _dev/tests/maintainer-verify.sh "$@"
