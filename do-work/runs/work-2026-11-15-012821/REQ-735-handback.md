# REQ-735 Hand-back — moved standing rules: which load in ordinary sessions

Branch `worktree-agent-REQ-735-moved-rules-always-on-scope`, base `547d4cc5`, tip `82db9eb2`.

## Commits

- `82db9eb2` [REQ-735] Move the numeric-target and redesign rules into the always-on file; action rules win over ask and model rules

## Files

- `skills/do-work/crew-members/communication-style.md` — modified (R1 and R2 bullets added after the bug-report bullet; ask-tool bullet gains "unless the action in use says how to ask")
- `skills/do-work/crew-members/coding-guardrails.md` — modified (R1 line deleted from § 1, R2 line deleted from § 3)
- `skills/do-work-toolbox/crew-members/coding-guardrails.md` — modified (same two deletions)
- `skills/do-work/crew-members/background-agents.md` — modified (model sentence gains "unless the action names another model", rewrapped from line 3 of the paragraph)
- `skills/do-work-toolbox/crew-members/background-agents.md` — modified (byte-identical paragraph)
- `skills/do-work-knowledge/crew-members/background-agents.md` — modified (byte-identical paragraph)
- `skills/do-work/actions/commit.md` — modified (one new bullet in § Rules after the `git add -A` bullet)

## Scope changes

None. Exactly the seven Scope files, +10/-11.

## Red-green evidence

RED at base `547d4cc5` (worktree before editing), `bash .../REQ-735-probe.sh`, exit 1, 7 FAIL lines:

```
FAIL: R1: 'unless the user says strict' expected only in skills/do-work/crew-members/communication-style.md, found: skills/do-work-toolbox/crew-members/coding-guardrails.md skills/do-work/crew-members/coding-guardrails.md
FAIL: R2: 'while redesigning' expected only in skills/do-work/crew-members/communication-style.md, found: skills/do-work-toolbox/crew-members/coding-guardrails.md skills/do-work/crew-members/coding-guardrails.md
FAIL: F3-ask: 'unless the action in use says how to ask' expected only in skills/do-work/crew-members/communication-style.md, found: nothing
FAIL: F5: 'change you did not make' expected only in skills/do-work/actions/commit.md, found: nothing
FAIL: moved line expected once in communication-style.md, found 0: - Read a number in the request (a size, a count, a time) as approximate unless the user says strict.
FAIL: moved line expected once in communication-style.md, found 0: - Never remove existing UI while redesigning unless the user asked for it.
FAIL: F3-model: the model paragraph lacks 'keep the session model unless the action names another model'
REQ-735 probe FAIL (7)
```

GREEN (run on the edited tree, then again at tip `82db9eb2`), exit 0:

```
PASS: R1 only in skills/do-work/crew-members/communication-style.md
PASS: R2 only in skills/do-work/crew-members/communication-style.md
PASS: R3 only in skills/do-work-knowledge/crew-members/background-agents.md skills/do-work-toolbox/crew-members/background-agents.md skills/do-work/crew-members/background-agents.md
PASS: F3-ask only in skills/do-work/crew-members/communication-style.md
PASS: F5 only in skills/do-work/actions/commit.md
PASS: moved line present once: - Read a number in the request (a size, 
PASS: moved line present once: - Never remove existing UI while redesig
PASS: model paragraph identical in skills/do-work-toolbox/crew-members/background-agents.md
PASS: model paragraph identical in skills/do-work-knowledge/crew-members/background-agents.md
PASS: F3-model clause in the model paragraph (identical in all three copies)
PASS: R3 has no copy or pointer in communication-style.md
PASS: ## Aliases unchanged
PASS: ## Examples unchanged
PASS: no machine, vendor, model or tool names in the seven changed files
PASS: shipped-package-reference-contract
REQ-735 probe PASS
```

## Verification

