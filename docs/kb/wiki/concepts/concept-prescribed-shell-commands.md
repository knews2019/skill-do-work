---
title: "Shell and Automation"
type: concept
topic_cluster: shell-and-automation
sources:
  - raw/processed/2026-09-04/REQ-258-split-the-prescribed-shell-behavior-suit.md
  - raw/processed/2026-09-04/REQ-460-make-outside-text-delimiter-containment.md
  - raw/processed/2026-09-04/REQ-543-reap-the-commit-hook-with-its-own-parent.md
  - raw/processed/2026-09-01/REQ-064-restore-blanked-archived-reqs-from-git-h.md
  - raw/processed/2026-09-01/REQ-066-clear-two-shellcheck-warnings-in-the-com.md
  - raw/processed/2026-09-01/REQ-072-go-utility-allocates-req-ids-and-version.md
  - raw/processed/2026-09-01/REQ-114-the-three-remaining-shell-logic-extracti.md
  - raw/processed/2026-09-01/REQ-152-review-fix-reject-reserved-just-recipe-c.md
  - raw/processed/2026-09-01/REQ-156-review-fix-handle-just-multiline-strings.md
  - raw/processed/2026-09-01/REQ-159-review-fix-complete-multiline-literal-st.md
  - raw/processed/2026-09-01/REQ-165-shell-block-lint-harness-for-shipped-act.md
  - raw/processed/2026-09-01/REQ-167-deduplicate-copy-pasted-shell-primitives.md
  - raw/processed/2026-09-01/REQ-171-addendum-promote-prescribed-shell-primit.md
  - raw/processed/2026-09-01/REQ-172-make-screenshot-source-cleanup-best-effo.md
  - raw/processed/2026-09-01/REQ-173-handle-first-line-bom-in-just-collision-.md
  - raw/processed/2026-09-01/REQ-187-no-single-local-maintainer-command-prove.md
  - raw/processed/2026-09-01/REQ-298-review-fix-sweep-the-unchecked-exit-stat.md
  - raw/processed/2026-09-01/REQ-409-implement-safe-cleanup-passes-and-explic.md
  - raw/processed/2026-09-01/REQ-413-implement-capture-file-answer-release-ve.md
  - raw/processed/2026-09-01/REQ-414-migrate-remaining-core-checks-publicatio.md
  - raw/processed/2026-09-01/REQ-419-add-flat-just-recipes-collision-validati.md
related:
  - page: entity-do-work-cli
    rel: extends
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Shell and Automation

Architectural overview and synthesis for the Shell and Automation subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 21 source documents:

- [[REQ-064-restore-blanked-archived-reqs-from-git-h]] — Restore blanked archived REQs from git history in cleanup
- [[REQ-066-clear-two-shellcheck-warnings-in-the-com]] — Clear two shellcheck warnings in the commit-hash guard fixture
- [[REQ-072-go-utility-allocates-req-ids-and-version]] — Go utility allocates REQ ids and version numbers and verifies release consistency
- [[REQ-114-the-three-remaining-shell-logic-extracti]] — The three remaining shell-logic extraction candidates, restated decay-free
- [[REQ-152-review-fix-reject-reserved-just-recipe-c]] — Review fix: Reject reserved Just recipe collisions without Just
- [[REQ-156-review-fix-handle-just-multiline-strings]] — Review fix: Handle Just multiline strings in collision scanning
- [[REQ-159-review-fix-complete-multiline-literal-st]] — Review fix: Complete multiline literal state in Just collision scanning
- [[REQ-165-shell-block-lint-harness-for-shipped-act]] — Shell-block lint harness for shipped action files
- [[REQ-167-deduplicate-copy-pasted-shell-primitives]] — Deduplicate copy-pasted shell primitives across action files
- [[REQ-171-addendum-promote-prescribed-shell-primit]] — Addendum: promote prescribed shell primitives to shipped, fixture-tested scripts
- [[REQ-172-make-screenshot-source-cleanup-best-effo]] — Make screenshot source cleanup best-effort
- [[REQ-173-handle-first-line-bom-in-just-collision-]] — Handle first-line BOM in Just collision scan
- [[REQ-187-no-single-local-maintainer-command-prove]] — No single local maintainer command proves shell plus both Go modules
- [[REQ-298-review-fix-sweep-the-unchecked-exit-stat]] — Review fix: sweep the unchecked-exit-status primitive across every shipped script
- [[REQ-409-implement-safe-cleanup-passes-and-explic]] — Implement safe cleanup passes and explicit destructive repairs
- [[REQ-413-implement-capture-file-answer-release-ve]] — Implement capture-file, answer, release, version, and changelog transactions
- [[REQ-414-migrate-remaining-core-checks-publicatio]] — Migrate remaining core checks, publication helpers, Git helpers, and surveys
- [[REQ-419-add-flat-just-recipes-collision-validati]] — Add flat Just recipes, collision validation, action delegation, and compatibility aliases
- [[REQ-258-split-the-prescribed-shell-behavior-suit]] — Split the prescribed shell behavior suite per script
- [[REQ-460-make-outside-text-delimiter-containment]] — Make outside-text delimiter containment condition-complete
- [[REQ-543-reap-the-commit-hook-with-its-own-parent]] — Reap the commit hook with its own parent

## Cross-References

See related system components and verification gates.

## Related Entities

- [[entity-do-work-cli]]
