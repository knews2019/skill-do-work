---
title: "ADR-016: Vendor the queue-kanban Go Tool Into the Skill"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-07-01
sources:
  - tools/queue-kanban/ (imported from g1w-game-find-the-difference @ 1.1.0)
  - actions/board.md (new action)
  - actions/version.md (update-path shipped-paths + diff exclusion)
  - README.md (Queue board section)
  - CLAUDE.md (Shipped Tooling section)
  - .gitattributes (export-ignore set — tools/ deliberately excluded from it)
related:
  - page: adr-013-harden-the-vendored-skill-distribution-model
    rel: extends
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: complements
created: 2026-07-01
updated: 2026-07-01
confidence: high
---

# ADR-016: Vendor the queue-kanban Go Tool Into the Skill

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-013-harden-the-vendored-skill-distribution-model]] (extends), [[adr-001-modular-action-prompts-and-companion-references]] (complements)

## Context

`queue-kanban` is a standalone Go module (~3,000 LOC, deps `goldmark` + `yaml.v3`, embedded `web/` frontend) that walks a repo's version-controlled `do-work/` Markdown tree and renders it as a Kanban board + completion calendar (`summary` / `generate` / `serve`). It was developed *inside consumer repos* — its Go module path was already `github.com/knews2019/skill-do-work/queue-kanban`, and it shipped its own `INSTALL.md` (a copy-into-another-repo prompt) and its own semver `CHANGELOG.md`, explicitly "independent of the do-work skill's own version" so divergent per-repo copies could reconcile by diffing that changelog.

That independence became the problem. The tool lived per-consumer and drifted from the skill on a separate update channel: two copies already existed on one machine (`g1w-game-find-the-difference` @ 1.1.0, `sa2-sentence-aligner2` @ 1.0.0, stale), and the skill's own updates (`update skill 0.93.1 → 0.99.5`) never touched the tool. Worse, the tool's parser mirrors the skill's **Schema Read Contract** (`actions/work-reference.md`: the `status`/`depends_on`/`domain` enums it buckets tickets on), but the parser and the contract lived in different repos — so a schema change and its parser update could not land together, and column-bucketing drift was invisible until a board rendered wrong.

Meanwhile the skill's distribution is already whole-tree: install and `do-work update` are a `git archive` tarball (`curl …/main.tar.gz | tar xz`) that ships **every tracked file not `export-ignore`'d** (ADR-013). So co-locating the tool needs no new pipe — committing it under a non-excluded path ships it automatically.

## Decision

**Vendor `tools/queue-kanban/` into the skill repo as shipped source, drive it with a new `do-work board` action, and fold its versioning into the skill's.** Concretely:

