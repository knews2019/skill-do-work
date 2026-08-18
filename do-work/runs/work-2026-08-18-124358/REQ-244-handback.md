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

---

# Addendum — Review Remediation

Written 2026-08-18T13:57:41Z (clock read at stamp time).

Remediation commit: **`c5f75e2`** — `[REQ-244] key the timestamp check on shape, and sweep the 12 sites it could not see`
Branch merged from `main` first (`23bad9d`), so this builds on REQ-241's and REQ-245's remediations; REQ-245's `forensics.md` edits merged clean with mine.

All five findings addressed. `maintainer-verify.sh` exits 0.

## Correction to the record (F4) — read this first

The original hand-back's Red-Green section asserts that every run used "the checker extracted from the committed `contract-regressions.sh`, not a hand-copy". **That was true of the three REDs and false of the GREEN.** The line I quoted —

```
$ python3 /tmp/req244/check.py "$PWD"
OK: 43 stamp sites, all cited
```

— came from `check.py`, the prototype I wrote before the check existed, not from the committed check, which printed nothing on success and could not have produced that line. The run happened and the count was right, but I presented a prototype's output under a sentence claiming provenance it did not have. In a REQ about agents writing values they did not read, that is the error the REQ exists to prevent, committed in the hand-back for it.

The committed check now prints its counts on success, so this class of confusion cannot recur: there is a real success line to quote, and it carries the numbers.

Every transcript in this addendum comes from `/tmp/req244b/shipped.py`, extracted verbatim from the committed `_dev/tests/contract-regressions.sh` with:

```
awk '/^if ! python3 - "\$repo_root" <<.PY.$/{n++} n==2 && !/^if ! python3/ && !/^PY$/{print} n==2 && /^PY$/{exit}' \
  _dev/tests/contract-regressions.sh > /tmp/req244b/shipped.py
```

## F2 — the check's own closed enumeration (the one that mattered)

The finding is correct and was the root cause of F1. The first check keyed on a literal `[…]` span containing the word "timestamp". That is an enumeration of *spellings*, and it locked in precisely the drift already removed while being blind to every spelling still present. My D-03 argued the reject keyed on the word rather than the two observed strings — true, and beside the point, because the bracket requirement was the enumeration one level up. Accepted without reservation.

**Recognition is now keyed on shape.** A site is a placeholder — any bracketed or quoted span short enough to name a value rather than be prose (≤ 30 characters of inner text, or wrapping a nested placeholder) — that is assigned to something, and that either denotes a clock value by name or **is the value of an `*_at` key**. The last clause is the rule's own stated trigger, quoted rather than invented: *"every `*_at` field in this schema, and any timestamp a future field adds."*

**Requirement stayed narrow while recognition went broad** — that split is what makes the check both complete and enforcing:

1. **Spelling** — an instant site's line must carry `<timestamp>` or `<now>`, the two spellings the rule's own sentence says mean it. A new spelling is still *recognized*, so it gets normalized rather than silently accepted.
2. **Citation** — unchanged. Same line, or the nearest non-blank line above the fence's opening delimiter. The reviewer upheld this definition and it is untouched.

**Names are excluded by shape, not by exception:** a `/` inside the placeholder, or a `-`/`/` glued to its outside, means it names a path or a compound. So `work-<YYYY-MM-DD-HHMMSS>`, `v<N>-<date>/`, and `raw/processed/<date>/` are skipped without any file or path list.

**Instant vs. date-only is now an explicit split.** Instants (`timestamp`, `now`, `iso`, or the `YYYY-MM-DDTHH…` shape, or any `*_at` value) must cite. Date-only placeholders (`date`, `today`, `YYYY-MM-DD` with no time part) are recognized and **counted but not required to cite** — see D-08 for why, and note this is what keeps `[today]`, `{today}` and `note.md`'s `[YYYY-MM-DD]` out of scope without a file list doing it.

### F2 probe set — every spelling the reviewer showed escaping

Each probe is a single-purpose file containing one already-cited baseline site (so the zero-site guard cannot be what fires) plus the probe line:

