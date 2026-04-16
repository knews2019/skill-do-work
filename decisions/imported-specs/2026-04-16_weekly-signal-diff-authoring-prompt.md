# Task: Add a `weekly-signal-diff` prompt to the do-work prompt library

You are working in the `skill-do-work` repository. Your job is to create a
new prompt file at `prompts/weekly-signal-diff.md` that follows the existing
prompt library conventions and integrates with the skill's BKB (Build
Knowledge Base) system.

**Read before you write.** Do not skip these — they contain the conventions
you must match:

1. `prompts/README.md` — how prompt files are shaped (header, `---`
   separator, body) and how the dispatcher in `actions/prompts.md` invokes
   them. Note the idempotency and resumability expectations.
2. `actions/prompts.md` — how `do work prompts run <name>` resolves,
   adopts the body as instructions, and passes args through.
3. `actions/build-knowledge-base.md` and `actions/bkb-reference.md` — the
   BKB sub-commands, source lifecycle (`kb/raw/inbox/` → triage → ingest →
   wiki page), topic cluster semantics, and the frontmatter contract for
   ingest sources.
4. `actions/present-work.md` — how deliverables are placed under
   `do-work/deliverables/` and how artifacts are structured.
5. `crew-members/karpathy.md` and `crew-members/general.md` — always-on
   behavioral guardrails. The prompt you write must not violate these.
6. `SKILL.md` — the help menu and top-level action surface, so your prompt
   integrates naturally with the rest of the skill.

If any of those files are missing, stop and report. Do not invent conventions.

## What the prompt does

`weekly-signal-diff` turns a week of AI-industry news into a **structural
diff** — what constraints shifted, who gained or lost leverage, what
dependencies got exposed, what business-model assumptions broke —
personalized by whatever the BKB already knows about the user. It is a
weekly ritual that produces both:

- An **inline digest** rendered in the chat, consumed in place, used to
  decide whether any shifts merit a new REQ or BKB query.
- A **durable deliverable** saved to
  `do-work/deliverables/weekly-signal-diff/<week-ending-YYYY-MM-DD>.md`
  and ingested back into BKB as a source under topic cluster
  `weekly-signal-diff` so next week's run can diff against it.

Both outputs carry the same facts. The deliverable is the canonical
artifact; the inline digest is a reading-surface presentation of the same
data.

## Design constraints

- **Library style, not a new action.** This is a prompt file under
  `prompts/`, not a new top-level action. It is invoked via
  `do work prompts run weekly-signal-diff` (or the shorthand
  `do work prompts weekly-signal-diff`).
- **Idempotent and resumable.** If the user re-runs within the same
  week-ending window, detect the existing deliverable and digest, and
  offer to append revisions rather than duplicate. State persistence
  lives in the deliverable file itself — do not invent a new state
  file.
- **BKB is the memory layer.** All personalization, prior-digest
  lookup, and durable storage flow through BKB. Do not reach for any
  other memory or retrieval system; if an alternative is genuinely
  needed, that is a scope change to flag, not a silent substitution.
- **Web retrieval is optional.** When the environment has web_search
  available, use it. When it does not, accept a source packet pasted by
  the user and say the diff is source-bounded. Do not fail closed.
- **Starter universe is a bootstrap layer.** The ten default categories
  (Frontier labs, Open model ecosystem, Search and answer interfaces,
  Developer tooling and agents, Cloud AI platforms, Data and model
  infrastructure, Enterprise software incumbents, Productivity and
  knowledge tools, Creative media generation, Robotics and embodied AI)
  are the scaffold. Re-rank and extend them using BKB context before
  writing the final diff. Every loaded lane gets a real paragraph of
  scan notes every week; lanes are never compressed to a quiet-lane
  one-liner, regardless of relevance or how noisy the week was.