- **Ships as source, built on demand.** The `*.go`, `go.mod`/`go.sum`, and embedded `web/` are committed under `tools/` (deliberately *not* `export-ignore`'d, so the tarball carries them). No `vendor/` tree and no compiled binary ship; the first `go build` fetches the two deps from the module proxy. The binary stays gitignored by the tool's nested `tools/queue-kanban/.gitignore`, which itself ships and keeps the binary ignored in every consumer.
- **`do-work board` is the entry point.** `actions/board.md` locates the tool under the skill dir, precondition-checks the Go toolchain (degrading gracefully when absent), builds, and dispatches `serve` (live at `:8090`) / `static` / `summary` with an explicit `--repo-root`. It is a read-only queue viewer.
- **Versioning folds into the skill.** The tool's independent `CHANGELOG.md` and `INSTALL.md` are removed; future tool changes get root `CHANGELOG.md` entries and normal skill version bumps. This ADR preserves the retired changelog (appendix below) for provenance.

## Alternatives

1. **Leave the tool per-consumer, keep installing via its own `INSTALL.md`.** Rejected — this is the status quo that produced divergent copies and a parser that drifts from the schema it depends on. The user's explicit goal was "the kanban go code is also deployed on do-work version update."
2. **Ship a prebuilt binary instead of source.** Rejected — platform-specific, bloats the text tarball, and defeats the `git archive` distribution model. Source + on-demand `go build` keeps the tarball small and portable.
3. **Vendor the Go deps (`go mod vendor`) for hermetic offline builds.** Considered and declined for now — it adds third-party Go source to every skill tarball for a benefit (offline first-build) that rarely bites, since a machine that has built the tool once has the deps cached. Revisit if consumers report proxy-less environments.
4. **Keep the tool's independent semver changelog after vendoring.** Rejected — once the skill is the single upstream, divergent-copy reconciliation (the whole reason for the separate changelog) no longer applies. One version to bump is simpler and matches every other shipped file.

## Consequences

A `do-work update` now carries the latest board into every consumer alongside the skill files, from one upstream — the drift and the stale second copy are designed out. The parser and the Schema Read Contract it tracks now live in one repo, so a schema change and its `model.go` update land in the same commit (CLAUDE.md's "Shipped Tooling" section records this lock-step as a rule). Costs: the skill is no longer pure markdown+shell — `do-work board` is the one action needing a compiler (Go), mitigated by a graceful precondition check that never blocks the rest of the skill; and the tarball grows by ~3,000 LOC of Go source. The update path was hardened to match (`actions/version.md` adds `tools/` to its shipped-paths dirty-check and excludes the gitignored `queue-kanban` binary from the fresh-upstream diff so it isn't flagged as a phantom customization). Follow-up outside this repo: consumer repos with a root-level `tools/queue-kanban/` should delete it and drive the board via `do-work board`, and the external skill-sync tool behind `skills-lock.json` should be confirmed to mirror the full tree (the plain `curl | tar` path already does).

## Appendix — Retired `tools/queue-kanban/CHANGELOG.md` (provenance)

The tool was independently versioned through 1.1.0 before being vendored in. Preserved here since the standalone changelog was removed:

- **1.1.0 — Board names its project (2026-06-30)** — stamps the project name (repo-root folder) into the page title/header so a tab or screenshot is self-identifying.
- **1.0.0 — The Portable Board (2026-06-30)** — first release packaged to install into any repo with a `do-work/` folder; added `INSTALL.md`, renamed `kanban-changelog.md` → semver `CHANGELOG.md`, made the repo-agnostic `go build` primary.
- **0.4.1 (2026-06-29)** — RECENTLY DONE window defaults to 24h instead of 7d.
- **0.4.0 (2026-06-29)** — live `serve` subcommand (re-walks the tree per request, rebuilds on mtime change); `generate` externalizes board JSON to a sibling `board-data.js` so the static artifact renders from `file://` with zero network.
- **0.3.1 (2026-06-29)** — drawer focus restored on REQ→UR navigation.
- **0.3.0 (2026-06-29)** — justfile recipes + status-based (not count-based) live-tree test.
- **0.2.1 (2026-06-29)** — documented in the backend prime.
- **0.2.0 (2026-06-29)** — system light/dark theme via `prefers-color-scheme`.
- **0.1.0 — MVP (2026-06-29)** — parser + board model + static `generate`; columns mirror the real do-work status vocabulary (Pending / Claimed / Needs-input-or-Blocked / Recently-done) plus a completion calendar. The versioned Markdown files are the database; the Go parser is the tool — no SQLite.

## References

- [tools/queue-kanban/](../../tools/queue-kanban/) — the vendored Go module (parser in `model.go`/`walk.go`, `//go:embed web/` frontend)
- [actions/board.md](../../actions/board.md) — the `do-work board` action
- [actions/version.md](../../actions/version.md) — update-path shipped-paths + binary diff exclusion
- [actions/work-reference.md](../../actions/work-reference.md) — the Schema Read Contract the parser must stay in lock-step with
- [CLAUDE.md](../../CLAUDE.md) — "Shipped Tooling (`tools/`)" conventions
- [[adr-013-harden-the-vendored-skill-distribution-model]] — the tarball/export-ignore distribution model this extends
