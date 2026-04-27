# Weekly Signal Diff

> Turn a week of AI-industry news into a structural diff — what constraints shifted, who gained or lost leverage, what dependencies got exposed — personalized via BKB. Produces an inline digest plus a durable deliverable ingested back into BKB so next week's run can diff against it.

**Aliases:** `wsd`, `signal-diff`

**When to use:**
- Weekly structural review of the AI industry — not a launch-tracker, benchmark roundup, or funding digest
- You want a diff against prior weeks' durable deliverables, not a one-shot news summary
- The output should feed BKB so momentum and calibration carry across weeks

**Inputs / flags:**
- `--week-ending=YYYY-MM-DD` — override today as the anchor (useful for backfills)
- `--source-packet=<path>` — path to a file of pasted headlines/links when `web_search` is unavailable
- `--topic=<string>` — weight a specific theme heavily (e.g., `"inference economics"`). The starter universe still anchors the scan; the theme gets extra weight on headline shifts.
- `--dry-run` — render the inline digest only; skip the deliverable write and the BKB ingest hand-off
- `--no-ingest` — write the deliverable but skip the BKB ingest hand-off

Pass any unrecognized args through untouched and mention them in the coverage note so the user sees they were ignored.

---

## Instructions for the executing agent

You are producing a weekly structural diff of AI-industry news, personalized by what BKB knows about the user. Your output has two forms: an **inline digest** rendered in chat for the user to consume in place, and a **durable deliverable** written to `do-work/deliverables/weekly-signal-diff/<week-ending>.md` with frontmatter suitable for BKB ingest. Both carry the same facts — the deliverable is the canonical artifact; the inline digest is a reading-surface presentation of the same data.

This prompt is **idempotent and resumable**. Before writing anything, check whether a deliverable already exists for the current `week-ending` date. If it does, append a timestamped revision section rather than overwriting — preserve history. State lives in the deliverable file; do not create a separate state file.

Memory layer note: **BKB is this skill's memory layer.** Personalization, prior-digest lookup, and durable storage all flow through BKB — do not reach for any other memory or retrieval system.

---

## Phase 1 — Establish frame

Confirm three things, either from the flags or by stating defaults explicitly in the coverage note:

- **Topic space** — default: the AI industry.
- **Freshness window** — default: the last 7 days, anchored to `--week-ending` (or today if the flag was omitted).
- **Audience lens** — default: operator. Acceptable alternatives if the user named one: investor, builder, content-prep.

If the user provided no frame, proceed with the defaults above and say so in the coverage note: *"Defaulting to a 7-day operator review because no frame was specified."*

## Phase 2 — Pull BKB context

Before any retrieval, query BKB for:

- active URs and pending REQs (via `do-work/user-requests/` and `do-work/queue/`) that reveal what the user is currently building or worried about
- the most recent 2–4 `weekly-signal-diff` wiki pages, if any — these ground the cross-week diff
- recurring entities or vendors mentioned in archived REQs, prime files, and past deliverables

Use `do-work bkb query` with **narrow** queries, not broad ones — one question per call, targeted at a specific concern. Examples:

- `do-work bkb query "What is the user currently shipping that depends on <vendor or platform>?"`
- `do-work bkb query "Which toolchain or supply-chain risks has the user flagged recently?"`

Extract a short **relevance profile** from the results — one paragraph covering: what the user is building, what they keep revisiting, what they are worried about. This profile shapes re-ranking in Phase 3.

If BKB is not initialized in this repo (`do-work bkb status` returns "not initialized" or `kb/` does not exist), say so explicitly in the coverage note and proceed using only the inlined starter universe. The diff is still valuable; it is just less personalized. Do not silently fail — announce the degraded state.

## Phase 3 — Build the watchlist

Start from the **10-lane core starter universe** inlined below. Then search the user's project (relative to the current working directory) for a personal sidecar named `weekly-signal-diff-personal.md`. Check these locations in order and take the first match:

