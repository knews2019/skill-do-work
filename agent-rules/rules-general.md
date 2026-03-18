# General Rules

## PRIME Files Philosophy

A prime file (`prime-<name>.md`) is a semantic index that lives in a utility or folder's root directory. It maps the architecture so AI agents don't reinvestigate from scratch every session. Think of it as a table of contents for a codebase area — it tells you what exists, where, and why.

Prime files always follow the `prime-*.md` naming convention. Projects typically register their prime files in CLAUDE.md (or equivalent root instructions file). To find relevant primes: check CLAUDE.md first, then glob for `prime-*.md` in the directories you're working in.

When reading or creating a prime file, adhere strictly to these rules:

- **Low Noise, High Value:** Keep them concise. A prime file that's too long defeats its purpose.
- **Pointers, Not Copies:** Do not copy-paste large blocks of code into the prime file. Instead, point to the code (`See src/utils/parser.ts for the core regex loop`) which acts as the source of truth.
- **No Volatile Metrics:** Do NOT include volatile data like test counts, exact line numbers, or pending invoice totals. These go stale immediately and create noise.
- **Multiple Aspects:** It is perfectly valid to have multiple prime files in the same folder if they describe different aspects (e.g., `prime-checkout-speed.md` and `prime-checkout-consolidation.md`).
- **Relative Paths:** Links in prime files use relative paths from the prime's own directory — never absolute paths.

### Satellite Docs

Prime files are the entry point. Satellite docs live alongside the prime and extend it:

- `implementation-history.md` — REQ traceability index organized by feature area
- `known-bugs-<name>.md` — bug documentation with code locations
- `lessons-learned/<topic>.md` — reusable lessons from past sessions

Every satellite doc must be linked from its prime file. If a satellite doc isn't linked from the prime, future sessions won't find it.

### Lessons Section

When a REQ captures lessons learned that are relevant to a prime file's domain, the prime file gets a `## Lessons` section with a link back to the archived REQ:

```markdown
## Lessons

- [REQ-042: Brief description of lesson](relative/path/to/do-work/archive/UR-NNN/REQ-042-slug.md#lessons-learned)
```

The link path must be relative from the prime file's directory. This is how institutional knowledge flows back into the codebase — the prime file becomes the durable index, and the archived REQ holds the detail.
