#!/usr/bin/env bash
# GREEN probe for REQ-723: the archive-collision helper skips assets/ folders.
# Checks whichever tree it runs in. Expected to FAIL at base: the new test does not exist yet.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"

fail() { printf 'REQ-723 probe FAIL: %s\n' "$*" >&2; exit 1; }
test_names='TestArchiveCollisionSkipsAssetsFolders TestArchiveCollisionFindsNestedRequestsWithoutFollowingDirectoryLinks TestArchiveCollisionDistinguishesMissingArchiveFromUnreadableArchive'

go_output="$(go test -count=1 -v \
  -run '^(TestArchiveCollisionSkipsAssetsFolders|TestArchiveCollisionFindsNestedRequestsWithoutFollowingDirectoryLinks|TestArchiveCollisionDistinguishesMissingArchiveFromUnreadableArchive)$' \
  ./internal/corehelpers/ 2>&1)" || fail "archive-collision tests failed: $go_output"
for test_name in $test_names; do
  grep -qF -- "--- PASS: $test_name " <<<"$go_output" || fail "$test_name did not PASS: $go_output"
done
printf 'REQ-723 probe PASS\n'
