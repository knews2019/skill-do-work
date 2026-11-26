---
id: REQ-760
title: 'Scan skill instructions for new maintenance opportunities'
status: completed
dispatch_at: 2026-11-25T23:29:58Z
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-11-25T23:27:30Z
  basis:
  - Route B
  - 3-file write set
  - 2 subsystems involved
  - 5 acceptance criteria
route: B
write_set: [ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/gallery-design.md, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.css, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.js, docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png]
created_at: 2026-11-25T22:57:02Z
user_request: UR-190
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
batch: skill-instruction-maintenance
related: [REQ-758, REQ-759]
claimed_at: 2026-11-25T23:27:19Z
builder_handback_at: 2026-11-25T23:46:54Z
integration_at: 2026-11-25T23:47:19Z
review_at: 2026-11-25T23:57:24Z
kb_status: pending
kb_entry: ''
commit: e791e736f261f5dc06da7b8553a304bc37ceacb2
heavy_verified_at: 2026-11-26T00:06:09Z
heavy_verified_revision: e791e736f261f5dc06da7b8553a304bc37ceacb2
completed_at: 2026-11-26T00:06:13Z
---
# Scan skill instructions for new maintenance opportunities
## What
Perform a fresh, read-only maintenance review of the do-work skill instructions. Trace one straightforward request through capture, verification, implementation, testing, and review. Present the findings in a short, ranked `do-work-toolbox ai-report` covering improvements beyond the five findings reviewed at `efcd3618`.

## Review brief
“Tell the model what only I or this project can tell it. Let it exercise judgment on the rest.” Look for an instruction an experienced engineer can infer, two instructions that conflict or could both be deleted, or a rule that adds a question, test run, or write-up without changing the outcome. Check what broke before a rule was added and what a second check actually caught.

For each new candidate give the current file and line, the exact deletion or narrowing, what could go wrong, and a recorded failure if one exists. Distinguish an observed incident from a plausible risk. Preserve the reason the user asked, consequential decisions, and evidence the work succeeded. For capture quality, use corrections to already-stated requirements and avoidable interruptions as the signals. Later changes of preference are separate.

## Constraints
The prior five topics are the baseline: passing mentions and repeated capture checks, routine verification and scores, factual-question formatting, the line-count heuristic, and reuse of a focused check at review. The requested third item is a new scan. Report any fresh findings for the maintainer to decide on. Do not apply the proposed cuts or queue them during the scan. If nothing additional earns a recommendation, say so. If command usage matters, people run `just run-kanban` and `just do-work-update` in their own project after installing the skill.

## Full Context
See `do-work/user-requests/UR-190/input.md` for the verbatim capture instruction.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the action-file prime and selected report procedures. Trace REQ-318 across the original input, phase commits and archive. Check historical evidence for both proposed cuts.
- [x] **[APPLY]:** Author only the report bundle and its preview within Scope. Preserve the user's reason for the ordering change, the builder's choices and the recorded acceptance evidence.
- [x] **[UNIFY]:** Review the five-file merge diff and design note, compare both presentation assets with the kit, check JS syntax and cited action lines, and view all four integrated Chrome renders plus the PNG preview.

*Builder evidence from the handback, checked against the integrated merge and recorded by the integrator. The handback predated takeover, so no builder-work duration was recorded.*
*Source: User request, third item. but also capture in 3 a new scan for new improvements*

## Triage

Route B. A fresh audit needs exploration of historical failures and the current workflow.

## Plan

**Planning not required** — Route B, exploration-guided implementation.

## Exploration

The capture and verify actions define the original-clause comparison and repair steps. The work testing and review companions govern the later handoffs. The archive includes small direct-route requests with documented acceptance evidence, including REQ-318's Timeline ordering. The five preceding audit topics are excluded by the request. The toolbox proposal report can carry ranked recommendations, source-line evidence and historical receipts without changing the source instructions.

## Scope

**Files I will touch:**
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html` (new) — local proposal.
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/gallery-design.md` (new) — design note.
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.css` (new) — local presentation stylesheet.
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.js` (new) — local presentation script.
- `docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png` (new) — preview of the rendered report.