```
=== F2 set: caught for a named reason, not the zero-site guard ===
caught  bracket-timestamp      | created_at: [timestamp]
          - probe.md:5: instant write site is written [timestamp]; the rule recognizes <timestamp> and <now>
          - probe.md:5: stamp write site cites no Timestamp rule on its own line or above its fence: created_at: [timestamp]
caught  iso-8601-timestamp     | **Scan date:** <ISO 8601 timestamp>
          - probe.md:5: instant write site is written <ISO 8601 timestamp>; the rule recognizes <timestamp> and <now>
          - probe.md:5: stamp write site cites no Timestamp rule on its own line or above its fence: **Scan date:** <ISO 8601 timestamp>
caught  brace-timestamp        | **Date:** {timestamp}
          - probe.md:5: instant write site is written {timestamp}; the rule recognizes <timestamp> and <now>
          - probe.md:5: stamp write site cites no Timestamp rule on its own line or above its fence: **Date:** {timestamp}
caught  quoted-iso             |   "started_at": "<iso>",
          - probe.md:5: instant write site is written "<iso>"; the rule recognizes <timestamp> and <now>
          - probe.md:5: stamp write site cites no Timestamp rule on its own line or above its fence: "started_at": "<iso>",
caught  iso-instant-shape      | created_at: <YYYY-MM-DDTHH:MM:SSZ>
          - probe.md:5: instant write site is written <YYYY-MM-DDTHH:MM:SSZ>; the rule recognizes <timestamp> and <now>
          - probe.md:5: stamp write site cites no Timestamp rule on its own line or above its fence: created_at: <YYYY-MM-DDTHH:MM:SSZ>
caught  bare-stamp-key         | blocked_at: <the moment it blocked>
          - probe.md:5: instant write site is written <the moment it blocked>; the rule recognizes <timestamp> and <now>
          - probe.md:5: stamp write site cites no Timestamp rule on its own line or above its fence: blocked_at: <the moment it blocked>
```

`bare-stamp-key` is not from the reviewer's list — it is the `*_at`-adjacency clause, proving a placeholder with no clock word in its name is still caught by the key it is assigned to.

### Controls — what must NOT be flagged

Same harness, baseline present so the run can legitimately succeed:

```
=== controls: must NOT be flagged (baseline cited site present) ===
not flagged  name-run-directory     | session_dir: "do-work/runs/work-<YYYY-MM-DD-HHMMSS>"
     Timestamp rule citation contract: 1 instant write sites cited, 0 date-only sites recognized.
not flagged  name-version-dir       | archive as: v<N>-<date>/
     Timestamp rule citation contract: 1 instant write sites cited, 0 date-only sites recognized.
not flagged  date-only-today        | **Date**: [today]
     Timestamp rule citation contract: 1 instant write sites cited, 1 date-only sites recognized.
not flagged  prose-quote-now        | Write "Now you can X; lives in Y subsystem" here.
     Timestamp rule citation contract: 1 instant write sites cited, 0 date-only sites recognized.
not flagged  powershell-format      | Use `"yyyy-MM-dd\THH:mm:ss\Z"` on Windows.
     Timestamp rule citation contract: 1 instant write sites cited, 0 date-only sites recognized.
not flagged  long-prose-bracket     | **RED prompt/case:** [Minimal prompt, repro, or example that should fail before the fix]
     Timestamp rule citation contract: 1 instant write sites cited, 0 date-only sites recognized.
```

`date-only-today` is the one that matters: `[today]` is *recognized* (the date-only count goes to 1) but not required to cite. Recognition and requirement are genuinely separate, not a detector that cannot see it.

The last three controls are real false positives I hit while building this and had to design out: a quoted prose sentence containing "now", the PowerShell format string in the rule's own home (`Get-Date` matches `\bdate\b`, `HH:mm` nearly matches the time shape), and a long bracketed descriptor. The length clause and case-sensitive format shapes are what exclude them.

## F1 — the 12 missed sites

All 12 fixed. The reviewer's method reproduced exactly.

### RED — widened check against the pre-remediation tree

Pre-remediation tree taken from git, not from a stash:

