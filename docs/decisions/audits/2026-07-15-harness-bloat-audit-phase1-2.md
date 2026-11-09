# do-work Harness-Bloat Audit — Phases 1–2

**Repo state audited:** commit `c97d32c`, clean tree. SHA-256 of all 98 audited files snapshotted to `scratchpad/phase1-hashes.txt` — phase 3 re-verifies each file against this snapshot and aborts any item whose file moved.

**Evidence legend** (every claim carries one):
- **VERIFIED** — observed directly in a file (cited `file:line`) or produced by a command run during this audit (word counts via `wc -w`, dates via `git log`).
- **INFERRED** — concluded from structure. In particular: *all* "when it loads" values describe **declared loading rules** (verified as text at the cited location). Whether an agent actually loads the file at runtime is never observable from the repo — every actual-load claim is INFERRED from the declared rule. Where a rule says "read only section X," a full-file load is additionally INFERRED from how agents typically consume single files.
- **USER_REPORTED** — carried over from the prior audit or the task description without independent verification. (None survive below — every prior-audit claim was re-checked and is now labeled VERIFIED, INFERRED, or corrected.)

---

## PHASE 1 — Map of the Active Path

### 1a. Load graph (declared rules, all VERIFIED at the cited location)

| Asset | Declared load rule | Where declared |
|---|---|---|
| SKILL.md | Entry point; body loads on every `do-work` invocation | Skill frontmatter contract |
| next-steps.md | "After every action completes … see next-steps.md for the full per-action reference" — every invocation, at end | SKILL.md:351-353 |
| actions/*.md (30 routed) | One per invocation, per the routing table → Action Dispatch table | SKILL.md:69-108, 295-330 |
| actions/*-reference.md (4) | By companion action, at named steps/sections | e.g. work.md:27,81,115 ("read the named section") |
| actions/kb-lessons-handoff.md | By work.md Lessons-Capture Phase and review-work.md Self-Validation | work.md:457 |
| actions/sample-archived-req.md | Pointer from work.md as example — load optional | work.md:592 |
| crew-members/general.md, karpathy.md | Always at implementation (work.md Step 6) | work.md:301-302 |
| crew-members/* (other 14) | JIT per each file's JIT_CONTEXT condition | each file line 3; work.md:303-312; capture.md:224,320; pipeline.md:149,175 |
| specs/*.md | work.md Step 3.7, on task-type/`suggested_spec` match | work.md:215-226 |
| interviews/*.md | interview action only | interview.md:54-80 |
| prompts/*.md | prompts action, only on user-requested `run <name>` (agent-verified: no other load path exists) | prompts.md sub-command `run` |
| docs/*-guide.md (22) | **Never.** No runtime file instructs reading a guide; all references are "a walkthrough exists" footnotes (agent-verified across SKILL.md, actions/, crew-members/, next-steps.md) | e.g. work.md:3 |
| CHANGELOG.md | First ~80 lines only, by version.md ("do NOT load the full file"; first 5 `## ` blocks) | version.md:35-36 |
| CLAUDE.md / decisions/ / _dev/ | Maintainer-side only; decisions/, _dev/ are export-ignored | .gitattributes |

### 1b. Context cost of the 5 most common invocations

Words of skill instruction loaded before user content is touched. Composition = declared rules (VERIFIED at citations above); the totals assume a compliant agent loads each declared file in full (INFERRED).

| Invocation | Composition | Orchestrator-context words | Notes |
|---|---|---|---|
| `do-work note` | SKILL 5,507 + note.md 1,027 + next-steps 1,667 | **8,201** | 87% of the load is router+overhead for a one-line file append |
| `do-work capture-request:` | SKILL 5,507 + capture.md 5,258 + prompt-injection.md 1,304 + clear-questions.md 741 + next-steps 1,667 | **14,477** | crew loads VERIFIED at capture.md:224, 320 |
| `do-work run` (1 REQ, Route B/C) | SKILL 5,507 + work.md 10,371 + work-reference.md 5,997 + general.md 670 + karpathy.md 826 + next-steps 1,667 | **25,038** | plus dispatched sub-contexts: review-work.md 4,566 (Step 7), kb-lessons-handoff.md 2,146 (Step 7.5), cleanup.md 2,798 (Step 10); plus conditional crew (security 1,644 / testing 1,451 / debugging 1,149 / domain file / caveman 357). work-reference is *declared* section-at-a-time (work.md phrasing "→ **Section**"); full-file load INFERRED |
| `do-work pipeline <request>` | SKILL 5,507 + pipeline.md 7,231 + anti-slop.md 902 + next-steps 1,667 = **15,307** orchestrator | **≈60k+ aggregate** | pipeline dispatches capture (~14.5k), run (~25k), verify, review, present — each a fresh sub-context with its own stack (INFERRED from dispatch table pipeline.md:139-142) |
| `do-work board` | SKILL 5,507 + board.md 1,552 + next-steps 1,667 | **8,726** | |

**Fixed floor on every invocation: 7,174 words (SKILL.md + next-steps.md) — ~9.5k tokens before any action file opens.** (VERIFIED word counts; always-load rule VERIFIED at SKILL.md:353.)

### 1c. SKILL.md internal anatomy (all VERIFIED via per-section `wc -w`)

| Section (lines) | Words | Job | Duplication |
|---|---|---|---|
| Frontmatter + argument-hint (1-9) | ~120 | Skill registration; argument-hint enumerates all 34 verbs | 6th partial enumeration of the action set |
| Actions bullet list (11-46) | 722 | One-line-per-action descriptions | Duplicates each action's blockquote + When-to-Use, and the help menu |
| Principles blockquotes (48-60) | 421 | Core concept, Trail of Intent, Capture≠Execute, two-windows | Unique — governs routing behavior |
| Routing priority table + parse rules (62-127) | 1,225 | Priority-ordered verb → route matching | Overlaps Verb Reference |
| Verb Reference (129-166) | 1,526 | Same verb→route mapping + disambiguation notes | ~60-70% overlaps the priority table; unique content: install normalization, dream-over-cleanup precedence, per-route payload notes |
| Help menu example (168-251) | 556 | Literal output spec for `do-work help` | Third full enumeration; only needed on the help route |
| Per-command help (253-274) | 109 | `<action> help` behavior | Only needed on help-suffix invocations |
| Payload preservation (276-289) | 83 | Never lose $ARGUMENTS | Unique |
| Action Dispatch (291-349) | 710 | Canonical name→path map + context-to-pass + subagent/foreground rules | CLAUDE.md names it the canonical mapping; `_dev/tests/contract-regressions.sh:63-69` greps this exact block |
| Next-steps pointer (351-353) | 36 | Load next-steps.md after every action | Unique |

**Correction to the prior audit (was USER_REPORTED "four enumerations"): there are FIVE full enumerations of the action set** (bullets, priority table, verb reference, help menu, dispatch table) plus the argument-hint. VERIFIED above.

### 1d. Per-file map — actions/ (39 files)

Load condition for all routed actions: per-action (loaded only when its verb routes there) — VERIFIED at SKILL.md:295-330. References load per companion action. "Evidence it still helps": last content change (`git log -1 --format=%cs`, VERIFIED); in-repo usage traces where they exist (this repo's own `do-work/archive/` holds 14 archived REQs incl. 11 `review_generated: true`, and `ai-reports/` holds 2 generated reports — VERIFIED; absence of traces here is weak evidence, since the skill is mostly used in consumer repos — INFERRED).

| File | Words | Job | Evidence it helps | Duplication |
|---|---|---|---|---|
| work.md | 10,371 | Queue orchestrator, 10 steps + 9 fractional | touched 07-14; archive traces | docs/work-guide (narrative, not verbatim); templates already extracted to work-reference |
| ai-report.md | 7,354 | HTML proof-of-work report generator | touched 07-14; 2 reports in ai-reports/ | docs/ai-report-guide (condensed); deconflicted from present-work in-file |
| pipeline.md | 7,231 | End-to-end orchestration, stateful | touched 07-13 | none — dispatches, never re-implements (VERIFIED pipeline.md:139-142) |
| bkb.md | 6,711 | KB wiki builder (init/triage/ingest/query/lint) | touched 06-30 | docs/bkb-guide; templates in bkb-reference |
| work-reference.md | 5,997 | Canonical schema + heavy templates for work.md; cited by ~12 actions | touched 07-14 | none — exists to prevent duplication |
| capture.md | 5,258 | UR+REQ capture, RED/GREEN inference | touched 07-07; archive traces | docs/capture-guide (351w) |
| install.md | 5,205 | Companion-tool installer (4 targets) | touched 07-13 | none |
| review-work.md | 4,566 | Post-work review, follow-up REQs | touched 07-11; 11 review_generated REQs | docs/review-work-guide (344w) |
| interview.md | 4,344 | Elicitation-template runner | touched 07-13 | docs/interview-guide; schemas in interview-reference |
| interview-reference.md | 4,257 | Interview schemas/formats companion | touched 06-30 | none (extraction target) |
| code-review.md | 3,707 | Standalone codebase review | touched 07-14 | docs/code-review-guide (399w) |
| dream.md | 3,471 | Memory-dir consolidation (destructive) | touched 06-30 | docs/dream-guide (1,001w) |
| present-work.md | 3,401 | Client-facing deliverables | touched 07-01 | docs/present-work-guide (236w) |
| inspect.md | 3,269 | Explain uncommitted changes | touched 06-30 | docs/inspect-guide (265w) |
| roadmap.md | 3,243 | Queue survey (read-only) | touched 07-10 | docs/roadmap-guide (411w) |
| prime.md | 3,182 | Prime-file create/audit | touched 07-14 | docs/prime-guide (644w) |
| cleanup.md | 2,798 | Archive consolidation; auto-runs after work loop | touched 07-14 | docs/cleanup-guide (349w) |
| ui-review.md | 2,779 | UI-quality audit | touched 06-30 | docs/ui-review-guide (294w) |
| deep-explore-reference.md | 2,570 | Personas/rubric companion | touched 07-14 | none |
| stray-check.md | 2,525 | Repo junk scan | touched 06-30 | docs/stray-check-guide (613w) |
| deep-explore.md | 2,450 | Multi-round concept exploration | touched 07-14 | none |
| version.md | 2,386 | Version/update/recap; parses CHANGELOG first ~80 lines | touched 07-14 | docs/version-guide (142w, near-zero overlap) |
| prompts.md | 2,271 | Prompt-library dispatcher (list/show/run) | touched 06-30 | docs/prompts-guide (572w) |
| tidy-repo.md | 2,258 | Safe layout reorganization | touched 07-13 | none |
| bkb-reference.md | 2,182 | BKB seed templates/crew defs | touched 06-30 | none |
| kb-lessons-handoff.md | 2,146 | Shared Lessons→KB handoff (work + review-work) | touched 06-30 | none — exists to dedupe the two callers |
| forensics.md | 1,936 | Pipeline diagnostics, 11 checks | touched 07-10 | docs/forensics-guide (305w) |
| commit.md | 1,904 | Atomic commit grouping | touched 06-30 | **internal**: stacks Checklist + Common-mistakes + CR + RF + VC — 5 overlapping guard sections (agent-VERIFIED) |
| tutorial.md | 1,855 | Interactive tutorials | touched 06-30 | **restates SKILL.md routing + concepts in Recipes/Concepts modes — its own Red Flags name this drift risk** (agent-VERIFIED) |
| quick-wins.md | 1,844 | Refactor/test-gap scan | touched 06-30 | some generic CR rows (agent-VERIFIED); carries the maintenance-marker contract (load-bearing) |
| verify-requests.md | 1,794 | Capture QA vs UR input | touched 06-30 | docs/verify-requests-guide (315w) |
| validate-feedback.md | 1,634 | External-findings triage | touched 06-30 | none; boilerplate strongly file-specific (agent-VERIFIED) |
| slop-check.md | 1,632 | Anti-slop artifact check | touched 07-01 | docs/slop-check-guide **heavily** restates it + the crew principles (agent-VERIFIED, worst-offender guide) |
| board.md | 1,552 | Kanban board build+run | touched 07-10 | boilerplate hard-won (stale-binary trap, gitignore walk-up) |
| clarify.md | 1,348 | Batch-review pending-answers | touched 07-11 | none |
| abandon.md | 1,245 | Cancel a REQ | touched 07-06 | none; CR rows hard-won (cancelled-vs-failed) |
| note.md | 1,027 | Append dated hint to notes.md | touched 07-13 | **internal**: CR/RF/VC rows map ≈1:1 onto its own Rules section (VERIFIED by direct read) |
| scan-ideas.md | 1,027 | Idea generation, read-only | touched 06-30 | **internal**: RF/VC restate its Rules (VERIFIED by direct read) |
| sample-archived-req.md | 599 | Golden example REQ | touched 06-30 | none — exemplar of work-reference schema |

Boilerplate census (agent-VERIFIED): 24 actions carry all three of Common Rationalizations/Red Flags/Verification Checklist; 10 carry RF+VC only; the 5 reference/sample files correctly carry none.

### 1e. Per-file map — crew-members/ (16 files)

All JIT-gated per their JIT_CONTEXT line-3 comments (VERIFIED); actual loading INFERRED from those declared rules. Referencing actions verified by grep (agent-VERIFIED).

| File | Words | Trigger | Referenced by | Duplication |
|---|---|---|---|---|
| general.md | 670 | ALWAYS at implementation | work, prime | none |
| karpathy.md | 826 | ALWAYS at implementation | 7 files incl. SKILL.md, specs | intentional contrast w/ maintenance.md |
| clear-questions.md | 741 | any interactive question | 7 files | none |
| prompt-injection.md | 1,304 | any third-party-content ingestion | 10 files | none |
| anti-slop.md | 902 | any human-facing artifact | 10 files | slop-check action applies it (layered, not dup) |
| background-agents.md | 2,020 | any fan-out | 8 files | none |
| maintenance.md | 786 | `maintenance: true` marker only | capture, work, quick-wins, work-reference | none |
| security.md | 1,644 | domain/heuristic security surface | work, code-review | mild layering w/ code-review's own checks |
| testing.md | 1,451 | tdd/domain/test-fail-loop | work + domain template | none |
| debugging.md | 1,149 | remediation / 2nd+ test attempt | work | none |
| ui-design.md | 1,043 | domain: ui-design | ui-review + template | none |
| backend.md | 793 | domain template only | (template only) | none |
| frontend.md | 538 | domain template only | (template only) | none |
| interviewer.md | 664 | interview action | interview | pairs w/ clear-questions (structure vs wording) |
| approach-directives.md | 459 | multi-REQ parallel dispatch | work | none |
| caveman.md | 357 | `caveman` frontmatter | work, work-reference | none |

### 1f. Per-file map — prompts/ (18), docs/ (22), root files

**prompts/** — all load ONLY on user-requested `prompts run <name>` (agent-VERIFIED: no other load path). Wiring beyond prompts.md/README/CHANGELOG: **zero for all 17** (agent-VERIFIED by grep). Relation to the task-queue mission:

| Group | Files | Words | Mission-related? |
|---|---|---|---|
| ADR log | architecture-decisions-log_create-or-expand | 2,470 | Yes (dev tooling) |
| Dark-code kit | audit, comprehension-gate, context-layer-generator | 3,608 | Yes-ish (comprehension-gate overlaps built-in code-review) |
| Prompt-kit coaching | step0–step6 (7 files) | 7,781 | No — personal AI-skills coaching (step3 partially spec-related) |
| Business/economics | business-vendor, economics ×2 | 2,492 | No |
| Tech-infra strategy | tech-inference, tech-infrastructure | 1,795 | No |
| News analysis | weekly-structural-diff | 756 | No |
| Index | README.md | 942 | — |

**docs/** — 22 guides, ~15.6k words, **never enter agent context** (agent-VERIFIED). Human-facing walkthroughs; overlap with action files ranges from near-zero (version-guide 142w, commit-guide 162w) to heavy (slop-check-guide 952w nearly duplicates the action's principle table, N/A rules, rewrite flow — agent-VERIFIED by diff-sampling 6 guides).

**Root files:** next-steps.md 1,667w (always-load, see 1b). CHANGELOG.md 26,193w / 162 entries — runtime reads only first ~80 lines (version.md:35, VERIFIED); the remaining ~24k words are pure shipped payload. Precedent for archiving exists: 0.76.0 deleted two archive changelogs; 0.76.1 (quoted in full by the survey) restored discoverability via a GitHub commit pointer (`bf15fe2`) for tarball installs with no `.git` — exactly the no-git-required pointer pattern phase 3 item 2 needs. README.md 1,745w (install/usage, human-only). AGENTS.md 2w (stub).

---

## PHASE 2 — Classification

Every asset in exactly one bucket. Items marked ⚠ diverge from the prior audit / task-description candidate list, with the reason.

### KEEP-ACTIVE

| Item | One-line justification (evidence) |
|---|---|
| SKILL.md: routing priority table, payload rules, principles blockquotes, dispatch mapping | The router's actual job; dispatch block is the canonical name→path map (CLAUDE.md) and is grep-asserted by `_dev/tests/contract-regressions.sh:63-69` (VERIFIED) |
| next-steps.md | The action-chaining mechanism, loaded once at end of action; 1,667w is tolerable and no cheaper design is proven (see PROBATION for the test) |
| All core queue actions: capture, work, work-reference, verify-requests, review-work, clarify, abandon, cleanup, commit, inspect, version, forensics, roadmap, note, board, pipeline, kb-lessons-handoff, sample-archived-req | Per-action loads — zero cost unless invoked (VERIFIED dispatch table); on-mission; usage traces in this repo's own archive |
| Scan/review family: code-review, quick-wins, scan-ideas, ui-review, validate-feedback, slop-check, present-work, stray-check, tidy-repo, prime, install, tutorial | Same per-action economics; deconfliction between siblings is explicit in each file (agent-VERIFIED) |
| deep-explore + reference ⚠ | Prior-audit RELOCATE candidate — but it feeds capture directly, is part of the ideation family with scan-ideas/quick-wins (staying), and was touched 2026-07-14 (actively developed). Relocation buys repo hygiene only, not context savings (per-action load). See PROBATION for the bounded test |
| All 16 crew-members | Already JIT-gated (VERIFIED JIT_CONTEXT lines); the two always-loads (general+karpathy, 1,496w combined) are the implementation guardrails the whole Trail-of-Intent design leans on |
| specs/ (5 files, 1,842w) | Already lazy (work.md:215-226); tiny |
| docs/ (22 guides) | Never in agent context (agent-VERIFIED); user-facing value; payload cost only |
| CHANGELOG.md head (newest ~20 entries) | Runtime contract: version.md reads first ~80 lines / 5 entries (VERIFIED version.md:35-36) |
| tools/queue-kanban | Compiled on demand, never context (VERIFIED board.md flow) |

### LAZY-LOAD

| Item | Trigger condition |
|---|---|
| SKILL.md help-menu example + per-command-help (665w combined) | Move to a new `actions/help.md`; load only on routing priority 1 (bare/`help`) and on `<action> help` invocations. It's an output spec, not routing logic — every non-help invocation pays 665w for nothing (VERIFIED sections 168-274). ⚠ Prior audit implied DELETE; it can't be deleted — `do-work help` needs a stable menu spec — but it can stop taxing the other 33 verbs |

(Nothing else qualifies: crew, specs, interviews, prompts, references, and docs are already lazy or never-loaded — the prior audit's "JIT discipline is real" holds, VERIFIED.)

### RELOCATE (phase 3 produces extraction plans only)

| Item | Words leaving | Evidence |
|---|---|---|
| prompts/ library — the 12 unwired, off-mission files (prompt-kit ×7, business/economics ×3, tech-infra ×2, weekly-diff) | ~12.8k | Zero references from any action (agent-VERIFIED grep); topics are personal/business coaching, not task-queue. Keep `actions/prompts.md` (the runner — it already supports project-local prompt dirs) and decide separately on the 4 dev-adjacent prompts (ADR log, dark-code kit) — recommendation: move the whole library for coherence, keep the runner |
| interview + interview-reference + interviews/ + crew-members/interviewer.md + docs/interview-guide | ~12.6k | Self-contained subsystem; only outbound touchpoint is optional bkb ingest (agent-VERIFIED); no other action depends on it |
| bkb + bkb-reference + docs/bkb-guide | ~10.5k | Soft coupling only: kb-lessons-handoff already degrades gracefully when `kb/` is absent (VERIFIED — defers to `pending`, points at `bkb init`, never blocks archival, per work.md:457 and CLAUDE.md contract), so extraction doesn't break the pipeline |
| dream + docs/dream-guide | ~4.5k | Pairs with bkb (consolidates the wiki bkb builds); goes wherever bkb goes |
| ai-report ⚠ | — | **Not recommended.** Candidate named in the task, but: it renders the REQ/UR trail (coupled to the schema that IS the skill's product), has real usage traces (2 reports in `ai-reports/`, VERIFIED), and is the most recently maintained action (07-14). Moved to PROBATION with a bounded test |
| deep-explore ⚠ | — | **Not recommended** — see KEEP-ACTIVE rationale; bounded test in PROBATION |

### HARDEN (work.md fractional steps — per-step verdict from full read of work.md)

| Step | Verdict | Reasoning |
|---|---|---|
| 2.0 Pre-claim archive collision (157-174) | **YES — fully scriptable** | Pure mechanics: extract REQ-NNN, glob two archive patterns, write status, print message. Zero judgment. Prose shrinks to a one-line pointer at the check |
| 3.5 Open Questions (191-213) | **NO — prose** ⚠ | Marking `- [~]` with decide-vs-escalate reasoning is judgment; only the D-XX counter-comment format is checkable, and that's not worth a check |
| 3.7 Spec loading (215-226) | **NO — prose** ⚠ | Half the trigger is judgment ("clearly indicates a task type"); the whole step is ~200 words — hardening saves nothing |
| 5.5 Scope declaration (265-281) | **PARTIAL** | Writing the declaration is judgment; the review-time drift check (declared file list vs Implementation Summary file list, set difference) is a testable condition → script; declaration prose stays |
| 5.75 Pre-flight (283-295) | **YES — fully scriptable** | `git status -uall`, run test baseline, check deps presence — all commands, all warnings-not-blockers. Script it; prose becomes a pointer |
| 6.25 Implementation Summary (332-348) | **NO for the writing; its validation folds into 6.3** | The manifest is generated prose; the "only do-work/ paths ⇒ not implemented" rule is checkable and belongs in the qualify check |
| 6.3 Qualify (350-374) | **PARTIAL** | Checks 1 (files exist / in diff), 4 (P-A-U boxes vs debug-artifact grep), and most of 5 (wiring grep) are scriptable; checks 2 (substantive), 3 (requirements traced), 6 (hollow data path) are genuinely judgment — prose stays for those, script feeds them evidence |
| 6.5 Testing (376-394) | **PROBATION** ⚠ | Test discovery is prime-driven judgment; the baseline-vs-regression comparison is only checkable if 5.75 records its baseline machine-readably — a design change, not a transcription. Bounded test below |
| 7.5 Lessons-capture (425-457) | **NO — prose** ⚠ | "What worked / what didn't / worth knowing" is irreducibly judgment |

Placement decision needed at approval: hardened checks must **ship** (they run in consumer repos), so they'd live in a shipped directory (`hooks/` or a new `tools/checks/`) as plain bash with graceful degradation, per the design-for-the-floor rule; `_dev/tests/contract-regressions.sh` gains maintainer-side assertions that the prose pointers stay in sync.

### DELETE

| Item | Words | Evidence |
|---|---|---|
| SKILL.md Actions bullet list (11-46) | 722 | Third-hand duplicate: each action's own blockquote + When-to-Use, and the help menu, carry the same content (VERIFIED) |
| SKILL.md Verb Reference (129-166) — **merge, not pure delete** | net ≈ −1,000–1,200 | ~60-70% duplicates the priority table (VERIFIED side-by-side); unique disambiguation (install normalization SKILL.md:159, dream-over-cleanup precedence :164, abandon ID-rule :136) must be folded into the priority table's Notes column before the section goes |
| CHANGELOG.md tail (entries 21-162) | ~23k payload | Runtime needs only first ~80 lines (VERIFIED version.md:35); 0.76.1 gives the exact archive+pointer pattern for tarball installs (VERIFIED, quoted in survey) |
| note.md CR+RF sections | ~230 | Rows map ≈1:1 onto its Rules section (VERIFIED by direct read) — the *lessons* survive in Rules; only the restatement goes |
| scan-ideas.md RF+VC restated rows | ~180 | Same 1:1 mapping onto Rules (VERIFIED by direct read) |
| commit.md redundant guard stack | ~150-250 | Carries FIVE overlapping guard sections (Checklist + Common-mistakes + CR + RF + VC, agent-VERIFIED); dedupe to the standard triad |
| quick-wins.md generic CR rows | ~60 | A few rows restate general refactoring wisdom (agent-VERIFIED); trim rows, keep section — its Rules carry the load-bearing maintenance-marker contract |

⚠ **Correction to the prior audit's boilerplate hypothesis:** the survey verified that in small actions the triad content is mostly *hard-won and file-specific* (note.md's reader-contract bug, board.md's stale-binary trap, validate-feedback's dishonest-pushback row, abandon's cancelled-vs-failed row). Wholesale stripping would delete real institutional memory. The defect is **triple-stating within a file**, so the DELETE scope above is restatement-dedupe in 4 named files — not section removal across all small actions. board, clarify, abandon, slop-check, forensics, verify-requests, validate-feedback keep their sections untouched.

### PROBATION (evidence too weak to change now — with the bounded test that settles it)

| Item | Why evidence is weak | Bounded test |
|---|---|---|
| ai-report relocation | Coupled to REQ/UR schema; actively used and maintained — but 7,354w in one action file is work.md-scale without a reference split | Count its schema touchpoints: if its only interface is "read an archived REQ/UR file," extraction is clean; separately, splitting a `ai-report-reference.md` (render-judge pass, SVG rules) mirrors the proven work/work-reference pattern regardless of relocation |
| deep-explore relocation | Feeds capture; recently maintained; but no in-repo evidence its briefs ever became REQs | Grep `do-work/archive/` + `decisions/` for REQs/URs originating from deep-explore outputs; zero hits over the archive's lifetime ⇒ promote to RELOCATE |
| next-steps.md always-load (1,667w) | It's the chaining UX; but no evidence suggestions are followed often enough to justify 1.7k words on every invocation | Instrument one week of real use: count invocations that are follow-ons of a suggestion. Alternative design if low: move each action's block into that action's file (per-action cost instead of always-cost) |
| Step 6.5 harden | Baseline comparison needs machine-readable baseline from 5.75 — untested design change | Prototype: have the 5.75 script emit `baseline.json` in one real run; if the 6.5 comparison consumes it without breaking floor-compat, promote to HARDEN |
| tutorial.md drift | Its Recipes/Concepts restate routing; its own Red Flags admit the risk — but no measured drift yet | Run its own Verification Checklist spot-check (recipes vs current routing table); actual drift found ⇒ add a contract-regression assertion or generate recipes from the routing table |
| docs/slop-check-guide.md near-duplication | Heaviest guide-vs-action overlap (agent-VERIFIED) but docs are never agent-loaded, so cost is maintenance-sync only | Diff its principle table against crew-members/anti-slop.md; byte-near-identical ⇒ slim guide to a pointer + examples |

---

## Corrections & pushback summary (things the candidate list got wrong)

1. **"3 of 4 enumerations go"** → there are **five**; the Action Dispatch table must stay (canonical map, regression-asserted), the help menu is LAZY-LOAD not DELETE, and the Verb Reference is a merge.
2. **Router under ~1,500 words is not reachable** without gutting routing correctness: priority table (1,225) + dispatch (710) + payload (83) + principles + frontmatter ≈ **2,400–2,900 realistic post-diet range** (aggressive variant merges dispatch into the routing table as a file-path column → ~2,100–2,300, and requires updating `contract-regressions.sh:63` in the same commit). Proposed revised target: **≤2,500 words**, ratcheted in phase 4.
3. **Boilerplate strip** → scoped down to restatement-dedupe in note, scan-ideas, commit, quick-wins (evidence above); the sections themselves survive.
4. **HARDEN list** → 2.0 and 5.75 fully; 5.5 and 6.3 partially; 3.5, 3.7, 6.25, 7.5 stay prose; 6.5 on probation.
5. **RELOCATE list** → bkb, dream, interview, prompts-library confirmed; **ai-report and deep-explore rejected for now** (probation tests defined).

## Phase-3 readiness (for after approval)

- Hash snapshot: `phase1-hashes.txt` (98 files, SHA-256, at `c97d32c`). Each phase-3 item re-verifies its target files against this before editing; mismatch ⇒ abort that item and report.
- Execution vehicle: maintenance REQs via `do-work capture-request:` with `maintenance: true` (loads crew-members/maintenance.md per work.md Step 6.5a), one REQ per approved bucket-item group, each with its own commit + changelog entry, on branch `claude/skill-bloat-assessment-xziqsb`.
- Known same-commit couplings: SKILL.md diet ⇔ `contract-regressions.sh` dispatch-block assertion; CHANGELOG truncation ⇔ version.md glob note + 0.76.1-style header pointer; any work.md step conversion ⇔ its Orchestrator Checklist line (work.md:544-564) and Common Rationalizations rows that cite the step.
