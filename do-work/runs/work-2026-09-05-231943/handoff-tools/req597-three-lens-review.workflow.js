export const meta = {
  name: 'req597-three-lens-review',
  description: 'Three independent adversarial reviews of REQ-597, then a synthesis with verified findings',
  phases: [
    { title: 'Review', detail: 'three reviewers with distinct lenses' },
    { title: 'Synthesize', detail: 'one synthesizer verifies and scores' },
  ],
}

const REPO = '/home/user/skill-do-work'
const RANGE = '804a8ba32129a3cd12a4aaa7e89346db1b95115c..d5cf28b996a6deb0a0df908cbe4aa722cf2a6ad8'
const MERGE = 'd5cf28b996a6deb0a0df908cbe4aa722cf2a6ad8'
const RUN = `${REPO}/do-work/runs/work-2026-09-05-231943`
const SP = '/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad'

const FINDINGS_SCHEMA = {
  type: 'object',
  properties: {
    findings: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          title: { type: 'string' },
          severity: { type: 'string' },
          file: { type: 'string' },
          line: { type: 'integer' },
          what_is_wrong: { type: 'string' },
          how_i_demonstrated_it: { type: 'string' },
          exact_command_and_output: { type: 'string' },
          suggested_fix: { type: 'string' },
        },
        required: ['title', 'severity', 'file', 'what_is_wrong', 'how_i_demonstrated_it', 'exact_command_and_output'],
      },
    },
    requirements_met: { type: 'array', items: { type: 'string' } },
    requirements_not_met: { type: 'array', items: { type: 'string' } },
    claims_in_the_record_i_checked: { type: 'array', items: { type: 'string' } },
    claims_that_did_not_hold: { type: 'array', items: { type: 'string' } },
    score_percent: { type: 'integer' },
  },
  required: ['findings', 'requirements_met', 'requirements_not_met', 'claims_in_the_record_i_checked', 'claims_that_did_not_hold', 'score_percent'],
}

const COMMON = `You are reviewing REQ-597 in the git repository at ${REPO}. READ-ONLY on that checkout: do not edit,
stage or commit anything there, and leave no file changed. FIRST THING: clone it into a directory of your
own under ${SP} (git clone ${REPO} ${SP}/review597-<yourlens>) and check out ${MERGE}; do every experiment
in that clone. Build the CLI from your clone (cd skills/do-work/tools/do-work-cli && go build -o
${SP}/review597-<yourlens>/dwcli .) and run prescribed blocks through your clone's launcher scripts.

The change is the merge range ${RANGE}: three shipped prose files, 23 insertions and 23 deletions —
skills/do-work-toolbox/actions/inspect.md (two prescribed blocks rewritten so they run at all, plus two
lead-in sentences), skills/do-work/actions/commit.md (exit and row prose), and
skills/do-work/docs/prescribed-shell-primitives.md (sixteen claims re-derived from the code, plus line 24
and line 137). Three builders produced it as three stacked commits; read git log ${RANGE} and each
commit message. The request (its Detailed Requirements list the sixteen claims by line and the caller
findings) is ${REPO}/do-work/working/REQ-597-*.md. The builders' hand-back is ${RUN}/REQ-597-handback.md
and the pre-build verification they worked from is ${RUN}/REQ-597-verification.json. Their fixtures are
reusable under ${SP}/req597-inspect/ (build-fixture.sh and the repo-* fixtures) and ${SP}/req597-guide/.

The repository's rule for this file class, learned the hard way across REQ-555, 595 and 596 (archived under
${REPO}/do-work/archive/): a guide sentence is checked by running the command it describes against a
fixture and comparing, not by reading the code that the sentence paraphrases; a sentence that generalizes
over two commands that behave differently is wrong; compressing a verified sentence is a rewrite; a claim
that was true of the draft and false of the shipped sentence is the author's.

Score out of 100. A claim you did not execute is not a finding. A finding must carry the exact command
and output that demonstrates it.`

