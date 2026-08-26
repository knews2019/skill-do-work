---
id: UR-074
title: 'Never show a bare REQ or UR number — autocomplete the title on display and on copy'
created_at: 2026-08-26T13:02:24Z
requests: [REQ-374, REQ-375, REQ-376]
word_count: 2990
---

# Never show a bare REQ or UR number — autocomplete the title on display and on copy

## Summary

A REQ body that reads `Read prime-cms.md, REQ-1679/REQ-1108 lessons` forces the reader to go
hunting: the numbers carry no meaning on their own. The same is true of a `UR-389` in the drawer
meta and of the bare ids in `Depends on`, `Blocked by`, `Unblocks` and `Overlapping write sets`.
On the board those ids are already clickable links, but the link text is the id alone; on copy
even the click is gone, because the clipboard receives the file's raw bytes.

The ask: wherever a REQ or UR id is **rendered** or **copied**, attach its title — first mention
expanded inline, plus a glossary of everything the document referenced. The REQ files on disk are
never rewritten; this is a presentation layer built from data the board already holds
(`requestsById` / `userRequestsById` carry a `title` on every record, and the board loads
`queue/`, `working/` and `archive/`, so an id as old as REQ-1108 still resolves).

The user also asked what else could improve here, and asked that new cross-references be written
self-describing in the first place rather than only patched at render time.

## Extracted Requests

| REQ | Title | Covers |
|---|---|---|
| REQ-374 | Title-bearing ticket links and a drawer glossary | "on display… it should always autocomplete the title of the detected request"; "a glossary and an inline mention when they appear for the first time"; "true for user requests and… for requests as well" |
| REQ-375 | Copy carries titles and a referenced-requests glossary | "even on copy… it should always autocomplete"; "don't change the rec file" |
| REQ-376 | Cross-Reference Convention for newly authored ids | "we shouldn't assume that people know what those numbers mean" applied at authoring time, so new references are self-describing outside the board too |

## Decisions Taken During Capture

Answered by the user through the planning questions, so these are validated requirements rather
than open questions:

- **D1 — Display shape:** expand the **first** mention of each id inline; later mentions stay bare.
  (Rejected: hover-tooltip-only, and every-mention expansion.)
- **D2 — Copy shape:** inline expansion **and** a glossary appendix. (Rejected: glossary alone.)
- **D3 — Scope:** board display and copy, **plus** an author-side convention for newly written
  cross-references. Existing REQ files are never rewritten.
- **D4 — Both id kinds:** the rule covers UR ids exactly as it covers REQ ids.

## Batch Constraints

- **Never rewrite a REQ file to add titles.** The user was explicit: "don't change the rec file".
  Every expansion is computed at render or copy time from board data.
- **Frontmatter stays bare ids.** `depends_on`, `related`, `addendum_to` and `user_request` are
  parsed, not read; annotating them would break YAML parsing and the board's own dependency graph.
- **One resolver, two consumers.** Display and copy must share a single id→title resolution path;
  two copies would drift.
- **Never guess.** An id that does not resolve to a board record stays bare and stays out of the
  glossary — the same posture as the existing ambiguous-REQ-segment rule in `board-detail.js`.
- **Accepted cost, stated up front:** `board-clipboard.js` currently documents "verbatim has to
  mean verbatim" for the clipboard payload. Inline annotation ends that contract. The break is
  client-side only — the Go payload keeps shipping exact file bytes and its two round-trip tests
  stay green — and the frontmatter fence is never touched, so a paste still parses as a REQ file.
  REQ-375 rewrites the stale comment in the same change.

## Adjacent Improvements Raised, Deliberately Not Captured

The user asked "in what other ways can we improve this". Four ideas were raised in the planning
reply and the user did not ask for any of them, so none became a REQ. Recorded here so the answer
is not lost with the session:

- **A dead cross-reference is silent today.** File paths the Go build cannot find render in the
  blocked accent (`.repo-file-missing`); an id naming a REQ that does not exist renders as ordinary
  prose. The same treatment would surface typos and references to never-captured work.
- **Card-level titles.** The `Depends on` count badge on a card carries no hover title — nearly free
  once `ticketTitleFor` exists.
