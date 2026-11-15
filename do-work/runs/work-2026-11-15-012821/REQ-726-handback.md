# REQ-726 hand-back (record the decision to adopt no external review canvas)

## Commits
- 47c505d9 [REQ-726] Record the decision to adopt no external review canvas

## Files
- docs/decisions/records/adr-023-adopt-no-external-review-canvas.md — created
- docs/decisions/log.md — modified (new section `## [2026-11-15] No external review canvas for ai-report bundles`)
- docs/decisions/_master_index.md — modified (cluster line renamed to "Pipeline Deliverables", "3 ADRs (1 accepted, 2 superseded)", ADR-023 added to related pages, `updated:`)
- docs/decisions/topics/_index_pipeline-deliverables.md — modified (title and opening widened, ADR-023 listed as accepted, `related:` entry, two sources, `updated:`)

## Scope changes
none

## Red-green evidence
- Base (changes stashed): `REQ-726 probe FAIL: no docs/decisions/records/adr-*review-canvas*.md record` (exit 1)
- Tip 47c505d9: `REQ-726 probe PASS (adr-023-adopt-no-external-review-canvas)` (exit 0)

## Verification
- GREEN probe: PASS (above).
- `bash _dev/tests/contract-regressions.sh`: exit 0, "Contract regression checks passed." One run, no rerun needed.
- `bash _dev/tests/audit-lockins.sh`: exit 0. The two new topic-index sources are a tracked path (`do-work/archive/UR-172/REQ-721-...md`, confirmed with `git ls-files --error-unmatch`) and a spaced `CHANGELOG.md (0.317.0 ...)` entry.
- `git diff --stat main...HEAD -- skills/ suite/ ai-reports/ do-work/`: empty. `git diff --stat main...HEAD`: only the four docs/decisions files.
- Read-through: every Detailed Requirements fact is present (decision; ~5,000 px strip tile, folded sections 38 to 46%; 127.0.0.1 only, no phone; same-address tile writes feedback; second-port tile blocked read + one write, two HTTP 403 refusals, live only after two clicks and one page at a time; star/vote/note move feedback out of chat; relation-map page as the direction; two reopen conditions). Tinka's repository is named once. The 38-46% and 127.0.0.1 / 403 facts were cross-checked in both bundles' index.html.

## P-A-U
- PLAN: copy adr-022's shape (frontmatter, topic line, Context/Decision/Consequences), shorter; link per D-01 to D-03.
- APPLY: wrote the record and three index edits; one commit by exact path.
- UNIFY: `git diff --stat main...HEAD` shows _master_index.md, log.md, adr-023 record, _index_pipeline-deliverables.md only; nothing under skills/, suite/, ai-reports/, do-work/; no debug artifacts.

## Decisions
- D-05 (DECIDE & STATE): `created:`/`updated:` are 2026-11-14, the UTC build date (2026-11-14T23:33Z), per D-02; `decided:` stays 2026-11-15. The log section header uses the decision date 2026-11-15. If the integrator prefers local date, change three `updated:`/`created:` lines.
- D-06 (DECIDE & STATE): added a `## When to Reopen` section to the record rather than burying the two reopen conditions in Consequences; adr-022 has extra sections (Alternatives, Instances) so the format allows it. `related: []` because no existing ADR covers ai-report review.
- D-07 (DECIDE & STATE): the record quotes the maintainer's "let's abandon tinka" (from the REQ's Builder Guidance) as the decision's trigger.

## Discovered Tasks
- `docs/decisions/_master_index.md` cluster lines do not list ADR-017 to ADR-022 (D-03 gap); a backfill is needed.
- `docs/decisions/_progress.md` still says "Next ADR number: ADR-017" (D-04); it should say ADR-024 after this lands.

## Integration seams
Used ADR number 023. If another branch landed an `adr-023-*` first, renumber this record to the next free number and update the links in log.md, _master_index.md and _index_pipeline-deliverables.md plus the record's title line. No sibling in this run edits docs/decisions/.

## Proposed CHANGELOG entry
none (no shipped file changes). `docs/decisions` is `export-ignore` in `.gitattributes` and not a source in `suite/modules.tsv`.

## Proposed lesson
none

## Lessons read
none required (no matching satellite)
