# Upstream suggestion for `knews2019/skill-do-work` — an opt-in `validate-feedback --capture [--run]` that chains triage → questions → capture → verify

**How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
consumer repos it was authored in. Observed against **v0.305.84**; all line numbers are from that tag.

---

## Request

`do-work-toolbox validate-feedback` triages well and then stops, by design: it is read-only and ends with a
suggested `do-work capture-request:` line. In practice users run the same four steps by hand after almost
every triage: answer the Discuss items in chat, retype the accepted finding ids into a capture, type
`do-work verify-requests` on the new UR, then type `do-work run`. Please make that chain one opt-in
invocation. Nothing in v0.305.60 to v0.305.87 touches validate-feedback, and no queued REQ covers it.

Keep Capture ≠ Execute exactly as it is today: without the flag or an explicit phrase, the action stays
read-only and its output is unchanged. No new REQ fields, no new statuses.

- **C1. The flag.** `validate-feedback --capture [--run]`, plus the phrases "then capture the accepted ones"
  and "capture and run" in the same invocation. `--run` is only valid together with `--capture`. Also
  accept a file path as the input, so a review or audit file keeps its source.
- **C2. Ask about Discuss items.** After the per-item verdicts, ask the user about each Discuss item with
  the interactive question tool, recommended option first, value and risk on each option.
- **C3. Capture once.** Capture the Accept items plus the Discuss items the user accepted as one UR (one
  user request record) with one REQ per finding. Each REQ cites its finding id and the source (pasted text,
  or the path of a review or audit file) in the provenance block the handoff already asks for
  (`validate-feedback.md:114`).
- **C4. Verify, then optionally run.** Run `verify-requests` on that UR automatically and print one combined
  report. With `--run`, and only when verification finds no gaps, continue into `do-work run` on just that
  UR.
- **C5. Routing.** Give `validate-feedback` a forward in the core `do-work` routing table, and have the suite
  installer detect the older standalone `do-validate-feedback` skill that still answers the same phrase in
  consumer repos, and offer to retire it.

## What happened, in consumer repos

Counts below come from session transcripts over the 30 days to 2026-10-09, re-counted in a second pass,
across four consumer repos plus the suite's own development checkout.

| Measure | Count |
| --- | --- |
| Prompts invoking `do-work validate-feedback` | 21, in about 18 sessions, 4 consumer repos + the suite checkout, 09-13 to 10-09 |
| Sessions that ran validate → capture (→ verify / run) by hand in one session | about 10, 2 to 5 prompts each |
| `do-work verify-requests` typed right after a capture | roughly 12 to 15 |
| Accept lists retyped by finding id | at least 1 verbatim (below), plus every "capture the accepted" follow-up |
| Pasted reviews that landed in the wrong repo | 1 |

Dated, anonymised snippets:

- 2026-09-17, suite checkout: after a triage, the user typed "capture F1, F2, F6, F9 and C1". The action
  had just printed those ids; the user copied them back by hand.
- 2026-09-24, a second consumer repo: "double check if all the features are added, use do-work
  verify-request as well after the request capture".
- 2026-09-28, four consumer repos: the same review was relayed by paste into `/do-work validate-feedback:`
  in each repo. In one of them the pasted text was the review written for a different repo, and triage ran
  against code the findings did not describe.
- 2026-09-28, a third consumer repo: validate-feedback, then "capture the requests", then verify-requests,
  three prompts for one chain.
- 2026-10-02, suite checkout: "capture the ones that are accepted".
- 2026-10-03, one consumer repo, two sessions: `/do-work capture-request` with pasted findings, then the
  next prompt was `do-work verify-requests UR-NNN` on the UR just captured.
- 2026-10-06, suite checkout: "do-work validate-feedback (check the following request, talk with me, use the
  ask tool and `capture-requests` after we got a common understanding)". This prompt describes C1 to C3.
- 2026-10-06, a second consumer repo: "review the changed code, are there bugs? run /do-work
  validate-feedback on the output of the code review".
- 2026-10-07, suite checkout: "capture to fix accepted".
- 2026-10-08, suite checkout, three sessions: "capture then, then run them", "capture them and run them",
  "capture and run".

