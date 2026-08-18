# REQ-238 Hand-Back — Point present-work at the Canonical Independent-Bytes Rationale

**Branch:** `worktree-agent-REQ-238-independent-bytes-pointer`
**Commit:** `3294a7f` — `[REQ-238] Point present-work at the canonical independent-bytes rationale`
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-238-independent-bytes-pointer`

## File Manifest

| File | Verb | What |
|---|---|---|
| `skills/do-work-toolbox/actions/present-work.md` | modified | Step 6's independent-bytes bullet (line 136): restatement replaced by local requirement + pointer; the caller's own policy clause kept verbatim |
| `_dev/tests/prescribed-shell-canonicalization.sh` | modified | one line added to the stale-pattern list |

`git diff --stat` against the branch point: 2 files changed, 2 insertions(+), 1 deletion(-). Nothing outside the write set was touched; `git status --short` showed exactly these two files before staging.

## The Pattern Phrase

**`follow every later in-place edit`**

Why this one:

- It is the rationale itself, not the requirement. Both wordings contain it verbatim — canonical `prescribed-shell-primitives.md:88` ("a snapshot linked to the canonical file would **follow every later in-place edit** of it") and the caller's copy ("would silently **follow every later in-place edit** of it"). A pattern that matched only one phrasing would be inert against a paraphrased re-copy of the same paragraph.
- It is narrow (REQ-230 D-04). It does not touch the neighbouring requirement vocabulary — `identical bytes`, `independent files`, `share storage` — all of which are legitimate for a caller to state about its own outputs. Widening to `share storage` in particular would have been wrong: `do-work-toolbox`-adjacent prose and future helper-policy text may legitimately say a helper's outputs do not share storage; only the *reason* has a single home.
- **Canonical-guide exclusion checked before choosing, not assumed.** The scan loop's first line is `[[ "$shipped_markdown" == "$canonical_guide" ]] && continue` — a full-path skip of `skills/do-work/docs/prescribed-shell-primitives.md`. So a pattern lifted straight out of the canonical sentence cannot fail the file that is supposed to own it. That is what made the exact rationale phrase available as the pattern rather than forcing a weaker near-miss.
- Placed immediately after `'container rather than a collision'` — the other publication-side rationale, added by REQ-230 — so the two publication entries sit together.

**Sized by grep before choosing.** Across `skills/` (the only tree the scan walks) exactly two files carry the phrase: the canonical guide (excluded) and `present-work.md`. Other hits are under `do-work/` — the queue REQ, REQ-230's archive record and hand-back — which the suite does not scan.

## RED (verbatim)

Pattern added to `_dev/tests/prescribed-shell-canonicalization.sh` **first**, prose still unchanged:

```
$ bash _dev/tests/prescribed-shell-canonicalization.sh; echo "EXIT=$?"
FAIL: skills/do-work-toolbox/actions/present-work.md restates canonical prescribed-shell rationale <follow every later in-place edit>; keep local intent and point at the guide.
EXIT=1
```

One failure, naming one file. That is the half REQ-230's builder called out: it proves simultaneously that the pattern fires on the real instance and that it is not over-broad anywhere else in the shipped markdown tree (the scan walks every `*.md` under `skills/`).

## GREEN (verbatim)

Same pattern still in the list, restatement now a pointer:

```
$ bash _dev/tests/prescribed-shell-canonicalization.sh; echo "EXIT=$?"
Prescribed shell primitive canonicalization checks passed.
EXIT=0
```

## The Prose Change

Before (`skills/do-work-toolbox/actions/present-work.md:136`):

> - The two outputs carry identical bytes but are independent files. A snapshot that shared storage with the canonical file would silently follow every later in-place edit of it, which is the opposite of what preserving a snapshot is for, so the helper verifies each output against the source separately and neither output can rewrite the other afterwards.

After:

> - The two outputs must carry identical bytes as independent files, for the reason the canonical [Portfolio summary publication](../../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication) contract gives. Meeting it here is why the helper verifies each output against the source separately and neither output can rewrite the other afterwards.

The **requirement** stays local (identical bytes, independent files — what this action must achieve). The **rationale** goes to its single home. The **caller's policy** survives word for word: *"the helper verifies each output against the source separately and neither output can rewrite the other afterwards"* — REQ-230's D-03 reason applies unchanged, that clause describes this helper's behaviour, not the shared primitive.

## Pointer Resolution — Checked by Hand

The canonicalization suite proves a restatement is *absent*; it never proves the pointer that replaced it *resolves* (REQ-230's review flagged exactly this). Verified manually:

- Path: from `skills/do-work-toolbox/actions/`, `../../do-work/docs/prescribed-shell-primitives.md` resolves to `skills/do-work/docs/prescribed-shell-primitives.md`. Resolved by `cd`-ing into the containing directory and printing `$PWD`, then `ls -l` on the link target through the relative path — the file exists (12435 bytes).
- Anchor: `grep -Fnx '## Portfolio summary publication'` matches at line 84 of that file, so `#portfolio-summary-publication` is a real GitHub-style anchor for a real heading. That heading is also one of the suite's own `required_heading` entries, so it cannot silently disappear.
- Depth cross-check: the identical `../../do-work/docs/…` prefix is already used by the two neighbouring pointers in the same file (lines 130 and 137), and the suite's `sibling_pointer` check requires exactly that string in `present-work.md`.

