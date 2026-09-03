#!/usr/bin/env bash
# Earned runtime incidents for the shipped core checks module.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
core_root="$repo_root/skills/do-work"
fail_count=0
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
# The quarantine merge moved into the typed protected-inventory command. Keep
# the behavior probes below; the retained shell file is launcher-only.

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
[ "$fail_count" -eq 0 ] || exit 1
printf 'core-checks contract probes passed.\n'