- **Search by referenced id.** The board filter box matches titles only. Matching "cites REQ-1679"
  would answer "what depends on this?" without opening anything.
- **The `[impact-rule-change]` title tag is noise inside an expanded inline mention.** Stripping a
  leading `[token] ` would read better, but that is a judgment call about the title convention
  rather than about linkification.

## Full Verbatim Input

> ````text
> [Screenshot: the do-work Kanban board at http://127.0.0.1:8090, project glw-game-find-the-difference,
> generated 2026-08-26 12:10 UTC. Columns Pending (3), Claimed (1), Needs input · Blocked (0),
> Recently done (18). The REQ-1685 detail drawer is open on the right showing its frontmatter meta
> table (status claimed, domain cms, user request UR-389, a long write set, route C, impact
> impact-rule-change, effort effort-substantive) and the start of its body, including the
> "AI Execution State (P-A-U Loop)" PLAN line reading "Read prime-cms.md, REQ-1679/REQ-1108 lessons,
> and the implementation guardrails."]
> 
> See below the request copied into clipboard and I just pasted it here. So one of the problems that I have with this is that anytime that I see request 1679, request 1108, I have to go spelunking and actually try to figure out what that is about. and this wouldn't necessarily have to be so so even on copy and on display before i'm going to a ticket it should always autocomplete just for rendering and copy sake don't change the rec file but it should always autocomplete the title of the detected request because otherwise these numbers become too cryptic and unless I actually go and hunt them down I wouldn't know what they are about but with the title I can get some idea. You're in plan mode let me know if this is good, feasible or if you have any other better solutions. Now the better solution doesn't have to be limited to this one so just go ahead and ask me if you have any other questions and explain to me how in what other ways can we improve this
> 
> 
> ---
> id: REQ-1685
> title: '[impact-rule-change] Close the preview-to-apply drift window across the four destructive card-mutation flows'
> status: claimed
> claimed_at: 2026-08-26T11:09:38Z
> route: C
> created_at: 2026-08-26T09:15:15Z
> user_request: UR-389
> domain: cms
> prime_files: [docs/prime-files/prime-cms.md]
> tdd: true
> suggested_spec:
> depends_on: []
> maintenance: false
> impact: impact-rule-change
> effort_estimate: effort-substantive
> estimate:
>   p50_active_minutes: 55
>   confidence: low
>   calculated_at: 2026-08-26T11:09:38Z
>   basis:
>     - Route C
>     - 8-file write set
>     - 3 subsystems involved
>     - 3 acceptance criteria
>     - persistence changes
>     - async lifecycle behavior
>     - cross-route regression gates
>     - full-suite verification
> write_set:
>   - server/modules/cms/publishing/cms-card-delete.js
>   - server/modules/cms/publishing/cms-unused-replacement-delete.js
>   - server/modules/cms/publishing/cms-takedown-flow.js
>   - server/modules/cms/publishing/cms-replacement-build-flow.js
>   - server/modules/cms/ingestion/photo-level-map.js
>   - server/modules/cms/ingestion/dedup-index.js
>   - server/modules/cms/core/cms-store.js
>   - server/modules/cms/core/cms-catalogue-store.js
>   - server/modules/cms/core/cms-sqlite-store.js
>   - server/modules/cms/core/cms-dev-feed-request.js
>   - server/modules/cms/core/cms-catalogue-commands-lifecycle.js
>   - server/modules/cms/ingestion/cms-catalogue-duplicate.js
>   - server/modules/cms/ingestion/cms-hand-form-ingest.js
>   - server/modules/cms/ingestion/ingestion-gate.js
>   - docs/prime-files/prime-cms.md
>   - tests/unit/cms-production/cms-card-delete.test.js
>   - tests/integration/cms-lifecycle/cms-unused-replacement-delete.test.js
>   - tests/integration/cms-lifecycle/cms-takedown-flow.test.js
>   - tests/unit/cms-production/cms-replacement-build-flow.test.js
>   - tests/unit/cms-core/cms-catalogue-transaction.test.js
>   - tests/unit/cms-core/cms-sqlite-store.test.js
> ---
> 
> # Close The Preview-To-Apply Drift Window Across The Four Destructive Card-Mutation Flows
> 
> ## What
> 
> Every multi-step destructive card flow takes a **flow-scoped** mutation lock that no per-record store
> shares, so a competing write can land between the in-lock re-preview and the destructive steps and be
> applied against state the admin never previewed. Decide the remedy — compare-and-swap revalidation
> immediately before each destructive write, versus shared/ordered locks — **once, for all four flows**,
> and implement it uniformly.
> 
> ## AI Execution State (P-A-U Loop)
> - [x] **[PLAN]:** Read `prime-cms.md`, REQ-1679/REQ-1108 lessons, and the implementation guardrails. Use exact store-local CAS: capture the affected file/SQLite rows during each flow's in-lock preview, pass those expectations to the real mutation primitive, and compare under that primitive's existing file lock or `BEGIN IMMEDIATE`. Optional expectations preserve unconditional legacy callers; CAS mismatches translate to each flow's existing drift refusal. No coordinator, shared lock, new lock nesting, or SQLite-to-ingestion import.
> - [x] **[APPLY]:** Wrote focused RED interleaving/primitive tests first, then added exact store-local CAS at the existing file locks and SQLite `BEGIN IMMEDIATE` boundaries. Wired the expectations through card deletion, unused-replacement deletion, takedown resolution, and replacement-build handoff; translated store drift to each flow's existing refusal while preserving omitted-argument behavior and committed-projection success semantics.
> - [x] **[UNIFY]:** Reviewed every scoped source, test, and authority-doc diff. The six focused suites pass (206 tests), `npm run lint` passes with zero warnings, and the full Jest suite passes (435 suites / 6500 tests; 9 suites / 40 tests skipped). `git diff --check` is clean and added lines contain no `console.log`, `debugger`, `TODO`, or `FIXME` artifacts.
> 
> ## Why
> 
> Raised as a P1 by the Codex reviewer on PR #69. The mechanism was verified against the code and holds; the
> scoping in the original report did not, which is why this is a REQ and not a fix on that branch.
> 
> ## Context
> 
> The four flows and their lock keys:
> 
> | Flow | Lock key | Module |
> |---|---|---|
> | card delete | `card-delete:<photoId>` | `server/modules/cms/publishing/cms-card-delete.js` |
> | unused-replacement delete | `unused-replacement-delete:<photoId>` | `server/modules/cms/publishing/cms-unused-replacement-delete.js:284` |
> | takedown resolution | `puzzle-repair-resolution:<photoId>` | `server/modules/cms/publishing/cms-takedown-flow.js` |
> | replacement build | `puzzle-repair:<photoId>` | `server/modules/cms/publishing/cms-replacement-build-flow.js` |
> 
> Meanwhile comments, catalogue entries, dev-feed requests, the photo↔level map and the dedup index each
> serialize on their own relative path. So the window is a property of the shape, present on all four.
> 
> **The absence of a shared lock is a recorded decision, not an oversight.** `server/modules/cms/core/cms-store.js:41-47`
> (REQ-1108): the catalogue-entries + dev-feed-requests pair — "the ONE cross-record pair whose correctness
> spans two records of a card" — moved into a single `node:sqlite` DB where the liveness check and the append
> run in one `BEGIN IMMEDIATE`. That "dissolved the REQ-1095 photo-scoped coordinator lock
> (`serializePhotoCardMutation`) and its LOCK ORDER contract." **Re-introducing cross-store lock sharing means
> reckoning with why it was removed and restoring a deadlock-safe lock order across four flows.** That is the
> main reason to prefer CAS — but make the call deliberately, and write it down wherever it lands.
> 
> Scope note: the fingerprint already covers the round trip from the admin's browser (the preview is re-run
> INSIDE the lock and compared there). This REQ is only about the residual window that spans the destructive
> steps themselves.
> 
> ## Detailed Requirements
> 
> - One decision — CAS revalidation vs shared/ordered locks — applied to **all four** flows, not just card
>   delete. A fix that closes one flow and leaves three reads as closure and is worse than none.
> - Narrowest concrete sub-case, and the natural first slice: `removePhotoLevelMapping`
>   (`server/modules/cms/ingestion/photo-level-map.js:525`) deletes `photos[photoId]` unconditionally, so a
>   mapping rewritten after the re-preview is erased rather than detected. Add an **optional**
>   `expectedLevelId` compare-and-swap parameter that defaults to today's unconditional behaviour, so the
>   existing unused-replacement caller is untouched.
> - Whatever is chosen, record it where the next reader will hit it: the REQ-1108 note in `cms-store.js` is
>   the existing home for this decision and should be extended rather than contradicted from elsewhere.
> 
> ## Constraints
> 
> - **Do not partially close.** Comments, catalogue entry, photo mapping and dedup signals share the same
>   window shape; a fix landing on one is not this REQ done.
> - No behaviour change for existing callers of any signature that gains a parameter — default to today's
>   semantics.
> - Deadlock safety is a hard requirement if the shared-lock route wins: four flows plus five record stores
>   need an explicit, written lock order.
> 
> ## Dependencies
> 
> None. Independent of REQ-1684.
> 
> ## Builder Guidance
> 
> **Certainty level: Mixed.** The defect is firm and verified; the remedy is a genuine open decision and the
> first real task is making it. Prefer CAS on the evidence above, but the REQ is not pre-committed to it —
> if the shared-lock route is better, take it and say why.
> 
> Start by writing the interleaving test, not the fix: it is what makes the two candidate remedies
> comparable, and it is the thing most likely to reveal that one flow differs from the other three.
> 
> ## Red-Green Proof
> 
> **RED prompt/case:** With a card previewed for deletion, let the in-lock re-preview pass, then land a
> competing write before the destructive steps run — the smallest being: rewrite the photo↔level mapping for
> that Photo ID to a different level, then let the delete proceed. Today `removePhotoLevelMapping` erases the
> newly written mapping without noticing. The same test shape with a comment appended instead of a remap
> shows the comment being deleted despite never appearing in the previewed counts.
> 
> **Why RED now:** The fingerprint is compared once, inside the flow-scoped lock, and nothing revalidates
> between that comparison and the individual destructive writes. No per-record store shares the flow's lock
> key, so nothing else serialises the competing write out.
> 
> **GREEN when:** The competing write is either serialised out or detected — the delete refuses with the
> existing `CARD_DELETE_PREVIEW_DRIFTED` code (or the flow's equivalent) rather than destroying unpreviewed
> state — and the same assertion holds for all four flows.
> 
> **Validation:** Inferred during capture, from a verified reviewer finding.
> 
> ## Full Context
> 
> See `do-work/user-requests/UR-389/input.md` for complete verbatim input.
> 
> ---
> *Source: Codex reviewer P1 on PR #69, thread discussion_r3861199188, verified against the code and rescoped.*
> 
> ---
> 
> ## Triage
> 
> **Route: C** - Complex
> 
> **Reasoning:** The request changes concurrency guarantees across four destructive flows and several stores, requires one shared architectural decision, test-first interleaving proof, and full regression verification.
> 
> **Planning:** Required
> 
> ## Plan
> 
> 1. Write RED interleaving tests for card delete, unused-replacement delete, takedown resolution, and replacement-build handoff. Each test will pause after the flow's last authoritative preview, perform a real competing store write, and prove today's destructive apply proceeds against facts that were never previewed.
> 2. Choose store-local compare-and-swap (CAS), not a resurrected cross-store coordinator. Add optional expected-value parameters to the affected store primitives, compare under each store's existing lock or SQLite transaction, preserve unconditional behavior when omitted, and translate mismatches into each flow's drift refusal.
> 3. Wire the CAS expectations through all four flows, extend the REQ-1108 note in `cms-store.js`, and verify focused interleaving coverage plus the repository's lint/full-suite gates.
> 
> **Decision:** CAS is the narrower and safer uniform remedy. REQ-1108 intentionally removed photo-scoped lock choreography after catalogue and dev-feed moved into one SQLite authority. Restoring a coordinator would require every writer across five stores to acquire a new ordered outer lock without re-entering existing non-reentrant locks. Store-local CAS instead compares the exact value being overwritten at the committing boundary and adds no new lock nesting.
> 
> **Requirement coverage:** Detailed requirement 1 maps to the four flow tests and four flow wirings; requirement 2 maps to optional `expectedLevelId` in `removePhotoLevelMapping`; requirement 3 maps to the REQ-1108 decision note in `cms-store.js`. The no-partial-close constraint is covered by the four-flow matrix and the card-delete mapping/comment/catalogue/dedup plane cases.
> 
> **Validation:** All stated requirements map to the three tasks and no task is orphaned. The plan has three tasks. Scope expansion beyond capture's seed is expected for store primitives and their focused tests; those files will be declared before implementation because flow-only rereads would recreate the same check-to-write race.
> 
> **Validation warning:** Exploration found that the replacement-build and takedown paths reach their first committing SQLite boundaries through duplicate/ingestion, lifecycle, and dev-feed command layers. Those required seams were added to Scope before implementation; a standalone preflight reread would not satisfy the four-flow requirement.
> 
> *Generated by Plan agent*
> 
> ## Exploration
> 
> The four orchestrators already re-preview inside flow-scoped locks, but the stores they later mutate use different keys or SQLite transactions. The atomic CAS therefore has to live in the store primitives themselves. `photo-level-map.js` and `dedup-index.js` can compare expected rows under their existing file locks. `cms-sqlite-store.js` can compare ordered raw rows inside the same `BEGIN IMMEDIATE` that deletes or replaces them; `cms-catalogue-store.js` is the wrapper seam for threading synchronous validators into ordinary catalogue/cascade writes.
> 
> Existing precedents are the restore kernel's synchronous `validateTargetCardState` callback and transaction test seam, plus `tests/helpers/controlled-interleaving.js` for positive-witness pause/release tests. Counts alone are insufficient because same-count replacements must drift; SQLite comparisons will use deterministic ordered row identities and raw `record_json`. Optional expectations must distinguish "omitted" from an expected missing row so existing callers retain unconditional/idempotent behavior.
> 
> No comment/activity module or export-barrel edit is needed. The known import cycle remains prohibited: `cms-sqlite-store.js` must not import `photo-level-map.js`.
> 
> *Generated by Explore agent*
> 
> ## Scope
> 
> **Files I will touch:**
> - `server/modules/cms/publishing/cms-card-delete.js` (modify) — pass exact CAS expectations through permanent card deletion
> - `server/modules/cms/publishing/cms-unused-replacement-delete.js` (modify) — pass exact CAS expectations through unused-build deletion
> - `server/modules/cms/publishing/cms-takedown-flow.js` (modify) — revalidate repair-resolution authority at each committing boundary
> - `server/modules/cms/publishing/cms-replacement-build-flow.js` (modify) — refuse a stale handoff before mint/apply
> - `server/modules/cms/ingestion/photo-level-map.js` (modify) — optional expected-level CAS under the map lock
> - `server/modules/cms/ingestion/dedup-index.js` (modify) — optional exact expected-row CAS under the dedup lock
> - `server/modules/cms/core/cms-store.js` (modify) — extend the REQ-1108 decision note
> - `server/modules/cms/core/cms-catalogue-store.js` (modify) — thread optional synchronous transaction validators
> - `server/modules/cms/core/cms-sqlite-store.js` (modify) — compare exact ordered raw rows inside the SQLite write transaction
> - `server/modules/cms/core/cms-dev-feed-request.js` (modify) — carry exact request expectations into stamp/resolve transactions
> - `server/modules/cms/core/cms-catalogue-commands-lifecycle.js` (modify) — carry expected card/request authority into retire/link commits
> - `server/modules/cms/ingestion/cms-catalogue-duplicate.js` (modify) — bind replacement duplication to the source/request expectation
> - `server/modules/cms/ingestion/cms-hand-form-ingest.js` (modify) — thread the expectation through the duplicate ingest path
> - `server/modules/cms/ingestion/ingestion-gate.js` (modify) — validate the replacement reservation at its committing catalogue boundary
> - `docs/prime-files/prime-cms.md` (modify) — document store-local CAS as the four-flow concurrency authority
> - `tests/unit/cms-production/cms-card-delete.test.js` (modify) — RED/GREEN mapping/comment/catalogue/dedup interleavings
> - `tests/integration/cms-lifecycle/cms-unused-replacement-delete.test.js` (modify) — RED/GREEN unused-build remap interleaving
> - `tests/integration/cms-lifecycle/cms-takedown-flow.test.js` (modify) — RED/GREEN resolution drift interleaving
> - `tests/unit/cms-production/cms-replacement-build-flow.test.js` (modify) — RED/GREEN handoff drift interleaving
> - `tests/unit/cms-core/cms-catalogue-transaction.test.js` (modify) — optional file-store CAS and back-compat coverage
> - `tests/unit/cms-core/cms-sqlite-store.test.js` (modify) — atomic exact-row snapshot/delete validation
> 
> **Files I will NOT touch:** HTTP/browser surfaces, lock implementation/order, unrelated CMS record modules, or public API behavior outside optional backward-compatible parameters.
> 
> **Acceptance criteria (restated from REQ):**
> - [ ] One CAS design closes the residual preview-to-apply window in all four destructive flows
> - [ ] `removePhotoLevelMapping` accepts optional `expectedLevelId`, compares under its own lock, and remains unconditional when omitted
> - [ ] Competing map/comment/catalogue/dedup writes are detected and preserved rather than destroyed
> - [ ] Existing callers without expectations keep today's behavior
> - [ ] The REQ-1108 note and CMS prime record why store-local CAS won over shared ordered locks
> - [ ] Focused interleaving tests, lint, and the full offline suite pass
> 
> ## Decisions
> 
> - **D-01 — DECIDE & STATE:** Use exact store-local CAS across all four flows and expand the scope through the actual committing command layers. Reasoning: a preflight-only reread or flow-local lock would leave the same check-to-write race; the required files are demanded by the REQ's "do not partially close" constraint, not optional scope growth.
> - **D-02 — Exact rows, optional expectations:** File-backed mutations compare the exact previewed row under their existing file lock; SQLite mutations compare raw affected rows in the writing `BEGIN IMMEDIATE`. Counts never authorize a write, so same-count replacement is drift. An omitted expectation keeps the pre-REQ unconditional/idempotent contract. Projection rebuild failures after a committed SQLite mutation remain committed successes.
> 
> ## Discovered Tasks
> 
> None.
> 
> ## Implementation Summary
> 
> **Files changed:**
> - `docs/prime-files/prime-cms.md` (modified)
> - `server/modules/cms/core/cms-catalogue-commands-lifecycle.js` (modified)
> - `server/modules/cms/core/cms-catalogue-store.js` (modified)
> - `server/modules/cms/core/cms-dev-feed-request.js` (modified)
> - `server/modules/cms/core/cms-sqlite-store.js` (modified)
> - `server/modules/cms/core/cms-store.js` (modified)
> - `server/modules/cms/ingestion/cms-catalogue-duplicate.js` (modified)
> - `server/modules/cms/ingestion/cms-hand-form-ingest.js` (modified)
> - `server/modules/cms/ingestion/dedup-index.js` (modified)
> - `server/modules/cms/ingestion/ingestion-gate.js` (modified)
> - `server/modules/cms/ingestion/photo-level-map.js` (modified)
> - `server/modules/cms/publishing/cms-card-delete.js` (modified)
> - `server/modules/cms/publishing/cms-replacement-build-flow.js` (modified)
> - `server/modules/cms/publishing/cms-takedown-flow.js` (modified)
> - `server/modules/cms/publishing/cms-unused-replacement-delete.js` (modified)
> - `tests/integration/cms-lifecycle/cms-takedown-flow.test.js` (modified)
> - `tests/integration/cms-lifecycle/cms-unused-replacement-delete.test.js` (modified)
> - `tests/unit/cms-core/cms-catalogue-transaction.test.js` (modified)
> - `tests/unit/cms-core/cms-sqlite-store.test.js` (modified)
> - `tests/unit/cms-production/cms-card-delete.test.js` (modified)
> - `tests/unit/cms-production/cms-replacement-build-flow.test.js` (modified)
> 
> **What was done:** Added exact store-local compare-and-swap expectations under the affected file locks and SQLite write transactions, then threaded them through permanent card deletion, unused-replacement deletion, repair resolution, and replacement-build handoff. Added witnessed interleaving regressions, preserved omitted-expectation compatibility and committed-projection semantics, and recorded the CAS-over-coordinator concurrency rule in the CMS authorities.
> 
> ## Qualification
> 
> Passed — 21 project files verified, all detailed requirements and six restated acceptance criteria traced to substantive diffs, P-A-U confirmed, store data paths inspected, and Scope/Implementation Summary matched exactly.
> 
> ## Testing
> 
> **Tests run:** `npm test -- --runInBand --testPathPatterns='cms-card-delete|cms-unused-replacement-delete|cms-takedown-flow|cms-replacement-build-flow|cms-catalogue-transaction|cms-sqlite-store'`; `npm run lint`; `npm test`; `git diff --check`
> **Result:** ✓ All passing — 206 focused tests; lint with zero warnings; 435 full-suite test suites / 6,500 tests passed (9 suites / 40 tests skipped by repository configuration); clean diff check.
> 
> **Red-green validation:**
> - Captured mapping case: card-delete and unused-build interleaving tests initially deleted a competing remap, then passed by refusing with their flow drift codes while preserving the new mapping.
> - Captured comment/row case: exact SQLite deletion test initially accepted a same-count comment replacement, then passed by refusing the whole transaction and preserving the competing row.
> - All-four-flow closure: replacement-build initially minted a stale duplicate and takedown initially retired a rewritten source; both now refuse at their transactional apply boundaries. The six-suite matrix passes together.
> - Primitive compatibility: optional photo-map/dedup/SQLite expectations initially did not exist; CAS cases now reject drift while omission cases retain unconditional behavior.
> 
> **New tests added:**
> - Interleaving regressions in the four flow suites named above
> - Exact file-store CAS/backward-compat cases in `tests/unit/cms-core/cms-catalogue-transaction.test.js`
> - Exact affected-row SQLite CAS/backward-compat cases in `tests/unit/cms-core/cms-sqlite-store.test.js`
> 
> *Verified by work action*
> 
> ## Review
> 
> **Overall: 67%** | 2026-08-26T12:00:08Z
> 
> | Dimension | Score |
> |-----------|-------|
> | Requirements | 72% |
> | Code Quality | 70% |
> | Test Adequacy | 65% |
> | Scope | 100% |
> | Risk | Low |
> | Acceptance | Partial |
> 
> **Important findings:**
> - Replacement drift is first checked only after the draft level, manifest row, Photo ID, and mapping are allocated; the refusal can leave an orphan while the test checks only catalogue count. — impact-user-visible → remediation required in REQ-1685
> - Card and unused-replacement deletion throw drift after earlier destructive steps can commit, bypassing their truthful partial-result/evidence paths and discarding `committedSteps`. — impact-user-visible → remediation required in REQ-1685
> - Replacement reuse/final stamp and takedown final resolution do not carry every relevant row expectation to each later committing boundary, so the documented four-flow guarantee is incomplete. — impact-user-visible → remediation required in REQ-1685
> 
> **Minor findings:** 1 — document the new public expectation parameters in their JSDoc contracts.
> **Acceptance:** Partial — low-level CAS works, but orchestration completeness and partial-state truthfulness must be repaired before this REQ can ship.
> **Suggested testing:** Assert all replacement allocation planes on drift; cover partial-result evidence after late delete drift; inject drift at replacement reuse/final stamp and takedown final resolution.
> **Follow-ups created:** None — the user requested all REQs implemented, so these core acceptance gaps remain in the current REQ's remediation loop.
> 
> *Reviewed by review-work action*
> 
> ## Remediation
> 
> **Attempt 1 review:** The CAS primitives passed, but the independent Route C review found three boundary-completeness defects. Remediation will add full-plane cleanup/truthful partial outcomes and carry expectations through every later commit, then rerun Summary → Qualification → Testing → Review.
> 
> ---
> 
> [Follow-up message 2:]
> 
> Basically, we shouldn't assume that people know what those numbers mean. This is why it's good to have a glossary and an inline mention when they appear for the first time. This is true for user requests and it is true for requests as well.
> 
> ---
> 
> [Follow-up message 3:]
> 
> Make sure to capture the intent using user requests and requests by a do work capture request.
> 
> ---
> 
> [Follow-up message 4:]
> 
> after the req is captured, run do-work verify-request
> ````

---
*Captured: 2026-08-26T13:02:24Z*