1. `./weekly-signal-diff-personal.md` (project root)
2. `./.claude/weekly-signal-diff-personal.md`
3. `./do-work/weekly-signal-diff-personal.md`
4. Any other location in the project — use `Glob` (or your environment's equivalent) with pattern `**/weekly-signal-diff-personal.md`, bounded to reasonable depth and ignoring `node_modules/`, `.git/`, and other dependency directories. Take the first match.

If a project-local copy is found, read its body and treat any lanes declared there as **full members of the watchlist**, weighted at least as heavily as the core 10. Personal lanes get the same full-coverage treatment as core lanes.

If no project-local copy is found, proceed with the core 10 only and say so in the coverage note ("no personal sidecar found in the project"). The skill ships a placeholder at its own `prompts/weekly-signal-diff-personal.md` that declares no real lanes — do not treat the shipped placeholder as a source of personal lanes; only project-local copies count.

Use the relevance profile to **promote** entities and sub-themes within each loaded lane — which frontier lab to watch first, which vendor to prioritize, which entity has a live dependency that makes it a near-term concern. Promotion changes ordering and sub-theme weighting inside a lane; **it does not shrink coverage of any lane**.

**Every loaded lane gets full coverage every week.** Every loaded lane is scanned and every loaded lane produces a real paragraph in the per-lane scan notes section of the output. No lane is demoted, compressed, or reduced to a one-liner based on relevance or noise. If a lane was quiet on structural shifts, its per-lane paragraph says what was scanned, what was considered, and why nothing hit the headline bar — it still fills a paragraph.

Add a **new lane** only when: an entity or theme appears across three or more recent captures, REQs, or wiki pages but doesn't fit any loaded lane; or the user explicitly asks to track something new via `--topic=`. Prefer adding it to a project-local `weekly-signal-diff-personal.md` (created or edited by the user in their own repo) rather than inlining into this library prompt. Do not add a lane for a single mention. Wait for the pattern.

### The core starter universe (10 lanes)

These lanes are the scaffold. They ship with the library and are the same for every user. Suggested entities are examples — swap them for competitors, upstream suppliers, or downstream customers that better match what BKB knows about the user's active work.

| # | Category | Suggested entities | What this lane reveals |
|---|---|---|---|
| 1 | Frontier labs | OpenAI, Anthropic, Google DeepMind, xAI | Model capability ceilings, pricing power, safety posture, government relationships |
| 2 | Open model ecosystem | Meta AI, Mistral, Hugging Face, Qwen, DeepSeek | Open-weights quality gap, self-hosting viability, licensing shifts |
| 3 | Search and answer interfaces | Perplexity, Glean, You.com, Arc Search | Distribution shifts away from traditional search, answer-engine economics |
| 4 | Developer tooling and agents | Cursor, Replit, Cognition, Windsurf, Aider | Agentic workflow maturity, IDE disruption, developer economic capture |
| 5 | Cloud AI platforms | Microsoft Azure AI, Google Cloud, AWS Bedrock, Oracle | Hyperscaler margin games, regional availability, enterprise lock-in |
| 6 | Data and model infrastructure | Databricks, Snowflake, Together AI, Fireworks, Modal | Inference economics, training data supply, compute arbitrage |
| 7 | Enterprise software incumbents | Salesforce, Atlassian, ServiceNow, Workday | Per-seat repricing exposure, agent-driven seat compression, AI feature monetization |
| 8 | Productivity and knowledge tools | Notion, Canva, Grammarly, Linear | Consumer SaaS AI-native competition, workflow integration depth |
| 9 | Creative media generation | Runway, ElevenLabs, Pika, Midjourney, Suno | Generation economics (Sora-scale burn rates), creator tool displacement |
| 10 | Robotics and embodied AI | Figure, Wayve, Physical Intelligence, 1X | Capital allocation into physical AI, labor-substitution timelines, supply chains |

### Personal-priority lanes (from a project-local sidecar, if present)

If Phase 3's search found a project-local `weekly-signal-diff-personal.md`, read its lanes table and load each row as a full-member lane in the watchlist. A personal lane has the same contract as a core lane: full weekly coverage, a per-lane scan-notes paragraph, weighted at least as heavily as any core lane. If no project-local copy was found, skip this section and say so in the coverage note.

### Re-ranking heuristics

#### Promote an entity or category when…

- it shows up repeatedly in the user's active URs, pending REQs, or recent archived work
- it affects a toolchain the user depends on (IDE, deploy platform, billing system, search tools, etc.)
- it competes with, supplies, or constrains something the user is building
- it appeared in the last 2–4 weekly digests and has unresolved momentum
- BKB returns matches for it across multiple topic clusters (signals genuine cross-cutting relevance)

Promotion changes which entities inside a lane get scanned first and which sub-themes get weighted most heavily in the headline structural shifts. **It does not shrink coverage of any other lane.**

#### Add a new lane when…

- an entity or theme appears across three or more recent captures, REQs, or wiki pages but doesn't fit any existing lane
- the user explicitly asks to track something new via `--topic=`

Do not add a lane for a single mention. Wait for the pattern.

## Phase 4 — Gather evidence

If `web_search` is available, run two passes:

1. **Broad 7-day sweep** — one search per lane (or per promoted entity within a lane), anchored to the freshness window. Collect headlines and short summaries.
2. **Targeted follow-ups** — for each strong candidate structural shift, run follow-up searches to confirm the claim, find counter-evidence, and locate the most authoritative source. Aim for the primary source (company post, regulator filing, financial statement, benchmark paper) over downstream coverage.

If `web_search` is **not** available, accept a pasted source packet via `--source-packet=<path>`. Read it, say explicitly in the coverage note that the diff is source-bounded, and proceed from only that material. Do not fail closed.

**Copyright discipline.** Paraphrase everything. Do not reproduce article paragraphs or strings of direct quotes. Cite URLs.

## Phase 5 — Ask the structural questions

For every candidate signal, ask:

- What constraint shifted?
- Who gained or lost leverage?
- What got cheaper, harder, faster, or more defensible?
- What dependency got exposed?
- What business-model or pricing assumption weakened?
- What changed in regulation, geography, or distribution?
- Why does this matter for the user's actual projects (grounded in the relevance profile)?

Discard candidates where none of these produce a meaningful answer. A launch that moves no constraint is not a structural shift — mention it in the relevant lane's scan notes so the reader sees it was considered, but do not promote it to the headline shifts.

## Phase 6 — Score and cut

Merge duplicates (two sources reporting the same shift collapse into one entry with both URLs cited). Label speculation as speculation. Select the **headline structural shifts** — however many genuinely emerged this week. There is no target count and no cap. The shift count floats with reality.

If the week is thin on headline shifts, say so plainly. Do not manufacture shifts. Three solid shifts beat seven padded ones; zero manufactured shifts beats a fake third.

The per-lane scan notes section still fills every week, regardless of how many shifts landed — that section is the audit trail of what was scanned.

## Phase 7 — Render the digest inline

Use this structure for the inline digest (the chat output):

### Top of mind this week

Hard cap: 5 bullets, 150 words. Name the 3–5 things the operator should hold in working memory this week — the synthesis, not the detail. Everything else in the digest is support material for mid-week re-reading. If the week is thin, give fewer bullets rather than padding.

### Actions this week

Two mandatory groups. Be concrete; if a group is empty, say so explicitly — that is a finding, not a hole.

**For the operator** — 1–3 things the operator could act on this week, formatted as `do-work capture request: <short description>` so they can capture any they want to pursue. Do not auto-capture.

**For clients** — 1–3 proactive client-outreach angles: which client archetype, what finding to raise, one-line draft of the outreach. If no shift this week has a client angle, state "No client-facing actions this week — purely structural."

### Coverage note

A short opening block, built from this template:

> This week's scan started from the 10-lane core starter universe [+ N personal lanes from `prompts/weekly-signal-diff-personal.md` | no personal sidecar loaded]. Coverage was reweighted using BKB context around [focus areas — e.g., the user's current top 2–3 active projects and critical toolchains]. Every loaded lane received a full scan; headline structural shifts are listed below and per-lane scan notes follow. Evidence: [`web_search` / source packet at `<path>` / source-bounded because BKB and web_search unavailable]. Frame: [operator / investor / builder / content-prep], window: `<YYYY-MM-DD>` to `<YYYY-MM-DD>`.

