# do-work: Development Portfolio

**Total versions shipped:** 20 releases (v0.1.0 — v0.20.0)
**Period:** 2026-01-27 — 2026-03-02
**Repository:** [knews2019/skill-do-work](https://github.com/knews2019/skill-do-work)

---

## Release History

### Foundation (v0.1.0 — v0.3.0)

Core system built from scratch. Task capture, work loop processing, REQ file lifecycle (pending → working → archived), git-aware commits, version management.

| Version | Name | What shipped |
|---------|------|-------------|
| 0.1.0 | Hello, World | Core capture + work loop + REQ lifecycle + git commits |
| 0.1.1 | Typo Patrol | Fixed install command username |
| 0.2.0 | Trust but Verify | Testing phase in work loop, orchestrator separation |
| 0.3.0 | Self-Aware | Version check, update check, docs |

### Structure & Safety (v0.4.0 — v0.6.0)

The system learned organization. User Request (UR) grouping, archive cleanup, immutability rules for in-flight and completed work.

| Version | Name | What shipped |
|---------|------|-------------|
| 0.4.0 | The Organizer | UR system, cleanup action, auto-archive consolidation |
| 0.5.0 | The Record Keeper | Changelog, version bump workflow |
| 0.6.0 | The Bouncer | Immutability for working/archive, addendum pattern |

### Quality & Verification (v0.7.0 — v0.9.5)

Added quality gates: verify checks capture accuracy, hints nudge users toward verification for complex requests, routing got reliable priority ordering, changelog display.

| Version | Name | What shipped |
|---------|------|-------------|
| 0.7.0 | The Nudge | Verify hint after complex captures |
| 0.8.0 | The Clarity Pass | UR+REQ pairing made unmissable, agent compatibility rules |
| 0.9.0 | The Rewind | Changelog display (reversed for terminal reading) |
| 0.9.1 | The Gatekeeper | Priority-ordered routing, keyword-first matching |
| 0.9.2 | The Front Door | Fixed SKILL.md frontmatter for CLI parsing |
| 0.9.3 | The Timestamp | Dates on every changelog entry |
| 0.9.4 | The Passport | Portable install/update commands |
| 0.9.5 | The Reinstall | Fixed update command (full reinstall over silent failure) |

### Portability & Scale (v0.10.0 — v0.12.7)

Separated capture from execution, moved actions to subagents, made the system tool-agnostic, and hardened the file system operations.

| Version | Name | What shipped |
|---------|------|-------------|
| 0.10.0 | The Hard Stop | Capture never auto-executes |
| 0.11.0 | The Delegate | Subagent dispatch for actions |
| 0.11.1 | The Safety Net | Subagent fallback for simpler tools |
| 0.12.0 | The Diet | 67-70% file size reduction, zero behavior change |
| 0.12.1 | The Passport Check | Removed hardcoded co-author metadata |
| 0.12.2 | The New Address | Upstream URLs updated to fork |
| 0.12.3 | The Time Traveler | Fixed changelog dates |
| 0.12.4 | The Right Address | Fixed ambiguous scan paths |
| 0.12.5 | The Deep Check | Content-aware duplicate detection |
| 0.12.6 | The Missed Spot | Fixed remaining bare path in duplicate scan |
| 0.12.7 | The Cold Start | First-run bootstrap guidance |

### Clarification & Review (v0.13.0 — v0.16.0)

Added code review, Open Questions system, clarification gates, non-blocking builder decisions, and complete question lifecycle.

| Version | Name | What shipped |
|---------|------|-------------|
| 0.13.0 | The Second Opinion | Post-work code review with 5 scoring dimensions |
| 0.14.0 | The Clarification Gate | Open Questions checkpoint, ambiguous-gap handling |
| 0.15.0 | The No-Block Build | Builder best-judgment, pending-answers follow-ups |
| 0.16.0 | The Full Loop | Verify resolves on spot, `do work answers` batch review, Builder Was Right path |

### Polish & Presentation (v0.17.0 — v0.20.0)

Renamed actions for clarity, expanded review to full acceptance testing, added client-facing presentation capabilities and institutional memory.

| Version | Name | What shipped |
|---------|------|-------------|
| 0.17.0 | The Name Tag | do.md → capture.md, consistent naming |
| 0.18.0 | The Clarity Pass | Action names clarified (capture requests, clarify questions), help menu |
| 0.19.0 | The Full Picture | Review → review work (requirements + acceptance testing + suggested testing), verify → verify requests, next-step suggestions |
| 0.19.1 | The Neighbor Check | Regression risk analysis, diff hygiene |
| 0.20.0 | The Pitch Deck | Present work action, Lessons Learned, knowledge-preserving diff hygiene |

---

## Cumulative Value Proposition

### What Was Built

A complete task management and delivery pipeline for AI-assisted development. Seven actions covering the full lifecycle from idea capture through client presentation:

1. **Structured capture** — natural language in, structured queue items out
2. **Intelligent triage** — simple/medium/complex routing adjusts the build pipeline
3. **Automated build** — plan, explore, implement, test per request
4. **Quality gates** — verify (capture quality) + review (code quality + acceptance testing)
5. **Clarification system** — Open Questions with defaults, non-blocking builder decisions, batch review
6. **Archive with memory** — self-contained UR folders, lessons learned, full traceability
7. **Client deliverables** — briefs, architecture diagrams, value propositions, video scripts

### Total Value Delivered

- **7 actions** covering capture → build → verify → review → present → cleanup → version
- **3-tier triage** (Route A/B/C) ensuring proportional effort per request
- **5-dimension code review** with requirements tracing and acceptance testing
- **Full lifecycle traceability** — every request has timestamps, triage decisions, implementation notes, test results, review scores, and lessons learned
- **Agent-agnostic design** — works with any AI coding tool that can read/write files
- **67-70% documentation reduction** (v0.12.0) with zero behavior change — lean and maintainable
- **20 releases** shipped over 5 weeks with consistent quality and changelog discipline

### Growth Opportunities

- **Dashboard/status view** — real-time queue status, archive stats, pending questions
- **Dependency-aware ordering** — process REQs in the right sequence based on declared dependencies
- **Parallel execution** — process independent REQs concurrently in multi-agent environments
- **Estimation calibration** — use triage + timestamp data to predict future implementation times
- **Template library** — reusable patterns for common request types
- **Cross-project portfolio** — aggregate work across multiple codebases

### Lessons Learned (Cross-Project)

- **Separation of concerns pays off.** Capture vs. execute, router vs. action, orchestrator vs. builder — every clean boundary prevented a class of bugs.
- **Markdown as infrastructure works.** No database, no server, no dependencies. Plain text files are readable, diffable, versionable, and portable.
- **Agent compatibility requires discipline.** Tool-specific APIs in action files break portability. The "design for the floor" principle (simplest agent that can read/write files) kept the skill working everywhere.
- **Quality gates should be automated and optional.** Mandatory gates frustrate users; absent gates produce bugs. Automated-but-skippable is the sweet spot.
- **Non-blocking decisions beat blocking questions.** The Open Questions system went through three iterations (blocking → recommended defaults → builder-decides-with-followup) before finding the right balance between thoroughness and velocity.
