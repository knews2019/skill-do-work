# Starter Universe

Bootstrap scaffold for the weekly structural scan. These categories and
suggested entities define the baseline sweep. Re-rank and extend using
BKB context before producing the final diff — the universe is a floor,
not a ceiling.

## Rules for using this file

- All categories below stay in scope every week. Demote, don't drop.
  A quiet lane gets a one-line "quiet this week" note; it does not
  disappear from the scan.
- Suggested entities are examples, not a contract. Swap them for
  competitors, upstream suppliers, or downstream customers that better
  match what BKB knows the user is building.
- When a lane has two or three strong candidate shifts, don't force a
  third entity just to fill the row. One well-sourced shift beats
  three padded ones.
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

### Demote — but do not drop — an entity or category when…

- it has low connection to the user's current work
- it generates plenty of headlines but the shifts are benchmark drama
  or launch hype rather than structural change
- the user's niche clearly lives elsewhere this week

Demotion means fewer entities scanned and a lower word count in the
output. It does not mean silence — "Robotics and embodied AI: quiet
this week, no structural shifts detected" is a valid line.

### Add a new lane when…

- an entity or theme appears across three or more recent captures,
  REQs, or wiki pages but doesn't fit any existing lane
- the user explicitly asks to track something new via `--topic=`

Do not add a lane for a single mention. Wait for the pattern.

## Coverage note template

Use wording like this at the top of the weekly diff so the reader can
see how the universe was personalized this week:

> This week's scan started from the 15-category starter universe
> (10 core AI + 5 personal-priority lanes). Coverage was reweighted
> using BKB context around [focus areas — e.g., "Chargebee billing,
> Claude Code tooling, Shopify mobile strategy"]. Lanes with no
> detected structural change this week: [list].

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
