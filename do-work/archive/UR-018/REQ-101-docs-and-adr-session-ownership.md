---
id: REQ-101
title: Docs + ADR — multi-checkout guide and the session-ownership decision record
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T21:28:46Z
completed_at: 2026-08-04T21:36:00Z
commit: e452989
kb_status: pending
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-096, REQ-097, REQ-099]
maintenance: false
related: [REQ-096, REQ-097, REQ-099]
batch: parallel-building
write_set: [docs/work-guide.md, decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md, decisions/log.md, decisions/topics/_index_workflow-orchestration.md, decisions/_master_index.md]
---

# Docs + ADR — Multi-Checkout Guide and Session-Ownership Decision Record

## What

User-facing documentation for the new model and a decision record for the re-grain. No ADR currently covers session ownership at all — the 0.161.0 exclusive-session decision was recorded only as an AI report plus REQ-069.

## Detailed Requirements

- **`docs/work-guide.md`:** a "several checkouts against one queue" section — how to claim from a workspace/clone/cloud session, the `assigned_to` earmark, the one-releaser rule, what happens on a double claim (ordinary merge conflict, fixed at merge, `queue-kanban verify` finds duplicates), and the automatic wave. Update the guide's existing "one queue owner, one REQ at a time" summary (~line 91) to match the new contract.
- **ADR in `decisions/records/`** (next number after adr-017) recording: the checkout→queue re-grain and claim-anywhere/one-releaser model; why the reserve *verb/status* stays dead while the advisory *field* returns (ratchet + router budget); capture-anywhere with fix-at-merge instead of id sharding; the auto-wave contract change (supersedes "nothing computes the set"); the static writer label vs. the liveness-machinery ban; links to UR-012/REQ-069 (the decision being partially reversed), UR-018, and the two acceptance-run artifacts (REQ-095, REQ-100).
- **`decisions/log.md`:** one-line entry (its last entry is 2026-07-01 — the log is stale; just append, don't backfill).
- Cross-references by file path per convention; no CLAUDE.md/AGENTS.md citations in shipped files.

## Red-Green Proof

**RED prompt/case:** `docs/work-guide.md:91` still tells users "one queue owner, and by default one REQ at a time" with no multi-checkout path; `decisions/records/` has no session-ownership ADR.
**Why RED now:** The docs describe the 0.161.0–0.166.0 contract; the re-grain decision lives only in plan files and this UR.
**GREEN when:** The guide documents the multi-checkout workflow accurately against the shipped contract text, and the ADR exists with the five decision points and links above.
**Validation:** User confirmed (approved plan, Phase 4).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 4).

---
*Source: approved plan, Phase 4*

## Triage

**Route: B** - Medium

**Reasoning:** The three deliverables and the ADR's five decision points are enumerated in the REQ. What needed discovery was the ADR file conventions (frontmatter shape, topic cluster, the `[[wikilink]]` + relative-path double citation, the Alternatives/Consequences/References sections), which topic cluster session ownership belongs to, and — the part that mattered — that the ADR indexes enumerate ADRs, so a new record is invisible until two more files change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided documentation

*Skipped by work action*

## Exploration

**ADR conventions**, read from `decisions/records/adr-016-vendor-queue-kanban-into-the-skill.md` as the most recent full-length example: YAML frontmatter with `title`/`type`/`status`/`topic_cluster`/`decided`/`sources`/`related` (each related entry a `page` + `rel` pair)/`created`/`updated`/`confidence`; a body opening with the topic-cluster line and a *See also* line that repeats the related pages; then Context → Decision → Alternatives → Consequences → References. References cite both `[[wikilink]]` form and relative paths, and `../../` resolves from `decisions/records/` to the repo root.

**`decisions/` is `export-ignore`d** (`.gitattributes:50`), so it is *not* a shipped path — which means the no-CLAUDE.md-citation rule does not bind here (ADR-016 cites it directly). The rule still binds `docs/work-guide.md`, which is shipped; verified zero mentions after writing.

**Topic cluster:** `workflow-orchestration` — its stated scope is "how pending work is stored and how the pipeline coordinates queue processing", which is where who-may-claim belongs. The alternative, `skill-architecture`, is about how the skill is structured and distributed.

