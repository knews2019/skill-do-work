---
id: UR-029
title: Capture four audited safety fixes
created_at: 2026-08-07T08:45:11Z
requests: [REQ-129, REQ-130, REQ-131, REQ-132]
word_count: 6
---

# Capture Four Audited Safety Fixes

## Summary

Capture four independently testable safety defects confirmed by a deep audit of the current do-work skill: secret-derived Git copies can become readable additions, malformed request history can be reported as successfully restored while remaining damaged, the commit-hash idempotency path can approve unrelated body edits, and the Kanban restart recipes can start while a listener remains on the requested port.

## Extracted Requests

| ID | Request |
|---|---|
| REQ-129 | Preserve secret rename quarantine while making Git copy records secret-safe. |
| REQ-130 | Recover malformed request files only from valid historical frontmatter. |
| REQ-131 | Reject non-metadata changes on commit-hash idempotent reruns. |
| REQ-132 | Harden and synchronize Kanban port shutdown behavior. |

## Batch Constraints

- Preserve newer upstream improvements rather than copying the consumer implementation wholesale.
- Keep the existing `XD` source plus `X` destination behavior for secret-shaped renames, even though the original consumer regression expected only the destination row.
- Add runnable RED-first regressions for every defect; the existing Bash 3.2 suites currently pass despite all four reproduced failures.
- Use the repository's existing shell test framework. Run relevant regressions, `bash -n`, ShellCheck, Bash 3.2 compatibility checks, queue-kanban Go tests, and `go vet ./...` as applicable.
- Keep the root recipe, installation template, and installed/generated behavior synchronized. Preserve Claude-specific metadata such as `argument-hint`.
- Follow the repository's release rules if implementation reaches an integrating commit.
- Do not stage, commit, or push unless the user explicitly asks afterward.

## Full Verbatim Input

run do-work capture-request on these issues

## Referenced Conversation Context

"These issues" refers to the four-finding audit immediately preceding this capture. The audit reproduced each failure against the current working tree, distinguished the stronger existing rename contract from the still-unsafe copy path, and confirmed that both existing shell regression suites remain green because the required cases are absent.

---
*Captured: 2026-08-07T08:45:11Z*
