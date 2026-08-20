# Capture Requests — Reference

> **Companion file to `actions/capture.md`.** Holds the REQ/UR templates, the schema-alias table, and the addendum-REQ template — the hard-contract field specifications that `actions/work.md`, `actions/roadmap.md`, and `../../do-work-board/tools/queue-kanban/model.go` all read. Each section below is pointed to by name from the matching step in `actions/capture.md`. Loading this file is only necessary when you reach the step that references it — and read only the named section. If this file is already in context from an earlier step this session, reuse it; don't re-read it at every reference site.

---

## REQ Title Convention

**Canonical home for the `title:` shape.** Every action that mints a REQ title follows this section — `actions/capture.md`, `actions/review-work.md` Step 10, `actions/work-reference.md`'s Builder-Decided Follow-up template and its **Discovered Tasks Classification (Step 8)** flow, and `../../do-work-toolbox/actions/code-review.md`. The condition is the rule — **any flow that mints a REQ carrying an `impact:` value follows this section** — so a new one inherits it without this list being re-counted.

```
title: "[<impact token>] <Kind prefix>: <brief description>"
```

- **The impact tag is a bracketed classification tag, not a kind prefix.** It composes with the existing `<Kind prefix>: ` conventions (`Addendum: `, `Review fix: `, `Confirm: `, `Code review: `) instead of competing for the same slot — a review-fix REQ that is also negligible reads `"[impact-negligible] Review fix: guard misses hex shorthand"`, never a double-colon title. Both parts are optional and independent: a plain capture with no kind prefix is `"[impact-negligible] Retitle the export button"`.
- **Emit the tag only when `impact:` is something other than the `impact-user-visible` default.** That mirrors the board's impact chip, which renders for every value except that default (`actions/work-reference.md` → Request File Schema), so the title and the card agree and the common case stays unadorned. An absent `impact:` reads as the default, so it gets no tag either.
- **The tag goes in `title:` and nowhere else.** Not in the filename — `REQ-NNN-slug.md` is unchanged, because a filename-borne token would mean renaming a file mid-pipeline whenever the verdict is revised, breaking every path already pointing at it. Not in the body `# H1` either: the H1 is Title-Cased prose, and Title-Casing a token would break the exact string the tag exists to be searched for.
- **Write the token in full, exactly as the `impact:` enum spells it** — `impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible`. Never an abbreviated form: the full prefixed token is what makes one `grep` return one axis, and it is what the board's search box matches (`../../do-work-board/tools/queue-kanban/web/board-filters.js` matches a case-insensitive substring on the title).
- **Double-quote the whole value, always.** A title carrying a colon or a leading `[` is not a plain YAML scalar, so an unquoted one makes strict YAML reject the whole block and the REQ is read only by the board parser's last-resort line-based recovery (`../../do-work-board/tools/queue-kanban/frontmatter.go`). That recovery is good — it exists so one bad line cannot cost a REQ its status, UR pointer, and dependencies, and it does hold — but it is a salvage path with a narrower contract than the parser proper, and it can corrupt the very titles this convention creates: recovery splits a value that opens with `[` and closes with `]` as a YAML flow list, so `title: [impact-negligible] Retitle export, again [v2]` is read back as `[impact-negligible] Retitle export again [v2]` — the comma silently eaten, no warning raised. Quoting costs two characters and takes the whole class off the table.
- **`impact:` is the source of truth; the title is a mirror.** When the two disagree the field wins, and every reader that acts on impact — Step 1's `--skip-impact-negligible` filter included — reads the field, never the title. The title tag exists so a human searching the board finds the REQ.

## Fold-First Rule

**Canonical home for the fold-before-mint contract.** The condition is the rule — **any flow about to mint a REQ from a finding, discovery, or triage result first scans the queue for an existing home, and creates a new REQ file only as the justified exception** — so `actions/capture.md` Step 2, `actions/review-work.md` Step 10, `actions/work-reference.md`'s Builder-Decided Follow-up and **Discovered Tasks Classification (Step 8)** flows, and `../../do-work-toolbox/actions/code-review.md` all follow this section, and a new flow inherits it without this list being re-counted. Half of all REQs ever queued were spawned by another REQ's review or build; the fix that keeps the queue draining is one file per root cause, not one file per facet.