1. RED probe at base: 7 FAIL lines, exit 1 (above).
2. GREEN probe at tip `82db9eb2`: `REQ-735 probe PASS`, exit 0.
3. `bash _dev/tests/shipped-package-reference-contract.sh`: exit 0 (1.3 s). `bash _dev/tests/contract-regressions.sh`: exit 0, "Contract regression checks passed." (27 s). One earlier attempt of contract-regressions exited 1 before running anything, because my output redirect named a scratch folder that did not exist yet; I created the folder and reran. That was not a test failure. No per-file time budget failed.
4. `git diff --stat 547d4cc5..HEAD`: exactly the seven Scope files, 10 insertions, 11 deletions.
5. Re-read each added or changed line against `communication-style.md` § Patterns to Avoid: no emphasis phrases, no analogy, no em-dash chain or semicolon, no bare identifier. The two moved bullets are byte-identical to their old lines.
6. Sweep for other restatements (`says strict`, `while redesigning`, `redesign`, `approximate unless`, `session model`, `numeric target`, `existing UI` over `skills/` and `_dev/`): no other shipped copy or test pins the moved lines. `ui-design.md:58` ("Don't redesign what isn't broken") is a different rule and stays, per the REQ. Only `CHANGELOG.md:19-20` (the 0.322.0 history entry) still names `coding-guardrails.md` as R1/R2's home; history is not rewritten.

## P-A-U

- PLAN: apply D-01 to D-05 as written: two byte-identical line moves, two precedence clauses (one in three lock-step copies), one new commit rule. No other wording.
- APPLY: one Python pass with an exact-once assertion per substitution over the seven files, then one commit by exact path.
- UNIFY: `git diff --stat 547d4cc5..HEAD` lists the seven Scope files only; checked each hunk: R1/R2 gone from both guardrails copies, present once in communication-style.md after the bug-report bullet; the three model paragraphs identical, with "cheapest model tier" still on one line; commit.md bullet sits right after the `git add -A` bullet, which is unchanged; no debug artifacts.

## Decisions

- **D-08 (DECIDE & STATE):** The F5 reason sentence is "Another session or the user may own it." It is the example the brief gave, one short sentence, and it names the condition (someone else owns the change) without naming a tool.
- **D-09 (DECIDE & STATE):** The model paragraph rewrap breaks after "unless the action", so lines 1 and 2 stay unchanged and the new line 3 is 81 characters, within the paragraph's existing width (line 1 is 86).

## Discovered Tasks

None.

## Integration seams

None with REQ-736 (it edits only `skills/do-work-toolbox/actions/help.md`). Headings edited:

- `skills/do-work/crew-members/communication-style.md` § Core Patterns
- `skills/do-work/crew-members/coding-guardrails.md` and `skills/do-work-toolbox/crew-members/coding-guardrails.md` § 1 (Think Before Coding, "Before implementing" list) and § 3 (Surgical Changes, "When editing existing code" list)
- the three `background-agents.md` copies § Match the Pattern to the Harness (last paragraph)
- `skills/do-work/actions/commit.md` § Rules

## Proposed CHANGELOG entry

Title: **Numeric-Target and Redesign Rules Apply in Every Session; Actions Can Override the Ask and Model Rules**

What changes for an installed copy: every session in an installed project now reads numbers in a request as approximate and keeps existing UI during a redesign. These two rules move from the implementation-only guardrails (`coding-guardrails.md`) into the always-on communication rules (`communication-style.md`). An action's own way of asking now wins over the general ask-tool rule, and a model an action names now wins over the general model rule in `background-agents.md`. The commit action gains a rule never to revert, stash or discard a change you did not make. The model-per-agent rule stays in `background-agents.md` only.

Bump: **minor**. This is impact-rule-change: `communication-style.md` loads in every session of every installed project, so behavior changes outside implementation, not only inside a do-work run (same reasoning as REQ-724 D-12), per `skills/do-work/actions/work-reference.md` → Changelog Entry Procedure.

## Proposed lesson

None. The lock-step and condition-preserving traps were answered up front by the probe and the pre-dispatch load-site check; nothing new failed.

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file): families `canonical-link-outlives-its-target` and `manifest-ownership-vs-edit-content` do not fire (no file moves, no manifest change).
- `_dev/primes/lessons-action-files.md`, fixed-string grep for `[family: alternate-writer-contract-drift]` and `[family: condition-preserving-prose-extraction]`: answered by the restatement sweep (Verification 6), the three-copy identity check, and the pre-dispatch evidence that Step 6 loads `communication-style.md` alongside `coding-guardrails.md`.
- `_dev/primes/prime-action-files.md` § Traps, `alternate-writer-contract-drift` line.
