#!/usr/bin/env bash
# Fixture execution proofs for atomic-download.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# atomic-download: a failed transfer may write private bytes but never change the final target.
atomic_bin="$fixture_root/atomic-bin"
mkdir -p "$atomic_bin"
printf '%s\n' '#!/usr/bin/env bash' 'output_path=""' 'while [ "$#" -gt 0 ]; do case "$1" in -o) output_path="$2"; shift 2 ;; *) shift ;; esac; done' 'printf partial > "$output_path"' 'exit 22' > "$atomic_bin/curl"
chmod +x "$atomic_bin/curl"
printf 'stable' > "$fixture_root/atomic-target"
PATH="$atomic_bin:$PATH" "$core_scripts/atomic-download.sh" https://example.invalid/fail "$fixture_root/atomic-target" >/dev/null 2>&1 && fail_case 'atomic-download partial-publication case accepted a failed transfer'
[ "$(cat "$fixture_root/atomic-target")" = stable ] || fail_case 'atomic-download partial-publication case changed the final target'
leaked_private_paths="$(find "$fixture_root" -name 'atomic-target.download.*' -print -quit)" \
  || fail_case 'atomic-download partial-publication case could not search the fixture tree for private scratch'
[ -n "$leaked_private_paths" ] \
  && fail_case 'atomic-download partial-publication case leaked private scratch'

# atomic-download: a rate-limited host answers 429 once and then succeeds. The fake curl
# below models curl's own internal retry loop, so it survives that 429 only if the caller
# allowed a retry — which is the whole point of the flag set. It also records the
# Authorization header it was handed, so the opt-in credential path is observable.
atomic_retry_bin="$fixture_root/atomic-retry-bin"
mkdir -p "$atomic_retry_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'output_path=""' \
  'retry_limit=0' \
  'authorization_header=""' \
  'while [ "$#" -gt 0 ]; do' \
  '  case "$1" in' \
  '    -o) output_path="$2"; shift 2 ;;' \
  '    --retry) retry_limit="$2"; shift 2 ;;' \
  '    -H) authorization_header="$2"; shift 2 ;;' \
  '    *) shift ;;' \
  '  esac' \
  'done' \
  'printf "%s" "$authorization_header" > "$ATOMIC_HEADER_LOG"' \
  'transfer_attempt=0' \
  'while :; do' \
  '  transfer_attempt=$((transfer_attempt + 1))' \
  '  printf "%s" "$transfer_attempt" > "$ATOMIC_ATTEMPT_LOG"' \
  '  if [ "$transfer_attempt" -gt 1 ]; then' \
  '    printf complete-payload > "$output_path"' \
  '    exit 0' \
  '  fi' \
  '  if [ "$transfer_attempt" -gt "$retry_limit" ]; then' \
  '    exit 22' \
  '  fi' \
  'done' \
  > "$atomic_retry_bin/curl"
chmod +x "$atomic_retry_bin/curl"

printf 'stale before retry\n' > "$fixture_root/atomic-retry-target"
GH_TOKEN='' GITHUB_TOKEN='' PATH="$atomic_retry_bin:$PATH" \
  ATOMIC_ATTEMPT_LOG="$fixture_root/atomic-attempts" \
  ATOMIC_HEADER_LOG="$fixture_root/atomic-header" \
  "$core_scripts/atomic-download.sh" https://example.invalid/rate-limited "$fixture_root/atomic-retry-target" >/dev/null 2>&1 \
  || fail_case 'atomic-download retry case did not survive a transient 429'
[ "$(cat "$fixture_root/atomic-retry-target")" = complete-payload ] \
  || fail_case 'atomic-download retry case did not publish the successful attempt'
[ "$(cat "$fixture_root/atomic-attempts")" = 2 ] \
  || fail_case 'atomic-download retry case did not let curl retry the rate-limited transfer'
[ -z "$(cat "$fixture_root/atomic-header")" ] \
  || fail_case 'atomic-download retry case sent an Authorization header with no token configured'