## P-A-U

### [PLAN]
Read, inside the worktree: `CLAUDE.md`; `_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale*; `skills/do-work/crew-members/maintenance.md` (this REQ is `maintenance: true`) alongside `coding-guardrails.md`, `general.md`, `communication-style.md`; `skills/do-work/docs/prescribed-shell-primitives.md` § *Portfolio summary publication*; the whole of `_dev/tests/prescribed-shell-canonicalization.sh`; `skills/do-work-toolbox/actions/present-work.md` Step 6. Read read-only from the main tree: REQ-238's brief body and REQ-230's archived record — its `## Implementation Summary`, `## Decisions`-equivalent reasoning, `## Testing` and `## Lessons Learned`, which set the accepted shape (pointer mirrors the canonical guide's own sentence; keep caller policy; add pattern before touching prose).

Approach: copy REQ-230's shape exactly rather than invent a second convention one line away. Order fixed up front — grep to size, read the exclusion mechanism, add the pattern, RED, then prose, then GREEN.

This is a maintenance pass, so `maintenance.md`'s "delete before you add" was the first question asked: the fix *is* a deletion (a restatement paragraph removed) plus one line in an existing mechanism. No new rule, no new test file — the added pattern is an entry in a list the suite already maintains, which is the shape this class of finding is pinned with.

### [APPLY]
Both edits landed inside the write set and nowhere else. `git status --short` before staging listed exactly `_dev/tests/prescribed-shell-canonicalization.sh` and `skills/do-work-toolbox/actions/present-work.md`. Staged by explicit path. Untouched, as instructed: `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `skills/do-work/docs/prescribed-shell-primitives.md`, and every file the three sibling builders own (`durations.go`, `durations_test.go`, `web/board-timeline.js`, `web/board.css`, `generate_test.go`). Nothing written under `do-work/` in either tree except this hand-back. No scratch files created.

### [UNIFY]

```
 _dev/tests/prescribed-shell-canonicalization.sh | 1 +
 skills/do-work-toolbox/actions/present-work.md  | 2 +-
 2 files changed, 2 insertions(+), 1 deletion(-)
```

- `skills/do-work-toolbox/actions/present-work.md` — one bullet rewritten. Checked: the pointer resolves and its anchor exists (above); the construction reads as one convention with the `#portfolio-summary-publication` pointer at line 130 and the `#verified-exact-publication` pointer at line 137; the caller's policy clause is intact; no other bullet in Step 6 changed.
- `_dev/tests/prescribed-shell-canonicalization.sh` — one entry added to `stale_patterns_file`. Checked: it is a `printf '%s\n'` argument list with the same `\`-continuation and single-quoting as its neighbours; the phrase contains no backtick or shell metacharacter needing care; matching is `grep -Fq`, so it is a literal, not a regex.

Linters and checks run, all unpiped:

- `shellcheck --severity=warning _dev/tests/prescribed-shell-canonicalization.sh` → exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` → RED exit 1 (pattern only), GREEN exit 0 (pattern + prose).
- `bash _dev/tests/maintainer-verify.sh` → **exit 0**, run directly with `echo $?` on its own line, never through `tail`/`head`. Full suite green including the shipped package reference contract, shell-block lint, the 42 named prescribed-shell script cases, the queue-kanban Go lanes and audit-metrics.

