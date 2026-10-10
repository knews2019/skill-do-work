#!/usr/bin/env bash
# GREEN checks for REQ-686 (integrator test gate): the lessons file says a stuck finalization journal is deleted by hand after a verified git revert, heavy_commands.go carries the one comment naming maintainer-verify.sh, the package is gofmt-clean and builds, and the fast citation contract holds. Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
lessons_text="$(<"$module_dir/lessons-do-work-cli.md")"
grep -q 'git revert.*by hand\|by hand.*git revert' <<<"$lessons_text" || { echo "lessons-do-work-cli.md has no line tying a performed git revert to deleting the journal by hand"; exit 1; }
heavy_source="$(<"$module_dir/internal/heavyverification/heavy_commands.go")"
grep -q 'maintainer-verify.sh' <<<"$heavy_source" || { echo "heavy_commands.go lacks the comment naming maintainer-verify.sh"; exit 1; }
unformatted="$(gofmt -l "$module_dir/internal/heavyverification")"
[ -z "$unformatted" ] || { printf 'gofmt -l lists: %s\n' "$unformatted"; exit 1; }
go build -C "$module_dir" ./internal/heavyverification/
bash "$root/_dev/tests/shipped-package-reference-contract.sh"
