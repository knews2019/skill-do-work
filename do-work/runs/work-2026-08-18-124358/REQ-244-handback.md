# REQ-244 Hand-back — Cite the Timestamp Rule at Every Stamp Write Site

Written 2026-08-18T13:16:47Z (clock read at stamp time).

## Branch

`worktree-agent-REQ-244-cite-the-timestamp-rule-at-every-stamp-write-site`
Implementation commit: `c417c58`

## Scope

Declared after the sweep, before any write. This is the list the diff should be judged against.

**Files I will touch (13, all modify — no new files):**

| File | new\|modify | Why |
|---|---|---|
| `_dev/tests/contract-regressions.sh` | modify | The lock-in check, added inside the existing REQ-078 timestamp block |
| `skills/do-work/actions/work-reference.md` | modify | Two fenced templates + one run-directory placeholder |
| `skills/do-work/actions/work.md` | modify | Uncited `completed_at` on the failure path (Step 8) |
| `skills/do-work/actions/clarify.md` | modify | Uncited `status_changed_at` on the unblock path |
| `skills/do-work/actions/review-work.md` | modify | Follow-up template `created_at` + report footer stamp |
| `skills/do-work/actions/forensics.md` | modify | Two `**Scan date:**` report templates |
| `skills/do-work/actions/roadmap.md` | modify | Two `**Scan date:**` report templates |
| `skills/do-work-toolbox/actions/code-review.md` | modify | Follow-up template `created_at` |
| `skills/do-work-toolbox/actions/present-work.md` | modify | `**Generated:**` portfolio header |
| `skills/do-work-toolbox/actions/deep-explore.md` | modify | Uncited `completed_at` in state.json write |
| `skills/do-work-toolbox/actions/deep-explore-reference.md` | modify | `session_dir` name placeholder |
| `skills/do-work-knowledge/actions/interview.md` | modify | Three `<now>` sites |
| `skills/do-work-knowledge/actions/interview-reference.md` | modify | Six `<now>` sites |

**Files I will NOT touch:**

- Everything the guardrails forbid: `do-work/` (except this hand-back), `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`.
- Every sibling-owned board path: `skills/do-work-board/tools/queue-kanban/{web/board-durations.js, generate_test.go, model.go, verify.go, web/board-cards.js, web/board-core.js, web/board.css, prime-do-kanban.md}`.
- **Date-only and local-date sites**, which the rule's own "Date-only stamps" paragraph governs and which the brief says not to convert: `actions/note.md`, `actions/memory.md`, `actions/memory-reference.md`, `actions/dream.md`, `actions/bkb*.md`, `actions/maintainability-audit*.md`, `actions/prime.md`, the `versions/v<N>-<YYYY-MM-DD>/` archive names in `interview*.md`, and the changelog/run-directory/report-slug names.
- **`skills/do-work-toolbox/actions/ui-review.md:216` (`**Date**: [today]`)** — a date-only report header, not a timestamp. Left alone for the same reason, and my check is keyed on the word *timestamp* so it does not catch it. Flagged rather than swept: see Pushback.
- No new `_dev/tests/` file. The REQ-078 timestamp block in `contract-regressions.sh` is the check's natural home ("Delete before you add"), and it is already the place that owns the other half of the same arrangement.

Diff matches Scope exactly: 13 files, all declared.

## P-A-U

