#!/bin/bash
# GREEN checks for REQ-676 (integrator test gate): the source-audit action and its guide exist, the router, both help menus, the README and the staged-skills list name it, and the fast citation contract holds. Fails before the build.
set -euo pipefail
test -f skills/do-work-toolbox/actions/source-audit.md
test -f skills/do-work-toolbox/docs/source-audit-guide.md
grep -q 'docs/source-audit-guide.md' skills/do-work-toolbox/actions/source-audit.md
sed -n '/^argument-hint:/p' skills/do-work-toolbox/SKILL.md | grep -q 'source-audit'
grep -Eq '^\|.*`source-audit`.*\| `\./actions/source-audit\.md` \|$' skills/do-work-toolbox/SKILL.md
grep -Eq '^  source-audit ' skills/do-work-toolbox/actions/help.md
grep -q 'source-audit' skills/do-work/actions/help.md
grep -q 'do-work-toolbox source-audit' README.md
grep -Eq '^  source-audit$' _dev/tests/staged-skills-contract.sh
bash _dev/tests/shipped-package-reference-contract.sh