Separately, 18 identical copies of a 62-line standalone `do-validate-feedback` skill sit in
`.claude/skills/` (and one `.agents/skills/`) across 9 consumer repos; several repos have a second, mirrored
checkout. It predates the toolbox action and has none of its later parts: no prompt-injection guard, no
Surface-cost check, no provenance-preserving capture handoff. The harness advertises it under its name,
which matches the phrase "validate feedback", while core `do-work` has no route for that phrase.

## Where the behaviour lives today

- `skills/do-work-toolbox/actions/validate-feedback.md:5`: "this action does NOT modify any files and does
  NOT create REQs … Accepted items become work through a separate, user-gated `do-work capture-request:`
  step (Capture ≠ Execute)."
- `validate-feedback.md:30`: input is "the pasted feedback"; "If `$ARGUMENTS` is empty, ask the user to paste
  the feedback". A file path is not an input form, so the source of a review or audit file is lost.
- `validate-feedback.md:78`: **Discuss** "has merit but the right path isn't clear-cut … Frame the trade-off."
  The trade-off is written into the report; the user is never asked to decide it.
- `validate-feedback.md:113-117`: the handoff block already says to keep "the accepted finding's **verbatim
  claim**, **original severity/source**, **Evidence**, and **Surface-cost** result together in the capture
  payload", then prints `do-work capture-request:` and `do-work run` as lines for the user to type.