```
$ git archive HEAD skills | tar -x -C /tmp/req244b/red        # HEAD = 23bad9d
$ python3 /tmp/req244b/shipped.py /tmp/req244b/red
timestamp citation failures:
- skills/do-work-knowledge/actions/interview-reference.md:133: instant write site is written "<iso>"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview-reference.md:133: stamp write site cites no Timestamp rule on its own line or above its fence: "started_at": "<iso>",
- skills/do-work-knowledge/actions/interview-reference.md:134: instant write site is written "<iso>"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview-reference.md:134: stamp write site cites no Timestamp rule on its own line or above its fence: "last_activity_at": "<iso>",
- skills/do-work-knowledge/actions/interview-reference.md:138: instant write site is written "<iso> | null"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview-reference.md:138: stamp write site cites no Timestamp rule on its own line or above its fence: "review_completed_at": "<iso> | null",
- skills/do-work-knowledge/actions/interview-reference.md:140: instant write site is written "<iso> | null"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview-reference.md:140: stamp write site cites no Timestamp rule on its own line or above its fence: "last_exported_at": "<iso> | null",
- skills/do-work-knowledge/actions/interview-reference.md:144: instant write site is written "<iso>"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview-reference.md:144: stamp write site cites no Timestamp rule on its own line or above its fence: "approved_at": "<iso>",
- skills/do-work-knowledge/actions/interview.md:120: instant write site is written "<iso>"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview.md:120: stamp write site cites no Timestamp rule on its own line or above its fence: "started_at": "<iso>",
- skills/do-work-knowledge/actions/interview.md:121: instant write site is written "<iso>"; the rule recognizes <timestamp> and <now>
- skills/do-work-knowledge/actions/interview.md:121: stamp write site cites no Timestamp rule on its own line or above its fence: "last_activity_at": "<iso>",
- skills/do-work-toolbox/actions/deep-explore-reference.md:326: instant write site is written "ISO 8601 timestamp"; the rule recognizes <timestamp> and <now>
- skills/do-work-toolbox/actions/deep-explore-reference.md:326: stamp write site cites no Timestamp rule on its own line or above its fence: "created_at": "ISO 8601 timestamp",
- skills/do-work-toolbox/actions/inspect.md:253: instant write site is written {timestamp}; the rule recognizes <timestamp> and <now>
- skills/do-work-toolbox/actions/inspect.md:253: stamp write site cites no Timestamp rule on its own line or above its fence: **Date:** {timestamp}
- skills/do-work-toolbox/actions/inspect.md:358: instant write site is written {timestamp}; the rule recognizes <timestamp> and <now>
- skills/do-work-toolbox/actions/inspect.md:358: stamp write site cites no Timestamp rule on its own line or above its fence: **Date:** {timestamp}
- skills/do-work-toolbox/actions/stray-check.md:92: instant write site is written <ISO 8601 timestamp>; the rule recognizes <timestamp> and <now>
- skills/do-work-toolbox/actions/stray-check.md:92: stamp write site cites no Timestamp rule on its own line or above its fence: **Scan date:** <ISO 8601 timestamp>
- skills/do-work-toolbox/actions/stray-check.md:127: instant write site is written <ISO 8601 timestamp>; the rule recognizes <timestamp> and <now>
- skills/do-work-toolbox/actions/stray-check.md:127: stamp write site cites no Timestamp rule on its own line or above its fence: **Scan date:** <ISO 8601 timestamp>
exit=1
```

Exactly the 12 sites, each failing both clauses. Nothing else in the corpus regressed into the widened net — the 42 sites already cited by the first commit stay green.

### GREEN

```
$ python3 /tmp/req244b/shipped.py "$PWD"
Timestamp rule citation contract: 54 instant write sites cited, 17 date-only sites recognized.
exit=0
```

54 instant sites, up from the 43 the first commit knew about: the 12 newly recognized, less one that the name-shape rule correctly reclassified as a directory name.

### What changed at each site

- **`interview.md:112,120,121`** — the fence caption became "Initial `session.json` — every `*_at` value is the current UTC instant (Timestamp rule, …):", and both values normalized to `"<timestamp>"`. The reviewer is right that this is the canonical template an agent copies; citing the prose at :165 while leaving the template bare was the wrong half to fix.
- **`interview-reference.md:126,133,134,138,140,144`** — same treatment on the "Full shape" schema; the `| null` variants become `"<timestamp> | null"`.
- **`deep-explore-reference.md:317,326`** — `"ISO 8601 timestamp"` → `"<timestamp>"`, cited in the paragraph above the fence. This was three lines below a line I had already edited in the same fence.
- **`stray-check.md:88,92,122,127`** and **`inspect.md:249,253,355,358`** — both report templates in each file, citation above each fence, placeholders normalized. Structurally identical to the `forensics.md`/`roadmap.md` sites, which is exactly why they should have been in the first sweep.

