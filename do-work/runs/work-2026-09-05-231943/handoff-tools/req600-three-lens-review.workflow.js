export const meta = {
  name: 'req600-three-lens-review',
  description: 'Three independent adversarial reviews of REQ-600, then a synthesis with verified findings',
  phases: [
    { title: 'Review', detail: 'three reviewers with distinct lenses' },
    { title: 'Synthesize', detail: 'one synthesizer verifies and scores' },
  ],
}

const REPO = '/home/user/skill-do-work'
const RANGE = '9e00a092cf29842506bea920137b52c952a62638..a25c7522566bea9d9d29c382e159b6a10157a9f1'
const MERGE = 'a25c7522566bea9d9d29c382e159b6a10157a9f1'
const REQ = `${REPO}/do-work/working/REQ-600-put-the-sigpipe-trap-where-shell-authors-read-it.md`
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

const COMMON = `You are reviewing REQ-600 in the git repository at ${REPO}. READ-ONLY on that checkout: do not
edit, stage or commit anything there, and leave no file changed. FIRST THING: clone it into a directory of
your own under ${SP} (git clone ${REPO} ${SP}/review600-<yourlens>) and check out ${MERGE}; do every
experiment and every mutation in that clone. The main checkout may briefly carry unrelated in-progress
edits, so read the files from your clone, not from ${REPO}.

The change is the merge range ${RANGE} (git diff --stat gives 5 files, +155/-77): a new section in
_dev/primes/prime-shell-commands.md; one prescribed block rewritten in
skills/do-work-knowledge/actions/memory-reference.md (lines 86-92); the scanner function
quiet_grep_pipeline_offenders lifted out of _dev/tests/quiet-grep-pipeline-audit.sh into a new
_dev/tests/quiet-grep-pipeline-scanner.sh that both the audit and _dev/tests/action-shell-blocks.sh now
source; and action-shell-blocks.sh running that scanner over every extracted Markdown fence and shipped
shell file, with a wiring fixture. The record is ${REQ}; read it whole, including Exploration, Scope,
Implementation Summary, Decisions and Discovered Tasks. The builder's hand-back is
${REPO}/do-work/runs/work-2026-09-05-231943/REQ-600-handback.md. The facts the prime section cites come
from the archived records of REQ-593 and REQ-594 under ${REPO}/do-work/archive/ (find them with ls).

The repository's own rule for shell: ${REPO}/_dev/primes/prime-shell-commands.md, and its naming rule:
${REPO}/skills/do-work/crew-members/coding-guardrails.md section 5 (two words minimum for anything with
reach, findable by plain-text search).

Score out of 100. A claim you did not execute is not a finding. A finding must carry the exact command
and output that demonstrates it.`

