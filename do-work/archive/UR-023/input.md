---
id: UR-023
title: Codex review of PR #130 — status normalization suppressed its own warning
created_at: 2026-08-05T19:41:00Z
requests: [REQ-115]
word_count: 0
---

# Codex Review of PR #130

## Full Verbatim Input

Two P2 findings from `chatgpt-codex-connector[bot]` on PR #130, reviewed commit `9f06e6c`:

1. **`tools/queue-kanban/frontmatter_cli.go` L185 — "Preserve warnings for status normalization."** When a caller runs `queue-kanban frontmatter get <file> status --normalize` (or `testing_status`) on a typo like `completedd`, this branch overwrites `recognized` to true, so the command exits 0, prints the typo to stdout, and emits no Schema Read Contract warning. The new CLI is intended to replace hand-rolled status reads; silently treating invalid status/testing_status values as recognized preserves the no-feedback path the command is meant to avoid.

2. **`do-work/archive/REQ-111-...md` L8 — "Point archived commit fields at reachable commits."** `e77383a` is not reachable from `8fb7a72` (same for `a6560bc` in REQ-112 and `66a7fcc` in REQ-110); a fresh checkout of the reviewed history will not have those objects, so following the archive record leads to an unknown revision.

## Disposition

Finding 1: **accepted** — reproduced, fixed by REQ-115.

Finding 2: **pushed back.** Verified against this repo's actual merge strategy rather than the synthetic PR ref Codex reviewed. All three hashes are reachable from the branch HEAD, and decisively REQ-109's archived hash `5f50fb7` — from an already-merged earlier PR — is reachable from `origin/main`. Main absorbs PRs as merge commits (`Merge pull request #131 from …`), not squashes, so branch commits survive into main and the archived `commit:` fields resolve. The finding would be correct in a squash-merge repo; it is not correct here. Rewriting those fields to some other commit would also violate the skill's own contract, which requires the *implementation* commit hash and validates it through `tools/checks/record-commit-hash.sh`.
