---
id: REQ-150
title: "Review fix: Repair live modular package references"
status: pending
domain: general
created_at: 2026-08-08T15:38:44Z
user_request: UR-031
addendum_to: REQ-144
review_generated: true
effort_estimate: normal
---

# Review Fix: Repair Live Modular Package References

## What
Make every URL and relative link published by the live modular packages resolve after installation and from the source archive. Close the entire shipped-reference class so a later package move cannot leave silent dangling guidance.

## Context
Found during review of REQ-144. The core version action still targets the deleted root action, the updater prime's archived-lesson links are file-relative from the wrong depth, and the installed core changelog links sidecars the package does not ship.

This is a standalone user-visible documentation and updater-reference defect rather than part of a broader sweep: one package-reference audit and executable guard closes the class.

## Requirements
- Point the live version action at the canonical modular upstream path.
- Repair or deliberately replace every broken relative lesson/history link in the installed core package.
- Decide and document how installed changelog history reaches repository-only sidecars without publishing dead links.
- Add a shipped-package reference check that resolves Markdown paths and validates raw upstream URLs against the live archive layout.
