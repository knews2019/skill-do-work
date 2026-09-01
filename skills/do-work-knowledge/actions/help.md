# Knowledge Help Action

> Prints the knowledge package menu and stops without reading or writing a knowledge store.

```text
do-work-knowledge — retained knowledge and memory workflows

  do-work-knowledge bkb [subcommand]       Build, query, lint, and maintain a compiled KB
  do-work-knowledge memory remember <text> Curate a fact into working memory
  do-work-knowledge memory forget <text>   Confirmed store removal and log redaction
  do-work-knowledge memory recall <query>  Attributed layered recall
  do-work-knowledge memory status          Memory engine health
  do-work-knowledge memory bootstrap       Consent-gated history import
  do-work-knowledge memory audit           Compare memory and BKB value
  do-work-knowledge dream [path]            Consolidate and prune a memory/wiki tree
  do-work-knowledge interview [template]   Structured operating-model elicitation
  do-work-knowledge prompts [subcommand]   Browse or run the shipped prompt library
  do-work-knowledge setup-memory           Explicitly scaffold memory and enable optional hooks
```

Fresh suite installs leave memory hooks disabled. `setup-memory` is the only enabling route.

Deterministic BKB, Dream, interview, and memory phases are directly runnable as flat Just recipes. Run `just --list` for the live inventory; the natural-language workflows retain semantic judgment and consent.
