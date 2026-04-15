# Master Index — Decisions

Top-level navigation for every Architecture Decision Record in this log. Entries are listed both by topic cluster and chronologically. See `README.md` for schema, `log.md` for the raw timeline.

**Total ADRs**: 10
**Status breakdown**: 10 accepted · 0 proposed · 0 superseded · 0 deprecated
**Date range**: 2026-02-01 → 2026-04-10

---

## By Topic Cluster

### Queue Model
How captured work is structured, how it flows, and where it lives on disk.
See [[topics/_index_queue-model]].

- [[adr/0002-ur-req-pairing]] — UR + REQ pairing is mandatory on every capture
- [[adr/0003-immutable-inflight-archived]] — `working/` and `archive/` are off-limits; follow-ups use `addendum_to`
- [[adr/0008-queue-canonical-path]] — pending REQs live at `do-work/queue/`, not `do-work/` root

### Platform Portability
The skill targets any agent that can read/write files and run shell commands — not a specific tool.
See [[topics/_index_platform-portability]].

- [[adr/0004-platform-agnostic-action-files]] — action files use generalized language; no tool-specific APIs
- [[adr/0005-subagent-dispatch-pattern]] — dispatch to subagents when available; fall back to direct read

### Routing & Dispatch
How a user utterance becomes an executed action.
See [[topics/_index_routing-dispatch]].

- [[adr/0006-priority-ordered-routing]] — SKILL.md routing table is numbered; first match wins

### Content Structure
How the skill's own source files are organized.
See [[topics/_index_content-structure]].

- [[adr/0007-crew-member-jit-loading]] — domain rules load just-in-time based on REQ frontmatter
- [[adr/0009-companion-reference-files]] — large actions split into primary + companion at the ~10k-token read limit

### Philosophy
The mental models the rest of the skill hangs on.
See [[topics/_index_philosophy]].

- [[adr/0001-capture-execute-boundary]] — capture hard-stops after writing files; user invokes work separately
- [[adr/0010-reqs-as-validated-intent]] — captured REQs are validated statements of user intent, not drafts

---

## By Date (Chronological)

| ADR | Date       | Version | Title                                  | Topic                |
|-----|------------|---------|----------------------------------------|----------------------|
| [[adr/0002-ur-req-pairing\|0002]]                | 2026-02-01 | 0.4.0  | UR + REQ pairing on every capture       | queue-model          |
| [[adr/0003-immutable-inflight-archived\|0003]]   | 2026-02-01 | 0.6.0  | Immutable working/ and archive/         | queue-model          |
| [[adr/0001-capture-execute-boundary\|0001]]      | 2026-02-03 | 0.10.0 | Capture ≠ Execute                       | philosophy           |
| [[adr/0004-platform-agnostic-action-files\|0004]]| 2026-02-03 | 0.8.0  | Platform-agnostic action files          | platform-portability |
| [[adr/0005-subagent-dispatch-pattern\|0005]]     | 2026-02-24 | 0.11.0 | Subagent dispatch with direct-read fallback | platform-portability |
| [[adr/0006-priority-ordered-routing\|0006]]      | 2026-02-04 | 0.9.1  | Priority-ordered routing table          | routing-dispatch     |
| [[adr/0007-crew-member-jit-loading\|0007]]       | 2026-04-07 | 0.50.0 | Crew members — JIT-loaded domain rules  | content-structure    |
| [[adr/0008-queue-canonical-path\|0008]]          | 2026-04-10 | 0.60.3 | Queue canonical path `do-work/queue/`   | queue-model          |
| [[adr/0009-companion-reference-files\|0009]]     | 2026-04-11 | 0.61.1 | Companion reference file pattern        | content-structure    |
| [[adr/0010-reqs-as-validated-intent\|0010]]      | 2026-04-08 | 0.51.3 | REQs as validated intent                | philosophy           |

---

## Status Legend

- **accepted** — in force today; behaviors downstream depend on it
- **proposed** — staged for a future release; not yet in force
- **superseded** — replaced by a newer ADR (`superseded_by:` points forward)
- **deprecated** — retired without replacement (rare)