- **Personalization lives in a project-local sidecar.** User-specific
  lanes, toolchain lists, and example focus areas must not be inlined
  into the library prompt. The skill ships a placeholder at
  `prompts/weekly-signal-diff-personal.md` that declares no real
  lanes — it is a template, not a live config. At Phase 3 the main
  prompt searches the user's project (relative to cwd) for a file
  named `weekly-signal-diff-personal.md` at prioritized locations
  (project root → `.claude/` → `do-work/` → glob for any other
  location) and uses the first match as the live sidecar. The
  shipped placeholder is never treated as a source of real lanes,
  only as a copy-me template. A project-local personal-sidecar lane
  has the same contract as a core lane — full weekly coverage,
  per-lane scan-notes paragraph, weighted at least as heavily as any
  core lane. If no project-local copy is found, the scan runs with
  the 10 core lanes only and says so in the coverage note.
- **Diff, not digest.** The weekly output leads with whatever
  structural shifts genuinely emerged — there is no target count and
  no cap. It then renders a mandatory per-lane scan-notes section
  covering every loaded lane (10 core + any from the personal
  sidecar), followed by the "what didn't change" calibration. Reject
  benchmark drama, launch hype, and funding theater as headline
  shifts unless they move a real constraint — but still mention them
  in the relevant lane's scan notes so the reader sees what was
  considered.
- **Copyright discipline.** When citing sources, paraphrase. Do not
  reproduce article paragraphs or strings of direct quotes.

## Prompt file shape

The file must follow the structure described in `prompts/README.md`:

```markdown
# Weekly Signal Diff

> One-line description suitable for `do work prompts list`.

**Aliases:** (optional)
**When to use:** 2–3 bullets
**Inputs / flags:** describe supported args (e.g., --source-packet, --week-ending, --dry-run)

---

<procedural body — the actual instructions the agent executes>
```

Everything above the `---` is metadata the dispatcher reads for `list`
and `show`. Everything below is adopted as operational instructions when
`run` fires.

## Body structure

Follow this order. Each phase should be explicit about what to do, what
to skip, and what to tell the user.

1. **Establish frame.** Confirm the topic space (default: AI industry),
   the freshness window (default: last 7 days), and the audience lens
   (operator / investor / builder / content-prep). If the user gave no
   frame, default to a 7-day operator review and say so in the coverage
   note.
2. **Pull BKB context.** Before any retrieval, query BKB for:
   - active URs and pending REQs (via `do-work/user-requests/` and
     `do-work/queue/`) that reveal what the user is currently building
     or worried about
   - the most recent 2–4 `weekly-signal-diff` wiki pages, if any
   - recurring entities or vendors mentioned in archived REQs and prime
     files
   Use `do work bkb query` with narrow queries, not broad ones. Extract
   a short **relevance profile** — what the user is building, what they
   keep revisiting, what they are worried about. This profile shapes
   re-ranking in the next step.
3. **Build the watchlist.** Start from the 10-lane core starter
   universe. Then search the user's project (relative to cwd) for a
   file named `weekly-signal-diff-personal.md`, checking prioritized
   locations in order: project root, `.claude/`, `do-work/`, and
   finally a bounded glob of the project tree. If a project-local
   copy is found, load its lanes as full-member lanes with the same
   contract as the core 10. Do not treat the skill's shipped
   placeholder at `prompts/weekly-signal-diff-personal.md` as a
   source of real lanes — it is a copy-me template, not a live
   config. If no project-local copy is found, proceed with the 10
   core lanes only and say so in the coverage note. Re-rank using
   the relevance profile. Promote entities and sub-themes inside
   each loaded lane that the user mentions repeatedly; swap
   suggested entities for competitors, upstream suppliers, or
   downstream customers that better match active work; add a new
   lane only when a theme appears across three or more recent
   captures, REQs, or wiki pages and doesn't fit any loaded lane
   (prefer adding to the project-local sidecar over inlining into
   the library prompt). **Every loaded lane stays in full scope
   every week** —
   promotion only re-orders and re-weights inside the headline-shift
   list and the per-lane scan notes; it never shrinks any lane's
   coverage floor.
4. **Gather evidence.** If web_search is available, run it in two
   passes: a broad 7-day sweep across the lanes, then targeted
   follow-ups on the strongest candidate shifts. If web_search is
   not available, accept a pasted source packet and proceed. Cite
   sources with URLs. Paraphrase — do not quote paragraphs.