phase('Review')
const reviews = await parallel([
  () => agent(`${COMMON}

YOUR LENS: the two action files — do the shipped blocks run, and is every new sentence true?

1. Extract the prescribed blocks from inspect.md and commit.md at ${MERGE} by program (not retyped), and
   run each from a fixture repository's root through the launcher <skill-root>/scripts/protected-inventory.sh
   exactly as the block spells it. Do the same at the parent of the first REQ-597 commit (git log ${RANGE}
   --reverse | head -1, then its parent) to reproduce the "exit 2 unknown option" the request reports.
2. Every sentence the diff adds or changes in these two files: run the case it describes. In particular:
   the wrapper rejects a root argument; from a subdirectory start prints the same rows and associate exits
   2 as if do-work/ were missing; without --quarantine-name the wrapper writes commit's file; associate
   exits 2 with a HELPER-USAGE finding when the quarantine is missing; the exit-status readings for 1 and 2
   in both files; the X-row and quarantine sentences. Report each as held / did not hold with output.
3. The sentences the diff did NOT change around the changed ones: are any of them now false given what
   the changed ones say? Read each step of inspect.md's protected-inventory flow end to end and commit.md's,
   executing as you go.
4. The two files describe the same wrapper. Do they now disagree with each other anywhere?`,
    { label: 'review:actions-executed', phase: 'Review', schema: FINDINGS_SCHEMA }),

  () => agent(`${COMMON}

YOUR LENS: the guide — is each re-derived claim true of the code, by measurement?

For every changed line in skills/do-work/docs/prescribed-shell-primitives.md in ${RANGE} (git diff
${RANGE} -- that file), take the new sentence and the old one, find the Go code it describes, and run a
fixture that distinguishes them. The request lists the sixteen by line number and the verification JSON
gives the drafts the builders were handed; the builders say they rejected several drafts as false — check
that what shipped is not one of them and is not a new compression of the same kind. Specific things to
measure: any sentence about rename/publication errors (os.Rename vs Root.Rename vs rename(2) errno
differences); anything about the merge-aware commit diff and the commit file listing (build a repository
with a root commit, an ordinary commit, a clean merge and a conflicted merge and run the real command);
lifecycle timing claims (run the timing command and read the categories it accepts); the verified-exact
publication rule and the portfolio summary (does a "retained script" exist for it — find it or prove it
does not); atomic download's occupancy behaviour on an existing file versus directory, dry-run versus live;
line 24 (what the launcher does with positionals) and line 137 (report image batch). A sentence true of
one of several commands it covers is a finding. Report each of the changed sentences as held / did not
hold, with the command and output.`,
    { label: 'review:guide-measured', phase: 'Review', schema: FINDINGS_SCHEMA }),

  () => agent(`${COMMON}

YOUR LENS: scope, completeness, guards and the record.

1. Completeness: the request's Detailed Requirements list the claims to correct by line. Map every one to
   a hunk in the diff; any listed claim with no hunk is a finding, and any hunk with no listed claim is
   widening (the builders say line 137 and line 24 are such; judge whether the request's own fourth
   requirement licenses them). Are the caller findings for commit.md and inspect.md all addressed?
2. Scope: exactly three files? Any hunk beyond the stated purpose? Does anything in the diff change a
   heading, the route column, the Mechanics column, or a pointer site that the guards pin
   (_dev/tests/prescribed-shell-canonicalization.sh: twelve headings and sixteen pointer sites;
   _dev/tests/audit-lockins.sh Findings 3 and 7)? Run all four guards in your clone: audit-lockins,
   prescribed-shell-canonicalization, quiet-grep-pipeline-audit, action-shell-blocks; report exit and last line.
3. The hand-back's claims, executed: "exit 2 unknown option" before; exit 0 and the row set after; the
   quarantine file mode 0600; the linked-worktree case writing under .git/worktrees/<name>/; the all-X
   fixture's exit 1; the gate line. Any number that does not reproduce is a finding.
4. The three commits: do the messages describe what each commit does? Is the stacking clean (no commit
   re-edits a line a previous one wrote, no leftover from a rejected draft)?
5. Release: three shipped files, so this is a release. Confirm no version or changelog was touched in the
   range (the finalizer owns that) and that nothing in the diff writes a model identifier or a session URL.
6. The builders' "found and not fixed" list (in the hand-back) names code defects: the shim loop at
   inventory.go:445-456 discarding prepared output, the launcher not passing --repo-root, atomic-download's
   dry-run/live asymmetry and an unchecked os.Stat at commands.go:891, finalization_apply.go:545 diff-tree
   without -m. Reproduce each in your clone and say which are real; the orchestrator is capturing them as
   requests and needs to know which claims hold.`,
    { label: 'review:scope-and-record', phase: 'Review', schema: FINDINGS_SCHEMA }),
])

const bodies = reviews.filter(Boolean)
log(`three lenses returned ${bodies.length} reports`)

phase('Synthesize')
const synthesis = await agent(`${COMMON}

Three independent reviewers have reported on REQ-597. Their reports:

${JSON.stringify(bodies, null, 1)}

VERIFY, then synthesize. Reproduce every finding yourself in your own clone before accepting it. Where
two reviewers disagree, run the experiment that settles it and say which you took and why.

Produce: an overall score out of 100, a per-dimension score (Requirements, Code Quality, Test Adequacy,
Scope, Risk, Acceptance), a verdict, the findings that survived verification ranked most severe first
with their reproduction and a suggested fix (exact replacement sentence where the finding is a sentence),
the findings that did NOT survive and why, a requirements checklist against the request's Detailed
Requirements, and — separately — which of the builders' "found and not fixed" code defects reproduced.
If the change is clean, say so plainly rather than inventing work.`, { label: 'synthesize:req597', phase: 'Synthesize', effort: 'high', schema: {
  type: 'object',
  properties: {
    overall_percent: { type: 'integer' },
    dimensions: {
      type: 'object',
      properties: {
        requirements: { type: 'integer' },
        code_quality: { type: 'integer' },
        test_adequacy: { type: 'integer' },
        scope: { type: 'integer' },
        risk: { type: 'string' },
        acceptance: { type: 'string' },
      },
      required: ['requirements', 'code_quality', 'test_adequacy', 'scope', 'risk', 'acceptance'],
    },
    verdict: { type: 'string' },
    confirmed_findings: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          title: { type: 'string' },
          severity: { type: 'string' },
          impact_token: { type: 'string' },
          file: { type: 'string' },
          line: { type: 'integer' },
          detail: { type: 'string' },
          reproduction: { type: 'string' },
          suggested_fix: { type: 'string' },
        },
        required: ['title', 'severity', 'impact_token', 'file', 'detail', 'reproduction'],
      },
    },
    rejected_findings: { type: 'array', items: { type: 'string' } },
    disagreements_and_how_settled: { type: 'array', items: { type: 'string' } },
    requirements_checklist: { type: 'array', items: { type: 'string' } },
    code_defects_reproduced: { type: 'array', items: { type: 'string' } },
    code_defects_not_reproduced: { type: 'array', items: { type: 'string' } },
    remediation_needed: { type: 'boolean' },
    remediation_brief: { type: 'string' },
  },
  required: ['overall_percent', 'dimensions', 'verdict', 'confirmed_findings', 'rejected_findings', 'disagreements_and_how_settled', 'requirements_checklist', 'code_defects_reproduced', 'code_defects_not_reproduced', 'remediation_needed', 'remediation_brief'],
} })

return synthesis