# Decisions Wiki — Build Progress

Working file. Tracks what's been created so a new session can resume mid-build without losing context.

**Branch**: `claude/create-decision-log-omF82`
**Started**: 2026-04-15
**Source of truth**: CHANGELOG.md v0.1.0 → v0.64.1

---

## Resume Prompt

If a session ends before this file is deleted, hand the next session this prompt:

> Resume the decisions wiki build on branch `claude/create-decision-log-omF82`. Read `decisions/_progress.md` first — it lists what's been created and what's left. The wiki layout, ADR schema, and topic cluster mapping are documented in `decisions/README.md`. The 10 planned ADRs, their titles, and which version introduced each are in `decisions/_master_index.md`. Source material for each ADR is in `CHANGELOG.md` — find the version cited in the progress file and read 10–40 lines around it. When each ADR is done, tick its checkbox in this progress file, commit with a message like `Add ADR-NNNN (title)`, and move on. After all ADRs land: (1) update `CLAUDE.md` Project Structure to register `decisions/`, (2) bump `actions/version.md` to 0.65.0 and add a CHANGELOG entry with codename "The Paper Trail", (3) delete this `_progress.md` file, (4) commit, (5) push the branch. Do not create a pull request.

---

## Layout

```
decisions/
├── README.md                    # Schema + how-to
├── _master_index.md             # Nav — by topic + chronological
├── _progress.md                 # THIS FILE — delete when complete
├── log.md                       # Append-only timeline
├── topics/
│   ├── _index_queue-model.md
│   ├── _index_platform-portability.md
│   ├── _index_routing-dispatch.md
│   ├── _index_content-structure.md
│   └── _index_philosophy.md
└── adr/
    ├── 0001-capture-execute-boundary.md
    ├── 0002-ur-req-pairing.md
    ├── 0003-immutable-inflight-archived.md
    ├── 0004-platform-agnostic-action-files.md
    ├── 0005-subagent-dispatch-pattern.md
    ├── 0006-priority-ordered-routing.md
    ├── 0007-crew-member-jit-loading.md
    ├── 0008-queue-canonical-path.md
    ├── 0009-companion-reference-files.md
    └── 0010-reqs-as-validated-intent.md
```

## Checklist

### Scaffolding
- [x] `decisions/README.md` — schema + how-to-add
- [x] `decisions/_master_index.md` — topic + chronological nav
- [x] `decisions/log.md` — timeline seed
- [x] `decisions/_progress.md` — this file

### Topic Indexes
- [x] `decisions/topics/_index_queue-model.md` (0002, 0003, 0008)
- [x] `decisions/topics/_index_platform-portability.md` (0004, 0005)
- [x] `decisions/topics/_index_routing-dispatch.md` (0006)
- [x] `decisions/topics/_index_content-structure.md` (0007, 0009)
- [x] `decisions/topics/_index_philosophy.md` (0001, 0010)

### ADR Pages

| Id | Slug | Version | Topic | Done |
|----|------|---------|-------|------|
| 0001 | capture-execute-boundary | 0.10.0 | philosophy | [x] |
| 0002 | ur-req-pairing | 0.4.0 / 0.8.0 | queue-model | [x] |
| 0003 | immutable-inflight-archived | 0.6.0 | queue-model | [x] |
| 0004 | platform-agnostic-action-files | 0.8.0 / 0.11.1 / 0.12.1 | platform-portability | [ ] |
| 0005 | subagent-dispatch-pattern | 0.11.0 / 0.11.1 | platform-portability | [ ] |
| 0006 | priority-ordered-routing | 0.9.1 | routing-dispatch | [ ] |
| 0007 | crew-member-jit-loading | 0.50.0 | content-structure | [ ] |
| 0008 | queue-canonical-path | 0.60.3 | queue-model | [ ] |
| 0009 | companion-reference-files | 0.61.1 / 0.64.1 | content-structure | [ ] |
| 0010 | reqs-as-validated-intent | 0.51.3 | philosophy | [ ] |

### Finalization
- [ ] Register `decisions/` in `CLAUDE.md` Project Structure
- [ ] Bump `actions/version.md` to 0.65.0
- [ ] Add CHANGELOG entry — `0.65.0 — The Paper Trail`
- [ ] Delete this `_progress.md`
- [ ] Final commit
- [ ] Push branch `claude/create-decision-log-omF82`

---

## Commit Cadence

One commit per logical unit so the branch is always resumable:

- [x] Commit 1: scaffolding (README, master index, log, progress file, queue-model index)
- [x] Commit 2: remaining 4 topic indexes
- [x] Commit 3: ADRs 0001–0003
- [ ] Commit 4: ADRs 0004–0006
- [ ] Commit 5: ADRs 0007–0010
- [ ] Commit 6: CLAUDE.md + version bump + CHANGELOG, delete this file

Version bump and CHANGELOG entry happen in the **final** commit only — the interim WIP commits on this feature branch don't each need a version bump.

## Notes

- ADR schema and relationship vocabulary live in `decisions/README.md`. Don't invent new `rel:` types without updating the README first.
- Each ADR's **References** section must cite the CHANGELOG version(s) where the decision landed, plus the primary action file(s) the decision is enforced in.
- When citing relationships between ADRs, frontmatter `related:` links are bidirectional — if 0001 relates to 0010, both files need the link.
