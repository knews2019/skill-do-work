#!/usr/bin/env bash
# GREEN checks for REQ-685 (integrator test gate): the no-root rollback path is gone from source, tests and the lock-in script, and the gittransaction package builds, vets, is gofmt-clean and passes. (The finalization package takes about 50 s and runs in the repository gate and the builder's own run.) Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
leftover="$(grep -rIl -e rollbackWithoutRoot -e openRollbackRoot "$module_dir/internal" "$root/_dev/tests" || true)"
[ -z "$leftover" ] || { printf 'deleted rollback names still referenced in:\n%s\n' "$leftover"; exit 1; }
unformatted="$(gofmt -l "$module_dir/internal/gittransaction")"
[ -z "$unformatted" ] || { printf 'gofmt -l lists: %s\n' "$unformatted"; exit 1; }
go vet -C "$module_dir" ./internal/gittransaction/
go test -C "$module_dir" -count=1 ./internal/gittransaction/
