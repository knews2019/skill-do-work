# Starter Universe

Bootstrap scaffold for the weekly structural scan. These categories and
suggested entities define the baseline sweep. Re-rank and extend using
BKB context before producing the final diff — the universe is a floor,
not a ceiling.

## Rules for using this file

- Every lane below stays fully in scope every week. Every lane
  receives a real paragraph of scan notes regardless of relevance to
  the user's current work; no lane is ever compressed to a one-line
  "quiet this week" note.
- Suggested entities are examples, not a contract. Swap them for
  competitors, upstream suppliers, or downstream customers that better
  match what BKB knows the user is building.
- When a lane has two or three strong candidate shifts, don't force a
  third entity just to fill the row. One well-sourced shift beats
  three padded ones — but the per-lane scan paragraph still gets
  written.
- Preserve baseline discovery. Personalization shapes weighting; it
  does not collapse the scan into only known favorites. Structural
  shifts often show up in lanes the user doesn't normally track.

## Core AI categories (baseline scan)

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

## Personal-priority lanes (added based on active work)

> The specific lanes below reflect one user's active work as of this
> document's authoring. At runtime, this content lives in
> `prompts/weekly-signal-diff-personal.md` (the personal sidecar),
> not inlined into the main library prompt. Treat this section as
> the seed content for the sidecar, not a contract — the sidecar
> gets edited as the user's projects shift.

These lanes are scoped to the user's current projects. They are
full members of the scan, not footnotes — structural shifts here are
weighted at least as heavily as shifts in the core AI categories.

| # | Category | Suggested entities | Why this lane matters to this user |
|---|---|---|---|
| 11 | Subscription and billing platforms | Chargebee, Stripe, Recurly, Zuora, Paddle | Directly relevant to Epoch Times subscription-technology work; pricing-model and dunning-flow changes affect checkout and billing integrations |
| 12 | Shopify ecosystem and commerce tooling | Shopify platform, Mobile Buy SDK, Tapcart, Shopney, Hydrogen, Gadget | Directly relevant to dsfantiquejewelry.com; mobile app strategy, storefront API changes, checkout extensibility shifts |
| 13 | Supply chain and npm/package security | Socket, Snyk, npm registry, GitHub Advanced Security, Sigstore | Live concern after axios npm compromise investigation; affects dev environment trust model |
| 14 | Multilingual NLP and alignment tooling | OPUS-MT, NLLB, Tatoeba, sentence-transformers, spaCy | Matches Romanian-Hungarian bilingual alignment workbench; any shift in quality, licensing, or availability feeds directly into that project |
| 15 | Claude Code ecosystem and agentic harnesses | Claude Code, Dorothy, Conductor, Overstory, OpenCode, Channels API | Matches daily tooling; sandbox, proxy, and multi-agent orchestration shifts directly affect how work gets done |

## Re-ranking heuristics

Apply these during step 3 of the weekly-signal-diff prompt (watchlist
construction).

### Promote an entity or category when…

- it shows up repeatedly in the user's active URs, pending REQs, or
  recent archived work
- it affects a toolchain the user depends on (e.g., Claude Code, LM
  Studio, Shopify, Chargebee)
- it competes with, supplies, or constrains something the user is
  building
- it appeared in the last 2–4 weekly digests and has unresolved
  momentum
- BKB returns matches for it across multiple topic clusters (signals
  genuine cross-cutting relevance)

Promotion changes which entities inside a lane get scanned first and
which sub-themes get weighted most heavily in the headline structural
shifts. Promotion does **not** shrink coverage of any other lane.
Every lane — promoted or not — still produces a full paragraph of
scan notes.

### Add a new lane when…

- an entity or theme appears across three or more recent captures,
  REQs, or wiki pages but doesn't fit any existing lane
- the user explicitly asks to track something new via `--topic=`

Do not add a lane for a single mention. Wait for the pattern.

## Coverage note template

Use wording like this at the top of the weekly diff so the reader can
see how the universe was personalized this week:

> This week's scan started from the 10-lane core starter universe
> [+ N personal lanes from `prompts/weekly-signal-diff-personal.md`
> | no personal sidecar loaded]. Coverage was reweighted using BKB
> context around [focus areas — e.g., the user's current top 2–3
> active projects and critical toolchains]. Every loaded lane
> received a full scan; headline structural shifts are listed below
> and per-lane scan notes follow.

## Notes for future expansion

- If the user takes on a new major project, add a lane for it rather
  than stretching an existing one. Lanes are cheap; overloaded lanes
  are expensive.
- If a lane stays quiet for 8+ consecutive weeks, consider whether
  the project it mapped to is still active. Propose dropping it in
  the weekly output; let the user decide.
- If two lanes keep producing overlapping shifts (e.g., "Developer
  tooling and agents" and "Claude Code ecosystem" frequently cite
  the same sources), consider merging them. Do not merge
  unilaterally — flag the overlap and let the user call it.
