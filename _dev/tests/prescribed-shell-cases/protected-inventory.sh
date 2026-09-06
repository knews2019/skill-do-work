#!/usr/bin/env bash
# Fixture execution proofs for protected-inventory.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# protected-inventory: once a secret path is quarantined it stays out of association.
inventory_repo="$fixture_root/inventory-repo"
fixture_repo_init "$inventory_repo"
mkdir -p "$inventory_repo/do-work/archive/UR-001"
printf '%s\n' '---' 'id: REQ-001' 'status: completed' '---' '' '## Implementation Summary' '' '- `safe.txt` — fixture' > "$inventory_repo/do-work/archive/UR-001/REQ-001-fixture.md"
printf 'base\n' > "$inventory_repo/safe.txt"
fixture_repo_commit_all "$inventory_repo" base
printf 'change\n' >> "$inventory_repo/safe.txt"
printf 'secret\n' > "$inventory_repo/.env.local"
inventory_output="$(cd "$inventory_repo" && "$core_scripts/protected-inventory.sh" start)" || fail_case 'protected-inventory start case returned nonzero'
grep -q $'X\t.env.local' <<<"$inventory_output" || fail_case 'protected-inventory start case did not quarantine the secret path'
association_output="$(cd "$inventory_repo" && "$core_scripts/protected-inventory.sh" associate)" || fail_case 'protected-inventory associate case returned nonzero'
grep -q $'REQ-001\tsafe.txt' <<<"$association_output" || fail_case 'protected-inventory associate case lost the safe owner'
grep -q '.env.local' <<<"$association_output" && fail_case 'protected-inventory associate case leaked the quarantined path'

# protected-inventory: --repo-root reaches the CLI from outside the repository, in both
# spellings and on either side of the mode token. This was the defect REQ-603 named first,
# and the launcher's own `set -u` crash on an empty argument array (bash 3.2) hid behind
# the same missing case: every call above passes no global flag at all.
outside_output="$(cd "$fixture_root" && "$core_scripts/protected-inventory.sh" --repo-root "$inventory_repo" associate)" \
  || fail_case 'protected-inventory --repo-root case returned nonzero from outside the repository'
[ "$outside_output" = "$association_output" ] \
  || fail_case 'protected-inventory --repo-root case did not produce the in-repository association'
outside_output="$(cd "$fixture_root" && "$core_scripts/protected-inventory.sh" associate "--repo-root=$inventory_repo")" \
  || fail_case 'protected-inventory --repo-root= case returned nonzero with the flag after the mode'
[ "$outside_output" = "$association_output" ] \
  || fail_case 'protected-inventory --repo-root= case did not produce the in-repository association'
(cd "$fixture_root" && "$core_scripts/protected-inventory.sh" --repo-root 2>"$fixture_root/inventory-usage-err") \
  && fail_case 'protected-inventory missing-value case did not exit nonzero'
[ -s "$fixture_root/inventory-usage-err" ] \
  || fail_case 'protected-inventory missing-value case exited silently'

prescribed_shell_finish
