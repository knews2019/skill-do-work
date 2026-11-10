# Work phase: review

> Companion to [work.md](./work.md); core queue procedures. Read when its phase routing selects this file. Existing step numbers remain stable.

### Step 7: Review

**Selected AI-Evaluation-Loop:** Follow `actions/ai-evaluation-loop-reference.md` instead of the ordinary review dispatch, score branches and one-remediation policy below. The parent owns all writes and the final result. Required reviewer failure/missing evidence stays unresolved; standalone review and unselected runs remain unchanged. Keep the applicable no-op validator, restatement sweep, timestamp and impact-routing contracts below. After the selected result, continue at Step 7.5.

**Ordinary review:** Run actions/review-work.md in **orchestrated mode** against this REQ.

For an already-green repository-gate repair, independent review reads the durable no-op evidence instead of an implementation diff and freshly invokes the same validator with the exact request path, finalization writer, and completion timestamp. It may accept the no-diff path only from typed success with `already_green_repair.review_allowed: true`; record the typed fingerprint, gate-evidence, canonical completion paths, staged paths, and decision in the ordinary `## Review`. Without that result the no-op may not complete.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?). **In worktree dispatch mode** the working tree is clean after the merge, so the review reads this REQ's merge range `<pre>..<merge_hash>` instead (`actions/review-work.md` Step 4, Get the Diff).

**Restatement sweep (MUST).** If this REQ's diff redefines something other text restates — a contract token, a schema field's semantics, a gate's wording, a prescribed command's output shape — or this REQ is its wave's last successful integration (passed below), the review runs the sweep defined in `actions/review-work.md` Step 6 (**Restatement Sweep**) and reports every stale restatement as a finding, including ones in files outside this REQ's declared Scope. Only an `impact-critical` result routes through Step 10's Fold-First Rule; every other result stays in the report and ends `→ report only`.

**How to run it:** Spawn an agent with actions/review-work.md file, the REQ path, and the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain has a matching file, load it; otherwise skip). In worktree dispatch mode, also pass this REQ's merge range `<pre>..<merge_hash>`; at a wave's last successful integration, also pass that fact and the other members' REQ ids (`actions/fan-out-reference.md` → **Restatement sweep and wave-end set-aside**). Or read actions/review-work.md file and follow its orchestrated-mode instructions in the current session.

**For ordinary review, what happens next depends on the review result:**

- **Acceptance = Pass AND overall ≥ 75%**: Append the Review section to the REQ and continue to archive as `completed`. Minor findings go in the report only.
- **Acceptance = Partial OR overall 50-74%**: Append Review and continue to archive as `completed`. Every finding is impact-stamped per `actions/review-work.md` Step 10. Only `impact-critical` findings create automatic follow-up work; every noncritical line ends `→ report only` and stays in the Review section.
- **Acceptance = Fail OR overall < 50%**: **Do NOT archive as completed.** Instead:
  1. Append the Review section to the REQ.
  2. Return to Step 6 (Implementation) with the review findings as context for the builder. Load `crew-members/debugging.md` for the remediation attempt — the builder needs structured debugging methodology, not just "try again."
  3. The builder gets **ONE remediation attempt**.
  4. Re-run Steps 6.25 through 7 (Summary → Qualification → Testing → Review) on the remediated code.
  5. If still failing after remediation: update frontmatter to `status: completed-with-issues`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), append a `## Remediation` section documenting both attempts, and route the remaining findings through `actions/review-work.md` Step 10. Only `impact-critical` findings auto-queue; the rest stay in the Review section. Then proceed to archive (Step 8) — the frontmatter is already set, so Step 8 should not overwrite it.

After the first review result is recorded, stamp `review_at: <now>` using the Timestamp rule, regardless of its verdict. If remediation runs, stamp `remediation_at: <now>` only after that builder hand-back is successfully integrated, then stamp `re_review_at: <now>` only after the post-remediation review result is recorded. All three are written only if the field is absent (`frontmatter set` enforces this; **Stamps are append-only**, `actions/work-reference.md`). A passing first review leaves both remediation fields absent.

The status `completed-with-issues` means the REQ was archived but has known unresolved problems. It counts toward UR completion for archiving purposes. Any critical follow-up remains queued; noncritical findings remain visible in the archived Review until a maintainer explicitly captures one. This status remains visible to recap and every completed-work presentation action; those readers inherit the Terminal-success status set from `actions/work-reference.md` rather than defining a caller-specific filter.

**Automatic review follow-ups are based on impact, not score or severity.** Apply `actions/review-work.md` Step 10 to every finding. Detect an existing `addendum_to` cycle per Step 8 before creating a critical follow-up.

