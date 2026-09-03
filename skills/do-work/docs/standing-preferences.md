# Standing Preferences — Nudges You Don't Need to Type

> **Part of the do-work skill.** A reference for the operating nudges people habitually paste at the start of a run — "keep writing lessons learned," "keep working until the queue is done," "commit often," "I'm AFK, don't block on questions." Most are already the skill's defaults. This page maps each common nudge to where that behavior already lives, so you can stop re-typing it.

If you find yourself pasting the same reminders every session, read this once. For each common nudge the table says whether it's already the default, how it works, and — for the two that are **deliberately not** the default — why.

## The map

| You often type... | Already built in? | How it works / where |
| --- | --- | --- |
| "when you see something, say something" / "keep writing lessons learned with improvement suggestions" | **Yes** | Every non-trivial REQ appends a `## Lessons Learned` section before archiving (`actions/work.md` Lessons-Capture Phase; standalone mode via `actions/review-work.md`). Out-of-scope things the builder notices go to `## Discovered Tasks`; only `impact-critical` findings auto-queue, while every noncritical line ends `→ report only` in the archived REQ (`crew-members/general.md` Discovered-Tasks Contract). Promote one explicitly with `do-work capture`, quoting the full finding line. For a free-form jot, `do-work-toolbox note "…"` (`../../do-work-toolbox/actions/note.md`). |
| "keep working until all the REQ's are done" | **Partly — bounded on purpose** | `do-work run` loops through every dependency-ready `pending` REQ, one at a time, until none remain (`actions/work.md` Loop-or-Exit). It exits when the current ready set is exhausted rather than running an unattended, infinite re-drain. A separate "loop until empty forever" runner was considered and **declined**: see `decisions/records/adr-014-*`. |
| "use background agents / workflows to manage context" | **Yes (when your tool supports it)** | Actions that fan work out use the durability pattern in `crew-members/background-agents.md`; `SKILL.md` dispatches `work` and `cleanup` to the background when subagents are available; each REQ gets a fresh agent and a context wipe between iterations (`actions/work.md` Loop-or-Exit). Subagents are a nice-to-have, never a requirement. |
| "it's safe to commit / `do-work commit` as a background agent" | **Commit: yes. Background: no, by design** | `do-work run` commits each finished REQ itself (`actions/work.md` Commit Phase); `do-work commit` batches loose changes into small atomic commits (`actions/commit.md`). Commit runs in the **foreground** on purpose — only `work` and `cleanup` are background-eligible in `SKILL.md`'s dispatch. |
| "I prefer frequent, well-defined commits" | **Yes** | The convention is **one atomic commit per REQ**, staging explicit files only — never `git add -A`/`.` (`actions/work.md` Commit Phase; `actions/commit.md`). That already is "frequent and well-defined." It was requested once as `UR-003` and folded into this convention. |
| "follow YAGNI" | **Yes** | The Simplicity-First guardrail is loaded on every build (`crew-members/coding-guardrails.md` §2 — the canonical home). It's stated once and referenced narrowly; over-repeating it would itself violate YAGNI. |
| "no cryptic or single-word names" / "names must be greppable" | **Yes** | The Naming-for-Reach guardrail is loaded on every build (`crew-members/coding-guardrails.md` §5 — the canonical home) and spot-checked at review (`actions/review-work.md` Coding-Guardrails Principle Check). It scopes the two-word minimum to names with *reach* and exempts idiomatic short locals, so it doesn't fight §3's match-existing-style rule. Language-specific vocabulary (which short locals your stack considers idiomatic) belongs in your own project's agent instructions — that layers on top of this guardrail and doesn't need to restate it. |
| "I'm AFK — turn questions into pending-answer, do as much as you can without getting blocked" | **Yes** | Builders never block on ambiguity. Each open question is recorded as a best-judgment decision (`- [~]` with reasoning), the builder keeps building every `pending` REQ, and a follow-up REQ with `status: pending-answers` is filed for you to review later via `do-work clarify` (`actions/work.md` Open-Questions step). A question whose real answerer is a named outside stakeholder skips your clarify queue entirely: it lands on that person's stakeholder REQ with a shareable HTML report, and `do-work stakeholder-answers` ingests their reply. Note the real status is spelled **`pending-answers`** (plural). |

## Two things worth knowing

- **"See something, say something" means capture — not fix.** When the builder notices an
  unrelated bug or bit of tech debt mid-task, the rule is to **log it in `## Discovered
  Tasks`, classify it, and end a noncritical line `→ report only` — not fix it inline or
  queue it automatically** (`crew-members/general.md` Discovered-Tasks Contract). Inline fixes blow the REQ's declared scope and inflate the
  diff so it can't be reviewed safely. If you actually want trivial adjacent issues fixed
  in place, that is a deliberate change to this guardrail — say so explicitly; it isn't the
  default.
- **The queue continuation is bounded on purpose.** "Keep going until it's all done" is
  honored *within* a run for every dependency-ready `pending` REQ, but there is no unattended
  loop that drains and re-drains newly generated work forever. That was a design decision
  (`adr-014`), not an oversight — the boundary keeps every REQ reviewed instead of racing to
  empty the board.

## So do I still need to paste my nudges?

Mostly no. Lessons-learned capture, discovered-tasks, YAGNI, per-REQ atomic commits,
background agents, and don't-block-on-questions are already the defaults. The only two you
*can't* get just by asking are an unbounded auto-drain and a backgrounded commit — and
those are deliberately not how the skill runs. Changing either is an ADR-level decision,
not a per-run reminder.
