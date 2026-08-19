# REQ-249 builder brief

**Route B.** Estimated 45 active minutes (P50, medium confidence).

**The decision is already made and is not yours to revisit:** the user confirmed the literal form. Every backticked cross-package citation must resolve as a real relative path from the citing file's own directory. Your job is the sweep and the rule statement, not the choice.

**The extent is larger than the REQ's capture counts.** Measured in the main tree at claim time: 126 backticked `../do-work*/...` citations and 14 `../../do-work*/...` ones across the four packages. Depth is per-file, not a blanket rewrite — from `skills/do-work/actions/` the correct form is `../../do-work-board/...`, from `skills/do-work-board/tools/queue-kanban/` it is `../../../do-work/...`, and a package-root `SKILL.md` citing `../do-work-board/...` is **already literal-correct**. Derive each one; do not sed the tree.

The REQ's third requirement asks you to answer whether these pointers should now become mechanically checkable, and to implement the check if the answer is yes. `_dev/tests/shipped-package-reference-contract.sh` already resolves markdown links and anchors across packages and is the natural home.

## How this build runs

You are a **worktree builder** dispatched by the do-work work pipeline (`skills/do-work/actions/work.md`). Read that file's Step 6 expectations if you need them; everything binding is repeated here.

**Your tree, your branch.** Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-249-decide-the-cross-package-citation-path-form`. That is a full checkout of this repository on branch `worktree-agent-REQ-249-decide-the-cross-package-citation-path-form`, cut from the integration tip `67dae6b`.

- Never write anything under `/home/user/skill-do-work` — the main tree belongs to the orchestrator. The one exception is your hand-back file, named below.
- Never read or write anything under `do-work/` in your own worktree. Your worktree carries a **stale snapshot** of the queue; the live queue lives in the main tree and is the orchestrator's alone. Your REQ body is inlined in this brief — that is your copy of it.
- Commit your work on your own branch, as many commits as the work naturally splits into. Do not bump `VERSION`, do not touch `CHANGELOG.md`, do not touch `skills/do-work/actions/version.md` — those are serial-only files the integrator owns, and a builder bumping them races every sibling.
- If you need one line of wiring in a file outside your write set — a shared registry entry, a pointer in someone else's doc — **do not edit it**. Hand back the exact line and where it goes as an *integration seam*; the orchestrator applies it inside the merge commit.
- If you discover you need a file outside your declared scope, stop and report it in your hand-back rather than writing it silently.

**Crew rules load from your own worktree** (they ship, so they are there at the same paths): read `skills/do-work/crew-members/general.md`, `skills/do-work/crew-members/coding-guardrails.md`, and `skills/do-work/crew-members/communication-style.md` before you write code. This REQ is `maintenance: true` — a deliberate pass on the skill's own operating instructions — so also read `skills/do-work/crew-members/maintenance.md` and carry its delete-before-you-add discipline alongside the guardrails. Read every path in the REQ's `prime_files` too — those primes encode prior mistakes in this exact area.

**P-A-U phasing is mandatory.** Your REQ body carries an `AI Execution State (P-A-U Loop)` block. Work it: [PLAN] a brief technical approach, [APPLY] inside declared scope only, [UNIFY] run `git diff --stat`, run the linters, verify no debug artifacts, and list each file you checked. The orchestrator audits the checked boxes against the diff — a checked [UNIFY] over a diff containing stray instrumentation is a qualification failure. Record the P-A-U evidence in your hand-back, not in a `do-work/` file.

**Log significant decisions as D-XX** in your hand-back with reasoning. A reversible, low-reach choice is DECIDE & STATE (reasoning only); an irreversible, taste-dependent, or genuinely contestable one is ESCALATE — add `Value:` and `Risk:` lines.

**Out-of-scope finds** go in a `## Discovered Tasks` list in your hand-back. Do not fix them inline.

## Environment notes for this checkout

This is a Linux container, not the maintainer's machine. Before you start, know:

