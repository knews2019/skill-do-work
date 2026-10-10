# REQ-686 hand-back

- Branch: worktree-agent-REQ-686-pushed-back-finding-doc-lines
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-686-pushed-back-finding-doc-lines
- Base: b629e5cd. Commit: adcdeb65 (one commit).

## File manifest
- skills/do-work/tools/do-work-cli/lessons-do-work-cli.md (modified): one sentence added on the `internal/finalization/` bullet (line 121): after a performed and verified `git revert` of the primary commit, delete the stuck journal `.git/do-work-finalization/REQ-N.json` (and its `REQ-N.payloads` folder) by hand, because `--discard-journal` only accepts phase `prepared`.
- skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go (modified): one comment line above `Handlers()`: "The three fast-stage commands serve _dev/tests/maintainer-verify.sh only; no action calls them."

## P-A-U
- [PLAN]: add the two lines in the two named files; no other change.
- [APPLY]: done. First try put the comment inside the const block; gofmt then wanted three const lines realigned, so I moved it above `Handlers()` (no reflow, one line).
- [UNIFY]: `git diff b629e5cd --stat`: 2 files, 2 insertions(+), 1 deletion(-) (the lessons line is replaced in place). Checks: REQ-686-probe.sh exit 0 (about 1.6s; includes both additions, gofmt, build, shipped-package-reference-contract PASS); `gofmt -l` empty; `go vet ./internal/heavyverification/` exit 0 (0.2s); `git diff --check` exit 0. Files checked: both manifest files, diff read in full. No rerun needed.

## Lessons file size (integration seam)
`wc -c lessons-do-work-cli.md` = 76181 bytes; tokens = (76181+3)//4 = 19046 (was 18987 in the REQ). Update the row in do-work/lessons-index.md.

## Caller evidence
- _dev/tests/maintainer-verify.sh:160 `decide-fast-stage --stage ...`
- _dev/tests/maintainer-verify.sh:186 `invalidate-fast-stage --stage ...`
- _dev/tests/maintainer-verify.sh:199 `record-fast-stage --stage ... --fingerprint ...`
- No other hit in *.sh, *.md, *.json, *.go outside do-work/ (archive and runs excluded). Claim holds.
- Sentence claim checked: finalization_commands.go (discardPreparedJournal comment and error text) removes only phase `prepared` journals.

## Decisions
- D-01 DECIDE & STATE: comment placed above `Handlers()` instead of inside the const block, to avoid a gofmt realignment of three lines.
- D-02 DECIDE & STATE: no actions/*.md home exists for FINALIZATION-PRIMARY-COMMIT, so the lessons bullet is the one home, as the brief said.

## Discovered Tasks
None.

## Lessons read
_dev/primes/prime-releases.md (via brief), REQ-686, brief. Did not read lessons-do-work-cli.md beyond line 121 (no code in its bullets touched).

## Anti-bloat check
Diff stat above. Added beyond what the REQ named: none.

## Proposed CHANGELOG entry
Title: Doc lines record the manual exit after a reverted finalization and the callers of the fast-stage commands
Why: a later reader no longer has to re-ask how to clear a stuck journal after a revert, or who calls the three fast-stage commands.
- lessons-do-work-cli.md: delete the stuck journal by hand after a verified `git revert` of the primary commit.
- heavy_commands.go: one comment says the fast-stage commands serve maintainer-verify.sh only.
(CHANGELOG 0.305.46 wording fix dropped by decision.)

## Proposed lesson bullet
none