## F3 — `dream.md` misclassification

Accepted; the classification claim in my Scope section was wrong. `skills/do-work-knowledge/actions/dream.md:59` writes a real instant (the `.lock` file, compared against a 5-minute mtime window), so filing the file under "date-only sites" was incorrect even though `dream.md` also contains genuine date-only sites at :162, :207 and :228.

Fixed in place:

```
-   - If absent: create it with the current UTC timestamp inside, then continue.
+   - If absent: create it with the current UTC instant inside (Timestamp rule, `../do-work/actions/work-reference.md`), then continue.
```

**The check cannot see this site and will not catch a regression of it** — it is a prose instruction with no placeholder at all, and the detector keys on placeholders. Stated plainly rather than left implied: see D-11.

## F5 — in-fence YAML comments removed

Accepted. The `lenientFrontmatterFields` path splitting on the first colon without stripping comments is a real hazard, and it is one this change introduced. All four moved to prose above their fence:

| File | Was | Now |
|---|---|---|
| `work-reference.md:620` | `created_at: <timestamp>   # …` | `Its `created_at` is the current UTC instant (Timestamp rule, above).` above the fence |
| `work-reference.md:892` | `session_ended: <timestamp>   # …` | `Its `session_ended` is the current UTC instant (Timestamp rule, above).` above the fence |
| `review-work.md:363` | `created_at: <timestamp>   # …` | citation appended to the sentence introducing the fence |
| `code-review.md:298` | `created_at: <timestamp>   # …` | citation appended to the sentence introducing the fence |

Every generated-artifact template now carries a bare `created_at: <timestamp>` / `session_ended: <timestamp>` line, so nothing this REQ added can reach `parseTimestamp` with a comment attached. As the finding notes, this removes a mechanism rather than adding one — the four templates now use the same pattern as the six report templates.

**Pre-existing instances of the same hazard, reported not fixed:** `capture-reference.md:16,29,142,168`, `work-reference.md:152,165,187,198,207`, and `estimate-reference.md:21` all carry trailing `#` comments on `*_at` lines inside fences. They predate REQ-244 and carry *literal* example values rather than placeholders, so they are a documentation schema rather than a copy-template in most cases — but `capture-reference.md:16` is inside the template capture copies, and it is the same shape F5 objects to. Out of scope here; worth a REQ if the lenient path is considered reachable in practice.

## Files Changed (remediation commit only)

```
 _dev/tests/contract-regressions.sh                 | 99 +++++++++++++++-------
 skills/do-work-knowledge/actions/dream.md          |  2 +-
 .../actions/interview-reference.md                 | 12 +--
 skills/do-work-knowledge/actions/interview.md      |  6 +-
 skills/do-work-toolbox/actions/code-review.md      |  4 +-
 .../actions/deep-explore-reference.md              |  4 +-
 skills/do-work-toolbox/actions/inspect.md          |  8 +-
 skills/do-work-toolbox/actions/stray-check.md      |  8 +-
 skills/do-work/actions/review-work.md              |  4 +-
 skills/do-work/actions/work-reference.md           |  8 +-
 10 files changed, 100 insertions(+), 55 deletions(-)
```

Three files are new to the REQ relative to the original Scope — `stray-check.md`, `inspect.md`, `dream.md` — each named in the remediation brief's scope grant. Nothing else was added. `git diff --name-only` against `do-work/`, `VERSION`, `version.md`, `CHANGELOG.md` and `skills/do-work-board/tools/queue-kanban/` returns nothing.

## Verification

```
maintainer-verify: aggregate contract suite
Timestamp rule citation contract: 54 instant write sites cited, 17 date-only sites recognized.
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
ok  	github.com/knews2019/skill-do-work/queue-kanban	15.887s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (5.62s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.821s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.440s
Maintainer verification passed.
0
```

Run unpiped from the worktree root; `echo $?` printed `0`. The check's success line is now visible in the suite output, which is the F4 fix working.

## Integration Seams

None. Merge the branch as-is.

## Decisions (continuing from D-06)

