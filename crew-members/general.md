# The Compass — General Crew Member

<!-- JIT_CONTEXT: This file is always loaded during implementation (Step 6) regardless of domain. It provides foundational conventions (PRIME file discipline, code hygiene, commit etiquette) that apply to every REQ. No domain tag gates it. -->

## PRIME Files Philosophy

When asked to read or create a Prime file (`prime-*.md`), adhere strictly to these rules:

- **Purpose:** Prime files are semantic indexes for a specific utility or folder. They prevent the AI from having to reinvestigate the entire architecture from scratch.
- **Low Noise, High Value:** Keep them concise.
- **Pointers, Not Copies:** Do not copy-paste large blocks of code into the prime file. Instead, point to the code (`See src/utils/parser.ts for the core regex loop`) which acts as the source of truth.
- **No Volatile Metrics:** Do NOT include volatile data like test counts, exact line numbers, or pending invoice totals. These go stale immediately and create noise.
- **Multiple Aspects:** It is perfectly valid to have multiple prime files in the same folder if they describe different aspects (e.g., `prime-checkout-speed.md` and `prime-checkout-consolidation.md`).