**Calibrate depth to route:** Route A gets a quick scan (skip dimensions that don't apply). Route B gets a standard review. Route C gets a thorough review comparing against the plan.

Append to the request file:

(append per the **Append to REQ File** template in `actions/review-work.md` — the file dispatched above, so it is already in context; review-work.md owns the Review section format)

### Step 7.5: Lessons-Capture Phase

> **Named entry point.** Other actions reference this as **work.md's Lessons-Capture Phase** (not by step number) — e.g. `actions/kb-lessons-handoff.md` and `actions/review-work.md`. The `7.5` is for internal navigation only; callers must use the phase name so they don't break if steps are renumbered.

Before archiving, capture what's worth remembering. This section is the institutional memory — when someone revisits this code in six months, the REQ file tells them what happened, what was tried, and why things ended up the way they did.

Append to the request file:

```markdown
## Lessons Learned

**What worked:** [1-2 bullets — approaches, patterns, or tools that paid off]
**What didn't:** [1-2 bullets — dead ends, failed approaches, and *why* they failed]
**Worth knowing:** [Anything the next person touching this code should know — gotchas, edge cases, non-obvious dependencies]
```

**Rules:**
- Keep it concise — pointers to code, not walls of text. The code is the source of truth.
- **Required for Routes B and C** — there's always something worth recording when exploration or planning was involved. **Optional for Route A** — skip if the change was straightforward with no unexpected discoveries, no failed approaches, and no gotchas worth noting. If anything surprised you (undocumented behavior, unexpected test failures, a file that wasn't where you expected), record it.
- "What didn't work" is the most valuable part — it prevents repeating mistakes.
- File lists are no longer needed here — they're covered by the mandatory Implementation Summary (Step 6.25).

**Write the `## Orientation` block (the hand-back's "what's being built"):** After the Lessons Learned section, append a short `## Orientation` section reporting the change at feature/subsystem altitude — "Now you can X; lives in Y subsystem" — not a file list. Use the REQ's `prime_files` to name the subsystem; flag `[MAP CHANGED]` only when the change alters the system's shape (new module, data flow, contract, or a renamed concept). Run a narrowed staleness spot-check on each touched prime per `../../do-work-toolbox/actions/prime.md` Step 2 / Step 6 (do its referenced paths still exist?) and flag any prime the change made stale. **Scale to reach:** a leaf REQ is one line; a map-changing REQ gets a short paragraph and a why-it-matters. When `prime_files` is empty, derive a one-line feature-altitude summary from the What / Implementation Summary instead — never a file list. This block feeds the **WHAT'S BEING BUILT** section of the Decision Brief (`actions/work-reference.md` → **Decision Brief (hand-back format)**). Crash recovery strips `## Orientation` on re-queue (it's orchestrator-generated).

**Update prime lesson satellites (deferred to Step 8):** After writing the Lessons Learned section, follow **Prime Discovery** in `crew-members/general.md` → PRIME Files Philosophy. For each discovered prime file relevant to this lesson, **collect a pending lesson write** — do NOT execute the write here. The REQ is still in `do-work/working/`, so any link pointing to its eventual archive location would either be broken or tempt a link to the transient working path.

The write target is the prime's **lesson satellite**, `lessons-<name>.md` beside `prime-<name>.md` — never the prime itself. A prime is read in full every time its area is touched; appending to it on every REQ is what turns a routing index into an archive (`crew-members/general.md` → PRIME Files Philosophy).

Choose one short kebab-case **failure-family slug** for each pending lesson. Record each pending write as a tuple: `{ primeFilePath, satellitePath, relativeLinkText, lessonSummary, familySlug }`. Hold them in memory (or a small scratch file under `do-work/working/`) until Step 8. The literal marker written to the satellite is `[family: <family-slug>]`; use the same slug for recurrence checks and the lessons index.

Compute each deferred prime-link path relative to the prime file's location (not the repo root) per `actions/work-reference.md` → **Deferred Prime-Link Path Computation (Step 7.5)**; Step 8 validates the resolved target against the exact archive path in the lifecycle plan before adding the prime or satellite to the finalization allowlist.

Only add a link when the lesson is relevant to that prime file's scope — don't spray every lesson into every prime file. If no relevant existing prime is discovered, skip this and clear the pending list.

**Knowledge-base handoff.** After the Lessons Learned section is written and prime-file links are in place, follow `actions/kb-lessons-handoff.md` to offer dropping a structured source document into the knowledge base's `raw/inbox/` (under `<project-root>/docs/kb/` by default) so the next `do-work-knowledge bkb triage` + `do-work-knowledge bkb ingest` cycle compiles the lessons into the wiki. The handoff asks the user before writing and records `kb_status` (plus `kb_entry` on success) back onto the REQ. In unattended work runs with no human in the loop, the handoff defaults to `kb_status: pending` — it never writes to the KB without consent. If the project has no knowledge base, the handoff points the user at `do-work-knowledge bkb init` and defers; it never blocks archival.
