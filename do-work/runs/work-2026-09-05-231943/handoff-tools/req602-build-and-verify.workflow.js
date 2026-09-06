export const meta = {
  name: 'req602-build-and-verify',
  description: 'Build REQ-602 (repoint fifteen lesson-satellite links, add a satellite link check) in an isolated worktree, then verify independently',
  phases: [
    { title: 'Build', detail: 'one builder in its own worktree' },
    { title: 'Verify', detail: 'one verifier re-runs every claim' },
  ],
}

const REPO = '/home/user/skill-do-work'
const BRANCH = 'claude/do-work-queue-drain-4ee2xl'
const WT = '/home/user/skill-do-work-worktrees/worktree-agent-REQ-602-satellite-links'
const WT_BRANCH = 'worktree-agent-REQ-602-satellite-links'
const SP = '/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad'
const REQ_GLOB = `${REPO}/do-work/working/REQ-602-*.md (or ${REPO}/do-work/queue/REQ-602-*.md if the claim has not landed yet; same content)`
const HANDBACK = `${REPO}/do-work/runs/work-2026-09-05-231943/REQ-602-handback.md`

const RULES = `Repository rules that bind you: ${REPO}/CLAUDE.md (delete before you add; programs beat prose;
state conditions, not lists), ${REPO}/_dev/primes/prime-shell-commands.md (read it whole before writing any
shell — in particular the quiet-grep-from-a-pipeline shape is forbidden and a guard scans for it),
${REPO}/_dev/primes/prime-action-files.md, ${REPO}/skills/do-work/crew-members/coding-guardrails.md section 5
(two-word names minimum for anything with reach, findable by plain-text search). The commit trailer every
commit must end with, verbatim, after a blank line:
Co-Authored-By: Claude Fable 5.1 <noreply@anthropic.com>
Claude-Session: https://claude.ai/code/session_01WAfirqbotH8v8zUnEyJC1d
Never write a model identifier anywhere else.`

phase('Build')
const build = await agent(`You are the builder for REQ-602 in the do-work repository. Read the request first: ${REQ_GLOB}.

SETUP. The main checkout at ${REPO} is on branch ${BRANCH}; do not edit files there except the hand-back
named below. Create your own worktree: cd ${REPO} && git worktree add -b ${WT_BRANCH} ${WT} ${BRANCH}
(if the worktree already exists, use it). Work, test and commit only inside ${WT}.

${RULES}

THE WORK, in this order:
1. Read ${WT}/_dev/tests/audit-lockins.sh end to end to learn how its findings are numbered, how each
   prints FAIL lines and exits, and how it is invoked by the gate. Read the three satellites named in the
   request's write_set and ${WT}/do-work/lessons-index.md (its header states how the tokens column is computed).
2. Add the check FIRST: every relative Markdown link in a maintainer lesson satellite (condition: a
   satellite is any _dev/primes/lessons-*.md; a relative link is a "](" target that is not an absolute
   URL and not a pure "#fragment") must resolve, with any "#fragment" stripped, from the satellite's own
   directory. Key it on the condition, not on a list of files or targets. Print one FAIL line per broken
   link naming the satellite, the line number and the target. Run it: it must exit 1 and name all fifteen
   broken links the request lists. Record the exact output in the hand-back (that is the red proof).
3. Repoint the fifteen links to the paths the request lists, changing nothing else in those bullets: not
   the text, not the family marker. Run the check: exit 0. Then a mutation: change one repointed link
   back, run the check, expect exit 1 naming exactly that link; restore it; run again, exit 0; diff the
   restore against the committed file to prove it is byte-identical.
4. The shipped satellite ${WT}/skills/do-work/tools/do-work-cli/lessons-do-work-cli.md uses canonical
   URLs of the form https://github.com/knews2019/skill-do-work/blob/main/<repo path>. For each such URL,
   check whether <repo path> exists in the tree (git ls-files). Report how many URLs there are, how many
   resolve, and list every one that does not. Do NOT edit that shipped file: it is a release-controlled
   file and the orchestrator will decide; put the list in the hand-back.
5. Recompute the three lessons-index.md rows for the satellites whose bytes changed: tokens =
   (bytes + 3) integer-divided by 4 (the header's own formula), families = the exact sorted set of
   [family: slug] markers, coverage as the header defines. Do it with a small program, not by hand, and
   show the before/after rows in the hand-back.
6. Run every guard that reads what you touched: bash _dev/tests/audit-lockins.sh, and the full fast gate
   from your worktree: DO_WORK_GATE_ROOT=${WT} bash ${SP}/gate.sh — it must end with the line
   "Maintainer verification passed." and exit 0. Record the exit and last line.
7. Commit on ${WT_BRANCH}: one commit for the check (red proof in the message is not needed, the
   hand-back carries it) and one for the repointing plus index rows, or one commit for all of it if you
   judge the split artificial — say which and why. Subject lines start with "[REQ-602] ". Do not merge;
   the orchestrator merges.
8. Write the hand-back to ${HANDBACK} (in the MAIN checkout's run directory): head hash, commits, the red
   output, the green output, the mutation evidence, the shipped-satellite URL report, the index rows
   before/after, the gate line, and a "Found and not fixed" list. Be exact: every number in it will be
   re-derived by a verifier.

Return a short JSON-like summary as your final text: branch head hash, the commits, the check's name in
audit-lockins.sh, and the count of shipped-satellite URLs that do not resolve.`,
  { label: 'build:req602', phase: 'Build', effort: 'high' })

log(`builder returned: ${String(build).slice(0, 300)}`)

phase('Verify')
const verify = await agent(`You are the independent verifier for REQ-602. The builder reports:

${build}

The builder's hand-back is ${HANDBACK}; the request is ${REQ_GLOB}; the worktree is ${WT} on branch
${WT_BRANCH}. READ-ONLY on ${WT} and ${REPO}: experiment in a clone under ${SP}/verify602 (git clone
${REPO} then check out ${WT_BRANCH}).

${RULES}

Re-derive, do not re-read: (1) at the parent of the builder's first commit, the fifteen listed links are
broken and no others in the three satellites are; (2) the new check exits 1 there naming exactly those
fifteen, and exits 0 at the branch head; (3) the check is keyed on the condition — add a new satellite
_dev/primes/lessons-zzz.md with one broken relative link in your clone and confirm it is caught; add a
link with a #fragment to an existing target and confirm it is NOT flagged; add an absolute https URL
and confirm it is not flagged; (4) no bullet text or family marker changed: diff the three satellites
between parent and head restricted to non-link characters, or show that every changed line differs only
in its link target; (5) the lessons-index.md rows: recompute tokens and families yourself and compare;
(6) the builder's shipped-satellite URL report: recompute it yourself; (7) the check's shell contains no
quiet grep fed from a pipeline and passes bash _dev/tests/quiet-grep-pipeline-audit.sh; (8) names in the
new code are two words and findable; (9) the gate claim: run DO_WORK_GATE_ROOT=<your clone> bash
${SP}/gate.sh and report exit and last line.

Return: a list of claims checked with pass/fail, every discrepancy with the exact command and output,
and a verdict: mergeable as is, or not, and why.`,
  { label: 'verify:req602', phase: 'Verify', effort: 'high' })

return { build, verify }