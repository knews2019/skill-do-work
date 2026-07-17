# Harness-Bloat Cleanup — Receipt (UR-003, phases 3–4)

Executed on branch `claude/skill-bloat-assessment-xziqsb`, versions 0.123.0 → 0.124.2,
one commit per REQ plus metadata commits. Every file edited was hash-verified against
`2026-07-15-phase1-hashes.txt` at phase-3 start (all 98 matched; no aborts).

## Before / after — words loaded before user content is touched

| Invocation | Before | After | Δ |
|---|---|---|---|
| Fixed floor (SKILL.md + next-steps.md, every invocation) | 7,174 | 4,063 | **−43%** |
| `do-work note` | 8,201 | 4,877 | **−41%** |
| `do-work board` | 8,726 | 5,615 | **−36%** |
| `do-work capture-request:` | 14,477 | 11,366 | **−21%** |
| `do-work pipeline` (orchestrator context) | 15,307 | 12,196 | **−20%** |
| `do-work run` (1 REQ, Routes B/C, orchestrator) | 25,038 | 21,763 | **−13%** |

Component deltas: SKILL.md 5,507 → 2,396 (−56%); work.md 10,371 → 10,207; note.md
1,027 → 814; scan-ideas.md 1,027 → 868; commit.md 1,904 → 1,629; quick-wins.md
1,844 → 1,788. New lazy-loaded `actions/help.md`: 760 words, loaded only on help routes.

**Shipped-payload delta:** live CHANGELOG 26,193 → ~2,900 words (144 entries moved to
the export-ignored `CHANGELOG-archive.md`; ~23.5k words leave every consumer tarball).

## Rules deleted vs relocated vs hardened

- **Deleted (≈3,900 instruction words + 23.5k payload words):** 2 of SKILL.md's 5
  action-set enumerations outright (Actions bullets, help menu*), 1 merged (Verb
  Reference → routing table Notes; every trigger verb and precedence rule preserved);
  CHANGELOG tail; intra-file guard restatements in note/scan-ideas/commit/quick-wins
  (every deleted row maps to a surviving rule — mappings in REQ-023).
  *help menu deleted from the always-loaded router; content lives on in `actions/help.md` (lazy).
- **Relocated: 0 files.** Per UR-003 Message 1, RELOCATE was plan-only this pass.
  Plans for prompts library (~19.8k words), interview (~12.6k), bkb+dream (~14.9k) in
  `2026-07-15-relocation-extraction-plans.md`, recommended order A→B→C.
- **Hardened: 2 steps fully, 2 partially** → `tools/checks/{archive-collision,preflight,scope-drift,qualify}.sh`
  (shipped, executable, script-missing fallbacks in prose, sync-asserted by the contract suite).
  preflight.sh also records `baseline.json` — groundwork for the Step 6.5 probation item.
- **Ratchets:** contract suite fails any commit pushing SKILL.md past 2,650 words
  (red-green demonstrated); CLAUDE.md now requires every new action to justify not
  being a sibling skill.

## Probation results (bounded tests, run this pass)

| Item | Test result | Disposition |
|---|---|---|
| tutorial.md drift | Every `do-work <verb>` in its recipes still routes post-diet | Resolved — keep; no drift found |
| docs/slop-check-guide duplication | Its table is NOT a byte-copy of crew-members/anti-slop.md (different form: 20-row table vs prose principles) | Resolved — keep as-is |
| ai-report relocation | 10 REQ/UR-schema touchpoints (threshold for clean extraction was <3) | Resolved — stays in do-work; optional future `ai-report-reference.md` split remains open |
| deep-explore relocation | No archived REQ originated from a deep-explore brief; but two maintenance REQs invested in its machinery (runs path, docs) | Still PROBATION — origination test says relocate, investment says keep; revisit after Plans A–C land |
| next-steps.md always-load (1,667 w) | Not testable in-repo (needs usage data: how often are suggestions followed) | Still PROBATION — untouched |
| Step 6.5 full harden | Groundwork shipped (baseline.json from preflight.sh); conversion needs validation in real runs | Still PROBATION — prose unchanged beyond the optional compare-against-baseline sentence |

## Deliberately NOT touched, and why

- **work.md steps 3.5, 3.7, 6.25, 7.5** — audit verdict: irreducibly judgment (decide-vs-escalate reasoning, spec-match judgment, generated manifests, lessons prose).
- **Guard sections in board, abandon, clarify, slop-check, forensics, verify-requests, validate-feedback** — survey verified them hard-won and file-specific; the original strip-all-small-actions hypothesis was falsified.
- **Action Dispatch table, principles blockquotes, argument-hint** — canonical name→path map (test-asserted), routing-governing philosophy, and skill-registration surface respectively.
- **docs/ guides (22 files)** — never enter agent context; user-facing value; payload cost only.
- **All 16 crew-members** — JIT gating verified real; the two always-loads (general+karpathy, 1,496 words) are the implementation guardrails the trail-of-intent design depends on.
- **Original ~1,500-word router target** — declined with evidence: full trigger-verb lists ARE routing behavior, and the dispatch table is contractually asserted; 2,396 was reached without sacrificing either. Budget ratcheted at 2,650.

## Traceability

UR-003 (verbatim request + capture notes) and its six REQs — REQ-019, REQ-020,
and REQ-021..024 (renumbered from REQ-015..018 on 2026-07-18 to resolve id
collisions with the earlier kanban-stream REQs) — (each with Triage,
P-A-U state, Implementation Summary, Testing, Lessons, and its implementation
commit hash) are archived in `do-work/archive/UR-003/`. Changelog entries:
0.123.0, 0.123.1, 0.123.2, 0.124.0, 0.124.1, 0.124.2.

Known review notes for the PR: shellcheck was unavailable in the execution
environment — the four `tools/checks/` scripts are hand-reviewed only; and the
`kb_status` line was dropped from REQ-021's frontmatter template during metadata
editing (harmless — the field is optional and handoff-owned).
