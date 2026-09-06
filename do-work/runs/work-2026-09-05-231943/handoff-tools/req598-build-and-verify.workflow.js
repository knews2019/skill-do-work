export const meta = {
  name: 'req598-build-and-verify',
  description: 'Build REQ-598 (decide the rollback handle once at its open; six-kind no-handle test; pin at zero) in a worktree, then verify independently',
  phases: [
    { title: 'Build', detail: 'one builder in its own worktree, following the judged plan step by step' },
    { title: 'Verify', detail: 'one verifier re-derives every claim in a fresh clone' },
  ],
}

const REPO = '/home/user/skill-do-work'
const BRANCH = 'claude/do-work-queue-drain-4ee2xl'
const WT = '/home/user/skill-do-work-worktrees/worktree-agent-REQ-598-decide-once'
const WT_BRANCH = 'worktree-agent-REQ-598-decide-once'
const SP = '/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad'
const CLI = 'skills/do-work/tools/do-work-cli'
const PKG = `${CLI}/internal/gittransaction`
const HANDBACK = `${REPO}/do-work/runs/work-2026-09-05-231943/REQ-598-handback.md`

const RULES = `Repository rules that bind you: ${REPO}/CLAUDE.md (delete before you add; machinery is not an
achievement), ${REPO}/skills/do-work/tools/do-work-cli/prime-do-work-cli.md and its satellite
lessons-do-work-cli.md, ${REPO}/_dev/primes/prime-shell-commands.md (for the lock-in edit: no quiet grep
fed from a pipeline; a guard counts shape, not text), ${REPO}/skills/do-work/crew-members/coding-guardrails.md
(section 5: two-word names minimum, findable by plain-text search). Every commit message ends, after a
blank line, with exactly these two trailer lines:
Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01WAfirqbotH8v8zUnEyJC1d
Never write a model identifier anywhere else. Do not touch CHANGELOG.md, VERSION or version.md: the
finalizer owns them.`

phase('Build')
const build = await agent(`You are the builder for REQ-598 in the do-work repository. Read the request whole first:
${REPO}/do-work/working/REQ-598-*.md — its Plan, Exploration, Scope and Pre-Flight sections describe the
judged plan (Plan B, decide the rollback handle once at its open) and the ordered steps. The judge's full
verdict is ${SP}/req598-verdict.json (read it: what_changes, the_no_handle_test, ordered_steps,
lock_in_change, behaviour_preserved_evidence_plan, risks). The verified tree is ${SP}/req598-final.patch
(applies cleanly on the current branch head) and the exact test is ${SP}/final_test_func.go.

SETUP. cd ${REPO} && git worktree add -b ${WT_BRANCH} ${WT} ${BRANCH}. Work, test and commit only in ${WT}.
Go module root: ${WT}/${CLI}. Do not edit anything under ${REPO} except the hand-back named below.

${RULES}

FOLLOW THE ORDERED STEPS, TAKING YOUR OWN EVIDENCE — the patch is the reference answer, not a substitute
for the steps. Concretely:
1. Clean-head record: in ${WT}/${CLI}, go test -count=1 ./internal/gittransaction/ (expect ok), go vet
   ./..., and note the current Finding 3 pin (8) by running bash _dev/tests/audit-lockins.sh from ${WT}.
2. Seam and test first: add openRollbackRoot and the syscall import and the test
   (from final_test_func.go, name TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest,
   placed before TestPrivatePreimageSnapshotRejectsReplacementDuringOpen). With ONLY that applied, run the
   test: it must fail with a nil-pointer panic at root.Mkdir in quarantineAndRollbackPrivate via
   rollbackFailure and ExecuteTransaction. Save the exact panic lines. Commit this as commit 1 ("[REQ-598]
   Pin the no-handle rollback: seam and six-kind test, red on the current code").
3. The restructure: apply the rest of the patch's git_transaction.go change (rollbackWithRoot,
   rollbackWithoutRoot, four helpers moved out, eight guards deleted, rootedOpenSnapshot's comment
   corrected, the rollbackFailure open through the seam). Run: the new test green; the whole package
   green with zero edits to existing tests; go vet; gofmt -l empty; GOOS=windows go vet ./...; go test
   -race ./internal/gittransaction/. Then the evidence the plan requires: (a) a temporary panic canary at
   the top of rollbackWithoutRoot — run the suite with -skip for the new test and show it never fires,
   then run the new test and show it fires; remove the canary; (b) git diff -U0 of git_transaction.go
   against the branch base restricted to removed lines matching the guard shape: exactly eight, list them;
   (c) rg for the guard shape on the tree exits 1. Commit as commit 2.
4. The lock-in: rewrite Finding 3 in ${WT}/_dev/tests/audit-lockins.sh per lock_in_change (pin at zero,
   scan pattern unchanged, rg exit 1 green, exit 0 fails listing sites, exit >1 fails loudly, header says
   why 8 -> 0 and what the scan cannot see). Prove it: passes at zero; re-add one guard in
   rollbackDirtyTracked and show the FAIL line naming the site; restore, diff the restore. bash -n and
   shellcheck --severity=warning on the file. Commit as commit 3.
5. Gate from the worktree: DO_WORK_GATE_ROOT=${WT} bash ${SP}/gate.sh — record exit and last line. If it
   is red, say exactly which test and whether it is one of the two heavyverification tests the plan calls
   environmental (they should NOT be red under this wrapper; if they are, report it, do not explain it away).
6. Hand-back to ${HANDBACK}: head hash, the three commits, the red panic text, the green results, the
   canary and -U0 evidence, the lock-in mutation, the gate line, every number a verifier needs, and a
   "Found and not fixed" list. Do not merge.

Return as final text: head hash, the three commit hashes, gate exit and last line, and any step where
your evidence differed from the plan's.`,
  { label: 'build:req598', phase: 'Build', effort: 'high' })