If any unrecognized flags were passed, mention them here so the user sees they were ignored.

### Headline structural shifts

For each shift (no fixed count):

- **Title** — one line naming what moved.
- **What changed** — 1–2 paragraphs of paraphrased description with source URLs cited inline.
- **Why it matters in general** — a short paragraph framing the shift in industry-wide terms. This framing is what makes future cross-week comparison work, so keep it even when the user already knows the context.
- **Why it matters to this user** — a visibly separated paragraph grounding the shift in the relevance profile (active projects, toolchains, constraints). Keep the general framing and the personal framing visually distinct — never collapse them into one.
- **Sources** — bullet list of URLs.
- **Speculation tag** — if any part of the entry is speculative, prefix that sentence with `[Speculation]`.

### Per-lane scan notes (every loaded lane, always)

One real paragraph per loaded lane — every core lane and every personal-sidecar lane. Every lane gets its paragraph regardless of whether a headline shift landed there. For each lane, the paragraph covers:

- What was scanned this week (entities, sub-themes)
- What moved (if anything) and where it landed — headline shift, possible future shift, or below the bar
- What didn't move and why that absence matters (or doesn't)
- Any cross-lane connection worth flagging

A lane with no headline shift still produces a paragraph — the paragraph simply describes what was scanned and why nothing hit the headline bar. Do **not** compress any lane to a one-liner.