- `bash _dev/tests/maintainer-verify.sh` is the canonical pass/fail gate and it **exits 0 at your branch point** — that is your baseline. Exit code zero is the only proof; never pipe it through `| tail` or judge it from a summary. It takes a few minutes.
- The toolchain was installed for this session: Go 1.26.1, ShellCheck 0.11.0, `just` 1.21.0. They are on `PATH`.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB `queue-kanban` binary into the source tree, which is gitignored (so nothing warns you) and multiplies through the installer probe's copies. Build with `go build -o /tmp/<name> .`.
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward and never compute one.
- Browser tooling may drop a `.playwright-mcp/` directory into the main tree without you issuing a write. It is gitignored and holds sibling agents' evidence — if you clean it, remove only your own files.

## Hand-back

When you are done, write your report to exactly this absolute path:

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-162355/REQ-249-handback.md
```

That is the one main-tree path you may write. Never stage it, never commit it — it is an orchestrator working file, not branch content.

Structure it as:

```markdown
# REQ-NNN hand-back

**Branch:** worktree-agent-REQ-249-decide-the-cross-package-citation-path-form
**Commits:** <short hashes, oldest first>

## What I built
## File manifest
- `path` (new|modified|deleted) — one factual line each
## P-A-U evidence
## Testing evidence
Red-green: the test name, the failure text before, the pass after. Quote real output — never a transcript from a prototype or from memory.
Full gate: the `maintainer-verify.sh` exit code you actually observed.
## Decisions (D-XX)
## Integration seams
Exact lines and where they go, or "none".
## Discovered Tasks
## Pushback
Anything in this brief you think is wrong, with your evidence.
```

**One standing warning, from the previous session's own record.** Five consecutive REQs shipped a mechanism that looked like it closed a class and closed only the instance — and in three of five the remaining hole was exactly where the real data lives. Assume your first fix has that shape and go looking for the hole before a reviewer does. A passing assertion is not evidence about the thing you did not sample.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-249
title: Decide the cross-package citation path form and sweep to match
status: claimed
created_at: 2026-08-18T13:54:59Z
claimed_at: 2026-08-18T16:09:27Z
route: B
status_changed_at: 2026-08-18T14:12:05Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: true
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-08-18T16:10:30Z
  basis:
    - Route B
    - 15-file write set
    - 4 subsystems involved
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
- _dev/primes/prime-action-files.md
- skills/do-work/**/*.md
- skills/do-work-board/**/*.md
- skills/do-work-knowledge/**/*.md
- skills/do-work-toolbox/**/*.md
- _dev/tests/shipped-package-reference-contract.sh
---

# Decide the Cross-Package Citation Path Form and Sweep to Match

## What

Two incompatible readings of the cross-package citation path coexist in shipped markdown, and nothing can tell them apart. Pick one and sweep.

## Context

Raised by REQ-244's builder as pushback rather than decided mid-REQ, verified independently by the orchestrator, and confirmed again by REQ-244's review. REQ-244 added eleven more citations in the prescribed form while the question was open, so the count grows until this is settled.

- `\`../do-work/actions/work-reference.md\`` is what `_dev/primes/prime-action-files.md:91` prescribes and what `actions/memory.md` uses. From `skills/do-work-toolbox/actions/` it resolves to `skills/do-work-toolbox/do-work/actions/...`, **which does not exist.** It coheres only read as *skill-root*-relative rather than relative to the citing file's own directory.
- `\`../../do-work/actions/work-reference.md\`` is what `skills/do-work-toolbox/actions/present-work.md:37` already uses. It resolves literally, in both the source and installed topologies.

Counts at capture, verified by REQ-244's review: 30 local-form citations in core, 3 sibling-form in do-work-board, 11 sibling-form in do-work-knowledge and do-work-toolbox — consistent by package, so whichever reading is intended, the corpus is at least internally regular.