**The indexes are closed enumerations.** `decisions/topics/_index_workflow-orchestration.md` lists every ADR in the cluster twice (frontmatter `related:` and a body bullet), and `decisions/_master_index.md` lists cluster members plus a hand-maintained count. A new ADR that updates neither is unreachable by navigation. Found while checking: **the master index's counts were already stale** — its header claimed 15 in-force decisions against 17 records on disk (18 after this one), and it predated ADR-016 and ADR-017 entirely.

**No prior record to amend.** `grep -rl "exclusive.session" decisions/` returns nothing: the 0.161.0 decision this partially reverses exists only as `do-work/archive/UR-012/REQ-069-…`. So ADR-018 is the first record of session ownership in either direction, which shapes its Context section — it has to state the prior contract rather than link to a record of it.

## Scope

**Files I will touch:**
- `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md` (new) — the decision record
- `docs/work-guide.md` — a new "Several checkouts against one queue" section
- `decisions/log.md` — one appended entry
- `decisions/topics/_index_workflow-orchestration.md` — cluster membership (frontmatter + body bullet)
- `decisions/_master_index.md` — cluster row

**Acceptance criteria (restated from REQ):**
- [ ] Guide section covers claiming from a workspace/clone/cloud session, the `assigned_to` earmark, the one-releaser rule, double-claim behavior, `verify`'s duplicate detection, and the automatic wave
- [ ] The guide's existing "one queue owner, one REQ at a time" summary matches the new contract
- [ ] ADR numbered 018, in `decisions/records/`
- [ ] ADR records all five decision points (re-grain; verb dead / field returns; capture-anywhere with fix-at-merge; auto-wave supersedes "nothing computes the set"; writer label vs. liveness ban)
- [ ] ADR links UR-012/REQ-069, UR-018, and both acceptance-run artifacts
- [ ] `decisions/log.md` gets one appended line, no backfill
- [ ] Cross-references by file path; no `CLAUDE.md`/`AGENTS.md` citation in any shipped file

## Pre-Flight

- **WARN — baseline suite red before any change:** the same 8 `chmod 500`-versus-root failures inherited by every REQ in this batch.
- Working tree clean outside `do-work/` at claim time.

## Implementation Summary