### What didn't change

Two or three cross-cutting assumptions that held steady despite the noise. This is the calibration section that prevents week-to-week drift.

### What changed from last week

New / rising / fading / resolved themes, referencing prior `weekly-signal-diff` wiki pages by slug (e.g., `[[weekly-signal-diff-2026-04-10]]`). If no prior digests exist, say so and note that next week will be the first true diff.

### Watch next

Entities, constraints, or questions worth monitoring over the next 1–4 weeks.

## Phase 8 — Write the deliverable

Save the same content (lightly reformatted as a Markdown document) to `do-work/deliverables/weekly-signal-diff/<week-ending>.md` with BKB-ready frontmatter:

```yaml
---
topic_cluster: weekly-signal-diff
week_ending: YYYY-MM-DD
sources:
  - https://...
  - https://...
created: YYYY-MM-DD
updated: YYYY-MM-DD
---
```

Frame, window bounds, and evidence mode live in the coverage note at the top of the body — BKB query can extract them from text when needed; they are not frontmatter fields.

**Idempotency rule.** If `do-work/deliverables/weekly-signal-diff/<week-ending>.md` already exists, do not overwrite. Append a new section titled `## Revision — <ISO 8601 timestamp>` with the fresh content and bump the `updated:` field in the frontmatter. Preserve history; the deliverable is its own state file.

Create the directory `do-work/deliverables/weekly-signal-diff/` if it doesn't exist.

If `--dry-run` is set, skip this phase entirely and tell the user the deliverable was not written.

## Phase 9 — Ingest hand-off (never auto-run)

Copy the deliverable to `kb/raw/inbox/` using a BKB-conventional name (e.g., `weekly-signal-diff-<week-ending>.md`). Then tell the user:

> Deliverable written to `do-work/deliverables/weekly-signal-diff/<week-ending>.md` and staged for BKB ingest at `kb/raw/inbox/weekly-signal-diff-<week-ending>.md`. To ingest, run:
>
> ```
> do-work bkb triage
> do-work bkb ingest
> ```
>
> Ingest is a lifecycle event — run it when you're ready.

Do **not** run those commands automatically. Ingest belongs to the user.

If `--no-ingest` or `--dry-run` is set, skip this phase and say so.

If BKB is not initialized (`kb/` does not exist), skip the copy step, say so in the close-the-loop summary, and tell the user they'd need to run `do-work bkb init` first if they want the ingest loop to work.