**Acceptance criteria:**
- Trace one straightforward request from capture through review and distinguish recorded failures from hypothetical risks.
- Rank only new, evidenced recommendations, or report that none qualify. Give current citations, exact cuts and consequences, preserving the original reason, decisions and success evidence.
- Deliver and render the report. Leave the recommendations for the maintainer to decide.

## Pre-Flight

- The builder's Route B `advance` preflight ran in the private checkout at dispatch revision `77c78cd721c24f0c49ed4f1901cd656a60575060`. The checkout was clean before authoring. The request-bound preflight itself ran after drafting, once the shared receipt was available. Its only warning, `PREFLIGHT-DIRTY`, named exactly the five report outputs. The shipped instruction tree used by the focused check was unchanged.
- Both typed gate records matched REQ-760 and `do-work/working/REQ-760-scan-skill-instructions-for-new-maintenance-opportunities.md`. The preflight was satisfied. Its private baseline recorded `bash _dev/tests/shipped-package-reference-contract.sh`, exit 0, `launched: true`.
- The direct canonical argv `["bash","_dev/tests/maintainer-verify.sh"]` had a satisfied green record with `exact_revision_match` at the dispatch revision, exit 0, digest `25aef6a61b17974fce3a08a63953c3eb1719dde37ce3ab556f066e3619e6bbb5`. The REQ-758 integration receipt identifies the shared 182-second run. This is inherited dispatch-revision evidence. Integrated testing uses the real merge range and a fresh canonical gate.

## Implementation Summary

