---
source_type: req_lesson
req_id: REQ-093
req_path: do-work/archive/UR-015/REQ-093-shipped-go-source-cites-the-export-ignored-maintainer-doc.md
date: 2026-08-04
domain: general
module: tools/queue-kanban
tags: [queue-kanban, go-source, export-ignored, shipped, sites]
---

# Lessons from REQ-093: Confirm: six shipped Go-source sites cite the export-ignored CLAUDE.md, and the suite's guard catches none of them

## What the REQ was about

REQ-088 fixed one line of `actions/memory-reference.md` that pointed a consumer at this repo's
`CLAUDE.md` — a file that is `export-ignore`d, so it never lands in a consumer install and the
pointer dangles. While verifying that REQ-088's one-site inventory was complete, a sweep of every
shipped path found **six more sites of the same defect**, all in shipped Go source:

| File | Line | The citation |
|---|---|---|
| `tools/queue-kanban/verify.go` | 123 | `covers CLAUDE.md § Before Every Commit items 1 and 2` |
| `tools/queue-kanban/verify.go` | 156 | `(CLAUDE.md § Before Every Commit item 2)` |
| `tools/queue-kanban/verify.go` | 186 | `the duplicate version numbers CLAUDE.md records as having already happened` |
| `tools/queue-kanban/verify_test.go` | 89 | `each direction has its own cause — CLAUDE.md § Before Every Commit` |
| `tools/queue-kanban/verify_test.go` | 141 | `CLAUDE.md records as having already happened` |
| `tools/queue-kanban/verify_test.go` | 171 | `CLAUDE.md § Before Every Commit, item 2` |

## Solution summary

Rewrote six comments in `verify.go` (123, 156, 186) and `verify_test.go` (89, 141, 171) to state the release invariant they describe instead of citing `CLAUDE.md § Before Every Commit`; neither file now mentions the maintainer doc. Dropped the unsourceable "CLAUDE.md standard:" attribution from `prompt-kit-step6-constraint-architecture.md:78`, keeping the rule. Replaced the suite's idiom-matching `self_citation_pattern` with an inverted check: it greps every `CLAUDE.md`/`AGENTS.md` mention across the shipped paths and filters out a 14-entry per-file `maintainer_doc_mention_allowlist`, keeping the original check's stderr/`fail_count` contract and adding a `FAIL:` message that names the allowlist as the third remedy. Updated `CLAUDE.md:122`, which described the guard as grepping "the common citation idioms" — stale the moment the guard stopped doing that.

## What worked

- Scoring a proposed guard against the real occurrences before shipping it, instead of reasoning about which idioms "should" appear. The measurement (old pattern 0/8, proposed widening 4/6 on the Go sites, inverted check 8/8) settled the design argument in one command, and the bare-prose negative probe demonstrated the failure the idiom approach would have shipped.

## What didn't work

- `git add` with a list containing one stale pathspec — the REQ-090 queue path that `git mv` had already moved. The whole invocation aborted, so the commit captured only the rename and left both the cancellation body and REQ-093's answers unstaged. Both halves of this trap are written down in `do-work/HANDDOWN-UR-015-016.md` (`git mv` stages content at move time; `git add` aborts on one bad pathspec) and it was still walked into, because the two combine into a single silent failure: the commit *succeeded*, reporting nothing wrong. `git status --porcelain -uall` after staging and before committing is the check that catches it — reading the commit's own `--name-status` afterwards is what caught it here, one step too late.

## Worth knowing

- `_dev/` is export-ignored, so the contract suite may cite `CLAUDE.md` in its own comments — the inverted check must never be pointed at `_dev/`, or it flags its own explanation. Relatedly, the per-file allowlist shape was chosen over per-directory for a concrete reason worth keeping: allowlisting `actions/` would have exempted `actions/memory-reference.md`, the file the whole thread started from.

## Back-reference

See `do-work/archive/UR-015/REQ-093-shipped-go-source-cites-the-export-ignored-maintainer-doc.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `2e29b36`.