log(`builder returned: ${String(build).slice(0, 300)}`)

phase('Verify')
const verify = await agent(`You are the independent verifier for REQ-598. The builder reports:

${build}

Hand-back: ${HANDBACK}. Request: ${REPO}/do-work/working/REQ-598-*.md. Judge's verdict:
${SP}/req598-verdict.json. Branch ${WT_BRANCH} in worktree ${WT}. READ-ONLY on ${WT} and ${REPO}: clone
${REPO} into ${SP}/verify598 and check out ${WT_BRANCH}; module root is ${SP}/verify598/${CLI}.

${RULES}

Re-derive, do not re-read:
1. RED: check out the builder's commit 1 (seam and test only) and run the new test: it must panic at
   root.Mkdir. Then at head: green. Then the whole package, -race, vet, gofmt, GOOS=windows vet.
2. Behaviour preservation, by your own means, not the builder's: (a) add the canary yourself and run
   the suite minus the new test; (b) compute the -U0 removed-lines set matching the guard shape between
   the branch base and head: exactly the eight sites REQ-558 kept (list them by function); (c) diff the
   bodies of the three moved loops against their pre-change text modulo the removed guard tokens — is
   anything else different? (d) differential: write a throwaway test that runs a six-kind fixture through
   ExecuteTransaction with the seam set to succeed, at base and at head, and compare Rollback.Actions,
   Rollback.Errors, outcome and file bytes.
3. Attack the new test: which single removals from rollbackWithoutRoot does it catch (remove the
   dirty-tracked unstage; remove the created-paths git rm --cached; remove one left-in-place error) — each
   must fail it; if one does not, that is a finding.
4. The lock-in: Finding 3 passes at zero; re-adding a guard fails naming the site; a guard written as
   "nil == root" and as "} else if root == nil {" are both caught; a comment containing "root == nil" is
   not; shellcheck clean; no quiet grep fed from a pipeline (run _dev/tests/quiet-grep-pipeline-audit.sh).
5. The record's and hand-back's numbers (+155/-81, +112, +27/-41, eight, six kinds, exact error list
   order) — recompute each.
6. Names for reach (openRollbackRoot, rollbackWithRoot, rollbackWithoutRoot, the moved helpers) against
   coding-guardrails section 5; every deleted comment sentence really described a thing that no longer
   exists; rootedOpenSnapshot's new comment true of the new code.
7. Gate in your clone: DO_WORK_GATE_ROOT=<clone> bash ${SP}/gate.sh; exit and last line.

Return: claims checked with pass/fail, every discrepancy with exact command and output, and a verdict:
mergeable as is, or not, and why.`,
  { label: 'verify:req598', phase: 'Verify', effort: 'high' })

return { build, verify }