- **The candidate set is every `pending` or `pending-answers` REQ in `do-work/queue/` — the whole queue, never filtered to the discovering UR.** A root cause is a property of the codebase, not of the UR that happened to surface it, and a same-UR filter lets sibling URs mint duplicate REQs for one root cause. Queue residency is the unclaimed check: a claimed REQ lives in `do-work/working/` and is never folded into.
- **Match on root cause, sweeps first.** `grep -rl "^sweep: true" do-work/queue/`, then compare the finding's root cause against each candidate's `sweep_key:` — append on a key match. When no key matches literally, compare against each candidate's `## What` root-cause statement: same rule (would one fix close both?) means same REQ. Then check non-sweep pending REQs' `## What` the same way. Title resemblance is forbidden as a match signal in both passes — two flows judging titles differently mint duplicates, recreating the runaway at half scale.
- **Fold into a sweep = one checklist line under `## Instances`, citing the discoverer:** `- [ ] [file/site]: [instance] (found by REQ-NNN / UR-NNN)`. The append never touches frontmatter, and the sweep's `user_request:` stays its original UR — the instance line's citation is the cross-UR attribution.
- **Escalation is the one exception — a fold never buries a more urgent verdict.** When the folded instance's judged impact outranks the sweep's (`impact-critical` outranks everything; `impact-user-visible` and `impact-rule-change` outrank `impact-negligible`), the fold also promotes the sweep's `impact:` to the instance's verdict and re-mirrors the title tag (REQ Title Convention above) — otherwise `--skip-impact-negligible` would skip a finding someone judged worth seeing. An `impact-critical` instance additionally flips a `pending-answers` sweep to `pending`: the critical pierce (`actions/review-work.md` Step 10) applies to folds exactly as it does to creation, because burying a security finding behind a consent checkbox is the wrong trade however it arrives. Escalation only, never downward, and always reported prominently.
- **Fold into a non-sweep pending REQ = convert it once.** When a finding shares the root cause a plain pending REQ already covers, add `sweep: true`, a `sweep_key:`, and an `## Instances` checklist seeded with that REQ's own original instance plus the new one. This one-time frontmatter edit is legal only while the REQ is `pending` or `pending-answers` and unclaimed — the same window in which capture already appends addenda to queued REQs — and after conversion the append-only rule above governs.
- **A prose-only discrepancy never mints a new REQ at all — it always has a home.** Prose-only means the fix changes no behavior, no checker's predicate, and no rule's stated condition; the example shapes (a stale count, a wrong cross-reference number, a comment describing a superseded mechanism) are illustrative, never the test (`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale). When no root-cause match exists, such a finding lands on the **standing prose-reconciliation sweep** — the queue's permanent `sweep: true` + `standing: true` REQ keyed `sweep_key: prose-only-discrepancy-reconciliation`, identified by that key and never by id, which drains on its own cadence and never closes — or is folded into the next commit that touches the file for another reason. **When the queue holds no REQ with that key, the minting flow creates it on demand** — the one sweep capture itself may mint from scratch — with `sweep: true`, `standing: true`, the canonical key, `impact: impact-negligible`, an empty `write_set`, and the discovering UR as its `user_request:`; a fresh install has no standing sweep until its first prose-only finding needs one, and that must never strand the finding. A capture that declines to mint must name where the finding went; declining without a destination loses findings, which is worse than over-capture. Three exemptions survive explicitly: an `impact-critical` finding is never deferred, at any depth, whatever its fix touches; a **contradiction** between two shipped instructions is not prose reconciliation — two rules that cannot both be followed change what an agent does, so they change behavior and earn a first-class REQ; and a prose discrepancy in a **user-facing** artifact contract (a changelog entry, a report template a consumer reads) is judged on its reader, not on its diff shape. This boundary governs what the minting flows *create*, never what already exists — it is no licence to cancel queued REQs retroactively; only `do-work abandon` removes user intent from the queue.
- **A new REQ states its exception in one line:** the fold-first scan found no pending REQ, sweep or not, in any UR, sharing the root cause — and, for a prose-only finding, why the standing sweep was not its home. The impact judgment is unchanged either way — a folded instance inherits its sweep's `impact:` (escalation above is the one exception); a created REQ carries its own judged verdict per the REQ Title Convention above.

## Request File Formats

### Simple REQ

```markdown
---
id: REQ-001
title: "Brief descriptive title"   # always double-quoted; prefix `[<impact token>] ` when impact: is anything other than impact-user-visible — REQ Title Convention above
status: pending
created_at: 2025-01-26T10:00:00Z  # current UTC instant, never local time with a Z suffix (Timestamp rule, actions/work-reference.md)
user_request: UR-001
domain: frontend  # choose one: frontend, backend, ui-design, general, security, testing, or cms
prime_files: []  # list paths to relevant prime-*.md files, or leave empty
tdd: true  # default true when a runnable RED test can be written in this project's harness; false otherwise (see heuristic below)
suggested_spec:  # optional — spec template name if one clearly matches (e.g., "api-endpoint", "bug-fix")
depends_on: []  # optional — list of REQ IDs that must complete before this REQ runs; honored by actions/work.md's selection scan
maintenance: false  # set true ONLY for a deliberate removal/narrowing of the skill's OWN operating instructions — see Step 1's maintenance assessment; loads crew-members/maintenance.md in actions/work.md Step 6
# impact: impact-user-visible   # EXPECTED on every REQ, but JUDGED, never copied: impact-critical | impact-user-visible | impact-rule-change | impact-negligible — whether anyone would ever notice the work, judged by Step 1's impact assessment. Deliberately commented out: an uncommented value gets copied more often than judged, and absence already reads as impact-user-visible, never as the user's stop signal (actions/work-reference.md → Request File Schema)
# effort_estimate: effort-mechanical   # OPTIONAL triage bit (effort-mechanical | effort-substantive; absent reads as effort-substantive) — the SIZE of the fix, a separate axis from impact. Capture MAY set it when the request is clearly a small mechanical fix, so the board chips it at a glance; never invent `effort-mechanical` for work whose size you haven't judged.
# assigned_to: "cloud-alpha"   # OPTIONAL advisory earmark — emit ONLY when the user names a session to reserve this work for (Step 1's earmark assessment); verbatim and YAML-quoted, never invented. The default work scan skips-and-reports it; explicit targeting clears it on claim (actions/work-reference.md → Request File Schema)
# External-condition fields — set ONLY when the task waits on something outside the queue (see Step 1's external-condition assessment). Omit all three for normal REQs.
# status: blocked            # use INSTEAD OF `status: pending` when the REQ cannot start until an external condition is met — distinct from depends_on (another REQ) and Open Questions (a question for the user)
# blocked_by: "..."          # free text naming the condition (always YAML-quoted), e.g. "LM Studio running locally", "designer answered on mockups"
# blocked_at: 2026-01-26T10:00:00Z   # stamp the moment it was captured blocked — current UTC instant, same Timestamp rule as created_at
# blocked_check: "..."       # OPTIONAL shell probe (YAML-quoted) — emit ONLY when the user supplies or explicitly confirms the command; never invent one
---

# [Brief Title]

## What
[1-3 sentences describing what is being requested]

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)
[User's stated reasoning — omit if not provided]

## Context
[Additional context, constraints, or details mentioned]

## Red-Green Proof
**RED prompt/case:** [Minimal prompt, repro, or example that should fail or be missing today]
**Why RED now:** [What is currently broken or absent]
**GREEN when:** [Observable result that proves the request is done]
**Validation:** [User confirmed / User adjusted / Inferred during capture]

## Assets
[Description of screenshots or links to saved files]

---
*Source: [original verbatim request]*
```

Include `## Red-Green Proof` when the request is behavior-changing and can be proven with a prompt, repro, or example. If `tdd: true`, this section is mandatory. The goal is proof of behavior, not implementation detail.

Treat defining the RED state as essential, high-value capture work. It is one of the most helpful things you can do for the downstream builder because it turns vague intent into a concrete failing proof target. Do not treat this as paperwork. Lean into it. Be eager to find the best RED case: the smallest, clearest prompt/repro/example that proves the behavior is missing now and will clearly turn GREEN later.

### Complex REQ (additional sections)

Complex requests use the same base format plus these sections:

```markdown
## Detailed Requirements
[Extract EVERY requirement from the original input that applies to THIS feature.
DO NOT SUMMARIZE — use the user's words. Include specific values, constraints,
conditions, edge cases, "must"/"should"/"never" statements.]

## Constraints
[Limitations, restrictions, batch-level concerns that apply to this REQ]

## Dependencies
[What this feature needs or what needs it — reference other REQ IDs]

## Builder Guidance
[Certainty level: Exploratory / Firm / Mixed. Scope cues like "keep it simple."
Any latitude given to the builder.]

## Open Questions
- [ ] [Question about ambiguity the user needs to clarify]
  Recommended: [best default based on context]
  Also: [alternative A], [alternative B]

Open Questions use checkbox syntax with recommended choices. Each question includes a `Recommended:` line (the best default if the user doesn't answer) and an `Also:` line with alternatives. The choices make questions answerable even when the question itself isn't fully understood — the user can just pick one.

`- [ ]` = unresolved, `- [x]` = answered (answer follows `→`), `- [~]` = deferred to builder (note follows `→`).

**Capture time is the optimal window for resolving these.** During capture (this action), use the ask tool if your environment provides one; otherwise use your environment's normal ask-user prompt/tool. Present Open Questions immediately. The user is here, engaged, and fleshing out the request — don't defer what you can clarify now. Only leave questions as `- [ ]` if you genuinely can't ask (e.g., batch processing, async capture).

Only add questions where the user's intent is genuinely unclear — don't add questions the builder can answer by reading the codebase.

## Full Context
See `do-work/user-requests/UR-NNN/input.md` for complete verbatim input.
```

**Additional frontmatter for complex requests:**
- `related: [REQ-006, REQ-007]` — other REQs in this batch
- `batch: auth-system` — batch name grouping related requests
- `addendum_to: REQ-005` — if this amends an in-flight/completed request
- `write_set: [src/auth/session.ts, src/auth/*.test.ts]` — repo-relative paths/globs this REQ expects to write
- `sweep: true` + `sweep_key: <root-cause-slug>` — a consolidation sweep (one REQ per root cause with an `## Instances` checklist; the key is the append discriminator). Set at sweep creation (`actions/review-work.md` Step 10 → Sweep consolidation) or by the **Fold-First Rule**'s one-time conversion of a pending REQ — never minted from scratch by capture, **with one exception: the standing prose-reconciliation sweep, which any minting flow creates on demand when the queue lacks its key** (Fold-First Rule). `standing: true` additionally marks that never-closing standing sweep — drains tick instances and leave it `pending`; it is never archived on empty and never holds its UR open

(`maintenance` is **not** complex-only — it lives in the base schema above. Step 1's *Maintenance assessment* is its authoritative home.)

**Populating `depends_on`.** When the request body mentions prior REQs that must complete first (e.g., "after REQ-486 lands", "depends on the auth refactor"), populate `depends_on` in the frontmatter with the REQ IDs. Don't rely on numeric ID ordering — actions/work.md honors `depends_on`, not ID-based heuristics. The optional prose `## Dependencies` section in REQ bodies remains for human readers; the frontmatter field is the source of truth for tooling (work-action selection, roadmap classification, upstream-failure detection).

**Populating `write_set`.** Seed it only when the request already names the files, or when the slice is inherently per-file ("rewrite each adapter in `src/adapters/`"). Otherwise **omit it** — the field is display-only (it feeds the board's *overlaps* badge; nothing schedules on it at any builder count — `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), so an invented set is strictly worse than absence: absence simply gets no overlaps badge (it reads as unknown, not conflict), while a wrong set makes the board badge misleadingly. Capture's value is a hint, never a commitment: the work pipeline's Scope step (`actions/work-reference.md` → Scope Declaration Template) firms it up and overwrites it — for Routes **B and C** only. A **Route A** REQ never runs that step (`actions/work.md` Step 5.5), so it keeps the capture-seeded value for the whole run. The field is optional on every REQ.

**Slicing convention.** When a single user request slices into multiple REQs with internal dependencies, the slicer should populate `depends_on` per the dependency graph it produced. actions/work.md then runs roots first, gates downstream REQs on their prerequisites, and supports `--wave N` for checkpointed execution one dependency depth at a time. A clean DAG in `depends_on` makes foundation-phase batches predictable; sloppy or missing `depends_on` returns to numeric-ID order and risks cascade misclassification. Prefer slice boundaries that also give each REQ its own files — split a surface per concern rather than leaving several REQs editing one shared block — because per-concern files keep each REQ's diff and review self-contained and its board *overlaps* badge quiet. That is a nudge, not a gate: when a file-clean boundary would distort the request, slice for intent and declare the unavoidable overlap in each REQ's `write_set` so the board's overlap badge can surface it.

`depends_on` is semantically distinct from `addendum_to`: `addendum_to: REQ-N` says "this REQ amends REQ-N" (used for follow-ups and review-generated remediation); `depends_on: [REQ-N, REQ-M]` says "this REQ requires REQ-N and REQ-M to be completed first." A REQ can carry both.

### Schema Aliases

Several fields above accept legacy aliases at read time so muscle-memory typos from sister tools don't silently drop information. The canonical key wins when multiple are present; capture always emits the canonical — aliases are read-only, never propagated on write.

| Canonical field | Aliases recognized | Read sites |
|---|---|---|
| `addendum_to` | `amends`, `parent`, `amendment_to` | capture's duplicate check, `actions/work.md`'s upstream-failure walk + cycle detection, `actions/roadmap.md` Blocked classification |
| `depends_on` | `dependencies` | capture's slicing convention, `actions/work.md`'s dependency-aware selection / cycle detection / `--wave` depth / upstream walk, `actions/roadmap.md` Ready/Blocked rubrics |
| `batch` | `batch_name` | `actions/roadmap.md` batch grouping; verify-requests cross-REQ summarization |
| `related` | `related_reqs` | `actions/roadmap.md` cross-REQ surfacing; verify-requests batch coverage |
| `suggested_spec` | `spec_hint`, `suggested-spec` | `actions/work.md`'s spec pre-load hint |

For enum-valued and boolean fields shared with `actions/work.md` (`status`, `domain`, `route`, `caveman`, `tdd`, `maintenance`, `error_type`, `kb_status`, `impact`, `effort_estimate`), capture honors the **normalize-and-warn contract** defined in the Schema Read Contract (in the companion `actions/work-reference.md`): invalid values trigger a warning and a documented default rather than silent acceptance. When writing the REQ files, if the captured value for any normalize-and-warn field doesn't match the canonical enum (after applying the contract's normalization), prompt the user to confirm the intended value before emitting the REQ — capture is the human-attention window for catching typos at the source. Never write a non-canonical value silently.

### UR input.md

Created for every invocation. For simple requests, it's minimal:

```markdown
---
id: UR-005
title: Add keyboard shortcuts
created_at: 2025-01-26T10:00:00Z  # current UTC instant, never local time with a Z suffix (Timestamp rule, actions/work-reference.md)
requests: [REQ-020]
word_count: 4
---

# Add keyboard shortcuts

## Full Verbatim Input

add keyboard shortcuts

---
*Captured: 2025-01-26T10:00:00Z*
```

For complex requests, add a Summary, an Extracted Requests table, and a Batch Constraints section before the Full Verbatim Input. The verbatim section must contain the COMPLETE, UNEDITED input — never summarize or clean it up.

## Addendum REQ Template

Written for **Addendum for in-flight/completed requests** (`actions/capture.md` Step 2) — a new UR + REQ pair where the REQ's `addendum_to` links back to the original. The fenced block shows the exact frontmatter and body shape, including the `## Prior Implementation` section used when the original is archived:

```markdown
---
id: REQ-021
title: "Addendum: dark mode sidebar support"   # a non-default impact prefixes the tag ahead of the kind prefix — "[impact-negligible] Addendum: …" (REQ Title Convention above)
status: pending
created_at: 2025-01-27T09:00:00Z  # current UTC instant (Timestamp rule, actions/work-reference.md)
user_request: UR-006        ← new UR created for this addendum
addendum_to: REQ-005        ← links back to the original request
# assigned_to: "cloud-alpha"   # OPTIONAL advisory earmark — same contract as the Simple/Complex REQ template above: emit ONLY when the user names a session to reserve this work for (Step 1's earmark assessment), verbatim and YAML-quoted, never invented
---

# Addendum: Dark Mode Sidebar Support

## What
Add sidebar support to the existing dark mode implementation (REQ-005).

## Context
Addendum to REQ-005, which is currently [in progress / completed].
The user wants the sidebar to also support dark mode.

## Prior Implementation
[For archived/completed originals: read the original REQ from the archive and
summarize what was built, key files modified, patterns used, and commit hash
(if available). Skip this section for in-flight originals — the builder will
encounter the work in progress naturally.]

## Requirements
- Sidebar must respect the dark mode theme
```