**REQ-243's checker cannot arbitrate this.** It resolves Markdown `[text](target)` link syntax; every one of these citations is a backticked path, which that checker never sees. So the correctness of cross-package pointers currently rests on convention alone.

## Open Questions

- [x] Shipped action files cite each other across skill packages with a backticked path, and there are two spellings in the tree that mean different things. One (`../do-work/actions/...`) is what the prime file tells writers to use, and it only makes sense if you read the `../` as "up to the skills folder" rather than as a real relative path — typed literally into a terminal from the citing file's folder, it points at nothing. The other (`../../do-work/actions/...`) is a real path that works from where the file actually sits. Both are in use today, the first far more than the second. Nothing checks either, because our new link checker only understands Markdown links and these are backticks. Which reading do you want to be the rule — the skill-root one that most files already use, or the literal one that a reader could paste and follow? Whichever you pick, the other spelling gets swept to match and the prime file gets updated to say so. → Confirmed: literal paths (`../../`), so a citation a reader pastes actually resolves and a future checker can verify it mechanically.
  Recommended: literal paths (`../../`), so a citation a reader pastes actually resolves and a future checker can verify it mechanically.
  Also: keep skill-root-relative (`../`), sweep `present-work.md` and `completed-work-presentation-reference.md` to match, and state the convention explicitly in the prime so nobody reads it as a filesystem path.

**Answered [2026-08-18]:** User confirmed the recommended literal form via `do-work clarify`. The rule is: every backticked cross-package citation must resolve as a real relative path from the citing file's own directory. Sweep every citation that doesn't (the Context counts are a starting set, not the extent — re-derive them), update `_dev/primes/prime-action-files.md:91` to prescribe the literal form, and answer the Requirements question in favor of mechanical checkability: backticked cross-package pointers should become checkable now that they resolve literally.

## Requirements

- One spelling is the documented rule, stated in `_dev/primes/prime-action-files.md`.
- Every shipped citation matches it — the sweep is the requirement, and the counts above are a starting set rather than the extent.
- Whether backticked cross-package pointers become mechanically checkable is part of the answer: if the literal form wins, say whether a checker should now read them.

---

## Triage

**Route: B** - Medium

**Reasoning:** The decision is already made; what is unknown is the extent — the capture counts are explicitly a starting set, so every citation across four shipped packages has to be found before any of them is rewritten.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `_dev/primes/prime-action-files.md` (modify) — § Cross-Referencing states the literal form as the rule
- Shipped markdown under `skills/do-work/`, `skills/do-work-board/`, `skills/do-work-knowledge/`, `skills/do-work-toolbox/` (modify) — every backticked cross-package citation swept to the spelling that resolves from the citing file's own directory
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — the mechanical check for backticked cross-package pointers, if the builder judges one is warranted

**Files I will NOT touch:**
- `skills/do-work-board/tools/queue-kanban/**` Go sources and `web/*.js` — REQ-248 is building there in this same wave.
- `CHANGELOG.md`, `VERSION`, `skills/do-work/actions/version.md` — serial-only, owned by the integrator.
- Markdown `[text](target)` links — REQ-243's checker already governs those; this REQ is about backticked paths.

**Acceptance criteria (restated from REQ):**
- [ ] One spelling is the documented rule, stated in `_dev/primes/prime-action-files.md`.
- [ ] Every shipped backticked cross-package citation resolves literally from its own file's directory — the capture counts are a starting set, not the extent (measured 126 local-form / 14 sibling-form at claim time, and package-root `SKILL.md` citations are already literal-correct).
- [ ] The answer states whether backticked cross-package pointers become mechanically checkable, and implements the check if the answer is yes.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 (recorded in `do-work/working/baseline.json`)
**Dependencies:** ⚠ this checkout needed Go 1.26.1, ShellCheck 0.11.0 and `just` installed before the baseline could run at all, and one pre-existing Linux-only test failure had to be fixed first (0.212.8) — see the REQ brief.

*Checked by work action*