- [x] **[PLAN]:** Read `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`, `CLAUDE.md`, the three always-on crew members, and the Timestamp rule itself (`skills/do-work/actions/work-reference.md:99-113`). Approach: (1) sweep all four skills for placeholder-shaped stamp sites and classify each as instant / date-only / directory name; (2) normalize every bracket-form placeholder to the two spellings the rule recognizes (`<timestamp>`, `<now>`) — never copying a clock command, since work-reference.md ~line 101 reserves that to itself; (3) add an inline citation at each uncited site, using the plain `actions/work-reference.md` form inside core and the `../do-work/actions/work-reference.md` sibling form from toolbox/knowledge per the `actions/memory.md` precedent the REQ names; (4) lock it in with a condition-keyed check placed in the existing REQ-078 timestamp block rather than a new file. No code written during this step.
- [x] **[APPLY]:** Applied exactly the planned scope through a fail-loud patch script (each replacement asserted to match exactly once, so a silent miss or a double-apply aborts). 26 replacements across 12 action files, plus the check in `contract-regressions.sh`. No file outside the declared Scope was written.
- [x] **[UNIFY]:** `git diff --stat` reviewed; 13 files, +123/−28. Verified per file:
  - `_dev/tests/contract-regressions.sh` — the only added block is the checker plus its rationale comment; ShellCheck clean (it ran over 50 tracked files inside `maintainer-verify`); the block uses the file's existing `if ! python3 - "$repo_root" <<'PY' … PY then … fi` idiom and increments `fail_count` like its neighbours.
  - `work-reference.md` — read lines 336, 627, 893 in context: run-directory placeholder now matches the `work-<YYYY-MM-DD-HHMMSS>` form the file's own dispatch table already uses at line 393; both fenced YAML templates carry a same-line comment citation.
  - `work.md`, `clarify.md` — citation inserted mid-parenthetical; the surrounding claims (`mandatory on every terminal flip`, `blocked_at is removed on this flip`) are preserved verbatim, and I checked `contract-regressions.sh` for assertions pinning those two lines (the `status_changed_at` assertion at :1824 targets Crash Recovery in work-reference.md, not these).
  - `review-work.md`, `forensics.md`, `roadmap.md`, `present-work.md` — rendered each fence and its introducing line; the templates still read as templates, and no citation text landed inside a rendered report body except the two YAML comments, which are valid YAML if copied.
  - `code-review.md`, `deep-explore.md`, `deep-explore-reference.md`, `interview.md`, `interview-reference.md` — all use the sibling citation form; `deep-explore-reference.md`'s `session_dir` now matches the `date +%Y%m%d-%H%M%S` its own action prescribes at `deep-explore.md:111`.
  - No debug artifacts: `git status --untracked-files=all` shows only the 13 tracked modifications; all scratch lives in `/tmp/req244/`.

## Files Changed

```
 _dev/tests/contract-regressions.sh                 | 91 ++++++++++++++++++++++
 .../actions/interview-reference.md                 | 12 +--
 skills/do-work-knowledge/actions/interview.md      |  6 +-
 skills/do-work-toolbox/actions/code-review.md      |  2 +-
 .../actions/deep-explore-reference.md              |  2 +-
 skills/do-work-toolbox/actions/deep-explore.md     |  2 +-
 skills/do-work-toolbox/actions/present-work.md     |  4 +-
 skills/do-work/actions/clarify.md                  |  2 +-
 skills/do-work/actions/forensics.md                |  8 +-
 skills/do-work/actions/review-work.md              |  6 +-
 skills/do-work/actions/roadmap.md                  |  8 +-
 skills/do-work/actions/work-reference.md           |  6 +-
 skills/do-work/actions/work.md                     |  2 +-
 13 files changed, 123 insertions(+), 28 deletions(-)
```

- **`_dev/tests/contract-regressions.sh`** — the lock-in check, inserted directly after the REQ-078 "only work-reference.md may spell a timestamp command" check, because it enforces the other half of that same arrangement. Two conditions, both keyed on shape: a bracket-form placeholder naming a timestamp fails anywhere in shipped `skills/*/actions/*.md`; every `<timestamp>`/`<now>` must carry `Timestamp rule` on its own line or on the nearest non-blank line above its fence's opening delimiter. Plus a zero-site guard.
- **`skills/do-work/actions/work-reference.md`** — `created_at` (Builder-Decided Follow-up Template) and `session_ended` (Session Checkpoint Template) normalized to `<timestamp>` with a same-line citation comment; `work-<timestamp>/` in Sole integrator respelled `work-<YYYY-MM-DD-HHMMSS>/`.
- **`skills/do-work/actions/work.md`** — the Step 8 failure path's `completed_at: <timestamp>` now cites the rule, matching its three already-cited siblings in the same file.
- **`skills/do-work/actions/clarify.md`** — the unblock path's `status_changed_at: <timestamp>` now cites the rule, matching the other four cited sites in the same file.
- **`skills/do-work/actions/review-work.md`** — the "Review Fix" follow-up template's `created_at` (the site the reported incident came from) and the report footer stamp.
- **`skills/do-work/actions/forensics.md`, `roadmap.md`** — both `**Scan date:**` report templates in each file; the placeholder was already spelled "timestamp", so normalizing to `<timestamp>` preserves its existing instant meaning rather than changing it.
- **`skills/do-work-toolbox/actions/present-work.md`** — `**Generated:** [UTC timestamp]` → `<timestamp>`, cited on the line introducing the template.
- **`skills/do-work-toolbox/actions/code-review.md`** — follow-up template `created_at`, sibling citation form.
- **`skills/do-work-toolbox/actions/deep-explore.md`** — the state.json `completed_at` write, in a skill that never loads the rule.
- **`skills/do-work-toolbox/actions/deep-explore-reference.md`** — `session_dir` respelled `<YYYYMMDD-HHMMSS>`, matching `deep-explore.md:111`.
- **`skills/do-work-knowledge/actions/interview.md`, `interview-reference.md`** — the nine-site `<now>` cluster (`started_at`, `last_activity_at` ×2, `last_validated_at` ×3, `approved_at` ×2, `review_completed_at`, `last_exported_at`), each with a sibling citation.

