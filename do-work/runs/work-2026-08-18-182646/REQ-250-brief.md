# REQ-250 builder brief

**Route B.** Estimated 25 active minutes (P50, medium confidence).

Your file, `_dev/tests/shipped-package-reference-contract.sh`, was just extended by REQ-249 (backticked-citation checker + a parser refactor extracting `mask_block_code`/`inline_code_regions` from `strip_markdown_code`). Read the current state before planning — your four instances concern the *Markdown-link* half, but the parser is now shared. Do not rework REQ-249's checker.

REQ-249's review also noted (Minor, report-only, relevant to your judgment): anchor fragments on *backticked* citations are discarded, never validated. That is out of your scope unless closing your instance 1 (same-file `#anchor` links) makes it nearly free — if so, note it as a Discovered Task rather than widening.

## How this build runs

You are a **worktree builder** dispatched by the do-work work pipeline. Everything binding is in this brief.

**Your tree, your branch.** Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-250-close-the-remaining-markdown-link-checker-gaps` — a full checkout on branch `worktree-agent-REQ-250-close-the-remaining-markdown-link-checker-gaps`, cut from integration tip `ad69e56`.

- Never write anything under `/home/user/skill-do-work` — the one exception is your hand-back file, named below.
- Never read or write `do-work/` in your own worktree (stale snapshot; your REQ body is inlined below).
- Commit on your own branch, in small increments so an interruption costs one step. Do not touch `VERSION`, `CHANGELOG.md`, or `skills/do-work/actions/version.md` — serial-only, integrator-owned.
- A needed one-line edit to a file outside your write set is an *integration seam*: hand back the exact line and where it goes; do not edit the file. A larger need: stop and report in your hand-back.
- Out-of-scope finds go in `## Discovered Tasks` in your hand-back — never fixed inline.

**Crew rules** (read from your own worktree before writing code): `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`. This REQ is `tdd: true` and `maintenance: true` — also read `crew-members/testing.md` and `crew-members/maintenance.md`; delete-before-you-add applies: two of the four instances may close better by deleting a claim than by adding code. Read every `prime_files` path too.

**P-A-U phasing is mandatory** — work the [PLAN]/[APPLY]/[UNIFY] block in your REQ body; record the evidence in your hand-back (the orchestrator transcribes it and audits it against the diff). Log significant choices as D-XX with reasoning (DECIDE & STATE vs ESCALATE with Value/Risk).

## Environment notes

- `bash _dev/tests/maintainer-verify.sh` exits 0 at your branch point — your baseline and your gate. Exit code is the only proof; never pipe it through `tail`.
- Toolchain present: Go 1.26.1, ShellCheck 0.11.0, `just`, Node 22, Chromium (Playwright, `/opt/pw-browsers/chromium`).
- **Never run bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — build to scratch (`go build -o /tmp/<name> .`).
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp; never carry or compute one.
- Build test fixtures in scratch space, never in this repo's own `do-work/`.

## Hand-back

Write your report to exactly this absolute path (the one main-tree write you may make; never stage or commit it):

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-182646/REQ-250-handback.md
```

Structure: `# REQ-NNN hand-back` with **Branch**, **Commits** (oldest first), then `## What I built`, `## File manifest` (one full path per line, `(new|modified|deleted)` + one factual line), `## P-A-U evidence`, `## Testing evidence` (real RED and GREEN output — never from a prototype or memory; the observed maintainer-verify exit code), `## Decisions (D-XX)`, `## Integration seams` (exact lines or "none"), `## Discovered Tasks`, `## Pushback`.

**Standing warning, now seven-for-seven in this repo:** every recent REQ shipped a mechanism that looked like it closed a class and closed only the instance — reviews keep finding the hole exactly where the real data lives. Assume your first fix has that shape and hunt the hole before the reviewer does.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-250
title: Close the remaining markdown link checker gaps
status: claimed
created_at: 2026-08-18T13:55:32Z
claimed_at: 2026-08-18T18:25:40Z
route: B
status_changed_at: 2026-08-18T13:55:32Z
user_request: UR-042
addendum_to: REQ-243
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-link-checker-unresolved-classes
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- _dev/tests/shipped-package-reference-contract.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T18:26:21Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Close the Remaining Markdown Link Checker Gaps

## What

REQ-243's anchor checker has four known limits. Three are false negatives — links it accepts that are broken — and one is a path-escape the same change escalated from a `stat` to a `read_text`.

## Context

All four came from REQ-243's independent review, which cleared the checker's core (skip rule condition-keyed, slug rule correct for every real anchor in this repo, corpus failure path non-vacuous) and named these as the remaining edges. They were routed here rather than folded in, because REQ-243's remediation was already carrying two more urgent findings.

## Instances

- [ ] **Same-file bare `#anchor` links are not validated at all.** `[see below](#renamed-section)` is silently broken today. The slugs are already cached by REQ-243's work, so this is roughly three lines.
- [ ] **`os.path.normpath` does not clamp `..` at the repo root, and this path is now read rather than stat'ed.** A link like `../../../../elsewhere.md#a` that happens to exist was previously only tested for existence; the anchor check now calls `read_text()` on it. Pre-existing in origin, escalated in effect. A `relative_to(repo_root)` guard closes it.
- [ ] **HTML tags and entities in heading text diverge from GitHub's slug, in the false-negative direction.** `## <kbd>Ctrl</kbd> and stuff` slugs to `kbdctrlkbd-and-stuff` here and `ctrl-and-stuff` on GitHub; `## Tom &amp; Jerry` gives `tom-amp-jerry` against `tom--jerry`. A link written to the checker's slug passes the suite and is broken in every renderer. No shipped heading trips this today — all 27 with suspicious characters slug correctly — so **check whether this is worth closing at all before closing it**; a stated limitation may be the better answer.
- [ ] **Blockquoted ATX headings are dropped.** `> # Quoted Heading` yields no anchor because the gate tests the code-masked line, which begins with `>`. GitHub does generate one. Fails loudly (spurious FAIL), so low severity — but the limitation is currently unstated.

## Requirements

- Each instance is either closed, or its limitation is stated in the file with the failure direction named. A silently-accepting limit and a loudly-failing one are different risks and the file should say which each is.
- `bash _dev/tests/maintainer-verify.sh` still exits 0, and the existing 27 corpus anchors still resolve.
- **`maintenance: true`: ask what can be removed before adding.** Two of these four may be better closed by deleting a claim than by adding code.

## Red-Green Proof

**RED prompt/case:** per instance — a same-file link to a non-existent anchor; a `..`-escaping relative target; a heading containing an HTML tag with a link written to GitHub's slug; a blockquoted heading with a link to it.
**Why RED now:** the first three pass the suite today and the fourth fails spuriously.
**GREEN when:** each case either fails for the right reason or is documented as an accepted limit, and the suite still exits 0.
**Validation:** Review findings on REQ-243; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** Four named edges of one checker, each with its instance stated and the fix sketched; the work is judgment about close-vs-document per instance, inside a single test file.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — the four checker edges closed or documented, with fixtures

**Files I will NOT touch:**
- Shipped markdown — this REQ changes the checker, not the corpus.
- The backticked-citation checker REQ-249 just added to the same file — coordinate around it, don't rework it.

**Acceptance criteria (restated from REQ):**
- [ ] Each instance is either closed, or its limitation is stated in the file with the failure direction named.
- [ ] `bash _dev/tests/maintainer-verify.sh` still exits 0 and the existing 27 corpus anchors still resolve.
- [ ] maintenance: true — ask what can be removed before adding; two of the four may close better by deleting a claim.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (wave-1 integration tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium all present

*Checked by work action*
