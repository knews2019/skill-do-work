#!/bin/bash
# GREEN checks for REQ-678 (integrator test gate): the release-check action and its guide exist, the action links the core Step 7 stage definitions and names the unassessed and historical labels, the router, both help menus, the README and the staged-skills list name it, and the fast citation contract holds. Fails before the build.
set -euo pipefail
test -f skills/do-work-toolbox/actions/release-check.md
test -f skills/do-work-toolbox/docs/release-check-guide.md
grep -q 'docs/release-check-guide.md' skills/do-work-toolbox/actions/release-check.md
grep -qF '../../do-work/actions/review-work.md' skills/do-work-toolbox/actions/release-check.md
for required_label in 'unassessed' 'historical' 'not ready'; do grep -qi "$required_label" skills/do-work-toolbox/actions/release-check.md; done
argument_hint_line="$(sed -n '/^argument-hint:/p' skills/do-work-toolbox/SKILL.md)"
grep -q 'release-check' <<<"$argument_hint_line"
grep -Eq '^\|.*`release-check`.*\| `\./actions/release-check\.md` \|$' skills/do-work-toolbox/SKILL.md
grep -Eq '^  release-check ' skills/do-work-toolbox/actions/help.md
grep -q 'release-check' skills/do-work/actions/help.md
grep -q 'do-work-toolbox release-check' README.md
grep -Eq '^  release-check$' _dev/tests/staged-skills-contract.sh
bash _dev/tests/shipped-package-reference-contract.sh