phase('Review')
const reviews = await parallel([
  () => agent(`${COMMON}

YOUR LENS: is the prime section true, and is it the section a shell author needs?

1. Read the new section "A Writer's SIGPIPE Death Reads as the Reader's Verdict" line by line. For every
   factual claim, find the evidence in REQ-593's and REQ-594's archived records or reproduce it. In
   particular: the measured window (the record says 0 of 50 misfires at 36 KB and 50 of 50 at 200 KB) —
   reproduce it in your clone with a small bash harness under set -o pipefail (a writer that emits N bytes
   then a marker, piped into grep -q for something early in the stream, repeated 50 times at each size) and
   report what you measured. If the section states a number the records do not carry or your measurement
   contradicts, that is a finding.
2. "Wrong in both directions": does the section actually show both directions, and is the negative
   matcher (! cmd | grep -q, or grep -q used to prove absence) correctly identified as the dangerous half?
   Demonstrate each direction once.
3. The two-half fix: is the herestring form the section prescribes correct under set -u and pipefail, and
   does the section say that the capture discards the producer's status? Test the one-line
   listing="$(...)" || fail form: does it actually fail on a producer failure?
4. The guard pointer: does what the section says about quiet-grep-pipeline-audit.sh (tracked *.sh only,
   grep-only reader set, Markdown not its input) and about action-shell-blocks.sh (now runs the same
   scanner over the fences) match what those scripts do at ${MERGE}? Read their code, not their comments.
5. Placement and shape: is it beside "Unchecked Exit Status Reads as Content" as the request requires, and
   does it read as a trap entry in this file's own style (compare the neighbours)? Does the prime's Lessons
   section need a line, or does the record's claim that lessons are appended on archive hold — check
   what the archive step actually does with _dev/primes/lessons-shell-commands.md (search the skills/ tree
   for who writes it).`,
    { label: 'review:prime-truth', phase: 'Review', schema: FINDINGS_SCHEMA }),

  () => agent(`${COMMON}

YOUR LENS: the scanner lift and the new lint — is nothing weaker, and does the new scan really run?

1. Byte-for-byte: extract quiet_grep_pipeline_offenders (function body plus its contract comment) from
   quiet-grep-pipeline-audit.sh at the parent commit 9e00a092 and from quiet-grep-pipeline-scanner.sh at
   ${MERGE}, and diff them. The record says they are identical. Report the diff.
2. Sourcing side effects: does sourcing quiet-grep-pipeline-scanner.sh run anything, set -e/-u/-o
   pipefail, change the caller's shell options, or leave variables behind? Source it in a clean bash and
   compare set -o and declare -p before and after.
3. Fixture parity: at the parent, the audit had 19 must-flag and 7 must-not-flag shapes. Confirm the
   same 26 shapes exist at the merge and each still produces the same verdict (run the audit; also run the
   audit's self-test if it has one).
4. Mutations, each in your clone and each reverted afterwards: (a) remove the scanner call from
   lint_shell_source in action-shell-blocks.sh — does the lint exit 1 and say why; (b) restore the
   pre-change memory-reference.md block from 9e00a092 — does the lint exit 1 at line 88 of the Markdown
   file, and is the reported line number right (count it); (c) delete --quiet from the option class in the
   scanner — does the audit exit 1; (d) NEW: add a fence to some shipped .md in your clone with the shape
   "cmd | grep -q x" at an indent of 3 spaces and again inside a fenced block nested in a list item — does
   the lint catch both, and at the right line; (e) add the shape to a shipped .sh under skills/ — is it
   caught by the lint AND by the audit, or only one.
5. The record says the scan runs "before the shellcheck gate, so it runs even where shellcheck is
   unavailable". Prove it: run the lint with shellcheck hidden from PATH (PATH without it, or a stub that
   does not exist) against mutation (b) and show the FAIL still appears.
6. The wiring fixture run_quiet_grep_wiring_fixture runs on every default invocation. Does it leave any
   temp file behind? Does it work when TMPDIR is unset or read-only? Is its FAIL text reachable (mutation a)?
7. Read the diff of action-shell-blocks.sh hunk by hunk: is every new line needed; is every new name
   two words and findable; does shellcheck --severity=warning pass on all three _dev/tests files?`,
    { label: 'review:lint-and-scanner', phase: 'Review', schema: FINDINGS_SCHEMA }),

  () => agent(`${COMMON}

YOUR LENS: the shipped block, scope, naming, and every number in the record.

1. The block at memory-reference.md:86-92. Run it as shipped in three environments in your clone: no
   ollama on PATH; an ollama stub that exits 1 with nothing on stdout; an ollama stub that prints a
   listing containing nomic-embed-text. Report the exit status of the grep line in each. Then judge the
   "|| true": the record says the producer's failure is the "no backend" answer the surrounding prose
   prescribes — read the surrounding prose (the Semantic Recall section) and say whether that is what it
   prescribes. Would the sibling prime section ("Unchecked Exit Status Reads as Content") accept this
   form, and does the in-line comment say why as that section requires?
2. Is the new block still consistent with the other two probes beside it (command -v embed, the API key
   test) in style, alignment and how a reader is told to use the three results?
3. Every number in the record and hand-back, executed: 166 .md files under skills/, 32 with a fence, 74
   fences, all bash, deepest indent 3, none unterminated, zero sh/shell fences; "the only other grep -q in
   any fence is memory-reference.md:142, a file-argument grep"; the supplementary grep for rg -q, head,
   sed -n ... q, awk ... exit, read, grep -m after a pipe found zero in the 74 blocks — rerun it; "+155/-77";
   "95 tracked shell files (94 plus the helper)"; "33 shipped .sh under skills/"; "ollama list prints
   under 2 KB on any realistic install" (estimate from ollama's own list format: bytes per row).
4. Scope: exactly the five files? Any hunk beyond the stated purpose? The request's Constraint says "no
   change to the guard, whose fixture pins 19 shapes" and the record reads that as fixture and behaviour,
   not bytes. Is that reading defensible, and does the Constraints section's second bullet in fact license
   the action-shell-blocks.sh change?
5. Naming for reach: quiet-grep-pipeline-scanner.sh, quiet_grep_pipeline_offenders,
   run_quiet_grep_wiring_fixture, ollama_models, quiet_grep_offender(s), quiet_grep_offender_line — each
   against coding-guardrails.md section 5.
6. Release: memory-reference.md is a shipped file, so this is a release. Is that stated correctly in the
   record, and are the changelog/version mirrors correctly left to the finalizer (look at how the last
   three "[REQ-NNN] release:" commits did it: git log --oneline -30)?
7. Discovered Tasks: the record says lessons-shell-commands.md has no entry for REQ-593 or REQ-594
   although both are archived. Verify: grep the file; then find what actually appends to it (search
   skills/do-work for the file name and for "lessons") and say whether the pipeline was ever supposed to
   append for those two — was it a step that was skipped, or a step that does not exist?`,
    { label: 'review:block-scope-record', phase: 'Review', schema: FINDINGS_SCHEMA }),
])

const bodies = reviews.filter(Boolean)
log(`three lenses returned ${bodies.length} reports`)

phase('Synthesize')
const synthesis = await agent(`${COMMON}

Three independent reviewers have reported on REQ-600. Their reports:

${JSON.stringify(bodies, null, 1)}

VERIFY, then synthesize. Reproduce every finding yourself in your own clone before accepting it. Where
two reviewers disagree, run the experiment that settles it and say which you took and why.

Produce: an overall score out of 100, a per-dimension score (Requirements, Code Quality, Test Adequacy,
Scope, Risk, Acceptance), a verdict, the findings that survived verification ranked most severe first
with their reproduction, the findings that did NOT survive and why, and a requirements checklist against
the request's Detailed Requirements and Scope acceptance criteria. If the change is clean, say so plainly
rather than inventing work.`, { label: 'synthesize:req600', phase: 'Synthesize', effort: 'high', schema: {
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
          detail: { type: 'string' },
          reproduction: { type: 'string' },
          suggested_fix: { type: 'string' },
        },
        required: ['title', 'severity', 'impact_token', 'detail', 'reproduction'],
      },
    },
    rejected_findings: { type: 'array', items: { type: 'string' } },
    disagreements_and_how_settled: { type: 'array', items: { type: 'string' } },
    requirements_checklist: { type: 'array', items: { type: 'string' } },
    remediation_needed: { type: 'boolean' },
    remediation_brief: { type: 'string' },
  },
  required: ['overall_percent', 'dimensions', 'verdict', 'confirmed_findings', 'rejected_findings', 'disagreements_and_how_settled', 'requirements_checklist', 'remediation_needed', 'remediation_brief'],
} })

return synthesis