**Files changed:**
- `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md` (new) — ADR-018. Context states the prior contract verbatim and why no record of it existed, plus the two things that made a naive re-grain dangerous (the committed checkpoint as recovery's input; the by-name liveness ban and its three-patch history). Decision carries **six** numbered points — the five the REQ named plus the checkpoint-stays-committed choice with the user's own reasoning. Alternatives runs to **seven**, each with its rejection reason, including the two that were nearly taken (keep the human confirmation gate; make `write_set` a scheduling input). Consequences names the costs plainly: the unspecified cases really are unspecified, the checkpoint is now a guaranteed conflict point, the label-less bullet is known-unsound and filed as REQ-104, agent behavior under concurrency is still unproven, and none of it exists on untracked-`do-work/` installs.
- `docs/work-guide.md` (modified) — new "Several checkouts against one queue" section before *Trigger aliases*: claiming from anywhere, the one-releaser rule, `assigned_to` with both override paths, the two conflict shapes and their resolutions (with the keep-both-sides rule and *why* each one-sided resolve loses data), capture collisions and the duplicate-id probe, the writer-label recovery rule, `--fan-out`, and a closing caveat that the whole section presumes a committed `do-work/`.
- `decisions/log.md` (modified) — one appended dated entry, no backfill of the stale gap.
- `decisions/topics/_index_workflow-orchestration.md` (modified) — ADR-018 in the frontmatter `related:` list and the body bullet list; cluster scope line extended to mention queue ownership; `updated` bumped.
- `decisions/_master_index.md` (modified) — ADR-018 added to the Workflow Orchestration row; the row's stale hand-maintained count replaced by a scope phrase, and the header's stale global count replaced with a pointer to read `records/` instead.

**What was done:** Wrote the decision record for the whole UR-018 re-grain, the user-facing multi-checkout guide, the log line, and the two index updates a new ADR requires to be reachable.

## Testing

### The guide is accurate against the shipped contract, not against the plan

Every load-bearing claim in the new section was checked against the text that implements it rather than against the approved plan:

```
$ grep -c 'one releaser per queue' actions/work-reference.md
1
$ grep -c 'claim held by' actions/work-reference.md
3
$ grep -c 'keeping every entry from both sides' actions/work-reference.md
1
$ grep -ci 'assigned to <name>' actions/work.md
1
$ grep -c 'fan-out' actions/work-reference.md actions/work.md
actions/work-reference.md:3
actions/work.md:7
```

The guide's conflict-shape claims (content conflict never a rename conflict; `CHECKPOINT.md` conflicts on every concurrent claim including disjoint ones; byte-identical claims leave the label as the only detector) are the observed results recorded in `do-work/archive/UR-018/REQ-095-two-clone-acceptance-run.md`, not predictions — which is why the section can state them flatly.

### No maintainer-doc citation in a shipped file

```
$ grep -c 'CLAUDE.md' docs/work-guide.md
0
```

Zero, as required. `decisions/` is `export-ignore`d (`.gitattributes:50`) and is not a shipped path, so the ADR's `CLAUDE.md` reference in its `sources:` list is the same allowance ADR-016 already uses.

### Every ADR reference resolves

```
$ cd decisions/records && for p in <every relative path in the References section>; do [ -e "$p" ] && echo OK || echo MISSING; done
OK  ../../actions/work-reference.md
OK  ../../actions/work.md
OK  ../../docs/work-guide.md
OK  ../../tools/queue-kanban/model.go
OK  ../../tools/queue-kanban/verify.go
OK  ../../do-work/archive/UR-012/REQ-069-exclusive-session-model-removes-concurrency-machinery.md
OK  ../../do-work/user-requests/UR-018/input.md
OK  ../../do-work/user-requests/UR-018/assets/approved-plan.md
OK  ../../do-work/archive/UR-018/REQ-095-two-clone-acceptance-run.md
OK  ../../do-work/archive/UR-018/REQ-100-live-wave-acceptance-run.md
OK  ../topics/_index_workflow-orchestration.md
```

Eleven for eleven, plus `do-work/queue/REQ-104-labelless-entry-authorship-heuristic.md` cited in Consequences, confirmed present. Broken relative links in a `decisions/` file are silent — nothing compiles them — so this check is the only thing standing between a written reference and a dangling one.

### ADR numbering and reachability

```
$ ls decisions/records/*.md | wc -l
18
$ grep -c adr-018 decisions/_master_index.md decisions/topics/_index_workflow-orchestration.md decisions/log.md
decisions/_master_index.md:1
decisions/topics/_index_workflow-orchestration.md:2
decisions/log.md:1
```

018 is the next number after 017. It appears in the topic index twice (frontmatter `related:` and body bullet, matching the file's own convention), once in the master index row, and once in the log — so it is reachable from every navigation entry point rather than only existing on disk.

### Contract suite

```
$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -c '^FAIL'
8
```

The pre-existing eight. Nothing in `docs/` or `decisions/` tripped a new assertion, including the shipped-path citation guard and the docs-exemption allowlist.

## Lessons Learned

**What worked:**
- Writing the ADR's Consequences from the batch's *findings* rather than from its plan. Four of the five costs listed (guaranteed checkpoint conflict, the unsound label-less bullet, unproven agent concurrency, nothing travels on untracked installs) are things REQ-095 and REQ-100 discovered; none of them were foreseeable when the plan was approved, and an ADR that omitted them would read as more confident than the evidence supports.
- Recording the two *nearly-taken* alternatives (keep the confirmation gate; schedule on `write_set`) with real reasons rather than listing only the obviously-wrong ones. A rejected alternative that nobody was tempted by teaches a later reader nothing.
- Checking every relative path in the References section by loop. Nothing validates links in `decisions/`, so a dangling reference is permanent and silent.

**What didn't:**
- Nearly shipping the ADR without touching either index. `decisions/records/` is not a scanned directory — the topic index and master index are hand-maintained lists, so a record absent from them is unreachable by navigation even though the file exists. Caught by grepping the indexes for `adr-017` to see whether *it* was listed.
- The master index's counts turned out to be **three behind** (header claimed 15 in-force decisions against 18 records, and the cluster row's "4 ADRs + 1 declined" predated nothing but was about to). Adding another hand-maintained number would have made the same failure again on ADR-019, so both counts were replaced with a scope phrase and a "read `records/`" pointer.

**Worth knowing:**
- `decisions/` is `export-ignore`d and therefore not a shipped path: it may cite `CLAUDE.md`, and it never ships to consumers. `docs/` **is** shipped and may not.
- A new ADR is **five** file changes, not one: the record, the topic index (twice — frontmatter and body), the master index row, and `decisions/log.md`.
- `decisions/log.md` was stale from 2026-07-01 to this entry. The REQ said append, don't backfill, which is right — a hand-written history reconstructed after the fact is worth less than an honest gap.

## Orientation

The multi-checkout model now has both halves of its documentation: a user-facing walkthrough in `docs/work-guide.md` covering how to claim from any checkout, how to earmark work, what the two conflict shapes look like and how to resolve them without losing a claim, and how `--fan-out` behaves — and `decisions/records/adr-018-…`, the first decision record this skill has ever had for session ownership, recording six decisions, seven rejected alternatives, and the costs the acceptance runs uncovered. `[MAP CHANGED]` — a new topic-cluster member and the first ADR in the repo about who owns a queue, so the workflow-orchestration cluster now covers ownership as well as queue mechanics. `prime_files` is empty; no prime went stale.

## Qualification

**Passed** — 5 files verified, 7 acceptance criteria traced, and the two claims most likely to be wrong checked mechanically rather than by reading: every ADR relative link resolved by loop, and every load-bearing guide claim grepped against the shipped text that implements it.

- **Substantive:** the ADR is a full-length record in the established format (Context / Decision / Alternatives / Consequences / References), not a stub; the guide section is a complete workflow walkthrough with its caveat stated.
- **Requirements traced:** all seven criteria map to a file or a `## Testing` check.
- **Accurate, not aspirational:** the guide describes shipped behavior — verified by grepping the implementing text — and the ADR's Consequences describe found costs, including one filed defect and one unproven claim.

## Review

**Overall: 94%** | 2026-08-04T21:36:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | n/a (documentation) |
| Test Adequacy | 90% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 3 minor
**Acceptance:** Pass — the guide documents the multi-checkout workflow accurately against the shipped contract, ADR-018 exists with all five required decision points (plus a sixth) and every required link, the log has its one appended line, and the record is reachable from every index.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Requirements checklist:** all seven `## Scope` criteria delivered; evidence per criterion in `## Testing`.

**Minor:**
- Scope grew by two files beyond the REQ's `write_set` (`decisions/topics/_index_workflow-orchestration.md`, `decisions/_master_index.md`), both declared before editing. Neither is optional: the indexes are hand-maintained lists, so an ADR absent from them is unreachable by navigation. The REQ's `write_set` said `decisions/records/*` and `decisions/log.md`, which is the shape of the omission — it assumed a new record is self-registering.
- The master index's stale counts were corrected in passing rather than filed. Two numbers in one line and one header sentence, in a file this REQ had to edit anyway; adding a third hand-maintained number instead would have re-created the exact failure on the next ADR. Replacing them with a "read `records/`" pointer is the durable fix, and it is the Closed-Enumerations rule applied to the indexes themselves.
- **Test Adequacy is 90% because documentation accuracy is verified by grep, not by execution.** Each guide claim was matched against the shipped text that implements it, and every ADR link resolved — but nothing here proves a *reader* follows the guide successfully, and the section's central instruction ("keep every entry from both sides") is exactly the kind of thing a hurried resolver skips. The one mitigation that exists is that the same rule is stated in `actions/work-reference.md`, where an agent resolving the conflict will meet it.

**Scope drift:** none against the declaration — five files declared, five touched.

**Restatement sweep (MUST):** run. This REQ *documents* rather than redefines, so the sweep ran in the reverse direction: every claim in the new guide section was checked against the shipped text it restates, since a guide that paraphrases a contract loosely is a stale restatement the moment either drifts. Two paraphrases were tightened to the contract's own wording (`one releaser per queue`, `claim held by …`). The sweep also caught what a *new ADR* makes stale — the two indexes and their counts — which is the same failure class one step out.

**Suggested additional testing:**
- After REQ-104 lands (the unsound label-less bullet), re-read the guide's crash-recovery paragraph: it currently describes the four-case classification correctly, and narrowing or dropping the label-less case will make one sentence of it wrong.
- A fresh reader walking the new section against a real second checkout, to catch the step it assumes rather than states. Every claim in it is verified true; none is verified *sufficient*.

*Reviewed by review-work action (pipeline mode, in-session)*
