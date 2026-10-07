---
id: REQ-644
title: '[impact-rule-change] Capture distinguishes a session earmark from work that needs the user as operator'
status: claimed
created_at: 2026-10-07T23:21:03Z
user_request: UR-140
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-643]
batch: earmark-placement
write_set: [skills/do-work/actions/capture.md, skills/do-work/docs/work-guide.md]
claimed_at: 2026-10-07T23:26:38Z
---
# Capture Distinguishes a Session Earmark From Work That Needs the User as Operator
## What
Capture's Earmark assessment says which of two shapes to write: `assigned_to` when another session will take the work (routing between sessions), and `status: blocked` with `blocked_by` naming the person plus `blocked_at` when the user has to be present as operator. The External-condition assessment lists the earmark as a fourth look-alike.
## Why
Today "leave this one for me" and "I have to be at the keyboard for this one" both read as earmarks. An earmark lands under Pending (today under Ready, after REQ-643 under Earmarked) and nothing releases it except a hand edit or explicit targeting. Work that needs the user belongs in Needs Input · Blocked, where `do-work clarify` already confirms the condition and releases the REQ. The maintainer had to flip two REQs by hand to get that.
## Verified Facts (from triage)
- `skills/do-work/actions/capture.md:107` (Earmark assessment) never mentions `blocked`.
- `capture.md:108` (External-condition assessment) already lists three look-alikes: a wait on another REQ is `depends_on`, a question for the user is `pending-answers`, an outside person's confirmation while work proceeds is `Answerer:`. The earmark is missing from that list.
- The release path exists: `skills/do-work/actions/clarify.md:13` and Step 5.5 confirm human-confirmable blocked conditions, and `skills/do-work/actions/work.md:97` keeps a `blocked` REQ out of the scan.
- `skills/do-work/docs/work-guide.md:113` introduces the field as "To say 'leave this one for me'", which is the ambiguous phrasing.
## Detailed Requirements
1. `capture.md:107`, Earmark assessment: add one sentence stating the distinction. Content to carry: `assigned_to` routes work between sessions; when the user must be present as operator, capture `status: blocked` with `blocked_by` naming the person (the user's words, Frontmatter Quoting) and `blocked_at`, which lands in Needs Input · Blocked and is released by `do-work clarify`. Point to the External-condition assessment rather than restating its mechanics.
2. `capture.md:108`, External-condition assessment: add the earmark as a fourth look-alike in the existing "Keep this distinct from the look-alikes" list: work reserved for another session is `assigned_to` (not blocked).
3. `docs/work-guide.md:113-119`, "Earmarking with `assigned_to`": one sentence saying the field is for another session or checkout, and that work waiting on you as operator is captured `blocked` with `blocked_by` naming you, so it shows under Needs Input · Blocked.
## Constraints
- Prose only. No new field, no new status, no Go change, no SKILL.md routing row.
- Keep the never-invent rule: no earmark in the user's words means no `assigned_to`; no person named means no `blocked_by`.
- Read `_dev/primes/prime-action-files.md` before editing; sweep any other restatement of the earmark rule the grep finds (family alternate-writer-contract-drift). Today's grep found only the three sites above.
## Dependencies
None. REQ-643 (board shows an earmarked pending REQ under Pending → Earmarked) is related; this sentence is true before and after it ships, so there is no `depends_on`.
## Builder Guidance
Certainty is high: the maintainer specified the sentence's content. Latitude is wording only. Keep each addition to one sentence per site.
## Red-Green Proof
**RED prompt/case:** A capture request reading "I need to be at the keyboard for this one, hold it for me". Follow `capture.md` Step 1 as written today.
**Why RED now:** The Earmark assessment is the only rule that matches, so capture seeds `assigned_to`; the REQ lands under Pending and nothing in the pipeline releases it.
**GREEN when:** `capture.md:107` says `assigned_to` is for routing between sessions and operator-present work is captured `blocked` with `blocked_by` naming the person and `blocked_at`; `capture.md:108` lists the earmark as a look-alike; `work-guide.md` says the same in one sentence. Checked by a prose read and `_dev/tests/shipped-package-reference-contract.sh`.
**Validation:** User confirmed. The maintainer wrote the distinction in the verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 6429 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` covers restated rules across `capture.md` and `work-guide.md`.
## Full Context
See `do-work/user-requests/UR-140/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: maintainer feedback "Board: a pending REQ earmarked for the user reads as 'Ready', which misleads", change 2, accepted by validate-feedback.*