**Files changed:**
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html` (new)
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/gallery-design.md` (new)
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.css` (new)
- `ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.js` (new)
- `docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png` (new)

**What was done:** Created a local proposal with two ranked instruction edits and an option to defer. The report traces REQ-318 from capture through review and links the observed outcomes to exact action lines and archived evidence. The edits are proposals for the maintainer.

**Provenance:** `65f68e0de0ba6f84d8385e7ca951a58cbcff128a..e791e736f261f5dc06da7b8553a304bc37ceacb2`. The merge's second parent is builder commit `1911eeeceef7cc6bd19855ce6127469cf719a71f`, based on `77c78cd721c24f0c49ed4f1901cd656a60575060`. The pre-merge queue-path guard was empty.

## Decisions

- **D-01 — Use the REQ-318 Timeline ordering case and Decision Matrix.** The separate capture, verification, implementation, and review commits make the trace independently checkable. Identical benefit, cost, and risk fields let the maintainer compare both edits and deferral. The report limits impact claims to observed outcomes and leaves time savings unmeasured.

- **D-02 — List the four bundle children individually in Scope and `write_set`.** The typed scope gate compares literal file paths. The original directory scope already covered these four files, and the committed preview retains its exact original path.

## Qualification

- The request-bound `qualify` and `scope-drift` records were satisfied against `65f68e0de0ba6f84d8385e7ca951a58cbcff128a..e791e736f261f5dc06da7b8553a304bc37ceacb2`. The first Scope comparison rejected the directory shorthand. D-02 records the correction to the same four files plus the declared preview.
- Read all five paths in the integrated range. Git diff hygiene, kit asset equality and JS syntax passed. The report and preview are usable locally. There were no shipped changes or debug artifacts.

## Testing

- **Focused:** The Route B test gate ran `bash _dev/tests/shipped-package-reference-contract.sh` against the merged checkout twice, including when recording the green canonical result. Both probes passed, exit 0 with no timeout. The request-specific baseline copy reproduced the builder's three reported fields for the same command. The baseline and both integrated probes were green.
- **Canonical:** `bash _dev/tests/maintainer-verify.sh` ran directly at `e791e736f261f5dc06da7b8553a304bc37ceacb2`. The first run exited 1 in the do-work-cli Go test stage. The required same-argv retry exited 0 in 174 seconds. Both expensive Go stages executed on the passing run. `advance` recorded the successful direct exit at the exact merge revision (argv digest `25aef6a61b17974fce3a08a63953c3eb1719dde37ce3ab556f066e3619e6bbb5`).
- **Report:** `ai-report-judge` passed with no findings in Google Chrome 158 at 1440 × 1000 and 390 × 844, light and dark. All four captures were inspected and had no horizontal overflow. Both assets match the shipped kit. JS syntax, all local references, unique IDs, and gallery metadata passed. The committed 1040 × 740 preview opened and matched a fresh capture of the merged report byte for byte.
- Historical Timeline measurements and test outcomes were checked against the archive and commits. The report proposes documentation edits only, so it adds no behavioral RED/GREEN test.

**Heavy verification plan:**
- Range: `65f68e0de0ba6f84d8385e7ca951a58cbcff128a..e791e736f261f5dc06da7b8553a304bc37ceacb2`
- Shared selection reason: coverage is uncertain for: ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.css, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.js, docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png.
- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Review

**Overall: 100%** | 2026-11-25T23:57:24Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** None.

**Acceptance:** Pass — verified the local report and its integrated preview. The report ranks two candidates beyond the five excluded subjects and gives each an exact change, evidence, and potential consequences. REQ-318 distinguishes the verifier's capture repair and sibling clarification from the builder's choices and the independent review's tie-break test finding.

The report's UR-190 citation gives its capture-time path at the dispatch revision. The current path after this UR closes is `do-work/archive/UR-190/input.md`.

The integrator compared the four cited action files to the dispatch revision, read their cited lines in the merge, and checked the opted-in verifier wording at line 85. The REQ-254 correction, REQ-264 gate warning, REQ-583 verified handback, and REQ-739 pending result match the archived sources. The original case's Timeline measurements are archival evidence. Chrome rendered both widths and color modes, and a fresh gallery capture matched the merged preview byte for byte.

**Restatement sweep:** Nothing redefined. The report proposes edits for the maintainer. REQ-761 was still in working, so REQ-760 was not the wave's last integration.

**Suggested testing:** None for the local artifact. The catalog check reports that the bundle awaits indexing.

**Follow-ups created:** None. The two candidate edits remain in the report.

*Reviewed by the serial integrator against the original UR and the exact merge range.*

## Lessons Learned

**What worked:**
- Pair action-line citations with archived incidents and revision diffs. REQ-254 documents a failed checkbox copy, REQ-583 documents a verified copy, and REQ-739 proves the pending KB outcome. Each source supports a different claim.

**What did not work:**
- Route B Scope originally used a bundle directory while the typed comparator expected literal paths. Listing its already-declared children cleared qualification without changing the deliverables.

**Worth knowing:**
- The builder's routed preflight was observed after drafting at the dispatch revision. The inherited baseline and the fresh main merge gate describe two different revisions.

## Orientation

The action-file maintenance review now has a local decision brief with two proposals and an archived request that shows what verification and review caught.

## Heavy Verification Plan

- Range: `65f68e0de0ba6f84d8385e7ca951a58cbcff128a..e791e736f261f5dc06da7b8553a304bc37ceacb2`
- Shared selection reason: coverage is uncertain for: ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.css, ai-reports/2026-11-26_0128_proposal-instruction-maintenance/presentation.js, docs/report-gallery/previews/ai-reports/2026-11-26_0128_proposal-instruction-maintenance/index.html.png.
- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Heavy Verification Result

- Target and execution revision: `e791e736f261f5dc06da7b8553a304bc37ceacb2`.
- The first CLI integrations run exited 1 after 95 seconds. `TestCaptureCollidesWithQueueKanbanReservationForSameNumber` received Go dependency-download notices before the expected number. The isolated test passed after the module cache warmed (0.748 seconds). The same typed CLI integration lane passed on retry.
- `queue-kanban-javascript`: executed, exit 0, 8 seconds.
- `queue-kanban-browser`: executed, exit 0, 90 seconds.
- `do-work-cli-integrations`: executed, exit 0, 71 seconds.
- `staged-skills`: executed, exit 0, 41 seconds.
- `updater`: executed, exit 0, 70 seconds.
- `installer`: executed, exit 0, 34 seconds.
