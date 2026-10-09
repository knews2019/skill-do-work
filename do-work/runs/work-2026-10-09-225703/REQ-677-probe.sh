#!/bin/bash
# GREEN checks for REQ-677 (integrator test gate): the journey-qa action and its guide exist, it cites ui-review's detection step, the router, both help menus, the README and the staged-skills list name it, and the fast citation contract holds. Fails before the build.
set -euo pipefail
test -f skills/do-work-toolbox/actions/journey-qa.md
test -f skills/do-work-toolbox/docs/journey-qa-guide.md
grep -q 'docs/journey-qa-guide.md' skills/do-work-toolbox/actions/journey-qa.md
grep -q 'actions/ui-review.md' skills/do-work-toolbox/actions/journey-qa.md
for result_class in 'product defect' 'test defect' 'unresolved'; do grep -qi "$result_class" skills/do-work-toolbox/actions/journey-qa.md; done
argument_hint_line="$(sed -n '/^argument-hint:/p' skills/do-work-toolbox/SKILL.md)"
grep -q 'journey-qa' <<<"$argument_hint_line"
grep -Eq '^\|.*`journey-qa`.*\| `\./actions/journey-qa\.md` \|$' skills/do-work-toolbox/SKILL.md
grep -Eq '^  journey-qa ' skills/do-work-toolbox/actions/help.md
grep -q 'journey-qa' skills/do-work/actions/help.md
grep -q 'do-work-toolbox journey-qa' README.md
grep -Eq '^  journey-qa$' _dev/tests/staged-skills-contract.sh
bash _dev/tests/shipped-package-reference-contract.sh