leaked_private_paths="$(find "$fixture_root" -name 'atomic-retry-target.download.*' -print -quit)" \
  || fail_case 'atomic-download retry case could not search the fixture tree for private scratch'
[ -n "$leaked_private_paths" ] \
  && fail_case 'atomic-download retry case leaked private scratch'

# atomic-download: an opt-in token becomes a bearer credential; GH_TOKEN wins over GITHUB_TOKEN.
printf 'stale before credential\n' > "$fixture_root/atomic-credential-target"
GH_TOKEN=primary-token GITHUB_TOKEN=fallback-token PATH="$atomic_retry_bin:$PATH" \
  ATOMIC_ATTEMPT_LOG="$fixture_root/atomic-credential-attempts" \
  ATOMIC_HEADER_LOG="$fixture_root/atomic-credential-header" \
  "$core_scripts/atomic-download.sh" https://example.invalid/private "$fixture_root/atomic-credential-target" >/dev/null 2>&1 \
  || fail_case 'atomic-download credential case returned nonzero'
[ "$(cat "$fixture_root/atomic-credential-header")" = 'Authorization: Bearer primary-token' ] \
  || fail_case 'atomic-download credential case did not send GH_TOKEN as a bearer credential'
printf 'stale before fallback credential\n' > "$fixture_root/atomic-fallback-target"
GH_TOKEN='' GITHUB_TOKEN=fallback-token PATH="$atomic_retry_bin:$PATH" \
  ATOMIC_ATTEMPT_LOG="$fixture_root/atomic-fallback-attempts" \
  ATOMIC_HEADER_LOG="$fixture_root/atomic-fallback-header" \
  "$core_scripts/atomic-download.sh" https://example.invalid/private "$fixture_root/atomic-fallback-target" >/dev/null 2>&1 \
  || fail_case 'atomic-download fallback-credential case returned nonzero'
[ "$(cat "$fixture_root/atomic-fallback-header")" = 'Authorization: Bearer fallback-token' ] \
  || fail_case 'atomic-download fallback-credential case did not fall back to GITHUB_TOKEN'

# atomic-download: a target occupied by a DIRECTORY must fail closed. `mv` treats a
# directory operand as a container rather than a collision, so the download nests
# inside it and exits zero — and the caller reads that status as proof the file
# landed. The canonical statement of the rule is the shipped guide's
# "Verified exact publication" section.
atomic_success_bin="$fixture_root/atomic-success-bin"
mkdir -p "$atomic_success_bin"
printf '%s\n' \
  '#!/usr/bin/env bash' \
  'output_path=""' \
  'while [ "$#" -gt 0 ]; do case "$1" in -o) output_path="$2"; shift 2 ;; *) shift ;; esac; done' \
  'printf complete-payload > "$output_path"' \
  'exit 0' \
  > "$atomic_success_bin/curl"
chmod +x "$atomic_success_bin/curl"

atomic_occupied_target="$fixture_root/atomic-occupied-target"
mkdir -p "$atomic_occupied_target"
printf 'occupant\n' > "$atomic_occupied_target/pre-existing.txt"
PATH="$atomic_success_bin:$PATH" \
  "$core_scripts/atomic-download.sh" https://example.invalid/payload "$atomic_occupied_target" >/dev/null 2>&1 \
  && fail_case 'atomic-download occupied-target case reported success for a publication that nested'
[ -d "$atomic_occupied_target" ] \
  || fail_case 'atomic-download occupied-target case did not leave the occupying directory in place'
[ "$(cat "$atomic_occupied_target/pre-existing.txt")" = occupant ] \
  || fail_case 'atomic-download occupied-target case disturbed the occupying directory contents'
leaked_private_paths="$(find "$atomic_occupied_target" -name '*.download.*' -print -quit)" \
  || fail_case 'atomic-download occupied-target case could not search the occupying directory'
[ -n "$leaked_private_paths" ] \
  && fail_case 'atomic-download occupied-target case abandoned its private file inside the occupant'

prescribed_shell_finish