## Red-Green Evidence

The check went in first, then failed for the reason it exists. All three REDs run the **checker extracted from the committed `contract-regressions.sh`**, not a hand-copy, so what fails is the shipped code.

### RED 1 — shipped checker against the pre-change action files

Pre-change tree materialized from git rather than by stashing (a stash of a clean file stashes nothing and the resulting green reads as proof):

```
$ git archive HEAD skills | tar -x -C /tmp/req244/red-tree
$ python3 /tmp/req244/shipped-check.py /tmp/req244/red-tree
timestamp citation failures:
- skills/do-work/actions/clarify.md:142: stamp write site cites no Timestamp rule on its own line or above its fence: - **Yes → unblock:** set `status: pending`, stamp `status_changed_at: <timestamp>` (blocked_at is removed on t
- skills/do-work/actions/forensics.md:216: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/forensics.md:257: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/review-work.md:365: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/review-work.md:425: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/roadmap.md:135: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/roadmap.md:244: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/work-reference.md:336: stamp write site cites no Timestamp rule on its own line or above its fence: **Sole integrator.** The builder never writes the main tree or its branch, **with exactly one exception: its o
- skills/do-work/actions/work-reference.md:627: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/work-reference.md:893: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work/actions/work.md:597: stamp write site cites no Timestamp rule on its own line or above its fence: Classify the failure and queue the right follow-up per `actions/work-reference.md` → **Failure Classification
- skills/do-work-knowledge/actions/interview-reference.md:297: stamp write site cites no Timestamp rule on its own line or above its fence: 5. Write a new empty `session.json` — fresh `session_id`, `template: <slug>`, `template_version: <current temp
- skills/do-work-knowledge/actions/interview-reference.md:315: stamp write site cites no Timestamp rule on its own line or above its fence: 3. On `confirm`: set `source_confidence: confirmed`, update `last_validated_at: <now>`, leave all other fields
- skills/do-work-knowledge/actions/interview-reference.md:316: stamp write site cites no Timestamp rule on its own line or above its fence: 4. On `edit`: enter an interactive edit — show current values, let the user override any field, produce a new
- skills/do-work-knowledge/actions/interview-reference.md:317: stamp write site cites no Timestamp rule on its own line or above its fence: 5. On `mark-stale`: set `status: stale`, update `last_validated_at: <now>`. The entry remains but is flagged i
- skills/do-work-knowledge/actions/interview-reference.md:322: stamp write site cites no Timestamp rule on its own line or above its fence: 6. Once all entries are processed in a layer, re-approve the layer as a whole — the layer-level approval gate
- skills/do-work-knowledge/actions/interview-reference.md:339: stamp write site cites no Timestamp rule on its own line or above its fence: 3. Write a new empty `session.json` — same shape as `fresh` (including `template_version: <current template ve
- skills/do-work-knowledge/actions/interview.md:165: stamp write site cites no Timestamp rule on its own line or above its fence: 6. **Persist on approval.** Write approved entries to `session.json` under `layers.<layer-id>.entries[]`. Set
- skills/do-work-knowledge/actions/interview.md:235: stamp write site cites no Timestamp rule on its own line or above its fence: 6. When every surfaced tension has a resolution (or was skipped), update `session.json`: set `review_completed
- skills/do-work-knowledge/actions/interview.md:257: stamp write site cites no Timestamp rule on its own line or above its fence: 2. **Stamp the export timestamp in-memory before rendering.** Set `session.last_exported_at = <now>` (ISO 8601
- skills/do-work-toolbox/actions/code-review.md:301: [timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
- skills/do-work-toolbox/actions/deep-explore-reference.md:323: stamp write site cites no Timestamp rule on its own line or above its fence: "session_dir": "do-work/runs/deep-explore-<sanitized-slug>-<timestamp>",
- skills/do-work-toolbox/actions/deep-explore.md:250: stamp write site cites no Timestamp rule on its own line or above its fence: Update state.json: set `writer_status: "done"`, `status: "complete"`, `completed_at: <timestamp>`, `surviving_
- skills/do-work-toolbox/actions/present-work.md:86: [UTC timestamp] is a placeholder spelling no rule governs; the Timestamp rule recognizes <timestamp> and <now>
exit=1
```

24 sites — the REQ's starting list plus five it did not name: `clarify.md:142`, `work.md:597`, `interview-reference.md:339`, and the two run-directory placeholders.

### RED 2 — mutation: one citation removed from the fixed tree

Proves the citation clause fails for its own reason, not merely because a bracket form is present. The placeholder is left untouched; only the citation is deleted:

```
$ python3 /tmp/req244/shipped-check.py /tmp/req244/mutant
timestamp citation failures:
- skills/do-work-toolbox/actions/deep-explore.md:250: stamp write site cites no Timestamp rule on its own line or above its fence: Update state.json: set `writer_status: "done"`, `status: "complete"`, `completed_at: <timestamp>`, `surviving_
exit=1
```

Exactly one failure, at the mutated site.

### RED 3 — anti-vacuity guard

A tree with `skills/*/actions/` present but no recognized placeholder anywhere. Without the guard this would pass silently, which is exactly what happens if the spellings are ever renamed:

```
$ python3 /tmp/req244/shipped-check.py /tmp/req244/blind
timestamp citation failures:
- no recognized stamp placeholder was found in any shipped action file — the spellings this check keys on were renamed and it has gone blind
exit=1
```

### GREEN — the same checker against the fixed tree

```
$ python3 /tmp/req244/check.py "$PWD"
OK: 43 stamp sites, all cited
exit=0
```

## Verification

```
maintainer-verify: aggregate contract suite
Maintainer verification self-test passed.
Suite manifest contract probes passed.
shipped package reference contract: PASS
Shell-block lint self-test passed.
Shell-block lint passed: 74 fenced blocks and 31 shipped shell files; ShellCheck enabled.
SessionStart hook behavior probes passed.
Prescribed shell primitive canonicalization checks passed.
Defensive-surface exact deletion regressions passed.
record-commit-hash and blanked-req-scan guard probes passed.
update-script behavior probes passed.
Prescribed shell script behavior probes passed (42 named script cases).
staged skills contract: PASS
suite installer behavior probes passed.
p50 estimator suite: all probes passed.
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	15.352s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (4.68s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	4.866s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.658s
Maintainer verification passed.
0
```

Run unpiped from the worktree root; `echo $?` printed `0`.

REQ-243's link+anchor checker (live in this worktree, reported above as `shipped package reference contract: PASS`) did not fire on any citation — expected, since it resolves Markdown `[text](target)` links and link definitions, and every citation I added is a backticked path, matching the convention `actions/memory.md` established. Worth knowing: **backticked cross-package pointers are not covered by that checker**, so their correctness rests on convention, not enforcement. See Lessons.

## Integration Seams

None. No shared registry, no cross-REQ text, nothing for the orchestrator to apply by hand. Merge the branch as-is.

One note for the integrator, not a seam: the diff touches `skills/do-work-toolbox/` and `skills/do-work-knowledge/` action markdown only — no overlap with REQ-242's or REQ-245's board paths.

## Decisions

- **D-01 — The check lives in the existing REQ-078 timestamp block, not a new `_dev/tests/` file.** That block already enforces "only the rule's home spells a clock command"; this is the same arrangement's other half ("every other site cites the rule"), and separating them would let one drift from the other. No new suite file was added.
- **D-02 — "Uncited site" is defined by position, not by a site list.** A recognized placeholder must carry `Timestamp rule` on its own line, or — when it sits inside a fenced block whose lines are copied verbatim into a generated artifact — on the nearest non-blank line above the fence's opening delimiter. Those are the only two positions a reader of the site actually sees; a citation elsewhere in the file is not *at* the site, which is the failure mode this REQ exists to close. The fence clause is a position, not an exemption: fenced YAML takes the same-line form as a trailing `#` comment (the `capture-reference.md:16` precedent), so it needs no carve-out at all.
- **D-03 — Bracket-form placeholders are rejected outright, not merely required to be cited.** `[timestamp]` and `[UTC timestamp]` are spellings no rule governs; the rule's own sentence names `<timestamp>` and `<now>` as the two forms that mean it. Keying the reject on the *word* rather than on the two observed strings means a future `[ISO timestamp]` is caught without anyone updating a list — and it deliberately does not catch date-only forms like `[today]`, which are a different shape.
- **D-04 — The two run-directory placeholders were respelled with their own calendar form.** `work-<timestamp>/` → `work-<YYYY-MM-DD-HHMMSS>/` and `deep-explore-…-<timestamp>` → `deep-explore-…-<YYYYMMDD-HHMMSS>`. These are names, not stamps, and the brief puts them out of scope for *citation* — but each file already states its own directory-name format elsewhere (`work-reference.md:393` and `deep-explore.md:111`), so the stamp token at those two sites was simply the wrong spelling for what they name. Respelling them makes the docs match themselves and, as a side effect, leaves `<timestamp>`/`<now>` meaning exactly one thing across the corpus — so the checker needs no path exception. This was the alternative to writing a path-context carve-out into the check, and deleting the ambiguity beat carving around it.
- **D-05 — Cross-package citations use the `../do-work/actions/work-reference.md` sibling form.** Both forms exist in the tree today (`present-work.md:37` uses `../../`), but the REQ names `actions/memory.md` as the precedent to follow, and `prime-action-files.md` prescribes the sibling-path form. Note this form is skill-root-relative, not directory-relative — it does not resolve literally from the citing file's own directory. That is the repo's established reading; I did not change it, and it is worth a maintainer decision at some point (see Pushback).
- **D-06 — A zero-site guard is part of the check.** If both recognized spellings are ever renamed, the check would otherwise pass by matching nothing. It now fails loudly instead.

## Lessons Learned

- **A rule that centralizes a command creates an obligation at every site that does not have it.** REQ-078 moved the clock command to one home for a real reason (Windows agents kept getting an unreachable fix). That move is only safe if every other site points home — and nothing enforced the second half for eight months, until a stamp was fabricated. When a future REQ centralizes something, the citation obligation at the periphery is part of the same change, not a follow-up.
- **The `<timestamp>` token was doing two jobs.** It marked stamps *and* appeared inside directory names, which is precisely why a mechanical check looked like it needed an exception. It did not need one — it needed the two jobs separated. Reach for "which of these uses is actually the odd one out" before writing a carve-out into a detector.
- **Backticked cross-package pointers are unenforced.** REQ-243's checker resolves Markdown link syntax; the repo's dominant convention for citing another action file is a backticked path, which that checker never sees. Two spellings of the sibling path (`../` and `../../`) coexist today because nothing can tell them apart. If cross-package pointer correctness matters, either the convention moves to Markdown links or a separate resolver has to read backticked paths — a candidate REQ, not something I did here.
- **`git archive HEAD <path> | tar -x -C <dir>` is a clean way to get a pre-change tree for a RED.** It never touches the working tree, so there is no stash to forget to pop and no risk of the vacuous "stashed a clean file" green.

## Pushback

Two things I did not do, flagged rather than silently absorbed:

1. **`skills/do-work-toolbox/actions/ui-review.md:216` — `**Date**: [today]` — is a real uncited clock-write site that this REQ's framing leaves uncovered.** It is date-shaped, so the REQ's "any `*_at:`/date-shaped placeholder" arguably reaches it, but the brief's hard constraint is not to convert date-only stamps to instants, and nothing states whether that header wants a UTC date or a deliberate local one. Converting it either way would be me deciding a question the REQ did not answer, so I left it and it is listed in Scope as untouched. If you want it governed, the cheap version is one line in the rule's Date-only paragraph plus a citation there — but somebody has to decide UTC vs. local for a human-facing report header first.

2. **The `../` vs `../../` split in cross-package backtick citations is a genuine inconsistency and I propagated the more common half.** `../do-work/actions/work-reference.md` does not resolve from `skills/do-work-toolbox/actions/`; `../../do-work/actions/work-reference.md` does, in both the source and installed topologies. The REQ told me to follow `actions/memory.md`, which uses `../`, and `prime-action-files.md` prescribes that form, so that is what I used — nine new citations' worth. If the intended reading is literal path resolution rather than "sibling package", then the prime and roughly a dozen existing citations are wrong together and want one sweep, not a per-REQ correction. I would rather that be decided deliberately than have me quietly pick the other form mid-REQ.