## Phase 10 — Close the loop

Print a short summary:

- Deliverable path (or "not written — dry run")
- Number of headline structural shifts
- Per-lane scan notes: "N / N lanes covered (10 core + M personal)" — always a full sweep of whatever lanes were loaded; this is a contract
- Ingest status: "staged in `kb/raw/inbox/`" / "skipped (flag)" / "BKB not initialized"
- Suggested next commands:
  - `do-work capture request: <...>` for any action items surfaced
  - `do-work bkb query "..."` to dig deeper on a specific shift
  - `do-work prompts run weekly-signal-diff --week-ending=YYYY-MM-DD` next week

## Rules

- Never manufacture shifts to pad the headline list. Three solid shifts beat seven padded ones; zero manufactured shifts beats a fake third.
- Never quote article paragraphs. Paraphrase everything; cite URLs.
- Never auto-ingest. Stop at the hand-off.
- **Every loaded lane gets full coverage every week.** No lane is demoted, compressed, or reduced to a one-liner based on relevance or noise. Promotion only re-orders the headline-shift list; per-lane scan notes are mandatory for every loaded lane every week (10 core + however many the personal sidecar declares).
- Keep "why this matters in general" and "why this matters to this user" visibly separated in every headline shift. Never collapse them.
- Speculation is allowed but must be labeled (`[Speculation]`).
- Thin weeks are thin on headline shifts. The per-lane scan notes section still fills every week.
- Top of mind is mandatory; hard cap enforced (5 bullets, 150 words); thin weeks produce fewer bullets, not padded ones.
- Actions section is mandatory, split operator vs. client; empty groups are stated explicitly, never omitted.

## Common Rationalizations

| If you're thinking… | STOP. Instead… | Because… |
|---|---|---|
| "There aren't enough shifts, I'll include this launch announcement" | Cut the count; say the week was thin; mention the launch in the relevant lane's scan notes if it's worth noting | Manufactured shifts poison the diff-over-time signal — future runs will try to diff against noise |
| "The user already knows this, I'll skip the general-audience framing" | Keep both layers visibly separated | The deliverable also serves as future input to itself; the general framing is what makes cross-week comparison work |
| "The starter universe is too generic, I'll compress Robotics and Creative media to a one-liner" | Give every lane a full paragraph in the per-lane scan notes | Compressing lanes silently destroys the baseline scan; structural shifts often surface in lanes the user doesn't normally track |
| "BKB isn't initialized, I'll skip BKB and pretend it worked" | Announce the degraded state in the coverage note and proceed | Silent degradation makes the user trust an output that wasn't personalized |
| "The deliverable already exists, I'll just overwrite it" | Append a timestamped revision section and bump the `updated:` frontmatter field | The deliverable is its own state file; overwriting destroys the history the diff-over-time signal depends on |

## Verification checklist (self-check before concluding)

- [ ] Deliverable exists at `do-work/deliverables/weekly-signal-diff/<week-ending>.md` with correct frontmatter (unless `--dry-run`).
- [ ] **Every loaded lane has a real paragraph in the per-lane scan notes section (10 core + any from the personal sidecar).** No lane is missing; no lane is reduced to a one-liner.
- [ ] No article paragraphs quoted verbatim; all citations are URLs + paraphrased summaries.
- [ ] Frontmatter fields populated: `topic_cluster`, `week_ending`, `sources`, `created`, `updated`.
- [ ] If a deliverable already existed for this `week-ending` date, a timestamped revision section was appended rather than overwritten.
- [ ] Coverage note at the top describes the frame, the freshness window, the personalization, and the evidence mode (web_search / source packet / source-bounded).
- [ ] "Top of mind this week" sits as the first subsection of the digest, holds 3–5 bullets, and stays under 150 words.
- [ ] "Actions this week" sits between "Top of mind" and "Coverage note", with both operator and client groups present (empty groups stated explicitly, not omitted).
- [ ] Ingest hand-off happened (file copied to `kb/raw/inbox/`) unless `--no-ingest`, `--dry-run`, or uninitialized BKB.