5. **Ask the structural questions** on every candidate signal:
   - What constraint shifted?
   - Who gained or lost leverage?
   - What got cheaper, harder, faster, or more defensible?
   - What dependency got exposed?
   - What business model or pricing assumption weakened?
   - What changed in regulation, geography, or distribution?
   - Why does this matter for the user's actual projects?
   Discard signals where none of these produce a meaningful answer.
6. **Score and cut.** Merge duplicates. Label speculation as
   speculation. Select the **headline structural shifts** that
   emerged this week — there is no target count and no cap; the
   number floats with how much genuinely moved. If the week was thin
   on headline shifts, say so — do not manufacture shifts.
7. **Render the digest inline.** Use this structure:
   - **Coverage note** — scope, window, how the watchlist was
     personalized, whether web_search or a source packet was used
   - **Headline structural shifts** — for each: what changed, why it
     matters in general, why it matters to this user (grounded in the
     relevance profile), sources with URLs. No fixed count.
   - **Per-lane scan notes** — one real paragraph per loaded lane
     (10 core + any from the personal sidecar), covering what was
     scanned, what moved, what didn't, and why. Lanes are never
     compressed to a one-liner.
   - **What didn't change** — 2–3 cross-cutting assumptions that held
     steady despite the noise
   - **What changed from last week** — new / rising / fading / resolved
     themes, referencing prior wiki pages by slug
   - **Watch next** — entities, constraints, or questions to monitor
   - **Actions** — optional follow-up REQs the user might capture; do
     not auto-capture them
8. **Write the deliverable.** Save the same content (with frontmatter
   suitable for BKB ingest — `topic_cluster: weekly-signal-diff`,
   `week_ending: <YYYY-MM-DD>`, `sources:` list, `created:` and
   `updated:` dates) to
   `do-work/deliverables/weekly-signal-diff/<week-ending>.md`. If the
   file already exists for the current week, append a revision section
   with a timestamp rather than overwriting — preserve history.
9. **Ingest into BKB.** Copy the deliverable to `kb/raw/inbox/` with a
   name matching BKB conventions, then tell the user to run
   `do work bkb triage` followed by `do work bkb ingest` so the page
   lands in the wiki under the `weekly-signal-diff` topic cluster. Do
   not run those commands automatically — they're the user's call,
   since ingest is a lifecycle event.
10. **Close the loop.** Print a short summary: deliverable path, number
    of shifts, whether ingest is pending, and suggested next commands
    (`do work capture request: ...` for any action items, or
    `do work bkb query "..."` to dig deeper on a specific shift).

## Arguments the prompt should accept

- `--week-ending=YYYY-MM-DD` — override the default of today. Useful
  for backfilling or re-running past weeks.
- `--source-packet=<path>` — path to a file containing pasted headlines
  or links. Use when web_search is unavailable or the user wants
  curated input.
- `--topic=<string>` — narrow the scope to a specific theme (e.g.,
  "inference economics", "enterprise SaaS repricing"). The starter
  universe still anchors the scan, but the diff weighs that theme
  heavily.
- `--dry-run` — produce the inline digest but do not write the
  deliverable or stage for ingest. Useful for previewing.
- `--no-ingest` — write the deliverable but skip the ingest hand-off
  step.

Pass any unrecognized args through untouched and mention them in the
coverage note so the user sees that they were ignored.

## Rules (enforced in the prompt body)

- Never manufacture shifts to pad the headline list. Three solid
  shifts beat seven padded ones; zero manufactured shifts beats a
  fake third.
- Never quote paragraphs. Paraphrase everything; cite URLs.
- Never auto-ingest. Stop at the hand-off.
- **Every loaded lane gets full coverage every week.** No lane is
  demoted, compressed, or reduced to a one-liner based on relevance
  or noise. Promotion only re-orders the headline-shift list;
  per-lane scan notes are mandatory for every loaded lane every week
  (10 core + however many the personal sidecar declares).