No debug artifacts, no commented-out code, no TODOs left behind.

## Decisions

- **D-01 — Pattern is the rationale phrase `follow every later in-place edit`, not the requirement vocabulary.** Reasoning above. Narrowest phrase that catches both existing phrasings of the reason while leaving a caller free to state *what* its outputs must be.
- **D-02 — Verified the canonical-guide exclusion is a path skip before choosing a phrase drawn from the canonical sentence.** The loop's `[[ "$shipped_markdown" == "$canonical_guide" ]] && continue` compares full paths, so the guide is exempt from every pattern, not from a curated subset. Had the exclusion been per-pattern, the phrase would have needed to be one the guide does not use — which would have made it weaker.
- **D-03 — Kept the caller's own policy clause verbatim.** "the helper verifies each output against the source separately and neither output can rewrite the other afterwards" is `present-work`'s policy, not the shared primitive. REQ-230 recorded this as its own D-03. Deleting a restatement is not the same as deleting policy, and `maintenance.md`'s "subtraction is not vandalism" points the same way.
- **D-04 — Pointed at `#portfolio-summary-publication`, the section that actually owns the sentence,** rather than at `#verified-exact-publication` (the neighbouring bullet's target). The canonical sentence lives in `## Portfolio summary publication` at line 84–88. The repeated pointer to the same section one paragraph above is deliberate, not redundant: the intro line tells the reader to read the contract before invoking, this bullet names where a specific claim's reason lives, and a bullet that a reader lands on directly should not depend on having read the intro.
- **D-05 — Changed "carry" to "must carry".** The old sentence's first clause was descriptive prose paired with the rationale that followed it. With the rationale gone, the clause has to stand as the requirement it always was, so it is stated as one. This is the only wording change beyond removing the restatement.

## Discovered Tasks

- [normal] **Nothing in the suite verifies that a markdown pointer resolves.** REQ-230's review raised this as a suggested-testing item and it is still open — this REQ is the second fix of the class, and it again traded a staleness risk for a broken-link risk that has no detector. Both fixes needed a by-hand check of the same two things: that the relative path normalizes to an existing file, and that the anchor text exists as a heading in it. A check over shipped markdown that resolves every relative `.md` link and every same-repo `#anchor` against the target's headings would be mechanical, would have caught a wrong relative depth in either REQ, and gets cheaper the more of this class ships. Worth its own decision, not a silent add-on.
- [minor] **The `## Portfolio summary publication` sentence is now referenced by name from a caller, which raises the cost of rewording it.** No action today — the suite's `required_heading` list already pins the heading, and the stale-pattern entry now pins the rationale. Recorded so that a future edit to that sentence knows a pointer depends on it.

## Integration Seams

None. Both files are exclusively this REQ's; neither is touched by the three sibling builders (`durations.go`/`durations_test.go`, `web/board-timeline.js`/`generate_test.go`, `web/board.css`/`generate_test.go`). No serial-only file (`VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, either `CHANGELOG.md`) was read into the change or modified — the version bump and changelog entry remain the orchestrator's at Step 9.

One note for the integrator, not a seam: RED is reproducible on the merged tree by holding the pattern fixed and checking out the pre-merge blob of `skills/do-work-toolbox/actions/present-work.md` only — `git checkout <pre-merge> -- skills/do-work-toolbox/actions/present-work.md`, run the suite, expect exit 1 naming that one file. REQ-230's lesson applies: do not use `git stash push` on a clean file to try to reproduce it, because stashing nothing yields a green run that reads as proof and is not.