- **D-07 — Recognition is broad, requirement is narrow.** The detector recognizes any clock-denoting placeholder in any bracket style, but only requires `<timestamp>`/`<now>` plus a citation. This is the split that fixes F2 without freezing the corpus: a spelling nobody anticipated is still *seen*, and being seen is what forces it to be normalized. A check that only recognizes what it requires can never find drift — it can only confirm the absence of drift it already knows about.
- **D-08 — Date-only sites are recognized and counted, but not required to cite.** The rule's Date-only paragraph governs UTC calendar dates while explicitly leaving local dates ungoverned ("changelog entry headings, run-directory names, report slugs… not governed here"). Seventeen sites — `{today}` report headers in `ui-review.md`/`quick-wins.md`/`code-review.md`, the `bkb.md` wiki-log dates, `note.md`'s bullet date — sit on the wrong side of an undecided question: nobody has ruled which are UTC and which are deliberately local. Requiring citations would have forced me to answer it across eight files outside the granted scope. Counting them puts the number in the suite output on every run, so the open question is visible rather than silently in scope. **This is the successor to my earlier Pushback #1, generalized from one file to a measured set.**
- **D-09 — Names are excluded by shape, not by exception.** A `/` inside the placeholder or a `-`/`/` glued outside it means the placeholder names a path or a compound. This is why `work-<YYYY-MM-DD-HHMMSS>` and `v<N>-<date>/` need no carve-out, and it is what let the D-04 respellings from the first commit stay useful rather than becoming a special case.
- **D-10 — A placeholder is short, or it is prose.** Inner text over 30 characters with no nested placeholder is not a placeholder. Without this the detector flags quoted prose containing "now", long bracketed descriptors, and the `(Get-Date).ToUniversalTime()` string in the rule's own home. It is a threshold and therefore the softest thing in the check; if it ever misfires the failure mode is a visible spurious FAIL, not a silent miss.
- **D-11 — The check governs placeholders, so a prose-only stamp instruction is outside it.** `dream.md:59` ("create it with the current UTC instant inside") is now cited, but nothing stops that citation being deleted. Detecting stamp writes described in prose means reading intent, which is a different and much less reliable check. Recorded as a known limit rather than papered over.
- **D-12 — The success line prints both counts.** It is the visible counterpart of the zero-site guard: the guard catches recognition collapsing to zero, the counts catch it collapsing to *fewer*, which is the failure that actually happened here. It is also what makes a GREEN transcript quotable, which is the durable fix for F4.

## Lessons Learned (addendum)

- **A lock-in check written immediately after a fix will tend to encode the fix, not the rule.** Mine matched the exact spellings I had just removed. The generalizing question is not "does this catch what I fixed" but "if this defect recurs in a form I have not seen, does the check see it?" — and the honest way to answer it is to write probe files for spellings *not* in the corpus, which is what the reviewer did and I had not.
- **Test the detector's negatives, not only its positives.** Three of my false positives (a quoted sentence containing "now", `Get-Date` matching `\bdate\b`, the PowerShell format string in the rule's own home) only appeared when I dumped every site the detector saw instead of only the failures. A detector that is never asked "what else did you match?" reports a clean pass that is really a narrow one.
- **A probe fixture with a single failing line can pass or fail for the wrong reason.** My first control run showed every control "caught" — because a file with no valid site trips the zero-site guard before any control logic runs. The controls only became meaningful once each fixture also contained one legitimately cited site. A guard that fires first will happily impersonate the check you meant to test.
- **When quoting evidence, quote the artifact under test.** The F4 error was not a wrong number; it was a right number from the wrong program, presented under a provenance claim. Making the shipped check print what a human would want to quote removes the temptation entirely — if the real thing produces a good transcript, nobody reaches for the prototype's.

## Pushback (addendum)

Nothing in the remediation brief is wrong; all five findings are accepted as stated, and F2 in particular identified a defect I had reasoned my way past in D-03.

One thing to weigh, not an objection: **the date-only count is 17 and will drift upward silently.** D-08 leaves those sites recognized-but-unrequired, which is right while local-vs-UTC is undecided, but it means the suite prints a number nobody owns. If that question gets answered — even just "report headers are deliberately local, log filenames are UTC" — the check can require citations on the UTC half in about ten lines, and the count stops being a parking space. Left as a REQ candidate rather than decided here.