- Never collapse general analysis into personal implications without
  marking the transition. Keep "why this matters in general" and "why
  this matters to this user" visibly separated.
- Speculation is allowed but must be labeled.
- Thin weeks are thin on headline shifts. The per-lane scan notes
  section still fills every week.

## Common rationalizations to guard against

Include a short "Common rationalizations" section in the prompt body,
following the pattern in `actions/prompts.md`. At minimum:

| If you're thinking… | STOP. Instead… | Because… |
|---|---|---|
| "There aren't enough shifts, I'll include this launch announcement" | Cut the count; say the week was thin | Manufactured shifts poison the diff-over-time signal |
| "The user already knows this, I'll skip the general-audience framing" | Keep both layers | The deliverable also serves as future input to itself; the general framing is what makes cross-week comparison work |
| "The starter universe is too generic, I'll compress Robotics and Creative media to a one-liner" | Give every lane a full paragraph in the per-lane scan notes | Demotion was removed from the design — compressing lanes silently destroys the baseline scan, and structural shifts often surface in lanes the user doesn't usually track |

## Verification checklist

Before concluding your work, verify:

- [ ] `prompts/weekly-signal-diff.md` exists and follows the header +
      `---` + body structure from `prompts/README.md`
- [ ] `prompts/README.md` has a new row in the "Available prompts"
      table describing `weekly-signal-diff`
- [ ] The body references BKB commands that actually exist (check
      `actions/build-knowledge-base.md` for the current sub-command
      surface)
- [ ] The deliverable path template is
      `do-work/deliverables/weekly-signal-diff/<week-ending>.md`
- [ ] The ingest hand-off uses `kb/raw/inbox/` (confirm that path
      against the BKB reference)
- [ ] Running `do work prompts list` after your change will show the
      new prompt with a coherent one-line description
- [ ] Running `do work prompts show weekly-signal-diff` will print the
      file verbatim without executing anything
- [ ] The body is written as instructions to an agent, not as prose
      about what the prompt does. Use imperatives: "Query BKB for...",
      "Write the deliverable to...", not "This prompt queries BKB..."
- [ ] The body mandates a **per-lane scan notes** section covering
      every loaded lane every week (10 core + any from the personal
      sidecar) — no lane is missing, no lane is reduced to a one-liner
- [ ] The body searches the user's project (cwd + prioritized
      locations + bounded glob) for a `weekly-signal-diff-personal.md`
      at Phase 3 and loads its lanes when a project-local copy is
      found. The shipped placeholder at
      `prompts/weekly-signal-diff-personal.md` declares no real
      lanes. No user-specific data (lane names, project names, vendor
      names) is hard-coded into either the library prompt or the
      shipped placeholder.

## Deliverables for this task

1. The new file `prompts/weekly-signal-diff.md` (library prompt, no
   user-specific data)
2. The new file `prompts/weekly-signal-diff-personal.md` (ships as
   a placeholder template — no real lanes; users create a
   project-local copy to personalize; see the "Personalization lives
   in a project-local sidecar" design constraint)
3. An update to `prompts/README.md` adding rows for both
   `weekly-signal-diff` and `weekly-signal-diff-personal` to the
   "Available prompts" table
4. A summary commit following the repo's commit conventions — use the
   `[prompts]` label since this is library work, not tied to a
   specific REQ. Do not push.
5. A short report back to the user: file paths, the one-line
   descriptions you chose, and the exact commands they can run to
   test it (e.g., `do work prompts show weekly-signal-diff` to
   confirm the file parses, then `do work prompts run
   weekly-signal-diff --dry-run` for a live preview).

Do not modify any other files. Do not edit `SKILL.md`, `CLAUDE.md`, or
the `actions/` directory. If you think any of those need to change to
accommodate this prompt, stop and report — don't change them silently.

If anything in these instructions conflicts with the repo's own
conventions (as read from `prompts/README.md`, `actions/prompts.md`,
or the BKB action files), the repo conventions win. Tell the user
about the conflict in your final report.