- `validate-feedback.md:122`, `:136`, `:154`: the Rules line ("Read-only … The capture handoff is a
  *suggestion* the user runs deliberately"), the rationalization row ("I'll capture the accepts to save the
  user a step → Stop after the report"), and the checklist item ("no REQs were created"). All three need an
  "unless `--capture` was given" clause.
- `skills/do-work/SKILL.md:20`: "Stop after capture unless the same user invocation explicitly requested
  execution too." The flag is that explicit request; C4 relies on this existing rule, not a new one.
- `skills/do-work/actions/capture.md:119`: "A review, build, triage, or consumer-report finding reaches
  capture only when the user invokes `do-work capture` and quotes the complete report-only finding line as
  the source". `--capture` is a user invocation; the fold-first scan named on the same line still runs, so
  findings that duplicate a queued REQ fold instead of minting a new one.
- `skills/do-work/next-steps.md:27`: after capture, suggest "`do-work verify-requests` before `do-work run`".
  `capture.md:261` makes the same suggestion for detailed requests. Both only suggest; the 12 to 15
  hand-typed follow-ups are the cost.
- `skills/do-work/actions/verify-requests.md:32` and `:38`: capture QA accepts one `UR-NNN`, so C4 needs no
  change in verify-requests.
- `skills/do-work/actions/work.md:105`: a `UR-NNN` token scopes the run to that UR's REQs and keeps
  `depends_on` gating, which an explicit `REQ-NNN` list would bypass.
- `skills/do-work/SKILL.md:61`: "Load `crew-members/clear-questions.md` before asking an interactive
  question." `crew-members/clear-questions.md:19-21` (say the consequence of each option) and `:27-29`
  (concrete options, never open-ended) already govern C2's wording.
- `skills/do-work-toolbox/actions/maintainability-audit.md:132`: the audit's loop footer is "`do-work-toolbox
  validate-feedback` → capture handoff → `do-work run` → re-audit", four hand-typed hops for every audit.
- `skills/do-work/SKILL.md:24-46`: the core routing table has no `validate-feedback` row, so
  "do-work validate-feedback: …" falls through to the capture fallback at `:46` ("unmatched descriptive
  multi-word input"). The forward exists only in `skills/do-work-toolbox/SKILL.md:18`.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go:743-791`
  (`reviewAndConfirm`) narrates module diffs and calls `warnAboutDirtyManagedPaths` (`:786`, defined at
  `:809`) before the `[y/N]` prompt. Nothing looks at sibling skill folders the suite does not own.
- `skills/do-work/actions/version.md:60`: the update "does not mutate … any other project configuration".
  So retiring an old skill must be an offer, not part of the install transaction.

## Proposed direction

**C1. `validate-feedback.md` Input.** Add `--capture` and `--run`, and the two phrases. Accept a file path as
`$ARGUMENTS` (read it, record the path as the source). Reject `--run` without `--capture` with a one-line
usage. Edit `:5`, `:122`, `:136`, `:154` to say "read-only unless `--capture` was given". Without the flag,
Steps 1 to 5 and the Output Format stay byte-for-byte as they are.

**Wrong-repo check before any write.** When `--capture` is set and more than half of the cited `file:line`
paths do not exist in this repo, stop before C2 and ask: "These findings cite paths that are not in this
repo. Continue capture here, or stop?" Recommended option: stop. This is the one incident on 2026-09-28; it
costs one existence check per cited path and only runs when a write is about to happen.

**C2. New Step 6, Discuss questions.** Only with `--capture`. One question per Discuss item, using the
interactive question tool after loading `crew-members/clear-questions.md`. Options, recommended first:
"Accept: capture as a REQ with remedy <one line>", "Park: `do-work-toolbox note` it for later", "Drop: no
work". Each option carries one line of value and one line of risk. Already done and Push back items are
never asked about. If no question tool is available, list the Discuss items with the same options and wait
for the answer in chat.

**C3. New Step 7, capture.** Build one payload from the Accept items plus the Discuss items the user
accepted, each as the existing `:113-114` provenance block, and run `skills/do-work/actions/capture.md` on
it. Result: one UR, one REQ per finding, each REQ's source line naming the finding id (F3, "Finding 3") and
the input kind (pasted text, or the file path from C1). If the set is empty, say so and stop without
capturing.

**C4. New Step 8, verify and optional run.** Run `skills/do-work/actions/verify-requests.md` with the new
`UR-NNN`. Print one combined report: the triage summary table, the UR and its REQs with titles, the verify
verdict. With
`--run` and a verify result with no gaps, continue to `do-work run UR-NNN` (UR token, so `depends_on` still
gates). If verify found gaps, stop and print the gaps with the exact command to run after fixing them.

**C5a. Core routing.** Add a row above the fallback at `skills/do-work/SKILL.md:46`:
`| validate-feedback, triage feedback | ../do-work-toolbox/actions/validate-feedback.md |`.

**C5b. Installer.** Add a read-only check beside `warnAboutDirtyManagedPaths`: if `.claude/skills/` or
`.agents/skills/` contains a folder named `do-validate-feedback` (keep a short list of known predecessor
names, starting with this one), narrate "`<path>` is an older standalone copy that answers the same phrase
as `do-work-toolbox validate-feedback`; remove it with `git rm -r <path>`". Never delete it inside the
transaction, per `version.md:60`.

## Acceptance check

- `validate-feedback` with no flag on a fixed sample paste produces the same report as v0.305.84 and writes
  no files.
- `validate-feedback --capture` on a paste with 2 Accept, 1 Discuss, 1 Push back items: one question is
  asked (about the Discuss item only); answering Accept produces exactly one new UR with 3 REQs, each REQ
  naming its finding id and source; `verify-requests` runs on that UR without a second prompt; one combined
  report is printed; no run starts.
- The same with `--run` and a clean verify ends in `do-work run UR-NNN`. With a verify gap it stops before
  the run and prints the gaps.
- A finding that duplicates a queued REQ folds into it through the existing fold-first scan.
- A paste whose cited paths are mostly absent from the repo stops before any question or write.
- `--run` without `--capture` is rejected with usage text.
- "do-work validate-feedback: …" routes to the toolbox action, not to capture.
- An install or update into a repo that has `.claude/skills/do-validate-feedback/` prints the retire line and
  leaves the folder in place.

## Out of scope

- A `--review N` front end that reviews the last N commits before triage. Feed the output of `/code-review`
  or `do-work-toolbox code-review` in as the paste instead.
- The "what earned this, and is the fix still cheaper than the surface it added" lens: Step 4 item 5
  (`validate-feedback.md:69`) and the Surface-cost field already apply it.
- Any change to verdict rules, to capture's REQ templates, or to how `do-work run` selects work.
- Capturing without the flag or phrase, for any verdict, including impact-critical ones.
