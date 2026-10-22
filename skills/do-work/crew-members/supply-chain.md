# The Quartermaster — Supply Chain Crew Member

<!-- JIT_CONTEXT: Loaded by `actions/work.md` Step 4 (the Route C plan agent) and Step 6 (substep 4b) when the REQ creates a new service, app, package, or deployable; adds, removes, or bumps a dependency; or touches install, build, container-image, or CI-install configuration. The trigger is heuristic — lean toward loading: the cost of an unneeded checklist is low, the cost of an unreviewed dependency is a supply-chain compromise. Also loaded by `../../do-work-toolbox/actions/code-review.md` Step 5 when the scoped files include a package manifest, lockfile, Dockerfile, or CI install step (and this file exists). This comment is the canonical statement of when this file loads; the loaders' conditions restate it. -->

> Fewer dependencies, no install-time code execution, minimal runtime images, with Go as the default for new services.

## Core Principle: Review Attention Is the Scarce Resource

Code is cheap now. A slug library or a retry helper got imported in 2018 because writing it cost a day; today it costs ten minutes and every line can be read. 200 lines you generated and reviewed beat 20,000 you imported and did not. Every dependency is code that runs with the app's privileges, on every developer machine at install time, and in every build.

## When Loaded at Planning

Architecture is decided in the plan, not in the diff. A Route C plan written with this file loaded names, under its architectural decisions:

1. Language and runtime, with the decision tier from **New Services Default to Go** below.
2. Every direct dependency, each with the reason it is kept and the alternative rejected (written in-repo, standard library, copied in).
3. The install method and the lockfile that pins it.
4. The final image base.
5. The SBOM step, or the reason none is needed.

Implementation then executes that list. A dependency that is not on it is a plan change, recorded as a DECIDE & STATE line (`crew-members/coding-guardrails.md` § The decide-vs-escalate gate), never a silent addition.

## Write Before You Import

Before adding a dependency, estimate the part you actually need. If you can write it and read it in one sitting, write it. Keep the serious ones: cryptography, TLS, parsers for formats you did not define, database drivers. In an existing PHP, Python, or Node app the REQ touches, cut dependencies instead of rewriting: remove trivial packages, copy small ones in, keep the serious ones. Record every added dependency as a DECIDE & STATE line naming the rejected alternative.

## New Services Default to Go

When the job is HTTP + JSON + database, a new service is written in Go. The standard library is the security argument, not the binary: `net/http`, `crypto/tls`, `encoding/json`, `database/sql`, `html/template`, and `crypto` cover the job with no install hooks, a checksum database that catches tampering, and minimum version selection that never silently upgrades.

- **Budget:** the standard library plus one database driver, plus a migrations tool if the project needs one. Anything else is written in-repo or escalated with value and risk.
- **Decision tier:** greenfield or polyglot project → DECIDE & STATE. Project that is single-stack in another language → ESCALATE: the value is the dependency count and the one-artifact deploy; the risk is a second runtime the team has to operate.
- **Deliverable:** a static binary. One artifact: hash it, sign it, put it in a scratch container. No interpreter to patch, no `node_modules` on the box, no `require()` of strings at runtime.

## Frontends Stay Hardened

npm is allowed, hardened:

- A lockfile with integrity hashes is committed; installs resolve from it, never from version ranges.
- Every install runs script-free (see **Installs Run No Code** below).
- One framework. No plugin sprawl: a date formatter, slug helper, or class-name joiner is written, not installed.
- A new version is adopted a few days after it is published, unless the bump is a security fix.
- The runtime image is distroless Node, or static files served by the backend.

## Installs Run No Code

Your own install commands follow the always-on rule in `crew-members/general.md` § Installs Run No Code. For the project's install, build, and CI steps the conditions are:

- A lockfile with integrity hashes is committed, and the install step refuses to resolve without it.
- The install step skips lifecycle scripts; a package that truly needs its build step gets that one build run explicitly and named in the hand-back.
- Prebuilt artifacts only: wheels, not source distributions; no compile-on-install.
- Cooldown: a new version is adopted a few days after release, unless the bump is a security fix.

Illustrative, not exhaustive: `npm ci --ignore-scripts`, `pip install --require-hashes --only-binary=:all: -r requirements.txt`, `composer install --no-scripts`, `go mod verify`. The condition is the rule; the package manager's own flag for it is what you look up.

## Runtime Image Is Minimal

The final image stage has no shell, package manager, or compiler. A Go binary goes in `scratch` or a distroless base. Interpreted apps get the same win: distroless images exist for Node and Python. A multi-stage build keeps the toolchain in the build stage only.

## Know What's Inside

Every build that produces an image or binary emits an SBOM (a dependency inventory) as a build artifact. A Go binary carries its own: `go version -m <binary>` prints the modules compiled in. Use the project's existing SBOM tool when one exists; otherwise choose one that reads the lockfile and record the choice in the plan.

## Scope Guard

Apply these rules to the surfaces the REQ touches. A project-wide gap the REQ did not touch (no lockfile anywhere, no SBOM step, a fat base image on an untouched service) goes to Discovered Tasks, not into this REQ.

## Quality Checks

Before marking UNIFY complete, verify:

| Criterion | What to check |
|-----------|---------------|
| Every direct dependency justified | Each addition has a reason and a rejected alternative recorded; nothing added "for later" |
| Lockfile current | The committed lockfile and its hashes match the manifest; no floating ranges resolved at install |
| Script-free install | The project's install step skips lifecycle scripts; any explicit per-package build is named |
| Prebuilt artifacts | No compile-on-install; wheels or binaries, not source distributions |
| Minimal final image | The runtime stage has no shell, package manager, or compiler |
| SBOM present | A build that produces an image or binary emits an inventory, or the plan says why not |
| Plan matches shipped | For a new service, the language and dependency list in the plan are what shipped |
