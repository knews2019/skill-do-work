# Hand-back: REQ-696 (sweep stale wording left by the wave-end review)

- Branch: `worktree-agent-REQ-696-stale-wording-sweep`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-696-stale-wording-sweep`
- Base commit: `fd9a2378`
- Commits: `e8381dc0` `[REQ-696] sweep stale stuck-routing, fence and sibling-route wording` (one commit, by exact path)

## File Manifest

- `skills/do-work/actions/forensics.md` (modified): line 10 no longer routes "stuck" to forensics; it points at `actions/status.md` for quiet in-flight work (REQ-690 F9 text).
- `README.md` (modified): last sentence of the "What happens if something goes wrong during processing?" answer now names `do-work status` for quiet in-flight work and `do-work forensics` for failed or broken work (REQ-690 F9 text).
- `skills/do-work/actions/clarify.md` (modified): line 106 fence rule now says one backtick longer, never shorter than three, no info string (REQ-688 F2 text).
- `skills/do-work/actions/capture.md` (modified): lines 136 and 138 of the queued-addendum example are now `> ` plus a bare three-backtick fence (REQ-688 F3).
- `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (modified): line 1 comment gains the validate-feedback forward-row exception sentence (REQ-692 M2). `# ` prefix kept; line 3 tabs unchanged (checked with `cat -vet`).
- `_dev/tests/staged-skills-contract.sh` (modified): line 792 fail message names the `./actions/` route and allows a `../$sibling_owner/` forward row (REQ-692 M2). `$sibling_owner` is set by the loop's `read` at line 776, so the message expands. Check logic unchanged.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the brief, the REQ (Triage D-01 to D-04), the three review excerpts, `shared-principles.md`, `prime-releases.md`, `prime-action-files.md` (Cross-Referencing: same-package citation `actions/status.md` is the correct form from `actions/forensics.md`), `lessons-releases.md`. Approach: six exact-string replacements with each old string asserted to occur once, then probe, `bash -n`, `git diff --check`.
- [x] **[APPLY]:** Applied the six replacements with the reviewers' exact text via one Python script (each `count(old) == 1` asserted). No other edits.
- [x] **[UNIFY]:**
  - `git diff fd9a2378 --stat`: 6 files changed, 7 insertions(+), 7 deletions(-) (README.md 2, tsv 2, staged-skills-contract.sh 2, capture.md 4, clarify.md 2, forensics.md 2).
  - `bash .../REQ-696-probe.sh` from worktree root: exit 0, wall under 1s, prints `REQ-696 GREEN probe: ok`.
  - `bash -n _dev/tests/staged-skills-contract.sh`: exit 0, wall under 1s.
  - `git diff fd9a2378 --check`: exit 0 (clean).
  - Each changed line read once in the diff: no stray whitespace; capture.md fences are exactly three backticks with no info string; tsv line 3 keeps its tabs.
  - Not run (per brief and D-04): `staged-skills-contract.sh` itself (heavy-only lane, integrator's heavy drain owns it), the repository gate.

## Proof Record

RED at base `fd9a2378` (before any edit):

```
README.md:176:... Run `do-work forensics` to diagnose stuck or failed work.
skills/do-work/actions/forensics.md:5:... feels broken, stuck, or produces confusing results ...
skills/do-work/actions/forensics.md:10:- User suspects something is stuck, broken, or producing confusing results
skills/do-work/actions/forensics.md:49:- stuck work, hollow completions, ...
skills/do-work/actions/forensics.md:159:- A REQ was flagged as `stuck` but ...
clarify.md:106: ... open a code fence longer than the longest backtick run anywhere in the text. ...
capture.md:136:> ````text
capture.md:138:> ````
tsv line 1: # TEST-ONLY historical inventory. Do not copy these retired aliases into runtime routing or guidance.
sh:792:    fail "core must not route sibling-owned action $public_action"
```

Probe at base: exit 1, 11 FAIL lines, `REQ-696 GREEN probe: 11 failure(s)` (matches pre-dispatch).

GREEN after edits: probe exit 0, `REQ-696 GREEN probe: ok`. `grep -n stuck forensics.md README.md` now hits forensics.md lines 5, 49 and 159 only.

## Decisions

- D-05: Applied all six edits in one Python pass that asserts each old string occurs exactly once, instead of six separate Edit calls, so a drifted site would stop the run instead of editing the wrong place. All six sites were still stale at base; none had been fixed by an unrelated commit. DECIDE & STATE.
- D-06: One commit for all six sites (brief prefers one commit; the change is one coherent sweep). DECIDE & STATE.

## Discovered Tasks

- `skills/do-work/actions/capture-reference.md:196` still restates the full fence rule after citing clarify.md Step 4. Now that clarify.md:106 states the complete rule, line 196 could cite it and drop the restatement (REQ-688 review F2 called this optional; the REQ does not name it). → report only
- `skills/do-work/docs/forensics-guide.md:3` ("detects stuck work") and `forensics.md:5` keep "stuck" as something forensics detects, per D-02. If a later reader finds them confusing next to the new routing, a status pointer could be added there. → report only

## Lessons Read

- `_dev/primes/lessons-releases.md` (whole file): families `canonical-link-outlives-its-target` and `manifest-ownership-vs-edit-content`. Neither applies (no record moved, no manifest edited).
- `_dev/primes/lessons-action-files.md`: not read (dropped for budget); every replacement text looked right, so the `restated-mechanism-unchecked` skim was not needed.

## Anti-bloat Check

`git diff fd9a2378 --stat`: 6 files changed, 7 insertions(+), 7 deletions(-), exactly the six write-set paths. Changes not named by the brief: none.

## Proposed CHANGELOG Entry

### Stuck Questions Route to Status, Fence Rule Stated in Full

Shipped prose no longer disagrees with rules changed in the previous run: "is it stuck" goes to `do-work status`, and the containment fence rule reads the same everywhere a session might copy it.

- `actions/forensics.md` and `README.md`: "stuck" questions point at `do-work status`; forensics stays for failed or broken work.
- `actions/clarify.md`: the outside-text containment fence is one backtick longer than the longest backtick run, never shorter than three, with no info string.
- `actions/capture.md`: the queued-addendum example uses a bare three-backtick fence.
- Test text only: the retired-trigger fixture header and the staged-skills contract fail message name the validate-feedback forward-row exception.

## Proposed Lesson Bullet

None. The work was a mechanical application of review-supplied text; nothing here would change what a later builder does.

## Integration Seams

None found. The six files are not in REQ-693, REQ-694 or REQ-695's write sets. The forensics line cites `actions/status.md` by path, which REQ-695 may edit but does not move.

## Check Wall Times

- RED probe: under 1s (exit 1). GREEN probe: under 1s (exit 0). `bash -n`: under 1s (exit 0). `git diff --check`: under 1s (exit 0). No reruns needed.
